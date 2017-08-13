package main

import (
	"crypto/tls"
	"net"
	"net/http"
)

func dialTLS4(network, addr string) (net.Conn, error) {
	return tls.Dial("tcp4", addr, new(tls.Config))
}

func dial4(network, addr string) (net.Conn, error) {
	return net.Dial("tcp4", addr)
}

func dialTLS6(network, addr string) (net.Conn, error) {
	return tls.Dial("tcp6", addr, new(tls.Config))
}
func dial6(network, addr string) (net.Conn, error) {
	return net.Dial("tcp6", addr)
}

var v4Client http.Client = http.Client{
	Transport: &http.Transport{
		Dial:    dial4,
		DialTLS: dialTLS4,
	},
}

var v6Client http.Client = http.Client{
	Transport: &http.Transport{
		Dial:    dial6,
		DialTLS: dialTLS6,
	},
}

func GetV4(addr string) (*http.Response, error) {
	return v4Client.Get(addr)
}
func GetV6(addr string) (*http.Response, error) {
	return v6Client.Get(addr)
}
