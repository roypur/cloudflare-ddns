package main

import (
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "strings"
    "sync"
    "math/big"
    "github.com/vishvananda/netlink"
    "errors"
)

const ipEndpoint = "https://icanhazip.com/"

const IP6_ADDR_LENGTH = 16

var ip6 net.IP
var ip4 net.IP
var success bool

var mutex sync.Mutex
var wg sync.WaitGroup

func main() {
    var conf Config = getConfig()

    var resp *http.Response
    var err error

    if conf.LocalNetwork {
        var v4Routes []netlink.Route
        var v6Routes []netlink.Route

        externV4 := net.ParseIP("1.2.3.4")
        externV6 := net.ParseIP("2a02::1")

        v6Routes, err = netlink.RouteGet(externV6)

        if (err == nil) && len(v6Routes) > 0 {
            ip6 = v6Routes[0].Src
        } else {
            fmt.Println("IPv6 address not found!")
        }

        v4Routes, err = netlink.RouteGet(externV4)

        if (err == nil) && len(v4Routes) > 0 {
            ip4 = v4Routes[0].Src
        } else {
            fmt.Println("IPv4 address not found!")
        }
    } else {
        resp, err = GetV6(ipEndpoint)

        if err == nil {
            content, err := ioutil.ReadAll(resp.Body)
            if err == nil {
                ip6 = net.ParseIP(string(trim(content)))
            } else {
                fmt.Println("IPv6 address not found!")
            }
        }

        resp, err = GetV4(ipEndpoint)

        if err == nil {
            content, err := ioutil.ReadAll(resp.Body)
            if err == nil {
                ip4 = net.ParseIP(string(trim(content)))
            } else {
                fmt.Println("IPv4 address not found!")
            }
        }
    }

    records, err := getRecords(conf.ApiEmail, conf.ApiKey, conf.Domain)

    var hasIp6 bool = (ip6 != nil)
    var hasIp4 bool = (ip4 != nil)

    if err == nil {
        if hasIp4{
            success = true
            wg.Add(len(conf.Ipv4))
        }
        if hasIp6{
            success = true
            wg.Add(len(conf.Ipv6))
            for hostName, host := range conf.Ipv6 {
                go findAndUpdate(records.Result, hostName, conf, host)
            }
        }
        if hasIp4{
            for _,host := range conf.Ipv4 {
                var empty Host
                go findAndUpdate(records.Result, host, conf, empty)
            }
        }
        wg.Wait()
        if success {
            fmt.Println("ddns-update: SUCCESS")
        } else {
            fmt.Println("ddns-update: FAIL")
        }
    } else {
        fmt.Println("Failed communicating with cloudflare")
    }
}

func findAndUpdate(records []Record, hostName string, conf Config, host Host) {
    defer wg.Done()
    var err error
    var isIp4 bool

    recType := "AAAA"
    if len(strings.TrimSpace(host.Addr)) == 0{
        isIp4 = true
        recType = "A"
    }

    for _, rec := range records {

        //lower case fqdn
        var fqdn string = strings.ToLower(strings.TrimSpace(hostName) + "." + strings.TrimSpace(conf.Domain))

        var recordName string = strings.TrimSpace(strings.ToLower(rec.Name))

        if (recordName == fqdn) && (rec.Type == recType) {
            if isIp4{
                rec.Content = ip4.String()
            }else{
                rec.Content, err = joinIP(host)
                if err != nil{
                    fmt.Printf("ERROR: %s\n", err)
                    mutex.Lock()
                    success = false
                    mutex.Unlock()
                    return
                }
            }
            err = update(conf.ApiEmail, conf.ApiKey, rec)

            if err != nil {
                mutex.Lock()
                fmt.Println("Failed to update " + recType + " record " + fqdn)
                success = false
                mutex.Unlock()
            }
        }
    }
}

func joinIP(host Host) (string, error) {

    var tmp string = host.Addr

    suffix := ""
    var addr []byte

    hostPrefixLength := host.HostPrefixLength

    if host.IsMac{
        var err error
        // remove colons and insert fffe in the middle
        for _, v := range tmp {

            if isHex(v) {
                suffix += string(v)
            }

            if len(suffix) == 6 {
                suffix += "fffe"
            }
        }

        addr, err = hex.DecodeString(suffix)

        if err == nil {
            addr[0] = addr[0] ^ 2
        } else {
            return "", errors.New("Invalid mac-address")
        }
        hostPrefixLength = 64
    }else{
        addr = net.ParseIP(host.Addr)
    }

    if host.PrefixLength > host.HostPrefixLength{
        return "", errors.New("prefix-length > host-prefix-length")
    }

    bigPrefix := big.NewInt(0)
    bigPrefix.SetBytes(ip6)

    for i:=0; i<(128-host.PrefixLength); i++{
        bigPrefix.SetBit(bigPrefix, i, 0)
    }

    bigHostAddr := big.NewInt(0)
    bigHostAddr.SetBytes(addr)

    for i:=(128-hostPrefixLength); i<128; i++{
        bigHostAddr.SetBit(bigHostAddr, i, 0)
    }

    bigIP := big.NewInt(0)
    bigIP.Or(bigHostAddr, bigPrefix)

    localPrefix := strings.TrimSpace(host.HostPrefix)

    if (len(localPrefix) % 2) == 1{
        localPrefix = "0" + localPrefix
    }

    local, err := hex.DecodeString(localPrefix)

    if err != nil{
        return "", errors.New(fmt.Sprintf("Invalid prefix-id: %s\n", localPrefix))
    }

    bigLocalPrefix := big.NewInt(0)
    bigLocalPrefix.SetBytes(local)
    bigLocalPrefix.Lsh(bigLocalPrefix, uint(hostPrefixLength))

    for i:=0; i<host.PrefixLength; i++{
        bigLocalPrefix.SetBit(bigLocalPrefix, i, 0)
    }

    bigIP.Or(bigLocalPrefix, bigIP)
    tmpBytes := bigIP.Bytes()

    ipBytes := make([]byte, IP6_ADDR_LENGTH)

    tmpLength := len(tmpBytes)
    for k,_ := range tmpBytes{
        ipBytes[IP6_ADDR_LENGTH-k-1] = tmpBytes[tmpLength-k-1]
    }

    return net.IP(ipBytes).String(), nil
}

func isHex(r rune) bool {
    // Upper and lower case HEX ascii values
    return (((r < 58) && (r > 47)) || ((r > 64) && (r < 71)) || ((r > 96) && (r < 103)))
}
