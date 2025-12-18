package engine

import "fmt"

// ParseError represents validation failures that happen before execution.
type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse error: %s", e.msg)
}

// RunError represents failures that happen during execution time.
type RunError struct {
	msg string
}

func (e *RunError) Error() string {
	return fmt.Sprintf("Run error: %s", e.msg)
}

func newParseError(msg string) error {
	return &ParseError{msg: msg}
}

func newRunError(msg string) error {
	return &RunError{msg: msg}
}
