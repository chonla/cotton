package executable

import (
	"cotton/internal/capture"
	"cotton/internal/execution"
	"errors"
	"net/http"
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
	// defer resp.Body.Close()

	_, err := http.DefaultClient.Do(ex.Request)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
