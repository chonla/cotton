package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// TestResult represents the result of a test run
type TestResult struct {
	FilePath       string
	Title          string
	Passed         bool
	Failures       []string
	Errors         []string
	SetupTitles    []string
	TeardownTitles []string
	Request        string
	Response       string
	Variables      map[string]string
	StartMs        int64
	FinishMs       int64
	ElapsedMs      int64
}

// TestResults represents the aggregated results of multiple test runs
type TestResults struct {
	Tests   []*TestResult
	Summary TestSummary
}

// TestSummary contains summary statistics
type TestSummary struct {
	StartMs   int64 `json:"start_time"`
	FinishMs  int64 `json:"finish_time"`
	ElapsedMs int64 `json:"time_elapsed"`
	Total     int   `json:"total"`
	Passed    int   `json:"passed"`
	Failed    int   `json:"failed"`
}

// TestResultJSON represents a test result for JSON serialization
type TestResultJSON struct {
	FilePath  string `json:"file_path"`
	Title     string `json:"title"`
	Result    string `json:"result"`
	StartMs   int64  `json:"start_time"`
	FinishMs  int64  `json:"finish_time"`
	ElapsedMs int64  `json:"time_elapsed"`
}

// JSONReport represents the complete JSON report structure
type JSONReport struct {
	StartMs   int64            `json:"start_time"`
	FinishMs  int64            `json:"finish_time"`
	ElapsedMs int64            `json:"time_elapsed"`
	Total     int              `json:"total"`
	Passed    int              `json:"passed"`
	Failed    int              `json:"failed"`
	TestCases []TestResultJSON `json:"test_cases"`
}

// TestRunner orchestrates the test execution
type TestRunner struct {
	failFast bool
	execFn   func(string, *TestContext) (*TestResult, error)
	timeout  time.Duration
}

// TestContext holds the shared state across multiple spec executions
type TestContext struct {
	variables map[string]string
	executor  *HTTPExecutor
	visited   map[string]bool
	debug     bool
}

// NewTestRunner creates a new test runner
func NewTestRunner() *TestRunner {
	return &TestRunner{timeout: defaultHTTPTimeout}
}

// SetFailFast configures whether the runner should stop after the first failure
func (tr *TestRunner) SetFailFast(enabled bool) {
	tr.failFast = enabled
}

// SetTimeout configures the HTTP timeout for requests
func (tr *TestRunner) SetTimeout(timeout time.Duration) {
	if timeout <= 0 || timeout > maxHTTPTimeout {
		tr.timeout = defaultHTTPTimeout
		return
	}
	tr.timeout = timeout
}

// NewTestContext creates a new test context with fresh variables and executor
func NewTestContext(timeout time.Duration) *TestContext {
	if timeout <= 0 || timeout > maxHTTPTimeout {
		timeout = defaultHTTPTimeout
	}

	return &TestContext{
		variables: make(map[string]string),
		executor:  NewHTTPExecutorWithTimeout(timeout),
		visited:   make(map[string]bool),
	}
}

