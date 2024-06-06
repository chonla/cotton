package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/request"
	"errors"
	"net/http"
	"slices"
)

// Test cases
type TestCase struct {
	Title       string
	Description string
	Request     *http.Request

	Captures   []*capture.Capture
	Setups     []*executable.Executable
	Teardowns  []*executable.Executable
	Assertions []*assertion.Assertion
}

func (t *TestCase) Execute() error {
	for _, setup := range t.Setups {
		setup.Execute()
	}

	if t.Request == nil {
		return errors.New("no request to be made")
	}

	t.Request.Close = true
	resp, err := http.DefaultClient.Do(t.Request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, teardown := range t.Teardowns {
		teardown.Execute()
	}

	return nil
}

func (t *TestCase) SimilarTo(anotherTestCase *TestCase) bool {
	return t.Title == anotherTestCase.Title &&
		t.Description == anotherTestCase.Description &&
		request.Similar(t.Request, anotherTestCase.Request) &&
		slices.EqualFunc(t.Captures, anotherTestCase.Captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		slices.EqualFunc(t.Setups, anotherTestCase.Setups, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		}) &&
		slices.EqualFunc(t.Teardowns, anotherTestCase.Teardowns, func(s1, s2 *executable.Executable) bool {
			return s1.SimilarTo(s2)
		})
}
