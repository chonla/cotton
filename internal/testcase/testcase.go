package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/clock"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/template"
	"cotton/internal/variable"
	"errors"
	"fmt"
	"slices"
)

type TestcaseOptions struct {
	Logger          logger.Logger
	RequestParser   httphelper.RequestParser
	InsecureRequest bool
	ClockWrapper    clock.ClockWrapper
}

// Test cases
type Testcase struct {
	options     *TestcaseOptions
	title       string
	description string
	reqRaw      string
	captures    []*capture.Capture
	setups      []*executable.Executable
	teardowns   []*executable.Executable
	assertions  []*assertion.Assertion
	variables   *variable.Variables
}

func NewTestcase(title, description, reqRaw string, options *TestcaseOptions) *Testcase {
	return &Testcase{
		options:     options,
		title:       title,
		description: description,
		reqRaw:      reqRaw,
		captures:    []*capture.Capture{},
		setups:      []*executable.Executable{},
		teardowns:   []*executable.Executable{},
		assertions:  []*assertion.Assertion{},
		variables:   variable.New(),
	}
}

func (t *Testcase) Title() string {
	if t.title == "" {
		return "Untitled"
	}
	return t.title
}

func (t *Testcase) Description() string {
	return t.description
}

func (t *Testcase) RawRequest() string {
	return t.reqRaw
}

func (t *Testcase) Captures() []*capture.Capture {
	// return clone of captures
	clones := []*capture.Capture{}
	for _, cap := range t.captures {
		clones = append(clones, cap.Clone())
	}
	return clones
}

func (t *Testcase) AddCapture(cap *capture.Capture) {
	t.captures = append(t.captures, cap.Clone())
}

func (t *Testcase) Assertions() []*assertion.Assertion {
	// return clone of assertions
	clones := []*assertion.Assertion{}
	for _, assert := range t.assertions {
		clones = append(clones, assert.Clone())
	}
	return clones
}

func (t *Testcase) AddAssertion(assert *assertion.Assertion) {
	t.assertions = append(t.assertions, assert.Clone())
}

func (t *Testcase) Setups() []*executable.Executable {
	// return clone of executables
	clones := []*executable.Executable{}
	for _, setup := range t.setups {
		clones = append(clones, setup.Clone())
	}
	return clones
}

func (t *Testcase) AddSetup(setup *executable.Executable) {
	t.setups = append(t.setups, setup.Clone())
}

func (t *Testcase) Teardowns() []*executable.Executable {
	// return clone of executables
	clones := []*executable.Executable{}
	for _, teardown := range t.teardowns {
		clones = append(clones, teardown.Clone())
	}
	return clones
}

func (t *Testcase) AddTeardown(teardown *executable.Executable) {
	t.teardowns = append(t.teardowns, teardown.Clone())
}

func (t *Testcase) Execute(passedVars *variable.Variables) *result.TestResult {
	testResult := &result.TestResult{
		Title:        t.title,
		Passed:       false,
		Assertions:   []*result.AssertionResult{},
		Error:        nil,
		EllapsedTime: nil,
	}
	watch := stopwatch.New(t.options.ClockWrapper)
	watch.Start()
	defer (func() {
		testResult.EllapsedTime = watch.Stop()
	})()

	if t.reqRaw == "" {
		testResult.Error = errors.New("no callable request")
		return testResult
	}

	initialVars := t.variables.MergeWith(passedVars)

	sessionVars := initialVars.Clone()
	if len(t.setups) > 0 {
		for _, setup := range t.setups {
			t.options.Logger.PrintSectionTitle("setup")
			execution, err := setup.Execute(sessionVars)
			if err != nil {
				testResult.Error = err
				return testResult
			}
			sessionVars = sessionVars.MergeWith(execution.Variables)
		}
	}

	t.options.Logger.PrintSectionTitle("test")
	t.options.Logger.PrintTestcaseTitle(t.Title())

	t.options.Logger.PrintVariables(sessionVars)

	reqTemplate := template.New(t.reqRaw)
	compiledRequest := reqTemplate.Apply(sessionVars)

	request, err := t.options.RequestParser.Parse(compiledRequest)
	if err != nil {
		testResult.Error = err
		return testResult
	}

	t.options.Logger.PrintRequest(compiledRequest)
	resp, err := request.Do(t.options.InsecureRequest)
	if err != nil {
		testResult.Error = err
		return testResult
	}

	t.options.Logger.PrintResponse(resp.String())

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
			// unexpected assertion
			testResult.Error = errors.New("unexpected assertion found")
			return testResult
		}
		var passed bool
		if assertion.Operator.IsArg2() {
			// unary assertion assertion
			passed, err = assertion.Operator.MustArg2().Assert(actual)
		} else {
			if assertion.Operator.IsArg3() {
				// binary assertion assertion
				passed, err = assertion.Operator.MustArg3().Assert(expected, actual)
			}
		}
		assertionResult := &result.AssertionResult{
			Title:    assertion.String(),
			Passed:   passed,
			Actual:   fmt.Sprintf("%v", actual),
			Expected: fmt.Sprintf("%v", expected),
			Error:    err,
		}
		testResult.Assertions = append(testResult.Assertions, assertionResult)
		t.options.Logger.PrintSectionTitle("assert")
		t.options.Logger.PrintAssertionResult(assertionResult)
		if err != nil {
			testResult.Error = err
			return testResult
		}
	}

	if len(t.teardowns) > 0 {
		for _, teardown := range t.teardowns {
			t.options.Logger.PrintSectionTitle("teardown")
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

func (t *Testcase) SimilarTo(anotherTestcase *Testcase) bool {
	return t.title == anotherTestcase.title &&
		t.description == anotherTestcase.description &&
		t.reqRaw == anotherTestcase.reqRaw &&
		slices.EqualFunc(t.captures, anotherTestcase.captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		slices.EqualFunc(t.setups, anotherTestcase.setups, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.teardowns, anotherTestcase.teardowns, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.assertions, anotherTestcase.assertions, func(a1, a2 *assertion.Assertion) bool {
			return a1.SimilarTo(a2)
		})
}
