package executable

import (
	"cotton/internal/capture"
	"cotton/internal/execution"
	"cotton/internal/logger"
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

func (ex *Executable) Execute(log logger.Logger) (*execution.Execution, error) {
	if log == nil {
		log = logger.NewNilLogger()
	}

	if ex.Request == nil {
		return nil, errors.New("no request to be made")
	}

	log.Printfln(" * %s", ex.Title)

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
