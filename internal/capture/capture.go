package capture

import (
	"cotton/internal/line"
)

type Capture struct {
	Name     string
	Selector string
}

func Try(mdLine line.Line) (*Capture, bool) {
	if caps, ok := mdLine.CaptureAll(`\s*\*\s+([^:]+)\s*:\s*` + "`([^`]+)`"); ok {
		return &Capture{
			Name:     line.Line(caps[1]).Trim().Value(),
			Selector: caps[2],
		}, true
	}
	if caps, ok := mdLine.CaptureAll(`\s*\*\s+([^:]+)\s*:\s*(.+)`); ok {
		return &Capture{
			Name:     line.Line(caps[1]).Trim().Value(),
			Selector: line.Line(caps[2]).Trim().Value(),
		}, true
	}
	return nil, false
}

func (c *Capture) SimilarTo(anotherCapture *Capture) bool {
	return c.Name == anotherCapture.Name &&
		c.Selector == anotherCapture.Selector
}
