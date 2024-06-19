package executable

import (
	"cotton/internal/capture"
	"cotton/internal/execution"
	"cotton/internal/logger"
	"cotton/internal/request"
	"cotton/internal/response"
	"cotton/internal/variable"
	"errors"
	"slices"
)

// For setups and teardowns
type Executable struct {
	Title   string
	Request request.Request

	Captures []*capture.Capture
}

func (ex *Executable) Execute(log logger.Logger) (*execution.Execution, error) {
	if log == nil {
		log = logger.NewNilLogger(false)
	}

	if ex.Request == nil {
		return nil, errors.New("no request to be made")
	}

	log.Printfln(" * %s", ex.Title)

	r, err := ex.Request.Do()
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	resp, err := response.New(r)
	if err != nil {
		return nil, err
	}

	vars := variable.New()
	for _, cap := range ex.Captures {
		value, err := resp.ValueOf(cap.Selector)
		if err != nil {
			return nil, err
		}
		vars.Set(cap.Name, value)
	}

	return &execution.Execution{
		Variables: vars,
	}, nil
}

func (ex *Executable) SimilarTo(anotherEx *Executable) bool {
	return ex.Title == anotherEx.Title &&
		slices.EqualFunc(ex.Captures, anotherEx.Captures, func(c1, c2 *capture.Capture) bool {
			return c1.SimilarTo(c2)
		}) &&
		ex.Request.Similar(anotherEx.Request)
}
