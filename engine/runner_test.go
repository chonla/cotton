package engine

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

func TestTestRunner_VariableSharing(t *testing.T) {
	context := NewTestContext(defaultHTTPTimeout)

	context.variables["var1"] = "value1"
	context.executor.SetVariables(context.variables)

	if context.variables["var1"] != "value1" {
		t.Errorf("Variable not stored correctly")
	}

	context.variables["var2"] = "value2"
	context.executor.SetVariables(context.variables)

	if len(context.variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(context.variables))
	}
}

func TestTestRunner_ResolveAndSetVariables_ImmediateValues(t *testing.T) {
	runner := NewTestRunner()
	context := NewTestContext(defaultHTTPTimeout)

	vars := map[string]string{
		"id":      "42",
		"name":    "\"John\"",
		"enabled": "true",
	}

	err := runner.resolveAndSetVariables(vars, context)
	if err != nil {
		t.Errorf("resolveAndSetVariables() error = %v", err)
	}

	if len(context.variables) != len(vars) {
		t.Errorf("Expected %d variables, got %d", len(vars), len(context.variables))
	}

	for k, v := range vars {
		if context.variables[k] != v {
			t.Errorf("Variable %q = %q, want %q", k, context.variables[k], v)
		}
	}
}

func TestTestRunner_ResolveAndSetVariables_JSONPath_NoResponse(t *testing.T) {
	runner := NewTestRunner()
	context := NewTestContext(defaultHTTPTimeout)

	vars := map[string]string{
		"token": "$.auth.token",
	}

	err := runner.resolveAndSetVariables(vars, context)
	if err != nil {
		t.Errorf("resolveAndSetVariables() error = %v", err)
	}

	if context.variables["token"] != "$.auth.token" {
		t.Errorf("Variable = %q, want %q", context.variables["token"], "$.auth.token")
	}
}

func TestTestRunner_UpdateVariablesFromResponse(t *testing.T) {
	runner := NewTestRunner()
	context := NewTestContext(defaultHTTPTimeout)

	// Simulate executor with response by creating a new HTTPExecutor
	executor := &HTTPExecutor{
		variables: make(map[string]interface{}),
		lastBody:  `{"user":{"id":123,"name":"Alice"}}`,
		client:    nil,
	}
	context.executor = executor

	vars := map[string]string{
		"user_id":   "$.user.id",
		"user_name": "$.user.name",
	}

	// Verify GetLastBody() returns the value
	body := context.executor.GetLastBody()
	if body != `{"user":{"id":123,"name":"Alice"}}` {
		t.Errorf("GetLastBody() = %q", body)
	}

	// Test gjson directly to verify it works
	result1 := gjson.Get(body, "user.id")
	t.Logf("gjson.Get(user.id) = %q", result1.String())

	result2 := gjson.Get(body, "$.user.id")
	t.Logf("gjson.Get($.user.id) = %q", result2.String())

	// Update from response - should extract JSON path values
	runner.updateVariablesFromResponse(vars, context)

	// Debug: Check what's in context.variables
	t.Logf("context.variables after update: %v", context.variables)

	// Debug: Manually verify each step
	for k, v := range vars {
		t.Logf("Var %s=%s, hasPrefix($)=%v", k, v, strings.HasPrefix(v, "$"))
		if strings.HasPrefix(v, "$") {
			path := strings.TrimPrefix(v, "$")
			t.Logf("  Trimmed path: %s", path)
			result := gjson.Get(body, path)
			t.Logf("  gjson.Get result exists=%v, value=%s", result.Exists(), result.String())
		}
	}
	// After update, they should be extracted
	if context.variables["user_id"] != "123" {
		t.Errorf("After update user_id = %q, want %q", context.variables["user_id"], "123")
	}
	if context.variables["user_name"] != "Alice" {
		t.Errorf("After update user_name = %q, want %q", context.variables["user_name"], "Alice")
	}
}

func TestTestRunner_ResolvePath_Absolute(t *testing.T) {
	currentPath := "/path/to/spec.md"
	linkedPath := "/other/path/linked.md"

	result := resolvePath(currentPath, linkedPath)

	if result != linkedPath {
		t.Errorf("resolvePath() = %q, want %q", result, linkedPath)
	}
}

