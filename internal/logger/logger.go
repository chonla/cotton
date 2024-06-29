package logger

import (
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/variable"
)

type LogLevel int

const (
	Compact LogLevel = iota
	Verbose
	Debug
)

type Logger interface {
	PrintSectionedMessage(section, message string) error
	PrintTestcaseSequence(index, total int) error
	PrintTestcaseTitle(title string) error
	PrintSectionTitle(sectionTitle string) error
	PrintExecutableTitle(title string) error
	PrintTestResult(passed bool) error
	PrintTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error
	PrintInlineTestResult(passed bool) error
	PrintInlineTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error
	PrintAssertionResults(assertionResults []*result.AssertionResult) error
	PrintAssertionResult(assertionResult *result.AssertionResult) error
	PrintRequest(req string) error
	PrintError(err error) error
	PrintTestsuiteResult(testsuiteResult *result.TestsuiteResult) error
	PrintDebugMessage(message string) error
	PrintVariables(variables *variable.Variables) error
}
