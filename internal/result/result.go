package result

import "cotton/internal/stopwatch"

type TestResult struct {
	Title        string
	Passed       bool
	Assertions   []*AssertionResult
	Error        error
	EllapsedTime *stopwatch.EllapsedTime
}

type AssertionResult struct {
	Title    string
	Passed   bool
	Expected string
	Actual   string
	Error    error
}

type TestsuiteResult struct {
	TestCount       int
	PassedCount     int
	ExecutionsCount int
	FailedCount     int
	SkippedCount    int
	EllapsedTime    *stopwatch.EllapsedTime
}
