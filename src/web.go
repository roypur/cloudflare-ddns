package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type LookupDialer struct {
	version uint
}

func (d LookupDialer) DialContext(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
	var sysDialer net.Dialer
	if d.version == 4 {
		return sysDialer.DialContext(ctx, "tcp4", addr)
	}
	return sysDialer.DialContext(ctx, "tcp6", addr)
}

func getRemoteAddr(version uint) (addr net.IP, err error) {
	transport := new(http.Transport)
	var notFoundError error

	var dialer LookupDialer
	dialer.version = version

	transport.DialContext = dialer.DialContext

	var client http.Client
	client.Transport = transport
	client.Timeout = time.Duration(TIMEOUT) * time.Second

	if version == 4 {
		notFoundError = errors.New("IPv4 address not found!")
	} else {
		notFoundError = errors.New("IPv6 address not found!")
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
	}

	return
}
