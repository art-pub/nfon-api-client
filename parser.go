package nfonapiclient

import (
	"encoding/json"
	"fmt"
	"log"
)

type NameValuePair struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}
type RelHrefPair struct {
	Name  string `json:"rel"`
	Value string `json:"href"`
}

type Items struct {
	Href    string
	Links   []string
	Data    []NameValuePair
	DataMap map[string]string
}

type SingleResult struct {
	Href    string
	Links   []string
	Data    []NameValuePair
	DataMap map[string]string
}

type MultiResult struct {
	Href     string
	Offset   int
	Total    int
	Links    []RelHrefPair
	LinksMap map[string]string
	Items    []Items
}

func SingleresultParser(body []byte) SingleResult {

	if !json.Valid([]byte(body)) {
		// handle the error here
		log.Printf("invalid JSON string: %s", string(body))
		return SingleResult{}
	}

	var parsed SingleResult
	json.Unmarshal(body, &parsed)

	parsed.DataMap = make(map[string]string)
	for _, d := range parsed.Data {
		parsed.DataMap[d.Name] = fmt.Sprintf("%v", d.Value)
	}

	return parsed
}

func MultiResultParser(body []byte) MultiResult {

	if !json.Valid([]byte(body)) {
		// handle the error here
		log.Printf("invalid JSON string: %s", string(body))
		return MultiResult{}
	}

	var parsed MultiResult
	json.Unmarshal(body, &parsed)

	parsed.LinksMap = make(map[string]string)
	for _, link := range parsed.Links {
		parsed.LinksMap[link.Name] = link.Value
	}

	for ii, i := range parsed.Items {
		i.DataMap = make(map[string]string)
		for _, d := range i.Data {
			i.DataMap[d.Name] = fmt.Sprintf("%v", d.Value)
		}
		parsed.Items[ii].DataMap = i.DataMap
	}

	return parsed
}
