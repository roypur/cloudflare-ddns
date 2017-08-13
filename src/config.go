package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type config struct {
	ApiKey   string            `json:"apiKey"`
	ApiEmail string            `json:"apiEmail"`
	Domain   string            `json:"domain"`
	Ipv4     []string          `json:"ipv4"`
	Ipv6     map[string]string `json:"ipv6"`
}

func getConfig() config {

	var fileName string = "config.json"

	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	fileContent, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var c config

	err = json.Unmarshal(fileContent, &c)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}
