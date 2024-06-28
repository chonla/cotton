package result

import (
	"cotton/internal/stopwatch"
)

type TestResult struct {
	Title        string
	Passed       bool
	Assertions   []*AssertionResult
	Error        error
	EllapsedTime *stopwatch.EllapsedTime
}
