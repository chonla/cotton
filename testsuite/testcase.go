package testsuite

import (
	"fmt"
	"strings"

	"github.com/chonla/yas/request"
	"github.com/chonla/yas/response"
	"github.com/kr/pretty"
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
	*Assert
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
	url := fmt.Sprintf("%s%s", tc.BaseURL, tc.Path)

	req, e := request.NewRequester(tc.Method)
	if e != nil {
		return e
	}
	req.SetHeaders(tc.Headers)
	resp, e := req.Request(url, tc.RequestBody)
	if e != nil {
		return e
	}

	res := response.NewResponse(resp)

	pretty.Println(res)

	// for k, v := range tc.Expectations {

	// }

	return nil
}
