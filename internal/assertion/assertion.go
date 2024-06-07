package assertion

import (
	"cotton/internal/line"
	"errors"
	"reflect"
	"strconv"
)

type AssertionOperator interface {
	Assert(actual, expected interface{}) (bool, error)
	Name() string
}

type Assertion struct {
	Selector string
	Value    interface{}
	Operator AssertionOperator
}

func Try(mdLine line.Line) (*Assertion, bool) {
	if caps, ok := mdLine.CaptureAll("\\s*\\*\\s+`([^`]+)`\\s*(.+)\\s*`([^`]+)`"); ok {
		op, err := New(line.Line(caps[2]).Trim().Value())
		if err == nil {
			value, err := parseValue(line.Line(caps[3]).Trim())
			if err == nil {
				return &Assertion{
					Selector: line.Line(caps[1]).Trim().Value(),
					Value:    value,
					Operator: op,
				}, true
			}
		}
	}
	return nil, false
}

func parseValue(mdLine line.Line) (interface{}, error) {
	if cap, ok := mdLine.Capture(`"(.+)"`, 1); ok {
		return cap, nil
	}
	if mdLine.LookLike(`^\d+$`) {
		return strconv.Atoi(mdLine.Value())
	}
	if mdLine.LookLike("null") {
		return nil, nil
	}
	return nil, errors.New("unexpected value")
}

func New(op string) (AssertionOperator, error) {
	if op == "==" {
		return &EqualAssertion{}, nil
	}
	return nil, errors.New("unrecognized assertion")
}

func (a *Assertion) SimilarTo(anotherAssertion *Assertion) bool {
	if anotherAssertion == nil {
		return false
	}

	return a.Selector == anotherAssertion.Selector &&
		reflect.TypeOf(a.Value) == reflect.TypeOf(a.Value) &&
		a.Value == anotherAssertion.Value &&
		a.Operator.Name() == anotherAssertion.Operator.Name()
}
