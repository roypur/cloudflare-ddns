package main

import (
	"encoding/json"
	"io/ioutil"
    "time"
)

type ConfigFile struct {
	Interval uint            `json:"interval"`
    Timeout  uint            `json:"timeout"`
	Token    string          `json:"token"`
	Domain   string          `json:"domain"`
	Ipv4     map[string]ConfigHost `json:"ipv4"`
	Ipv6     map[string]ConfigHost `json:"ipv6"`
}
type ConfigHost struct {
	LocalMode        bool   `json:"local"`
	Addr             string `json:"addr"`
	PrefixLength     int    `json:"prefix-length"`
	HostPrefixLength int    `json:"host-prefix-length"`
	HostPrefix       string `json:"prefix-id"`
	IsMac            bool   `json:"ismac"`
}
type Config struct {
	Interval time.Duration
    Timeout time.Duration
	Token    string
	Domain   string
	Hosts    []Host
}
type Host struct {
	Mask             uint
    Name             string
	Addr             string
	PrefixLength     int
	HostPrefixLength int
	HostPrefix       string
	IsMac            bool
}

func getConfig(filename string) (conf Config, err error) {
    content := make([]byte, 0, 0)
	content, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	var cf ConfigFile
	err = json.Unmarshal(content, &cf)
	if err != nil {
        return
	}

    createHost := func(name string, configHost ConfigHost, flag uint)(host Host) {
        host.Name = name
        host.Addr = configHost.Addr
        host.PrefixLength = configHost.PrefixLength
        host.HostPrefixLength = configHost.HostPrefixLength
        host.HostPrefix = configHost.HostPrefix
        host.IsMac = configHost.IsMac
		host.Mask = REMOTE_FLAG | flag
		if configHost.LocalMode {
			host.Mask = LOCAL_FLAG | flag
		}
        return
    }

    clen := len(cf.Ipv4) + len(cf.Ipv6)
    conf.Hosts = make([]Host, clen, clen)

    if cf.Interval < 60 {
        conf.Interval = time.Minute
    } else {
        conf.Interval = time.Duration(cf.Interval) * time.Second
    }
    if cf.Timeout < 1 {
        conf.Timeout = time.Second
    } else {
        conf.Timeout = time.Duration(cf.Timeout) * time.Second
    }

    conf.Token = cf.Token
    conf.Domain = cf.Domain

    for k, v := range cf.Ipv4 {
        clen--
        conf.Hosts[clen] = createHost(k, v, IP4_FLAG)
	}
	for k, v := range cf.Ipv6 {
        clen--
        conf.Hosts[clen] = createHost(k, v, IP6_FLAG)
	}
	return
}
