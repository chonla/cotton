package logger

import (
	"cotton/internal/line"
	"cotton/internal/result"

	"fmt"

	"github.com/fatih/color"
)

type TerminalLogger struct {
	debug bool
}

func NewTerminalLogger(debug bool) Logger {
	return &TerminalLogger{
		debug: debug,
	}
}

func (c *TerminalLogger) PrintTestCaseTitle(title string) error {
	_, err := fmt.Printf("%s ... ", title)
	return err
}

func (c *TerminalLogger) PrintTestResult(passed bool) error {
	var val string
	if passed {
		val = color.New(color.FgGreen).Add(color.Bold).Sprint("PASSED")
	} else {
		val = color.New(color.FgRed).Add(color.Bold).Sprint("FAILED")
	}
	_, err := fmt.Println(val)
	return err
}

func (c *TerminalLogger) PrintAssertionResults(assertionResults []result.AssertionResult) error {
	if !c.debug {
		return nil
	}

	for _, assertionResult := range assertionResults {
		err := c.PrintAssertionResult(assertionResult)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TerminalLogger) PrintAssertionResult(assertionResult result.AssertionResult) error {
	var val string
	if assertionResult.Passed {
		val = color.New(color.FgGreen).Add(color.Bold).Sprint("PASSED")
	} else {
		val = color.New(color.FgRed).Add(color.Bold).Sprint("FAILED")
	}
	_, err := fmt.Printf("* %s ... %s\n", assertionResult.Title, val)
	if err == nil && !assertionResult.Passed {
		errMsg := color.New(color.FgRed).Add(color.Bold).Sprint(assertionResult.Error)
		_, err = fmt.Printf("  %s\n", errMsg)
	}
	return err
}

func (c *TerminalLogger) Print(args ...any) (int, error) {
	return fmt.Print(args...)
}

func (c *TerminalLogger) Println(args ...any) (int, error) {
	return fmt.Println(args...)
}

func (c *TerminalLogger) Printf(format string, args ...any) (int, error) {
	return fmt.Printf(format, args...)
}

func (c *TerminalLogger) Printfln(format string, args ...any) (int, error) {
	return fmt.Printf(fmt.Sprintf("%s%s", format, line.DetectLineSeparator()), args...)
}
