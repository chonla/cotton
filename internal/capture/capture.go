package capture

import (
	"cotton/internal/line"
)

type Captured struct {
	Name    string
	Locator string
}

func Try(mdLine line.Line) (*Captured, bool) {
	if captured, ok := mdLine.CaptureAll(`\s*\*\s+([^=]+)=` + "`([^`]+)`"); ok {
		return &Captured{
			Name:    captured[1],
			Locator: captured[2],
		}, true
	}
	if captured, ok := mdLine.CaptureAll(`\s*\*\s+([^=]+)=(.+)`); ok {
		return &Captured{
			Name:    captured[1],
			Locator: captured[2],
		}, true
	}
	return nil, false
}
