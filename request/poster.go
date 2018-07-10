package request

import (
	"bytes"
	"net/http"
)

// Poster is a requester doing POST
type Poster struct {
	*Requester
}

// Request do actual request
func (p *Poster) Request(url, body string) (*http.Response, error) {

	bodyBytes := []byte(body)
	r, e := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if e != nil {
		return nil, e
	}
	for k, v := range p.Requester.headers {
		r.Header.Set(k, v)
	}

	p.Requester.LogRequest(r)

	return p.Requester.client.Do(r)
}
