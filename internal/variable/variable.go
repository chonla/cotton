package variable

import (
	"cotton/internal/line"
)

type Variable struct {
	Name  string
	Value interface{}
}

func Try(mdLine line.Line) (*Variable, bool) {
	if caps, ok := mdLine.CaptureAll(`\s*\*\s+([^:]+)\s*:\s*(.+)`); ok {
		value, err := line.Line(caps[2]).Trim().ReflectJSValue()
		if err != nil {
			// return as it is
			value = line.Line(caps[2]).Trim().Value()
		}
		v := &Variable{
			Name:  line.Line(caps[1]).Trim().Value(),
			Value: value,
		}
		return v, true
	}
	return nil, false
}
