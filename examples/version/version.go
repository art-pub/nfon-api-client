package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	nfonapiclient "github.com/art-pub/nfon-api-client"
)

// run this with go run version.go https://API_ROOT_URL API_KEY API_SECRET
func main() {

	if len(os.Args) < 4 {
		panic("missing API parameters")
	}

	apiRootURL := os.Args[1]
	apiKey := os.Args[2]
	apiSecret := os.Args[3]

	if apiRootURL == "" || apiKey == "" || apiSecret == "" {
		panic("missing API parameters")
	}

	var req = http.Request{
		Method:     http.MethodGet,
		Body:       io.NopCloser(strings.NewReader("")),
		RequestURI: "/api/version",
	}

	_, s, failed := nfonapiclient.Request(&req, nfonapiclient.ApiConfig{BaseURL: apiRootURL, Public: apiKey, Secret: apiSecret}, false)

	if !failed {
		println("Status is " + string(s))
	} else {
		println("Something failed. Please check your API parameters")
	}

	d := nfonapiclient.SingleresultParser(s)
	println("version is '" + d.DataMap["version"] + " - " + d.DataMap["buildTime"] + " - " + d.DataMap["host"] + "'")
}
