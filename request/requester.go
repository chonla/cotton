package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

// RequesterInterface is http request interface
type RequesterInterface interface {
	Request(string, string) (*http.Response, error)
	SetHeaders(map[string]string)
}

// Requester is base class for all requester
type Requester struct {
	headers map[string]string
	client  *http.Client
}

// NewRequester creates a new request
func NewRequester(method string) (RequesterInterface, error) {
	var req RequesterInterface
	var e error
	switch method {
	case "POST":
		req = &Poster{
			&Requester{
				headers: map[string]string{},
				client:  &http.Client{},
			},
		}
	default:
		e = errors.New("unsupported http method")
	}
	return req, e
}

// SetHeaders set header values to request
func (r *Requester) SetHeaders(h map[string]string) {
	for k, v := range h {
		r.headers[k] = v
	}
}

// LogRequest prints request data
func (r *Requester) LogRequest(req *http.Request) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

	fmt.Printf("%s\n", magenta("<<--"))
	fmt.Printf("%s %s %s\n", green(req.Method), req.URL, req.Proto)
	for k, v := range req.Header {
		for _, t := range v {
			fmt.Printf("%s: %s\n", yellow(k), t)
		}
	}
	fmt.Println()
	if req.Method == "POST" {
		bodyCopy, _ := req.GetBody()
		body, _ := ioutil.ReadAll(bodyCopy)
		r.LogBody(string(body))
	}
}

// LogBody prints request body
func (r *Requester) LogBody(body string) {
	blue := color.New(color.FgBlue).SprintFunc()

	prettyBody := r.prettifyJSON(body)

	fmt.Printf("%s\n", blue(prettyBody))
}

func (r *Requester) prettifyJSON(jsonString string) string {
	jsonObj := map[string]interface{}{}
	json.Unmarshal([]byte(jsonString), &jsonObj)
	prettyBody, _ := json.MarshalIndent(jsonObj, "", "    ")
	return string(prettyBody)
}
