package response

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

// Response is response from http
type Response struct {
	Proto       string
	Status      string
	StatusCode  int
	Header      map[string][]string
	Body        string
	printDetail bool
}

// NewResponse creates a new parsed response
func NewResponse(resp *http.Response, detail bool) *Response {
	headers := map[string][]string{}

	for k, v := range resp.Header {
		if headers[k] == nil {
			headers[k] = []string{}
		}
		for _, t := range v {
			headers[k] = append(headers[k], t)
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	r := &Response{
		Proto:       resp.Proto,
		Status:      resp.Status,
		StatusCode:  resp.StatusCode,
		Header:      headers,
		Body:        string(body),
		printDetail: detail,
	}

	return r
}

// LogResponse prints response body
func (r *Response) LogResponse() {
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Printf("%s %s\n", magenta("-->>"), white("RESPONSE"))
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
	fmt.Printf("%s\n", magenta("-->>"))
}
