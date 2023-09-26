package nfonapiclient

import (
	"io/ioutil"
	"log"
	"net/http"
)

type ApiConfig struct {
	BaseURL string
	Public  string
	Secret  string
}

// Do a request and send it to the NFON admin portal REST API.
//
func Request(r *http.Request, config ApiConfig, debugging bool) (*http.Response, []byte, bool) {
	// add service portal API prefix
	apiRequestPath := config.BaseURL + r.RequestURI
	if debugging {
		log.Printf("[DEBUG] apiRequestPath is '%s'", apiRequestPath)
	}

	// set the API KEY and SECRET
	apiKey := config.Public
	apiSecret := config.Secret

	// prepare data for authentication
	body, _ := ioutil.ReadAll(r.Body)
	if debugging {
		log.Printf("[DEBUG] (incoming) body for this request is '%s'", body)
	}

	// request the serviceportal API authentication
	date, contentType, contentMD5, signature := GetAuthentication(r.Method, body, r.RequestURI, apiSecret, r.Header.Get("Content-Type"))

	client, req := GetServiceportalClient(r.Method, apiRequestPath, body, apiKey, signature, date, contentMD5, contentType, r.Header)
	// *** MAGIC *** MAGIC *** MAGIC *** //
	resp, err := client.Do(req)
	if err != nil {
		if debugging {
			log.Printf("[ERROR] " + err.Error())
		}
		return nil, []byte(err.Error()), true
	}

	// read the body from the NFON admin portal REST API response
	apiBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("[ERROR] " + err.Error())
		return nil, nil, true
	}

	if debugging {
		log.Printf("[DEBUG] (returned) body is:\n%s", apiBody)
	}

	return resp, apiBody, false
}
