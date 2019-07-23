package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Interval         int               `json:"interval"`
	ApiKey           string            `json:"api-key"`
	ApiEmail         string            `json:"api-email"`
	Domain           string            `json:"domain"`
	Ipv4             map[string]Host   `json:"ipv4"`
	Ipv6             map[string]Host   `json:"ipv6"`
}

type Host struct {
        mask             uint
        LocalMode        bool              `json:"local"`
	Addr             string            `json:"addr"`
	PrefixLength     int               `json:"prefix-length"`
	HostPrefixLength int               `json:"host-prefix-length"`
	HostPrefix       string            `json:"prefix-id"`
	IsMac            bool              `json:"ismac"`
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
        var c Config
        err = json.Unmarshal(fileContent, &c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

        for k,v := range c.Ipv4 {
            v.mask = IP4_REMOTE_FLAG
            if v.LocalMode {
                v.mask = IP4_LOCAL_FLAG
            }
            c.Ipv4[k] = v
        }

        for k,v := range c.Ipv6 {
            v.mask = IP6_REMOTE_FLAG
            if v.LocalMode {
                v.mask = IP6_LOCAL_FLAG
            }
            c.Ipv6[k] = v
        }
	return c
}
