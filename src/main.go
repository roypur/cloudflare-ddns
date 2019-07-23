package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"
)

const IP_LOOKUP_ENDPOINT = "https://icanhazip.com/"

const IP6_ADDR_LENGTH = 16
const TIMEOUT = 5

const IP4_LOCAL_FLAG uint =  2
const IP6_LOCAL_FLAG uint =  4
const IP4_REMOTE_FLAG uint = 8
const IP6_REMOTE_FLAG uint = 16

var localIP4 net.IP
var localIP6 net.IP
var remoteIP4 net.IP
var remoteIP6 net.IP

var success bool

var mutex sync.Mutex
var wg sync.WaitGroup

func main() {
    var conf Config = getConfig()
    if conf.Interval > 0 {
        for {
            time.Sleep(time.Duration(conf.Interval) * time.Second)
            loop(conf)
        }
    } else {
        loop(conf)
    }
}

func loop(conf Config) {
    var err error
    mask := uint(0)

    localIP4, err = getLocalAddr(4)
    if err == nil {
        mask |= IP4_LOCAL_FLAG
    } else {
        fmt.Println("Can't find local IPv4")
    }

    localIP6, err = getLocalAddr(6)
    if err == nil {
        mask |= IP6_LOCAL_FLAG
    } else {
        fmt.Println("Can't find local IPv6")
    }

    remoteIP4, err = getRemoteAddr(4)
    if err == nil {
        mask |= IP4_REMOTE_FLAG
    } else {
        fmt.Println("Can't find remote IPv4")
    }

    remoteIP6, err = getRemoteAddr(6)
    if err == nil {
        mask |= IP6_REMOTE_FLAG
    } else {
        fmt.Println("Can't find remote IPv6")
    }

    records, err := getRecords(conf.ApiEmail, conf.ApiKey, conf.Domain)

    if err == nil {
        for _,host := range conf.Ipv4 {
            if (host.mask & mask) > 0 {
                success = true
                wg.Add(1)
            }
        }
        for _,host := range conf.Ipv6 {
            if (host.mask & mask) > 0 {
                success = true
                wg.Add(1)
            }
        }
        for hostName, host := range conf.Ipv4 {
            if (host.mask & mask) > 0 {
                go findAndUpdate(records.Result, hostName, conf, host)
            }
        }
        for hostName, host := range conf.Ipv6 {
            if (host.mask & mask) > 0 {
                go findAndUpdate(records.Result, hostName, conf, host)
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

    ip4b := make([]byte, IP6_ADDR_LENGTH)
    ip6b := make([]byte, IP6_ADDR_LENGTH)

    if host.LocalMode {
        for k,v := range localIP4 {
            ip4b[k] = v
        }
        for k,v := range localIP6 {
            ip6b[k] = v
        }
    } else {
        for k,v := range remoteIP4 {
            ip4b[k] = v
        }
        for k,v := range remoteIP6 {
            ip6b[k] = v
        }
    }
    ip4 := net.IP(ip4b)
    ip6 := net.IP(ip6b)

    recType := "AAAA"
    if (host.mask & (IP4_LOCAL_FLAG | IP4_REMOTE_FLAG)) > 0 {
        isIp4 = true
        recType = "A"
    }

    for _, rec := range records {

        //lower case fqdn
        var fqdn string = strings.ToLower(strings.TrimSpace(hostName) + "." + strings.TrimSpace(conf.Domain))

        var recordName string = strings.TrimSpace(strings.ToLower(rec.Name))

        if (recordName == fqdn) && (rec.Type == recType) {
            if isIp4 {
                rec.Content = ip4.String()
            } else {
                rec.Content, err = joinIP(host, ip6)
                if err != nil {
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

func joinIP(host Host, ip6 net.IP) (string, error) {

    var tmp string = host.Addr

    suffix := ""
    var addr []byte

    hostPrefixLength := host.HostPrefixLength

    if host.IsMac {
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
    } else {
        addr = net.ParseIP(host.Addr)
    }

    if host.PrefixLength > host.HostPrefixLength {
        return "", errors.New("prefix-length > host-prefix-length")
    }

    bigPrefix := big.NewInt(0)
    bigPrefix.SetBytes(ip6)

    for i := 0; i < (128 - host.PrefixLength); i++ {
        bigPrefix.SetBit(bigPrefix, i, 0)
    }

    bigHostAddr := big.NewInt(0)
    bigHostAddr.SetBytes(addr)

    for i := (128 - hostPrefixLength); i < 128; i++ {
        bigHostAddr.SetBit(bigHostAddr, i, 0)
    }

    bigIP := big.NewInt(0)
    bigIP.Or(bigHostAddr, bigPrefix)

    localPrefix := strings.TrimSpace(host.HostPrefix)

    if (len(localPrefix) % 2) == 1 {
        localPrefix = "0" + localPrefix
    }

    local, err := hex.DecodeString(localPrefix)

    if err != nil {
        return "", errors.New(fmt.Sprintf("Invalid prefix-id: %s\n", localPrefix))
    }

    bigLocalPrefix := big.NewInt(0)
    bigLocalPrefix.SetBytes(local)
    bigLocalPrefix.Lsh(bigLocalPrefix, uint(hostPrefixLength))

    for i := 0; i < host.PrefixLength; i++ {
        bigLocalPrefix.SetBit(bigLocalPrefix, i, 0)
    }

    bigIP.Or(bigLocalPrefix, bigIP)
    tmpBytes := bigIP.Bytes()

    ipBytes := make([]byte, IP6_ADDR_LENGTH)

    tmpLength := len(tmpBytes)
    for k, _ := range tmpBytes {
        ipBytes[IP6_ADDR_LENGTH-k-1] = tmpBytes[tmpLength-k-1]
    }

    return net.IP(ipBytes).String(), nil
}

func isHex(r rune) bool {
    // Upper and lower case HEX ascii values
    return (((r < 58) && (r > 47)) || ((r > 64) && (r < 71)) || ((r > 96) && (r < 103)))
}
