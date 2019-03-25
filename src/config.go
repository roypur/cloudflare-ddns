package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ApiKey       string            `json:"apiKey"`
	ApiEmail     string            `json:"apiEmail"`
	Domain       string            `json:"domain"`
        LocalNetwork bool              `json:"local"`
	Ipv4         []string          `json:"ipv4"`
	Ipv6         map[string]Host   `json:"ipv6"`
}
type Host struct {
    Addr string `json:"addr"`
    PrefixLength int `json:"prefix-length"`
    HostPrefixLength int `json:"host-prefix-length"`
    HostPrefix string `json:"prefix-id"`
    IsMac bool `json:"ismac"`
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

	return c
}
