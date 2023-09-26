package nfonapiclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthentication(t *testing.T) {
	// given
	method := "GET"
	body := []byte("")
	path := "customers"
	apiSecret := "4f5f5402-da77-410a-9aad-fd5ef74f746e"
	date := "Sun, 31 Dec 2023 12:00:00 GMT"
	contentType := ContentTypeJson

	// when
	date, contentType, contentMD5, signature := getAuthenticationWithDate(date, method, body, path, apiSecret, contentType)

	// then
	exppectedDateRegExp := "^[a-zA-Z]{3}, [0-9]{2} [a-zA-Z]{3} [0-9]{4} [0-9]{2}:[0-9]{2}:[0-9]{2} GMT$"
	expectedContentType := "application/json"
	expectedContentMD5 := "d41d8cd98f00b204e9800998ecf8427e"
	expectedSignature := "UpujUjJOS5jywvmHSyw9dTl0SmM="

	assert.Regexp(t, exppectedDateRegExp, date, "date expected '"+exppectedDateRegExp+"', but got '"+date+"'")
	assert.Equal(t, contentType, expectedContentType, "contentType expected "+expectedContentType+", but got '"+contentType+"'")
	assert.Equal(t, contentMD5, expectedContentMD5, "contentMD5 expected "+expectedContentMD5+", but got '"+contentMD5+"'")
	assert.Equal(t, signature, expectedSignature, "signature expected "+expectedSignature+", but got '"+signature+"'")
}

func TestGetServiceportalClient(t *testing.T) {
	// given
	method := "POST"
	apiRequestPath := "/api/customers"
	body := []byte("")
	apiKey := "6bc2fe16-c241-49fa-af41-01a8642dd885"
	signature := "fe8ubWPbe42fDvEsiqA/XImzVIk="
	date := "Sun, 31 Dec 2023 12:00:00 GMT"
	contentMD5 := "d41d8cd98f00b204e9800998ecf8427e"
	contentType := ContentTypeJson
	header := http.Header{"Accept": []string{"*"}}

	// when
	client, req := GetServiceportalClient(method, apiRequestPath, body, apiKey, signature, date, contentMD5, contentType, header)
	if nil == client || nil == req {
		client = nil // dummy to prevent go-staticcheck warning
	}

	// then
	expectedAuthorization := "NFON-API 6bc2fe16-c241-49fa-af41-01a8642dd885:fe8ubWPbe42fDvEsiqA/XImzVIk="
	assert.Equal(t, req.Method, method, "method expected '"+method+"', got '"+req.Method+"'")
	assert.Equal(t, req.Header.Get("Authorization"), expectedAuthorization, "Authorization expected '"+expectedAuthorization+"', got '"+req.Header.Get("Authorization")+"'")
	assert.Equal(t, req.Header.Get("x-nfon-date"), date, "x-nfon-date expected '"+date+"', got '"+req.Header.Get("x-nfon-date")+"'")
	assert.Equal(t, req.Header.Get("Date"), date, "date expected '"+date+"', got '"+req.Header.Get("Date")+"'")
	assert.Equal(t, req.Header.Get("Content-MD5"), contentMD5, "Content-MD5 expected '"+contentMD5+"', got '"+req.Header.Get("Content-MD5")+"'")
	assert.Equal(t, req.Header.Get("Content-Type"), contentType, "Content-Type expected '"+contentType+"', got '"+req.Header.Get("Content-Type")+"'")
	assert.Equal(t, req.URL.String(), apiRequestPath, "apiRequestPath expected '"+apiRequestPath+"', got '"+req.URL.String()+"'")
}
