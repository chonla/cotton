package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/result"
	"cotton/internal/template"
	"cotton/internal/variable"
	"errors"
	"fmt"
	"slices"
)

type TestCaseOptions struct {
	Logger        logger.Logger
	RequestParser httphelper.RequestParser
}

// Test cases
type TestCase struct {
	options     *TestCaseOptions
	title       string
	description string
	reqRaw      string
	captures    []*capture.Capture
	setups      []*executable.Executable
	teardowns   []*executable.Executable
	assertions  []*assertion.Assertion
}

func New(title, description, reqRaw string, options *TestCaseOptions) *TestCase {
	return &TestCase{
		options:     options,
		title:       title,
		description: description,
		reqRaw:      reqRaw,
		captures:    []*capture.Capture{},
		setups:      []*executable.Executable{},
		teardowns:   []*executable.Executable{},
		assertions:  []*assertion.Assertion{},
	}
}

func (t *TestCase) Title() string {
	return t.title
}

func (t *TestCase) Description() string {
	return t.description
}

func (t *TestCase) RawRequest() string {
	return t.reqRaw
}

func (t *TestCase) Captures() []*capture.Capture {
	// return clone of captures
	clones := []*capture.Capture{}
	for _, cap := range t.captures {
		clones = append(clones, cap.Clone())
	}
	return clones
}

func (t *TestCase) AddCapture(cap *capture.Capture) {
	t.captures = append(t.captures, cap.Clone())
}

func (t *TestCase) Assertions() []*assertion.Assertion {
	// return clone of assertions
	clones := []*assertion.Assertion{}
	for _, assert := range t.assertions {
		clones = append(clones, assert.Clone())
	}
	return clones
}

func (t *TestCase) AddAssertion(assert *assertion.Assertion) {
	t.assertions = append(t.assertions, assert.Clone())
}

func (t *TestCase) Setups() []*executable.Executable {
	// return clone of executables
	clones := []*executable.Executable{}
	for _, setup := range t.setups {
		clones = append(clones, setup.Clone())
	}
	return clones
}

func (t *TestCase) AddSetup(setup *executable.Executable) {
	t.setups = append(t.setups, setup.Clone())
}

func (t *TestCase) Teardowns() []*executable.Executable {
	// return clone of executables
	clones := []*executable.Executable{}
	for _, teardown := range t.teardowns {
		clones = append(clones, teardown.Clone())
	}
	return clones
}

func (t *TestCase) AddTeardown(teardown *executable.Executable) {
	t.teardowns = append(t.teardowns, teardown.Clone())
}

func (t *TestCase) Execute(initialVars *variable.Variables) *result.TestResult {
	testResult := &result.TestResult{
		Title:      t.title,
		Passed:     false,
		Assertions: []result.AssertionResult{},
		Error:      nil,
	}

	if t.reqRaw == "" {
		testResult.Error = errors.New("no callable request")
		return testResult
	}

	t.options.Logger.PrintTestCaseTitle(t.title)

	sessionVars := initialVars.Clone()

	if len(t.setups) > 0 {
		t.options.Logger.PrintBlockTitle("Setups")
		for _, setup := range t.setups {
			execution, err := setup.Execute(sessionVars)
			if err != nil {
				testResult.Error = err
				return testResult
			}
			sessionVars = sessionVars.MergeWith(execution.Variables)
		}
	}

	reqTemplate := template.New(t.reqRaw)
	compiledRequest := reqTemplate.Apply(sessionVars)

	request, err := t.options.RequestParser.Parse(compiledRequest)
	if err != nil {
		testResult.Error = err
		return testResult
	}

	t.options.Logger.PrintBlockTitle("Execute test")
	resp, err := request.Do()
	if err != nil {
		testResult.Error = err
		return testResult
	}

	for _, cap := range t.captures {
		value, err := resp.ValueOf(cap.Selector)
		if err != nil {
			testResult.Error = err
			return testResult
		}
		sessionVars.Set(cap.Name, value)
	}

	for _, assertion := range t.assertions {
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

	if len(t.teardowns) > 0 {
		t.options.Logger.PrintBlockTitle("Teardowns")
		for _, teardown := range t.teardowns {
			execution, err := teardown.Execute(sessionVars)
			if err != nil {
				testResult.Error = err
				return testResult
			}
			sessionVars = sessionVars.MergeWith(execution.Variables)
		}
	}

	testResult.Passed = true
	testResult.Error = nil
	return testResult
}

func (t *TestCase) SimilarTo(anotherTestCase *TestCase) bool {
	return t.title == anotherTestCase.title &&
		t.description == anotherTestCase.description &&
		t.reqRaw == anotherTestCase.reqRaw &&
		slices.EqualFunc(t.captures, anotherTestCase.captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		slices.EqualFunc(t.setups, anotherTestCase.setups, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.teardowns, anotherTestCase.teardowns, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.assertions, anotherTestCase.assertions, func(a1, a2 *assertion.Assertion) bool {
			return a1.SimilarTo(a2)
		})
}