// RunSpecs runs multiple spec files and aggregates the results
func (tr *TestRunner) RunSpecs(specFiles []string) (*TestResults, error) {
	startTime := time.Now()

	results := &TestResults{
		Tests: make([]*TestResult, 0, len(specFiles)),
		Summary: TestSummary{
			StartMs: startTime.UnixMilli(),
		},
	}

	var encounteredErr error

	for _, specFile := range specFiles {
		exec := tr.executeSpec
		if tr.execFn != nil {
			exec = tr.execFn
		}

		context := NewTestContext(tr.timeout)
		testStartTime := time.Now()
		result, err := exec(specFile, context)
		testFinishTime := time.Now()

		if err != nil {
			// Continue with other tests even if one fails
			if result == nil {
				result = &TestResult{
					FilePath:  specFile,
					Title:     specFile,
					Passed:    false,
					Errors:    []string{fmt.Sprintf("Failed to execute spec: %v", err)},
					StartMs:   testStartTime.UnixMilli(),
					FinishMs:  testFinishTime.UnixMilli(),
					ElapsedMs: testFinishTime.Sub(testStartTime).Milliseconds(),
				}
			}
			if encounteredErr == nil {
				encounteredErr = err
			}
		}

		// Set timing information if not already set
		if result.StartMs == 0 {
			result.StartMs = testStartTime.UnixMilli()
			result.FinishMs = testFinishTime.UnixMilli()
			result.ElapsedMs = testFinishTime.Sub(testStartTime).Milliseconds()
		}
		result.FilePath = specFile

		results.Tests = append(results.Tests, result)

		// Honor fail-fast: stop as soon as a test fails (failures or errors)
		if tr.failFast && !result.Passed {
			break
		}
	}

	// Calculate summary
	finishTime := time.Now()
	results.Summary.FinishMs = finishTime.UnixMilli()
	results.Summary.ElapsedMs = finishTime.Sub(startTime).Milliseconds()
	results.Summary.Total = len(results.Tests)

	for _, test := range results.Tests {
		if test.Passed {
			results.Summary.Passed++
		} else {
			results.Summary.Failed++
		}
	}

	var runErr error
	if encounteredErr != nil {
		runErr = encounteredErr
	} else if results.Summary.Failed > 0 {
		runErr = newRunError("one or more specs failed")
	}

	return results, runErr
}

// RunSpec is the main entry point to execute a test case spec
func (tr *TestRunner) RunSpec(specFile string) (*TestResult, error) {
	context := NewTestContext(tr.timeout)
	return tr.executeSpec(specFile, context)
}

// executeSpec parses and runs a single spec file, updating the shared context
func (tr *TestRunner) executeSpec(specFile string, context *TestContext) (*TestResult, error) {
	startTime := time.Now()

	if context.visited[specFile] {
		return nil, newRunError(fmt.Sprintf("cycle detected involving %s", specFile))
	}
	context.visited[specFile] = true
	defer delete(context.visited, specFile)

	parser, err := NewMarkdownParser(specFile)
	if err != nil {
		return nil, newParseError(fmt.Sprintf("failed to read spec '%s': %v", specFile, err))
	}

	spec, err := parser.Parse()
	if err != nil {
		return nil, newParseError(fmt.Sprintf("failed to parse spec '%s': %v", specFile, err))
	}

	result := &TestResult{
		FilePath:       specFile,
		Title:          spec.Title,
		Failures:       []string{},
		Errors:         []string{},
		SetupTitles:    []string{},
		TeardownTitles: []string{},
		Variables:      make(map[string]string),
		StartMs:        startTime.UnixMilli(),
	}

	// 1. Execute Setup Specs sequentially
	for _, setupFile := range spec.SetupLinks {
		setupPath := resolvePath(specFile, setupFile)
		setupResult, err := tr.executeSpec(setupPath, context)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error in setup spec '%s': %v", setupFile, err))
			return result, err
		}
		if setupResult != nil {
			result.SetupTitles = append(result.SetupTitles, setupResult.Title)
		}
		if !setupResult.Passed {
			result.Errors = append(result.Errors, fmt.Sprintf("Setup spec '%s' did not pass", setupFile))
			return result, fmt.Errorf("setup spec failed")
		}
	}

	// 2. Resolve and set variables for the current spec
	if err := tr.resolveAndSetVariables(spec.Variables, context); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to resolve variables: %v", err))
		return result, err
	}

	// Copy current variables to result for debug output
	for k, v := range context.variables {
		result.Variables[k] = v
	}

	// 3. Execute the main request for the current spec
	result.Request = spec.RequestBody
	resp, err := context.executor.ExecuteRequest(spec.RequestBody)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result, err
	}
	defer resp.Body.Close()

	// Store response for debug output
	result.Response = context.executor.GetLastBody()

	// After the request, update context variables from the response body for subsequent steps
	tr.updateVariablesFromResponse(spec.Variables, context)

	// 4. Evaluate expectations ONLY if the spec is a test case
	if spec.IsTestCase {
		evaluator := NewEvaluator()
		failures, err := evaluator.EvaluateExpectations(resp, spec.Expectations)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to evaluate expectations: %v", err))
			return result, err
		}
		result.Failures = failures
	}

	// 5. Execute Teardown Specs sequentially
	for _, teardownFile := range spec.TeardownLinks {
		teardownPath := resolvePath(specFile, teardownFile)
		teardownResult, err := tr.executeSpec(teardownPath, context)
		if teardownResult != nil {
			result.TeardownTitles = append(result.TeardownTitles, teardownResult.Title)
		}
		if err != nil {
			// Log teardown error but don't fail the main test
			fmt.Printf("Warning: Error in teardown spec '%s': %v\n", teardownFile, err)
		}
	}

	finishTime := time.Now()
	result.FinishMs = finishTime.UnixMilli()
	result.ElapsedMs = result.FinishMs - result.StartMs
	result.Passed = len(result.Failures) == 0 && len(result.Errors) == 0
	return result, nil
}

