package logger

import "cotton/internal/result"

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
	PrintInlineTestResult(passed bool) error
	PrintAssertionResults(assertionResults []*result.AssertionResult) error
	PrintAssertionResult(assertionResult *result.AssertionResult) error
	PrintRequest(req string) error
	PrintError(err error) error
	PrintTestsuiteResult(testsuiteResult *result.TestsuiteResult) error
	PrintDebugMessage(message string) error
}
