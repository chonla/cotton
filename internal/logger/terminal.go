package logger

import (
	"cotton/internal/result"

	"fmt"

	"github.com/fatih/color"
)

type TerminalLogger struct {
	level LogLevel
}

func NewTerminalLogger(level LogLevel) Logger {
	return &TerminalLogger{
		level: level,
	}
}

func (c *TerminalLogger) PrintTestcaseTitle(title string) error {
	var err error
	if c.level == Compact {
		_, err = fmt.Printf("%s ... ", title)
	} else {
		_, err = fmt.Printf("%s\n", title)
	}

	return err
}

func (c *TerminalLogger) PrintExecutableTitle(title string) error {
	if c.level == Compact {
		return nil
	}

	_, err := fmt.Printf("  - %s\n", title)
	return err
}

func (c *TerminalLogger) PrintBlockTitle(title string) error {
	if c.level == Compact {
		return nil
	}

	_, err := fmt.Printf("* %s\n", title)
	return err
}

func (c *TerminalLogger) PrintTestResult(passed bool) error {
	if c.level == Compact {
		return nil
	}

	var val string
	if passed {
		val = color.New(color.FgGreen).Add(color.Bold).Sprint("PASSED")
	} else {
		val = color.New(color.FgRed).Add(color.Bold).Sprint("FAILED")
	}
	_, err := fmt.Printf("* Test result: %s\n", val)
	return err
}

func (c *TerminalLogger) PrintInlineTestResult(passed bool) error {
	if c.level != Compact {
		return nil
	}

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
	if c.level == Compact {
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
	_, err := fmt.Printf("  - %s ... %s\n", assertionResult.Title, val)
	if err == nil && !assertionResult.Passed {
		errMsg := color.New(color.FgRed).Add(color.Bold).Sprint(assertionResult.Error)
		_, err = fmt.Printf("  %s\n", errMsg)
	}
	return err
}

func (c *TerminalLogger) PrintRequest(req string) error {
	if c.level != Debug {
		return nil
	}

	val := color.New(color.FgWhite).Sprint(req)
	_, err := fmt.Printf("\n%s\n\n", val)
	return err
}
