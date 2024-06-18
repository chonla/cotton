package logger

import "cotton/internal/result"

type NilLogger struct {
	debug bool
}

func NewNilLogger(debug bool) Logger {
	return &NilLogger{
		debug: debug,
	}
}

func (c *NilLogger) PrintTestCaseTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintTestResult(passed bool) error {
	return nil
}

func (c *NilLogger) PrintAssertionResults(assertions []result.AssertionResult) error {
	return nil
}

func (c *NilLogger) PrintAssertionResult(assertion result.AssertionResult) error {
	return nil
}

func (c *NilLogger) Print(args ...any) (int, error) {
	return 0, nil
}

func (c *NilLogger) Println(args ...any) (int, error) {
	return 0, nil
}

func (c *NilLogger) Printf(format string, args ...any) (int, error) {
	return 0, nil
}

func (c *NilLogger) Printfln(format string, args ...any) (int, error) {
	return 0, nil
}
