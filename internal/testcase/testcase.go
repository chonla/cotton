package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/executable"
	"net/http"
)

// Test cases
type TestCase struct {
	Title       string
	Description string
	Request     *http.Request

	Setups     []*executable.Executable
	Teardowns  []*executable.Executable
	Assertions []*assertion.Assertion
}
