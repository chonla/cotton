package request

import (
	"errors"
	"net/http"
)

// RequesterInterface is http request interface
type RequesterInterface interface {
	Request(string, string) (*http.Response, error)
	SetHeaders(map[string]string)
}

// Requester is base class for all requester
type Requester struct {
	headers map[string]string
	client  *http.Client
}

// NewRequester creates a new request
func NewRequester(method string) (RequesterInterface, error) {
	var req RequesterInterface
	var e error
	switch method {
	case "POST":
		req = &Poster{
			&Requester{
				headers: map[string]string{},
				client:  &http.Client{},
			},
		}
	default:
		e = errors.New("unsupported http method")
	}
	return req, e
}

// SetHeaders set header values to request
func (r *Requester) SetHeaders(h map[string]string) {
	for k, v := range h {
		r.headers[k] = v
	}
}
