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
		return New(line.Line(caps[1]).Trim().Value(), caps[2]), true
	}
	if caps, ok := mdLine.CaptureAll(`\s*\*\s+([^:]+)\s*:\s*(.+)`); ok {
		return New(line.Line(caps[1]).Trim().Value(), line.Line(caps[2]).Trim().Value()), true
	}
	return nil, false
}

func New(name, selector string) *Capture {
	return &Capture{
		Name:     name,
		Selector: selector,
	}
}

func (c *Capture) SimilarTo(anotherCapture *Capture) bool {
	return c.Name == anotherCapture.Name &&
		c.Selector == anotherCapture.Selector
}

func (c *Capture) Clone() *Capture {
	return &Capture{
		Name:     c.Name,
		Selector: c.Selector,
	}
}
