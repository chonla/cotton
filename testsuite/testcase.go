package testsuite

import (
	"fmt"
	"strings"

	"github.com/chonla/yas/assertable"
	"github.com/chonla/yas/request"
	"github.com/chonla/yas/response"
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
	white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	url := fmt.Sprintf("%s%s", tc.BaseURL, tc.Path)

	fmt.Printf("%s\n", white("================================================================================"))
	fmt.Printf("Testcase: %s\n", white(tc.Name))
	fmt.Printf("%s\n", white("================================================================================"))

	req, e := request.NewRequester(tc.Method)
	if e != nil {
		return e
	}
	req.SetHeaders(tc.Headers)
	resp, e := req.Request(url, tc.RequestBody)
	if e != nil {
		fmt.Printf("%s: %s\n", red("Error"), e)
		return e
	}

	as, e := assertable.NewAssertable(response.NewResponse(resp))
	if e != nil {
		return e
	}

	e = as.Assert(tc.Expectations)

	return e
}
