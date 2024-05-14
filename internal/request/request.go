package request

import (
	"bufio"
	"net/http"
	"strings"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Proto   string
}

func Parse(request string) (*Request, error) {
	reqReader := bufio.NewReader(strings.NewReader(request))
	reqHttpReader, err := http.ReadRequest(reqReader)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method: reqHttpReader.Method,
		URL:    "",
		Proto:  reqHttpReader.Proto,
	}, nil
}
