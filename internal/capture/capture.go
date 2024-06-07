package capture

import (
	"cotton/internal/line"
)

type Capture struct {
	Name    string
	Locator string
}

func Try(mdLine line.Line) (*Capture, bool) {
	if caps, ok := mdLine.CaptureAll(`\s*\*\s+([^=]+)\s*:\s*` + "`([^`]+)`"); ok {
		return &Capture{
			Name:    line.Line(caps[1]).Trim().Value(),
			Locator: caps[2],
		}, true
	}
	if caps, ok := mdLine.CaptureAll(`\s*\*\s+([^=]+)\s*:\s*(.+)`); ok {
		return &Capture{
			Name:    line.Line(caps[1]).Trim().Value(),
			Locator: line.Line(caps[2]).Trim().Value(),
		}, true
	}
	return nil, false
}

func (c *Capture) SimilarTo(anotherCapture *Capture) bool {
	return c.Name == anotherCapture.Name &&
		c.Locator == anotherCapture.Locator
}
