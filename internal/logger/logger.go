package logger

import "cotton/internal/result"

type LogLevel int

const (
	Compact LogLevel = iota
	Verbose
	Debug
)

type Logger interface {
	PrintTestcaseTitle(title string) error
	PrintExecutableTitle(title string) error
	PrintTestResult(passed bool) error
	PrintInlineTestResult(passed bool) error
	PrintBlockTitle(header string) error // printing Setups or Teardowns or Assertions or Captures or whatever
	PrintAssertionResults(assertionResults []result.AssertionResult) error
	PrintAssertionResult(assertionResult result.AssertionResult) error
	PrintRequest(req string) error
}
