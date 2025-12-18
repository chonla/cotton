package cli

import (
	"testing"

	"github.com/chonla/cotton/engine"
)

func TestResultFormatter_PrintResult_Passed(t *testing.T) {
	result := &engine.TestResult{
		Title:    "API Test",
		Passed:   true,
		Failures: []string{},
		Errors:   []string{},
	}

	formatter := NewResultFormatter()
	if formatter == nil {
		t.Errorf("NewResultFormatter() returned nil")
	}

	formatter.PrintResult(result)
}

func TestResultFormatter_PrintResult_WithFailures(t *testing.T) {
	result := &engine.TestResult{
		Title:    "API Test",
		Passed:   false,
		Failures: []string{"expectation 1 failed", "expectation 2 failed"},
		Errors:   []string{},
	}

	formatter := NewResultFormatter()
	if formatter == nil {
		t.Errorf("NewResultFormatter() returned nil")
	}

	formatter.PrintResult(result)
}

func TestResultFormatter_PrintResult_WithErrors(t *testing.T) {
	result := &engine.TestResult{
		Title:    "API Test",
		Passed:   false,
		Failures: []string{},
		Errors:   []string{"request failed", "invalid response"},
	}

	formatter := NewResultFormatter()
	if formatter == nil {
		t.Errorf("NewResultFormatter() returned nil")
	}

	formatter.PrintResult(result)
}

func TestResultFormatter_PrintResult_WithBothFailuresAndErrors(t *testing.T) {
	result := &engine.TestResult{
		Title:    "API Test",
		Passed:   false,
		Failures: []string{"expectation failed"},
		Errors:   []string{"error occurred"},
	}

	formatter := NewResultFormatter()
	if formatter == nil {
		t.Errorf("NewResultFormatter() returned nil")
	}

	formatter.PrintResult(result)
}

func TestNewResultFormatter(t *testing.T) {
	formatter := NewResultFormatter()
	if formatter == nil {
		t.Errorf("NewResultFormatter() returned nil")
	}

	if _, ok := interface{}(formatter).(*ResultFormatter); !ok {
		t.Errorf("NewResultFormatter() returned wrong type")
	}
}
