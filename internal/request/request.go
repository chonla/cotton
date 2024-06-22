package request

import (
	"bytes"
	"net/http"
)

type Request interface {
	Data() []byte
	Do() (*http.Response, error)
	SimilarTo(anotherRequest Request) bool
	String() string
}

type HTTPRequest struct {
	request      *http.Request
	plainRequest []byte
}

func (r *HTTPRequest) SimilarTo(anotherRequest Request) bool {
	return bytes.Equal(r.plainRequest, anotherRequest.Data())
}

func (r *HTTPRequest) Do() (*http.Response, error) {
	return http.DefaultClient.Do(r.request)
}

func (r *HTTPRequest) String() string {
	return string(r.plainRequest)
}

func (r *HTTPRequest) Data() []byte {
	return r.plainRequest
}