// resolveAndSetVariables resolves any JSON path values and updates the executor's context
func (tr *TestRunner) resolveAndSetVariables(vars map[string]string, context *TestContext) error {
	lastBody := context.executor.GetLastBody()

	for k, v := range vars {
		if strings.HasPrefix(v, "$") {
			// This is a JSON path
			if lastBody == "" {
				// Cannot resolve path without a response yet
				context.variables[k] = v
			} else {
				// Use GJSON to extract the value (v already includes $ prefix)
				result := gjson.Get(lastBody, v)
				if result.Exists() {
					context.variables[k] = result.String()
				} else {
					// Path not found, keep original for potential later resolution
					context.variables[k] = v
				}
			}
		} else {
			// Immediate value
			context.variables[k] = v
		}
	}

	// Update the executor with the resolved set of variables
	context.executor.SetVariables(context.variables)
	return nil
}

// updateVariablesFromResponse updates the context with new variables defined in the current spec
// that source their value from the most recent response
func (tr *TestRunner) updateVariablesFromResponse(vars map[string]string, context *TestContext) {
	lastBody := context.executor.GetLastBody()
	if lastBody == "" {
		return
	}

	for k, v := range vars {
		if strings.HasPrefix(v, "$") {
			// Strip the JSONPath prefix for GJSON compatibility
			var path string
			if strings.HasPrefix(v, "$.") {
				path = strings.TrimPrefix(v, "$.")
			} else {
				path = strings.TrimPrefix(v, "$")
			}
			result := gjson.Get(lastBody, path)
			if result.Exists() {
				context.variables[k] = result.String()
			}
		}
	}

	// Re-set variables in executor to include newly captured ones
	context.executor.SetVariables(context.variables)
}

// resolvePath returns an absolute path for a linked spec file
func resolvePath(currentSpecFile, linkedFile string) string {
	if filepath.IsAbs(linkedFile) {
		return linkedFile
	}
	return filepath.Join(filepath.Dir(currentSpecFile), linkedFile)
}

// WriteJSONReport writes the test results as a JSON report to a file
func (tr *TestResults) WriteJSONReport(filename string) error {
	// Build test cases list
	testCases := make([]TestResultJSON, 0, len(tr.Tests))
	for _, test := range tr.Tests {
		resultStr := "PASSED"
		if !test.Passed {
			resultStr = "FAILED"
		}
		testCases = append(testCases, TestResultJSON{
			FilePath:  test.FilePath,
			Title:     test.Title,
			Result:    resultStr,
			StartMs:   test.StartMs,
			FinishMs:  test.FinishMs,
			ElapsedMs: test.ElapsedMs,
		})
	}

	// Build complete report
	report := JSONReport{
		StartMs:   tr.Summary.StartMs,
		FinishMs:  tr.Summary.FinishMs,
		ElapsedMs: tr.Summary.ElapsedMs,
		Total:     tr.Summary.Total,
		Passed:    tr.Summary.Passed,
		Failed:    tr.Summary.Failed,
		TestCases: testCases,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON report: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}
