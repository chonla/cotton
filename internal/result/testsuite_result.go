package result

import "cotton/internal/stopwatch"

type TestsuiteResult struct {
	TestCount       int
	PassedCount     int
	ExecutionsCount int
	FailedCount     int
	SkippedCount    int
	EllapsedTime    *stopwatch.EllapsedTime
	Start           int64
	Stop            int64
	AppVersion      string
	TestResults     []*TestResult
}
