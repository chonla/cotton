package request

import (
	"net/http"
)

// Optioner is a requester doing OPTION
type Optioner struct {
	*Requester
}

// Request do actual request
func (o *Optioner) Request(url, body string) (*http.Response, error) {
	url = o.EscapeURL(url)

	r, e := http.NewRequest("OPTION", url, nil)
	if e != nil {
		return nil, e
	}
	for k, v := range o.Requester.headers {
		r.Header.Set(k, v)
	}
	for _, cookie := range o.Requester.cookies {
		r.AddCookie(cookie)
	}

	o.Requester.LogRequest(r)

	return o.Requester.client.Do(r)
}
