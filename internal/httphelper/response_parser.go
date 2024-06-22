package httphelper

import "net/http"

type ResponseParser interface {
	Parse(responseString string) (*HTTPResponse, error)
	From(response *http.Response) (*HTTPResponse, error)
}
