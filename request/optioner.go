package request

import (
	"net/http"
)

// Optioner is a requester doing OPTION
type Optioner struct {
	*Requester
}

// Request do actual request
func (g *Optioner) Request(url, body string) (*http.Response, error) {

	r, e := http.NewRequest("OPTION", url, nil)
	if e != nil {
		return nil, e
	}
	for k, v := range g.Requester.headers {
		r.Header.Set(k, v)
	}

	g.Requester.LogRequest(r)

	return g.Requester.client.Do(r)
}
