package request

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httputil"
)

type Request interface {
	Data() []byte
	Do() (*http.Response, error)
	Similar(anotherRequest Request) bool
	String() string
}

type HTTPRequest struct {
	request      *http.Request
	plainRequest []byte
}

func New(req *http.Request) (*HTTPRequest, error) {
	if req == nil {
		return nil, errors.New("unexpected nil request")
	}
	req.Close = true
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	return &HTTPRequest{
		request:      req,
		plainRequest: reqBytes,
	}, nil
}

func (r *HTTPRequest) Similar(anotherRequest Request) bool {
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
