package logger

import (
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/variable"
)

type NilLogger struct {
	level LogLevel
}

func NewNilLogger(level LogLevel) Logger {
	return &NilLogger{
		level: level,
	}
}

func (c *NilLogger) PrintTestcaseTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintTestcaseTitleWithPath(title, path string) error {
	return nil
}

func (c *NilLogger) PrintExecutableTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintTestResult(passed bool) error {
	return nil
}

func (c *NilLogger) PrintSectionedMessage(section, message string) error {
	return nil
}

func (c *NilLogger) PrintTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error {
	return nil
}

func (c *NilLogger) PrintInlineTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error {
	return nil
}

func (c *NilLogger) PrintInlineTestResult(passed bool) error {
	return nil
}

func (c *NilLogger) PrintAssertionResults(assertions []*result.AssertionResult) error {
	return nil
}

func (c *NilLogger) PrintAssertionResult(assertion *result.AssertionResult) error {
	return nil
}

func (c *NilLogger) PrintRequest(req string) error {
	return nil
}

func (c *NilLogger) PrintResponse(resp string) error {
	return nil
}

func (c *NilLogger) PrintError(fileContext string, err error) error {
	return nil
}

func (c *NilLogger) PrintTestcaseSequence(index, total int) error {
	return nil
}

func (c *NilLogger) PrintTestsuiteResult(testsuiteResult *result.TestsuiteResult) error {
	return nil
}

func (c *NilLogger) PrintSectionTitle(sectionTitle string) error {
	return nil
}

func (c *NilLogger) PrintDebugMessage(message string) error {
	return nil
}

func (c *NilLogger) PrintDetailedDebugMessage(messages ...string) error {
	return nil
}

func (c *NilLogger) PrintVariables(variables *variable.Variables) error {
	return nil
}
