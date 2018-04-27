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
                rec.Content = joinIP(host)
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

func joinIP(host Host) string {

    var tmp string = host.Addr

    mac := ""

    if host.IsMac{
        // remove colons and insert fffe in the middle
        for _, v := range tmp {

            if isHex(v) {
                mac += string(v)
            }

            if len(mac) == 6 {
                mac += "fffe"
            }
        }
    }else{
        //Do some magic here please
        mac = host.Addr
    }

    tmpAddr, err := hex.DecodeString(mac)

    if err == nil {
        tmpAddr[0] = tmpAddr[0] ^ 2
    } else {
        return "Invalid mac-address"
    }

    addr := make([]byte, IP6_ADDR_LENGTH+1)

    addrLen := len(addr)
    tmpAddrLen := len(tmpAddr)
    addr[0] = 1
    for k,_ := range tmpAddr{
        addr[addrLen-k-1] = tmpAddr[tmpAddrLen-k-1]
    }

    tmpPrefix := append([]byte{1}, ip6...)

    bigPrefix := big.NewInt(0)
    bigPrefix.SetBytes(tmpPrefix)

    for i:=host.PrefixSize; i<128; i++{
        bigPrefix.SetBit(bigPrefix, int(i), 0)
    }

    bigHostAddr := big.NewInt(0)
    bigHostAddr.SetBytes(addr)

    bigIP := big.NewInt(0)
    bigIP = bigIP.Or(bigHostAddr, bigPrefix)

    ret := net.IP(bigIP.Bytes()).String()
    fmt.Println(bigIP.Bytes())
    fmt.Println(len(bigHostAddr.Bytes()))
    fmt.Println(len(bigIP.Bytes()))
    return ret
}

func isHex(r rune) bool {
    // Upper and lower case HEX ascii values
    return (((r < 58) && (r > 47)) || ((r > 64) && (r < 71)) || ((r > 96) && (r < 103)))
}
