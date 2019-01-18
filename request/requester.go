package request

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

// RequesterInterface is http request interface
type RequesterInterface interface {
	Request(string, string) (*http.Response, error)
	SetHeaders(map[string]string)
	SetCookies([]*http.Cookie)
}

// Requester is base class for all requester
type Requester struct {
	headers     map[string]string
	cookies     []*http.Cookie
	client      *http.Client
	insecure    bool
	printDetail bool
}

// NewRequester creates a new request
func NewRequester(method string, insecure, detail bool) (RequesterInterface, error) {
	var client *http.Client
	if insecure {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
			},
		}
	} else {
		client = &http.Client{}
	}
	var req RequesterInterface
	reqer := &Requester{
		headers:     map[string]string{},
		cookies:     []*http.Cookie{},
		client:      client,
		insecure:    insecure,
		printDetail: detail,
	}
	var e error
	switch method {
	case "POST":
		req = &Poster{Requester: reqer}
	case "PUT":
		req = &Putter{Requester: reqer}
	case "PATCH":
		req = &Patcher{Requester: reqer}
	case "GET":
		req = &Getter{Requester: reqer}
	case "DELETE":
		req = &Deleter{Requester: reqer}
	case "OPTION":
		req = &Optioner{Requester: reqer}
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

// SetCookies set cookie values to request
func (r *Requester) SetCookies(c []*http.Cookie) {
	for _, cookie := range c {
		for _, cc := range r.cookies {
			if cc.Name == cookie.Name {
				continue
			}
		}
		r.cookies = append(r.cookies, cookie)
	}
}

// LogRequest prints request data
func (r *Requester) LogRequest(req *http.Request) {
	if !r.printDetail {
		return
	}

	green := color.New(color.FgGreen).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()

	insecureLabel := ""
	if strings.ToLower(req.URL.Scheme) == "https" && r.insecure {
		insecureLabel = fmt.Sprintf("%s", red(" (insecure)"))
	}

	fmt.Printf("%s %s%s\n", magenta("<<--"), white("REQUEST"), insecureLabel)
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
	fmt.Printf("%s\n", magenta("<<--"))
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

// EscapeURL to do escape parameters in given URL and rebuild the proper URL
func (r *Requester) EscapeURL(u string) string {
	uSegment, e := url.Parse(u)
	if e != nil {
		return u
	}
	q := uSegment.Query()
	uSegment.RawQuery = q.Encode()
	return uSegment.String()
}
