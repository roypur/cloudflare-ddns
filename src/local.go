package main

import (
	"errors"
	"net"
	"strings"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

const UDP_ENDPOINT4 = "1.2.3.4"
const UDP_ENDPOINT6 = "2a02::1"

type LookupDialer struct {
	version uint
}

type DiscoveredAddress struct {
    Addr net.IP
    Err error
    Mask uint
}

func getLocalAddr(version uint, _ time.Duration) (addr net.IP, err error) {
    fa := new(net.IPNet)
    fb := new(net.IPNet)

    notFoundError := errors.New("Local IPv6 address not found!")
    if version == 4 {
        notFoundError = errors.New("Local IPv4 address not found!")
    }

    _, fa, err = net.ParseCIDR("169.254.0.1/16")
    if err != nil {
        return
    }

    _, fb, err = net.ParseCIDR("fe80::1/10")
    if err != nil {
        return
    }

    currentEndpoint := UDP_ENDPOINT6
    currentNetwork := "udp6"
    if version == 4 {
        currentEndpoint = UDP_ENDPOINT4
        currentNetwork = "udp4"
    }

    laddr := new(net.UDPAddr)
    raddr := new(net.UDPAddr)
    var conn net.Conn
    var host string

    raddr.IP = net.ParseIP(currentEndpoint)
    conn, err = net.DialUDP(currentNetwork, laddr, raddr)

    if err != nil {
        err = notFoundError
        return
    }

    defer conn.Close()

    raw := conn.LocalAddr().String()
    if strings.ContainsRune(raw, '%') {
        err = notFoundError
        return
    }

    host, _, err = net.SplitHostPort(raw)
    if err == nil {
        tmp := net.ParseIP(host)
        if fa.Contains(tmp) || fb.Contains(tmp) {
            err = notFoundError
            return
        }
        addr = tmp
    } else {
        err = notFoundError
    }
    return
}

func (d LookupDialer) DialContext(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
	var sysDialer net.Dialer
	if d.version == 4 {
		return sysDialer.DialContext(ctx, "tcp4", addr)
	}
	return sysDialer.DialContext(ctx, "tcp6", addr)
}

func getRemoteAddr(version uint, timeout time.Duration) (addr net.IP, err error) {
    transport := new(http.Transport)
    var notFoundError error

    var dialer LookupDialer
    dialer.version = version

    transport.DialContext = dialer.DialContext

    var client http.Client
    client.Transport = transport
    client.Timeout = timeout

    if version == 4 {
        notFoundError = errors.New("Remote IPv4 address not found!")
    } else {
        notFoundError = errors.New("Remote IPv6 address not found!")
    }

    resp := new(http.Response)
    content := make([]byte, 0)

    resp, err = client.Get(IP_LOOKUP_ENDPOINT)

    if err == nil {
        content, err = ioutil.ReadAll(resp.Body)
        if err == nil {
            addr = net.ParseIP(string(trim(content)))
        } else {
            err = notFoundError
        }
    } else {
        err = notFoundError
    }
    return
}

func getAddr(version uint, local bool, timeout time.Duration)(chan DiscoveredAddress) {
    lookupFunc := getRemoteAddr
    if local {
        lookupFunc = getLocalAddr
    }

    ch := make(chan DiscoveredAddress)
    go func() {
        var res DiscoveredAddress
        res.Addr, res.Err = lookupFunc(version, timeout)
        tmp := make([]byte, IP6_ADDR_LENGTH, IP6_ADDR_LENGTH)
        for k,v := range res.Addr {
            tmp[k] = v
        }
        res.Addr = net.IP(tmp)
        ch <-res
    }()
    return ch
}
