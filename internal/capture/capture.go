package capture

import (
	"cotton/internal/line"
)

type Captured struct {
	Name    string
	Locator string
}

func Try(mdLine line.Line) (*Captured, bool) {
	if captured, ok := mdLine.CaptureAll(`\s*\*\s+([^=]+)\s*=\s*` + "`([^`]+)`"); ok {
		return &Captured{
			Name:    line.Line(captured[1]).Trim().Value(),
			Locator: captured[2],
		}, true
	}
	if captured, ok := mdLine.CaptureAll(`\s*\*\s+([^=]+)\s*=\s*(.+)`); ok {
		return &Captured{
			Name:    line.Line(captured[1]).Trim().Value(),
			Locator: line.Line(captured[2]).Trim().Value(),
		}, true
	}
	return nil, false
}
