package response

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

// Response is response from http
type Response struct {
	Proto      string
	Status     string
	StatusCode int
	Header     map[string][]string
	Body       string
}

// NewResponse creates a new parsed response
func NewResponse(resp *http.Response) *Response {
	headers := map[string][]string{}

	for k, v := range resp.Header {
		if headers[k] == nil {
			headers[k] = []string{}
		}
		for _, t := range v {
			headers[k] = append(headers[k], t)
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	r := &Response{
		Proto:      resp.Proto,
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     headers,
		Body:       string(body),
	}

	r.LogResponse()

	return r
}

// LogResponse prints response body
func (r *Response) LogResponse() {
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Printf("%s\n", magenta("-->>"))
	fmt.Printf("%s %s\n", r.Proto, green(r.Status))
	for k, v := range r.Header {
		for _, t := range v {
			fmt.Printf("%s: %s\n", yellow(k), t)
		}
	}
	fmt.Println()

	if r.Body != "" {
		fmt.Printf("%s\n", blue(r.Body))
	}
}
