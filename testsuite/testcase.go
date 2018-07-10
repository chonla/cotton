package testsuite

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chonla/yas/request"
	"github.com/fatih/color"
)

// TestCase holds a test case
type TestCase struct {
	Name         string
	Method       string
	BaseURL      string
	Path         string
	ContentType  string
	RequestBody  string
	Headers      map[string]string
	Expectations map[string]string
}

// NewTestCase creates a new testcase
func NewTestCase(name string) *TestCase {
	return &TestCase{
		Name:         name,
		Headers:      map[string]string{},
		Expectations: map[string]string{},
	}
}

// SetContentType set a corresponding content type
func (tc *TestCase) SetContentType(ct string) {
	switch strings.ToLower(ct) {
	case "json":
		ct = "application/json"
	default:
		ct = "application/json"
	}
	tc.ContentType = ct
}

// Run executes test case
func (tc *TestCase) Run() error {
	// client := &http.Client{}

	green := color.New(color.FgGreen).SprintFunc()
	// blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	url := fmt.Sprintf("%s%s", tc.BaseURL, tc.Path)
	fmt.Printf("%s: %s\n", green(tc.Method), url)

	req, e := request.NewRequester(tc.Method)
	if e != nil {
		return e
	}
	req.SetHeaders(tc.Headers)
	resp, e := req.Request(url, tc.RequestBody)

	// var req *http.Request
	// var e error
	// if tc.Method == "GET" {
	// 	req, e = http.NewRequest(tc.Method, url, nil)
	// } else {
	// 	fmt.Printf("%s\n", blue(tc.RequestBody))
	// 	body := []byte(tc.RequestBody)
	// 	req, e = http.NewRequest(tc.Method, url, bytes.NewBuffer(body))
	// }
	// for _, header := range tc.Headers {
	// 	req.Header.Set(header.Key, header.Value)
	// }
	// if e != nil {
	// 	return e
	// }

	// resp, e := client.Do(req)

	b, e := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if e != nil {
		return e
	}

	fmt.Printf("->\n%s\n", magenta(string(b)))
	return nil
}
