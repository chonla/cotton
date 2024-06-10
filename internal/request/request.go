package request

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Request struct {
	request      *http.Request
	plainRequest []byte
}

func New(req *http.Request) (*Request, error) {
	if req == nil {
		return nil, errors.New("unexpected nil request")
	}
	req.Close = true
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	return &Request{
		request:      req,
		plainRequest: reqBytes,
	}, nil
}

func (r *Request) Similar(anotherRequest *Request) bool {
	fmt.Println("-------------------")
	fmt.Println("Expected")
	fmt.Println(string(r.plainRequest))
	fmt.Println("-------------------")
	fmt.Println("Actual")
	fmt.Println(string(anotherRequest.plainRequest))
	fmt.Println("-------------------")
	return bytes.Equal(r.plainRequest, anotherRequest.plainRequest)
}

func (r *Request) Do() (*http.Response, error) {
	return http.DefaultClient.Do(r.request)
}
