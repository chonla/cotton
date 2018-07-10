package request

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/fatih/color"
)

// Poster is a requester doing POST
type Poster struct {
	*Requester
}

// Request do actual request
func (p *Poster) Request(url, body string) (*http.Response, error) {
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Printf("%s\n", blue(body))
	bodyBytes := []byte(body)
	r, e := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if e != nil {
		return nil, e
	}
	for k, v := range p.Requester.headers {
		r.Header.Set(k, v)
	}
	return p.Requester.client.Do(r)
}
