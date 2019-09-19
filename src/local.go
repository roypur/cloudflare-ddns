package main

import (
	"errors"
	"net"
	"strings"
)

const UDP_ENDPOINT4 = "1.2.3.4"
const UDP_ENDPOINT6 = "2a02::1"

func getLocalAddr(version uint) (addr net.IP, err error) {
	fa := new(net.IPNet)
	fb := new(net.IPNet)

	notFoundError := errors.New("IPv6 address not found!")
	if version == 4 {
		notFoundError = errors.New("IPv4 address not found!")
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
	}

	return
}
