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
	ch      chan string
}

func (d LookupDialer) DialContext(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
	var sysDialer net.Dialer
	if d.version == 4 {
		conn, err = sysDialer.DialContext(ctx, "tcp4", addr)
	} else {
		conn, err = sysDialer.DialContext(ctx, "tcp6", addr)
	}

	local, _, _ := net.SplitHostPort(conn.LocalAddr().String())
	d.ch <- local

	return
}

func CheckRedirect(req *http.Request, via []*http.Request) (err error) {
	if len(via) > 5 {
		err = errors.New("stopped after 5 redirects")
	}
	return
}

func GetAddr(addr string, version uint) (internal net.IP, external net.IP, err error) {
	transport := new(http.Transport)
	var notFoundError error

	var dialer LookupDialer
	dialer.version = version
	dialer.ch = make(chan string, 7)

	transport.DialContext = dialer.DialContext

	var client http.Client
	client.Transport = transport
	client.CheckRedirect = CheckRedirect
	client.Timeout = time.Duration(TIMEOUT) * time.Second

	if version == 4 {
		notFoundError = errors.New("IPv4 address not found!")
	} else {
		notFoundError = errors.New("IPv6 address not found!")
	}

	resp := new(http.Response)
	content := make([]byte, 0)

	resp, err = client.Get(addr)

	if err == nil {
		tmp := <-dialer.ch

		content, err = ioutil.ReadAll(resp.Body)
		if err == nil {

			internal = net.ParseIP(tmp)
			external = net.ParseIP(string(trim(content)))
		} else {
			err = notFoundError
		}
	}

	return
}
