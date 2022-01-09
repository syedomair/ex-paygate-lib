package mockserver

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

func MockTestServer(method string, url string, jsonInput []byte) (*httptest.ResponseRecorder, *http.Request) {
	var jsonStr = []byte(jsonInput)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	response := httptest.NewRecorder()
	return response, req
}
