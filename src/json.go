package main

import "strconv"

type Zone struct {
    CFSuccess
    Id   string `json:"id"`
    Name string `json:"name"`
}

type ZoneList struct {
    CFSuccess
    Result []Zone `json:"result"`
}

type RecordList struct {
	CFSuccess
	Result []Record `json:"result"`
}

type Record struct {
	Name     string `json:"name"`
	RecordId string `json:"id"`
	Type     string `json:"type"`
	Ttl      int    `json:"ttl"`
	Content  string `json:"content"`
	ZoneId   string `json:"zone_id"`
}

type CFSuccess struct {
	Success bool      `json:"success"`
	Errors  []CFError `json:"errors"`
}
type CFError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s CFSuccess) Error() string {

	var errorString string

	for _, v := range s.Errors {
		errorString += "error-code: " + strconv.Itoa(v.Code) + ", error-message: " + v.Message + "\n"
	}

	return errorString
}
