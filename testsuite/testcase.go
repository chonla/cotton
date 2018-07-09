package testsuite

import (
	"strings"
)

// TestCase holds a test case
type TestCase struct {
	Name         string
	Method       string
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
