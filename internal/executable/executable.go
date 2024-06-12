package executable

import (
	"cotton/internal/capture"
	"cotton/internal/console"
	"cotton/internal/execution"
	"cotton/internal/request"
	"errors"
	"slices"
)

// For setups and teardowns
type Executable struct {
	Title   string
	Request *request.Request

	Captures []*capture.Capture
}

func (ex *Executable) Execute(logger console.Console) (*execution.Execution, error) {
	if logger == nil {
		logger = console.NewNilConsole()
	}

	if ex.Request == nil {
		return nil, errors.New("no request to be made")
	}

	logger.Printfln(" * %s", ex.Title)

	_, err := ex.Request.Do()
	if err != nil {
		return nil, err
	}
	return &execution.Execution{}, nil
}

func (ex *Executable) SimilarTo(anotherEx *Executable) bool {
	return ex.Title == anotherEx.Title &&
		slices.EqualFunc(ex.Captures, anotherEx.Captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		ex.Request.Similar(anotherEx.Request)
	// request.Similar(ex.Request, anotherEx.Request)
}
