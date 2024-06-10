package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/request"
	"errors"
	"fmt"
	"slices"

	"github.com/kr/pretty"
)

// Test cases
type TestCase struct {
	Title       string
	Description string
	Request     *request.Request

	Captures   []*capture.Capture
	Setups     []*executable.Executable
	Teardowns  []*executable.Executable
	Assertions []*assertion.Assertion
}

func (t *TestCase) Execute() error {
	if t.Request == nil {
		return errors.New("no request to be made")
	}

	for _, setup := range t.Setups {
		_, err := setup.Execute()
		if err != nil {
			return err
		}
	}

	resp, err := t.Request.Do()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, teardown := range t.Teardowns {
		_, err := teardown.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TestCase) SimilarTo(anotherTestCase *TestCase) bool {
	fmt.Println("Compare Captures")
	fmt.Println(slices.EqualFunc(t.Captures, anotherTestCase.Captures, func(c1, c2 *capture.Capture) bool {
		return c1.SimilarTo(c2)
	}))
	fmt.Println(pretty.Diff(t.Captures, anotherTestCase.Captures))
	fmt.Println("Compare Setups")
	fmt.Println(slices.EqualFunc(t.Setups, anotherTestCase.Setups, func(s1, s2 *executable.Executable) bool {
		return s1.SimilarTo(s2)
	}))
	fmt.Println(pretty.Diff(t.Setups, anotherTestCase.Setups))
	fmt.Println("Compare Teardowns")
	fmt.Println(slices.EqualFunc(t.Teardowns, anotherTestCase.Teardowns, func(s1, s2 *executable.Executable) bool {
		return s1.SimilarTo(s2)
	}))
	fmt.Println(pretty.Diff(t.Teardowns, anotherTestCase.Teardowns))
	fmt.Println("Compare Assertions")
	fmt.Println(slices.EqualFunc(t.Assertions, anotherTestCase.Assertions, func(a1, a2 *assertion.Assertion) bool {
		return a1.SimilarTo(a2)
	}))
	fmt.Println(pretty.Diff(t.Assertions, anotherTestCase.Assertions))

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
