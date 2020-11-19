package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
	Interval int             `json:"interval"`
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
	Interval int
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

func getConfig() Config {
	var fileName string = "config.json"

	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var cf ConfigFile
	err = json.Unmarshal(fileContent, &cf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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

    var c Config
    clen := len(cf.Ipv4) + len(cf.Ipv6)
    c.Hosts = make([]Host, clen, clen)
    c.Interval = cf.Interval
    c.Token = cf.Token
    c.Domain = cf.Domain

    for k, v := range cf.Ipv4 {
        clen--
        c.Hosts[clen] = createHost(k, v, IP4_FLAG)
	}
	for k, v := range cf.Ipv6 {
        clen--
        c.Hosts[clen] = createHost(k, v, IP6_FLAG)
	}
	return c
}
