package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
    "errors"
)

const cloudflareEndpoint string = "https://api.cloudflare.com/client/v4/"

func getZone(email string, key string, domain string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", cloudflareEndpoint+"zones/", nil)

	var zones ZoneList

	if err == nil {
		req.Header.Set("X-Auth-Key", key)
		req.Header.Set("X-Auth-Email", email)
		resp, err := client.Do(req)
		if err == nil {
			content, err := ioutil.ReadAll(resp.Body)

			if err == nil {
				err = json.Unmarshal(trim(content), &zones)
			}
		}
	}

	// firecords.Errornd correct zone
	for _, v := range zones.Result {
		var apiName string = strings.TrimSpace(strings.ToLower(v.Name))
		domain = strings.TrimSpace(strings.ToLower(domain))

		if apiName == domain {
			return v.Id, err
		}
	}

	return "", err
}

func getRecords(email string, key string, domain string) (RecordList, error) {

	// fetching domain id
	zoneId, err := getZone(email, key, domain)
	client := &http.Client{}
	req, err := http.NewRequest("GET", cloudflareEndpoint+"zones/"+zoneId+"/dns_records/", nil)

	var records RecordList
    resp := new(http.Response)
	if err == nil {
		req.Header.Set("X-Auth-Key", key)
		req.Header.Set("X-Auth-Email", email)
		resp, err = client.Do(req)
		if err == nil{
            var content []byte
            if resp.StatusCode > 400{
                err = errors.New("Request to cloudflare not authorized.")
            }else{
			    content, err = ioutil.ReadAll(resp.Body)
			    if err == nil {
				    err = json.Unmarshal(trim(content), &records)
                    if (err == nil) && !records.Success{
                        err = errors.New("Request to cloudflare not authorized.")
                    }
			    }
            }
		}
	}
	return records, err
}

func update(email string, key string, rec Record) error {
	jsonContent, err := json.Marshal(rec)
	client := &http.Client{}

	if err == nil {
		req, err := http.NewRequest("PUT", cloudflareEndpoint+"zones/"+rec.ZoneId+"/dns_records/"+rec.RecordId, bytes.NewBuffer(jsonContent))

		if err == nil {
			req.Header.Set("X-Auth-Key", key)
			req.Header.Set("X-Auth-Email", email)
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			var resp *http.Response
			resp, err = client.Do(req)
			var content []byte
			content, err = ioutil.ReadAll(resp.Body)
			var success CFSuccess
			err = json.Unmarshal(trim(content), &success)

			if !success.Success {
				return success
			}
		}
	}
	return err
}

func trim(data []byte) (tmp []byte) {
	var isFirst bool = true

	var first int
	var last int

	for k, v := range data {
		if (v > 40) && (v < 126) {
			last = k + 1
			if isFirst {
				isFirst = false
				first = k
			}
		}
	}
	tmp = data[first:last]
	return tmp
}
