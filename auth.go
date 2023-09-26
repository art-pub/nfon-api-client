package nfonapiclient

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const ContentTypeJson = "application/json"

// get the NFON Admin Portal REST API signature for a given API secret.
func GetAuthentication(method string, body []byte, path string, apiSecret string, contentType string) (string, string, string, string) {
	loc, _ := time.LoadLocation("GMT")
	date := time.Now().In(loc).Format("Mon, 02 Jan 2006 15:04:05 MST")

	// set default content type if empty
	if contentType == "" {
		contentType = ContentTypeJson
	}

	return getAuthenticationWithDate(date, method, body, path, apiSecret, contentType)
}

// get signature for current timestamp/date
func getAuthenticationWithDate(date string, method string, body []byte, path string, apiSecret string, contentType string) (string, string, string, string) {
	// calculate signature, etc.
	// log.Printf("[DEBUG] (request) content-type is '%s'", contentType)

	// md5 for content-md5
	h := md5.New()
	io.WriteString(h, string(body))
	// attention: content-MD5 is normally base64 encoded, wrong implementation in NFON Admin Portal REST API!
	contentMD5 := fmt.Sprintf("%x", h.Sum(nil))
	// contentMD5 := base64.StdEncoding.EncodeToString(h.Sum(nil))

	stringToSign := method + "\n" + contentMD5 + "\n" + contentType + "\n" + date + "\n" + path
	// stringToSign = r.Method + "\n" + date + "\n/api/" + path
	// log.Printf("[DEBUG] String to sign: '%s'", stringToSign)

	mac := hmac.New(sha1.New, []byte(apiSecret))
	mac.Write([]byte(stringToSign))
	macString := mac.Sum(nil)
	// macStringHex := fmt.Sprintf("%x", macString)
	// log.Printf("[DEBUG] sha1 hmac is '%s', in hex '%s'", macString, macStringHex)

	// signature := base64.StdEncoding.EncodeToString([]byte(macStringHex))
	signature := base64.StdEncoding.EncodeToString([]byte(macString))
	return date, contentType, contentMD5, signature
}

func GetServiceportalClient(method string, apiRequestPath string, body []byte, apiKey string, signature string, date string, contentMD5 string, contentType string, header http.Header) (*http.Client, *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest(method, apiRequestPath, bytes.NewBuffer(body))
	if nil != err {
		log.Printf("[ERROR] could not create client!")
		return client, nil
	}

	req.Header.Add("Authorization", "NFON-API "+apiKey+":"+signature)
	req.Header.Add("x-nfon-date", date)
	req.Header.Add("date", date)

	// Content-MD5
	if contentMD5 != "" /* && http.MethodGet != r.Method */ {
		req.Header.Add("Content-MD5", contentMD5)
	}

	// Content-Length
	if contentLength := len(body); 0 < len(body) {
		req.Header.Add("Content-Length", strconv.Itoa(contentLength))
	}

	// Content-Type
	if contentType != "" /* && http.MethodGet != r.Method */ {
		req.Header.Add("Content-Type", contentType)
	}

	// copy and log headers
	for name, values := range header {
		switch name {
		case "Accept", "X-Forwarded-For":
			for _, value := range values {
				req.Header.Add(name, value)
				// log.Printf("[DEBUG] header %s:%s\n", name, value)
			}
		}
	}

	return client, req
}
