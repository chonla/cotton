package logger

import "cotton/internal/result"

type Logger interface {
	PrintTestCaseTitle(title string) error
	PrintTestResult(passed bool) error
	PrintAssertionResults(assertionResults []result.AssertionResult) error
	PrintAssertionResult(assertionResult result.AssertionResult) error

	Print(args ...any) (int, error)
	Println(args ...any) (int, error)
	Printf(format string, args ...any) (int, error)
	Printfln(format string, args ...any) (int, error)
}
