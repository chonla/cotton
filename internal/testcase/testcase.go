package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/logger"
	"cotton/internal/request"
	"cotton/internal/response"
	"cotton/internal/result"
	"cotton/internal/variable"
	"errors"
	"fmt"
	"slices"
)

// Test cases
type TestCase struct {
	Title       string
	Description string
	Request     *request.HTTPRequest

	Captures   []*capture.Capture
	Setups     []*executable.Executable
	Teardowns  []*executable.Executable
	Assertions []*assertion.Assertion
}

func (t *TestCase) Execute(log logger.Logger) *result.TestResult {
	if log == nil {
		log = logger.NewNilLogger(false)
	}

	testResult := &result.TestResult{
		Title:      t.Title,
		Passed:     false,
		Assertions: []result.AssertionResult{},
		Error:      nil,
	}

	if t.Request == nil {
		testResult.Error = errors.New("no request to be made")
		return testResult
	}

	log.PrintTestCaseTitle(t.Title)

	sessionVars := variable.New()

	for _, setup := range t.Setups {
		execution, err := setup.Execute(sessionVars, log)
		if err != nil {
			testResult.Error = err
			return testResult
		}
		sessionVars = sessionVars.MergeWith(execution.Variables)
	}

	r, err := t.Request.Do()
	if err != nil {
		testResult.Error = err
		return testResult
	}
	defer r.Body.Close()

	resp, err := response.New(r)
	if err != nil {
		testResult.Error = err
		return testResult
	}

	// capture response
	for _, cap := range t.Captures {
		value, err := resp.ValueOf(cap.Selector)
		if err != nil {
			testResult.Error = err
			return testResult
		}
		sessionVars.Set(cap.Name, value)
	}

	for _, assertion := range t.Assertions {
		actual, err := resp.ValueOf(assertion.Selector)
		if err != nil {
			testResult.Error = err
			return testResult
		}
		expected := assertion.Value
		if assertion.Operator.IsArg1() {
			testResult.Error = errors.New("unexpected assertion found")
			return testResult
		}
		var passed bool
		if assertion.Operator.IsArg2() {
			passed, err = assertion.Operator.MustArg2().Assert(actual)
		} else {
			if assertion.Operator.IsArg3() {
				passed, err = assertion.Operator.MustArg3().Assert(expected, actual)
			}
		}
		testResult.Assertions = append(testResult.Assertions, result.AssertionResult{
			Title:    assertion.String(),
			Passed:   passed,
			Actual:   fmt.Sprintf("%v", actual),
			Expected: fmt.Sprintf("%v", expected),
			Error:    err,
		})
		if err != nil {
			testResult.Error = err
			return testResult
		}
	}

	for _, teardown := range t.Teardowns {
		execution, err := teardown.Execute(sessionVars, log)
		if err != nil {
			testResult.Error = err
			return testResult
		}
		sessionVars = sessionVars.MergeWith(execution.Variables)
	}

	testResult.Passed = true
	testResult.Error = nil
	return testResult
}

func (t *TestCase) SimilarTo(anotherTestCase *TestCase) bool {
	return t.Title == anotherTestCase.Title &&
		t.Description == anotherTestCase.Description &&
		t.Request.Similar(anotherTestCase.Request) &&
		slices.EqualFunc(t.Captures, anotherTestCase.Captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		slices.EqualFunc(t.Setups, anotherTestCase.Setups, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.Teardowns, anotherTestCase.Teardowns, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.Assertions, anotherTestCase.Assertions, func(a1, a2 *assertion.Assertion) bool {
			return a1.SimilarTo(a2)
		})
}
