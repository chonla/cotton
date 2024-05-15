package executable

import (
	"cotton/internal/capture"
	"net/http"
)

// For setups and teardowns
type Executable struct {
	Title   string
	Request *http.Request

	Captures []*capture.Captured
}
