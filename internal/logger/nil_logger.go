package logger

import "cotton/internal/result"

type NilLogger struct {
	level LogLevel
}

func NewNilLogger(level LogLevel) Logger {
	return &NilLogger{
		level: level,
	}
}

func (c *NilLogger) PrintTestCaseTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintExecutableTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintBlockTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintTestResult(passed bool) error {
	return nil
}

func (c *NilLogger) PrintInlineTestResult(passed bool) error {
	return nil
}

func (c *NilLogger) PrintAssertionResults(assertions []result.AssertionResult) error {
	return nil
}

func (c *NilLogger) PrintAssertionResult(assertion result.AssertionResult) error {
	return nil
}

func (c *NilLogger) PrintRequest(req string) error {
	return nil
}
