package request

import (
	"bytes"
	"net/http"
)

// Putter is a requester doing PUT
type Putter struct {
	*Requester
}

// Request do actual request
func (p *Putter) Request(url, body string) (*http.Response, error) {
	url = p.EscapeURL(url)

	bodyBytes := []byte(body)
	r, e := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
	if e != nil {
		return nil, e
	}
	for k, v := range p.Requester.headers {
		r.Header.Set(k, v)
	}

	p.Requester.LogRequest(r)

	return p.Requester.client.Do(r)
}
