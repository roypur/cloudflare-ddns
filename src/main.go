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
    "os"
)

const IP_LOOKUP_ENDPOINT = "https://icanhazip.com/"

const IP6_ADDR_LENGTH = 16

const IP4_FLAG uint = (1 << 1)
const IP6_FLAG uint = (1 << 2)
const LOCAL_FLAG uint = (1 << 3)
const REMOTE_FLAG uint = (1 << 4)

var lookup = make(map[uint]net.IP)
var success bool

var mutex sync.Mutex
var wg sync.WaitGroup

func main() {
    if len(os.Args) != 2 {
        fmt.Printf("%s <config.json>\n")
        return
    }
    conf, err := getConfig(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }
	if conf.Interval > 0 {
        time.Sleep(time.Minute)
        loop(conf)
		for {
			time.Sleep(time.Duration(conf.Interval))
			loop(conf)
		}
	} else {
		loop(conf)
	}
}

func loop(conf Config) {
	var err error
	mask := uint(0)

    dc := make(map[uint]chan DiscoveredAddress)
    dc[IP4_FLAG | LOCAL_FLAG] = getAddr(4, true, conf.Timeout)
    dc[IP6_FLAG | LOCAL_FLAG] = getAddr(6, true, conf.Timeout)
    dc[IP4_FLAG | REMOTE_FLAG] = getAddr(4, false, conf.Timeout)
    dc[IP6_FLAG | REMOTE_FLAG] = getAddr(6, false, conf.Timeout)

    for flag,ch := range dc {
        data := <-ch
        if data.Err == nil {
            lookup[flag] = data.Addr
            mask |= flag
        } else {
            fmt.Println(data.Err)
        }
    }
    records, err := getRecords(conf.Token, conf.Domain, conf.Timeout)

    if err == nil {
        for _, host := range conf.Hosts {
            if (host.Mask & mask) > 0 {
                success = true
                wg.Add(1)
            }
        }
        for _, host := range conf.Hosts {
            if (host.Mask & mask) > 0 {
                go findAndUpdate(records.Result, host.Name, conf, host)
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
    addr, ok := lookup[host.Mask]
    if ok {
        recType := "AAAA"
        if (host.Mask & IP4_FLAG) > 0 {
            isIp4 = true
            recType = "A"
        }

        for _, rec := range records {

            //lower case fqdn
            var fqdn string = strings.ToLower(strings.TrimSpace(hostName) + "." + strings.TrimSpace(conf.Domain))

            var recordName string = strings.TrimSpace(strings.ToLower(rec.Name))

            if (recordName == fqdn) && (rec.Type == recType) {
                if isIp4 {
                    rec.Content = addr.String()
                } else {
                    rec.Content, err = joinIP(host, addr)
                    if err != nil {
                        fmt.Printf("ERROR: %s\n", err)
                        mutex.Lock()
                        success = false
                        mutex.Unlock()
                        return
                    }
                }
                err = update(conf.Token, rec, conf.Timeout)

                if err != nil {
                    mutex.Lock()
                    fmt.Println("Failed to update " + recType + " record " + fqdn)
                    success = false
                    mutex.Unlock()
                }
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
