package main

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
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

	go func() {
		d.ch <- local
	}()

	return
}

func (d LookupDialer) DialTLS(network string, addr string) (conn net.Conn, err error) {
	if d.version == 4 {
		conn, err = tls.Dial("tcp4", addr, new(tls.Config))
	} else {
		conn, err = tls.Dial("tcp6", addr, new(tls.Config))
	}

	local, _, _ := net.SplitHostPort(conn.LocalAddr().String())

	go func() {
		d.ch <- local
	}()

	return
}

func (d LookupDialer) Dial(network string, addr string) (conn net.Conn, err error) {
	if d.version == 4 {
		conn, err = net.Dial("tcp4", addr)
	} else {
		conn, err = net.Dial("tcp6", addr)
	}

	local, _, _ := net.SplitHostPort(conn.LocalAddr().String())

	go func() {
		d.ch <- local
	}()

	return
}

func GetAddr(addr string, version uint) (internal net.IP, external net.IP, err error) {
	transport := new(http.Transport)
	var notFoundError error

	var dialer LookupDialer
	dialer.version = version
	dialer.ch = make(chan string)

	transport.DialContext = dialer.DialContext
	transport.DialTLS = dialer.DialTLS
	transport.Dial = dialer.Dial

	var client http.Client
	client.Transport = transport

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
