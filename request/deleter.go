package request

import (
	"net/http"
)

// Deleter is a requester doing DELETE
type Deleter struct {
	*Requester
}

// Request do actual request
func (g *Deleter) Request(url, body string) (*http.Response, error) {

	r, e := http.NewRequest("DELETE", url, nil)
	if e != nil {
		return nil, e
	}
	for k, v := range g.Requester.headers {
		r.Header.Set(k, v)
	}

	g.Requester.LogRequest(r)

	return g.Requester.client.Do(r)
}
