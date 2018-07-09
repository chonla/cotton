package testsuite

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	Expectations []Expectation
}

// Expectation is a set of expectation
type Expectation struct {
	Key   string
	Value string
}

// NewExpectation creates a new expectation
func NewExpectation(key, value string) Expectation {
	return Expectation{
		Key:   key,
		Value: value,
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
	client := &http.Client{}
	url := fmt.Sprintf("%s%s", tc.BaseURL, tc.Path)
	fmt.Println(url)
	req, e := http.NewRequest(tc.Method, url, nil)
	if e != nil {
		return e
	}
	resp, e := client.Do(req)

	b, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		return e
	}

	pretty.Println(string(b))
	return nil
}
