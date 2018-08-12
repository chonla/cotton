package request

import (
	"net/http"
)

// Deleter is a requester doing DELETE
type Deleter struct {
	*Requester
}

// Request do actual request
func (d *Deleter) Request(url, body string) (*http.Response, error) {
	url = d.EscapeURL(url)

	r, e := http.NewRequest("DELETE", url, nil)
	if e != nil {
		return nil, e
	}
	for k, v := range d.Requester.headers {
		r.Header.Set(k, v)
	}

	d.Requester.LogRequest(r)

	return d.Requester.client.Do(r)
}
