package executable

import (
	"cotton/internal/capture"
	"cotton/internal/execution"
	"cotton/internal/request"
	"errors"
	"net/http"
	"slices"
)

// For setups and teardowns
type Executable struct {
	Title   string
	Request *http.Request

	Captures []*capture.Capture
}

func (ex *Executable) Execute() (*execution.Execution, error) {
	if ex.Request == nil {
		return nil, errors.New("no request to be made")
	}

	ex.Request.Close = true

	_, err := http.DefaultClient.Do(ex.Request)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ex *Executable) SimilarTo(anotherEx *Executable) bool {
	return ex.Title == anotherEx.Title &&
		slices.EqualFunc(ex.Captures, anotherEx.Captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		request.Similar(ex.Request, anotherEx.Request)
}
