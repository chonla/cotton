package logger

import (
	"cotton/internal/line"

	"fmt"

	"github.com/fatih/color"
)

type TerminalLogger struct{}

func NewTerminalLogger() Logger {
	return &TerminalLogger{}
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
