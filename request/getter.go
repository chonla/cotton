package request

import (
	"net/http"
)

// Getter is a requester doing POST
type Getter struct {
	*Requester
}

// Request do actual request
func (g *Getter) Request(url, body string) (*http.Response, error) {

	r, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return nil, e
	}
	for k, v := range g.Requester.headers {
		r.Header.Set(k, v)
	}

	g.Requester.LogRequest(r)

	return g.Requester.client.Do(r)
}