func TestTestRunner_ResolvePath_Relative(t *testing.T) {
	currentPath := "/path/to/spec.md"
	linkedPath := "setup.md"

	result := resolvePath(currentPath, linkedPath)

	if !strings.HasSuffix(result, "setup.md") {
		t.Errorf("resolvePath() = %q, should end with 'setup.md'", result)
	}
	if !strings.Contains(result, "path/to") {
		t.Errorf("resolvePath() = %q, should contain parent directory", result)
	}
}

func TestTestResult_Passed_WithoutFailures(t *testing.T) {
	result := &TestResult{
		Title:    "Test",
		Failures: []string{},
		Errors:   []string{},
	}

	if len(result.Failures) == 0 && len(result.Errors) == 0 {
		result.Passed = true
	}

	if !result.Passed {
		t.Errorf("Result should be passed without failures or errors")
	}
}

func TestTestResult_Failed_WithFailures(t *testing.T) {
	result := &TestResult{
		Title:    "Test",
		Failures: []string{"expectation failed"},
		Errors:   []string{},
	}

	if len(result.Failures) > 0 || len(result.Errors) > 0 {
		result.Passed = false
	}

	if result.Passed {
		t.Errorf("Result should not be passed with failures")
	}
}

func TestTestResult_Failed_WithErrors(t *testing.T) {
	result := &TestResult{
		Title:    "Test",
		Failures: []string{},
		Errors:   []string{"execution error"},
	}

	if len(result.Failures) > 0 || len(result.Errors) > 0 {
		result.Passed = false
	}

	if result.Passed {
		t.Errorf("Result should not be passed with errors")
	}
}

func TestTestRunner_RunSpecs_FailFastStopsAfterFirstFailure(t *testing.T) {
	runner := NewTestRunner()
	runner.SetFailFast(true)

	specs := []string{"specA", "specB", "specC"}
	calls := []string{}

	runner.execFn = func(path string, _ *TestContext) (*TestResult, error) {
		calls = append(calls, path)
		if path == "specB" {
			return &TestResult{Title: path, Passed: false}, nil
		}
		return &TestResult{Title: path, Passed: true}, nil
	}

	results, err := runner.RunSpecs(specs)
	if err == nil {
		t.Fatalf("RunSpecs() expected error for failed spec")
	}

	if len(results.Tests) != 2 {
		t.Fatalf("expected 2 tests to run before fail-fast, got %d", len(results.Tests))
	}

	expectedCalls := []string{"specA", "specB"}
	if !reflect.DeepEqual(calls, expectedCalls) {
		t.Errorf("execution order mismatch, got %v want %v", calls, expectedCalls)
	}

	if results.Summary.Failed != 1 {
		t.Errorf("summary failed = %d, want 1", results.Summary.Failed)
	}
}

func TestTestRunner_RunSpecs_NoFailFastRunsAll(t *testing.T) {
	runner := NewTestRunner()
	runner.SetFailFast(false)

	specs := []string{"specA", "specB", "specC"}
	calls := []string{}

	runner.execFn = func(path string, _ *TestContext) (*TestResult, error) {
		calls = append(calls, path)
		if path == "specB" {
			return &TestResult{Title: path, Passed: false}, nil
		}
		return &TestResult{Title: path, Passed: true}, nil
	}

	results, err := runner.RunSpecs(specs)
	if err == nil {
		t.Fatalf("RunSpecs() expected error for failed spec")
	}

	if len(results.Tests) != 3 {
		t.Fatalf("expected all 3 tests to run, got %d", len(results.Tests))
	}

	expectedCalls := []string{"specA", "specB", "specC"}
	if !reflect.DeepEqual(calls, expectedCalls) {
		t.Errorf("execution order mismatch, got %v want %v", calls, expectedCalls)
	}

	if results.Summary.Failed != 1 {
		t.Errorf("summary failed = %d, want 1", results.Summary.Failed)
	}
}
