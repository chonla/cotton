package logger

import (
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/variable"
	"strings"

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
		_, err = fmt.Printf("%s", title)
	} else {
		_, err = fmt.Printf("%s\n", title)
	}

	return err
}

func (c *TerminalLogger) PrintTestcaseTitleWithPath(title, path string) error {
	var err error
	if c.level == Compact {
		_, err = fmt.Printf("%s (%s)", color.New(color.FgHiWhite).Sprint(title), color.New(color.FgWhite).Sprint(path))
	} else {
		_, err = fmt.Printf("%s (%s)\n", color.New(color.FgHiWhite).Sprint(title), color.New(color.FgWhite).Sprint(path))
	}

	return err
}

func (c *TerminalLogger) PrintSectionTitle(sectionTitle string) error {
	if c.level == Compact {
		return nil
	}

	val := color.New(color.FgWhite).Sprintf("[%s] ", sectionTitle)
	_, err := fmt.Print(val)

	return err
}

func (c *TerminalLogger) PrintExecutableTitle(title string) error {
	if c.level == Compact {
		return nil
	}

	_, err := fmt.Println(title)
	return err
}

func (c *TerminalLogger) PrintSectionedMessage(section, message string) error {
	if c.level == Compact {
		return nil
	}

	c.PrintSectionTitle(section)
	_, e := fmt.Println(message)
	return e
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
	_, err := fmt.Println(val)
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
	openBracket := color.New(color.FgWhite).Sprint("[")
	closeBracket := color.New(color.FgWhite).Sprint("]")
	_, err := fmt.Printf(" %s%s%s ", openBracket, val, closeBracket)
	return err
}

func (c *TerminalLogger) PrintAssertionResults(assertionResults []*result.AssertionResult) error {
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

func (c *TerminalLogger) PrintAssertionResult(assertionResult *result.AssertionResult) error {
	if c.level == Compact {
		return nil
	}

	var val string
	if assertionResult.Passed {
		val = color.New(color.FgGreen).Sprint("passed")
	} else {
		val = color.New(color.FgRed).Sprint("failed")
	}
	_, err := fmt.Printf("%s...%s\n", assertionResult.Title, val)
	if err == nil && !assertionResult.Passed {
		errMsg := color.New(color.FgRed).Add(color.Bold).Sprint(assertionResult.Error)
		_, err = fmt.Printf("  %s\n", errMsg)
	}
	return err
}

func (c *TerminalLogger) PrintRequest(req string) error {
	if c.level != Debug && c.level != DetailedDebug {
		return nil
	}

	c.PrintSectionTitle("request")
	fmt.Println("")
	val := color.New(color.FgBlue).Sprint(req)
	_, err := fmt.Println(val)
	return err
}

func (c *TerminalLogger) PrintResponse(resp string) error {
	if c.level != Debug && c.level != DetailedDebug {
		return nil
	}

	c.PrintSectionTitle("response")
	fmt.Println("")
	val := color.New(color.FgBlue).Sprint(resp)
	_, err := fmt.Println(val)
	return err
}

func (c *TerminalLogger) PrintError(fileContext string, err error) error {
	section := color.New(color.FgRed).Sprintf("[error] ")
	file := color.New(color.FgWhite).Sprintf("%s", fileContext)
	val := color.New(color.FgRed).Sprintf("%s", err)
	message := fmt.Sprintf("%s%s in %s", section, val, file)
	_, err = fmt.Println(message)
	return err
}

func (c *TerminalLogger) PrintTestcaseSequence(index, total int) error {
	val := color.New(color.FgWhite).Sprintf("[testcase %d/%d] ", index, total)
	_, err := fmt.Print(val)
	return err
}

func (c *TerminalLogger) PrintTestsuiteResult(testsuiteResult *result.TestsuiteResult) error {
	_, err := fmt.Println(color.New(color.FgWhite).Sprintf("-------------------------"))
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s\n", c.buildFieldValue("Testcases executed", fmt.Sprintf("%d/%d", testsuiteResult.ExecutionsCount, testsuiteResult.TestCount)))
	if err != nil {
		return err
	}
	passedPercentage := 0.0
	failedPercentage := 0.0
	skippedPerecentage := 0.0
	if testsuiteResult.TestCount > 0 {
		passedPercentage = float64(testsuiteResult.PassedCount) * 100.0 / float64(testsuiteResult.TestCount)
		failedPercentage = float64(testsuiteResult.FailedCount) * 100.0 / float64(testsuiteResult.TestCount)
		skippedPerecentage = float64(testsuiteResult.SkippedCount) * 100.0 / float64(testsuiteResult.TestCount)
	}
	_, err = fmt.Printf("%s\n", c.buildFieldValue("Passed", fmt.Sprintf("%d (%0.2f%%)", testsuiteResult.PassedCount, passedPercentage)))
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s\n", c.buildFieldValue("Failed", fmt.Sprintf("%d (%0.2f%%)", testsuiteResult.FailedCount, failedPercentage)))
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s\n", c.buildFieldValue("Skipped", fmt.Sprintf("%d (%0.2f%%)", testsuiteResult.SkippedCount, skippedPerecentage)))
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s\n", c.buildFieldValue("Ellapsed", testsuiteResult.EllapsedTime))
	return err
}

func (c *TerminalLogger) buildFieldValue(label string, value interface{}) string {
	labelData := color.New(color.FgWhite).Sprintf("%s: ", label)
	valueData := color.New(color.FgHiWhite).Sprintf("%v", value)
	return fmt.Sprintf("%s%s", labelData, valueData)
}

func (c *TerminalLogger) PrintDebugMessage(message string) error {
	if c.level != Debug && c.level != DetailedDebug {
		return nil
	}

	return c.PrintSectionedMessage("debug", message)
}

func (c *TerminalLogger) PrintDetailedDebugMessage(messages ...string) error {
	if c.level != DetailedDebug {
		return nil
	}

	wrappedMessage := ""

	switch len(messages) {
	case 0:
		wrappedMessage = ""
	case 1:
		wrappedMessage = color.New(color.FgWhite).Sprintf("%s", messages[0])
	default:
		wrappedDetailed := color.New(color.FgHiWhite).Sprintf("%s", strings.Join(messages[1:], " "))
		wrappedLabel := color.New(color.FgWhite).Sprintf("%s", messages[0])
		wrappedMessage = fmt.Sprintf("%s %s", wrappedLabel, wrappedDetailed)
	}

	return c.PrintSectionedMessage("detailed-debug", wrappedMessage)
}

func (c *TerminalLogger) PrintTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error {
	if c.level == Compact {
		return nil
	}

	val := color.New(color.FgWhite).Sprint(ellapsedTime)
	_, err := fmt.Println(val)
	return err
}

func (c *TerminalLogger) PrintInlineTimeEllapsed(ellapsedTime *stopwatch.EllapsedTime) error {
	if c.level != Compact {
		return nil
	}

	val := color.New(color.FgWhite).Sprintf("(%s)", ellapsedTime)
	_, err := fmt.Println(val)
	return err
}

func (c *TerminalLogger) PrintVariables(variables *variable.Variables) error {
	if c.level != Debug {
		return nil
	}

	names := variables.Names()
	for _, name := range names {
		value, _ := variables.ValueOf(name)
		c.PrintSectionedMessage("variable", fmt.Sprintf("%s=%v", name, value))
	}
	return nil
}
