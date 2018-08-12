package request

import (
	"bytes"
	"net/http"
)

// Patcher is a requester doing PATCH
type Patcher struct {
	*Requester
}

// Request do actual request
func (p *Patcher) Request(url, body string) (*http.Response, error) {
	url = p.EscapeURL(url)

	bodyBytes := []byte(body)
	r, e := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyBytes))
	if e != nil {
		return nil, e
	}
	for k, v := range p.Requester.headers {
		r.Header.Set(k, v)
	}

	p.Requester.LogRequest(r)

	return p.Requester.client.Do(r)
}
