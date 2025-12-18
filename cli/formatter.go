package cli

import (
	"fmt"

	"github.com/chonla/cotton/engine"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

// ResultFormatter handles formatting and output of test results
type ResultFormatter struct {
	debug      bool
	quiet      bool
	monochrome bool
}

// NewResultFormatter creates a new result formatter
func NewResultFormatter() *ResultFormatter {
	return &ResultFormatter{}
}

// SetDebug sets debug mode
func (f *ResultFormatter) SetDebug(debug bool) {
	f.debug = debug
}

// SetQuiet sets quiet mode
func (f *ResultFormatter) SetQuiet(quiet bool) {
	f.quiet = quiet
}

// SetMonochrome sets monochrome mode
func (f *ResultFormatter) SetMonochrome(monochrome bool) {
	f.monochrome = monochrome
}

// PrintResults prints the test results to console
func (f *ResultFormatter) PrintResults(results *engine.TestResults) {
	if f.quiet {
		// In quiet mode, print nothing
		return
	}

	// Print individual test results
	for _, result := range results.Tests {
		f.PrintResult(result)
	}

	// Print summary
	fmt.Printf("\n%s\n", f.colorize("=== Test Summary ===", ColorBlue))
	fmt.Printf("Total: %d\n", results.Summary.Total)
	fmt.Printf("Passed: %s\n", f.colorize(fmt.Sprintf("%d", results.Summary.Passed), ColorGreen))
	fmt.Printf("Failed: %s\n", f.colorize(fmt.Sprintf("%d", results.Summary.Failed), ColorRed))
	fmt.Printf("Time: %dms\n", results.Summary.ElapsedMs)
}

// PrintResult prints a single test result to console
func (f *ResultFormatter) PrintResult(result *engine.TestResult) {
	if f.quiet {
		return
	}

	fmt.Printf("\n%s\n", f.colorize("=== Test Result ===", ColorBlue))
	fmt.Printf("Test: %s\n", result.Title)

	if result.Passed {
		fmt.Printf("Status: %s\n", f.colorize("PASSED", ColorGreen))
	} else {
		fmt.Printf("Status: %s\n", f.colorize("FAILED", ColorRed))
	}

	if f.debug {
		// In debug mode, print setup and teardown info
		if len(result.SetupTitles) > 0 {
			fmt.Println("\nSetup:")
			for _, title := range result.SetupTitles {
				fmt.Printf("  - %s\n", title)
			}
		}

		if len(result.TeardownTitles) > 0 {
			fmt.Println("\nTeardown:")
			for _, title := range result.TeardownTitles {
				fmt.Printf("  - %s\n", title)
			}
		}

		// Print request/response if available
		if result.Request != "" {
			fmt.Println("\nRequest:")
			fmt.Println(result.Request)
		}

		if result.Response != "" {
			fmt.Println("\nResponse:")
			fmt.Println(result.Response)
		}

		// Print variables
		if len(result.Variables) > 0 {
			fmt.Println("\nVariables:")
			for k, v := range result.Variables {
				fmt.Printf("  %s = %s\n", k, v)
			}
		}
	}

	if len(result.Failures) > 0 {
		fmt.Println(f.colorize("\nFailures:", ColorRed))
		for _, failure := range result.Failures {
			fmt.Printf("  - %s\n", failure)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println(f.colorize("\nErrors:", ColorRed))
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}
	}
}

// colorize applies color to text if monochrome mode is disabled
func (f *ResultFormatter) colorize(text, color string) string {
	if f.monochrome {
		return text
	}
	return color + text + ColorReset
}
