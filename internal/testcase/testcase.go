package testcase

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"errors"
	"net/http"
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
	return nil

}
