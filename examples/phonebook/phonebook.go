package main

import (
	"fmt"
	"io"
	"net/http"
	nfonapiclient "nfon-api-client"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		panic("missing API parameters")
	}

	apiRootURL := os.Args[1]
	apiKey := os.Args[2]
	apiSecret := os.Args[3]
	customer := os.Args[4]

	if apiRootURL == "" || apiKey == "" || apiSecret == "" {
		panic("missing API parameters")
	}

	var req = http.Request{
		Method:     http.MethodGet,
		Body:       io.NopCloser(strings.NewReader("")),
		RequestURI: "/api/customers/" + customer + "/phone-books?_pagesize=3",
	}

	_, s, failed := nfonapiclient.Request(&req, nfonapiclient.ApiConfig{BaseURL: apiRootURL, Public: apiKey, Secret: apiSecret}, false)

	if !failed {
		println("API returns " + string(s))
	} else {
		println("Something failed. Please check your API parameters")
	}

	d := nfonapiclient.MultiResultParser(s)
	fmt.Printf("%v\n", d)
	if d.Total > 0 {
		fmt.Printf("First result is: %s: %s", d.Items[1].DataMap["displayName"], d.Items[1].DataMap["displayNumber"])
	}
}
