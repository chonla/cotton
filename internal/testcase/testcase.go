package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"net/http"
)

// Test cases
type TestCase struct {
	Title       string
	Description string
	Request     *http.Request

	Captures   []*capture.Captured
	Setups     []*executable.Executable
	Teardowns  []*executable.Executable
	Assertions []*assertion.Assertion
}
