package httphelper

import (
	"crypto/tls"
	"errors"
	"net/http"
)

type Request interface {
	Do(insecure bool) (*HTTPResponse, error)
	String() string
	Specs() map[string]interface{}
}

type HTTPRequest struct {
	req             *http.Request
	reqRaw          string
	method          string
	path            string
	protocol        string
	protocolVersion string
	headers         map[string]string
	body            string
}

func (r *HTTPRequest) Do(insecure bool) (*HTTPResponse, error) {
	if r.req == nil {
		return nil, errors.New("empty request")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	resp, err := client.Do(r.req)
	if err != nil {
		return nil, err
	}
	defer func() { resp.Body.Close() }()

	responseParser := &HTTPResponseParser{}
	httpResponse, err := responseParser.From(resp)
	if err != nil {
		return nil, err
	}
	return httpResponse, nil
}

func (r *HTTPRequest) String() string {
	return r.reqRaw
}

func (r *HTTPRequest) Specs() map[string]interface{} {
	return map[string]interface{}{
		"method":          r.method,
		"path":            r.path,
		"headers":         r.headers,
		"body":            r.body,
		"protocol":        r.protocol,
		"protocolVersion": r.protocolVersion,
	}
}
