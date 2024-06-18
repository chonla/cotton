package assertion

import (
	"cotton/internal/line"
	"cotton/internal/response"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/samber/mo"
)

type UndefinedOperator interface {
	Name() string
}

type UnaryAssertionOperator interface {
	Assert(actual *response.DataValue) (bool, error)
	Name() string
}

type BinaryAssertionOperator interface {
	Assert(expected interface{}, actual *response.DataValue) (bool, error)
	Name() string
}

type Assertion struct {
	Selector string
	Value    interface{}
	Operator mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator]
}

func Try(mdLine line.Line) (*Assertion, bool) {
	// binary assertion operator
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
	// unary assertion operator
	if caps, ok := mdLine.CaptureAll("\\s*\\*\\s+`([^`]+)`\\s*(.+)"); ok {
		op, err := New(line.Line(caps[2]).Trim().Value())
		if err == nil {
			return &Assertion{
				Selector: line.Line(caps[1]).Trim().Value(),
				Value:    nil,
				Operator: op,
			}, true
		}
	}
	return nil, false
}

func parseValue(mdLine line.Line) (interface{}, error) {
	if cap, ok := mdLine.Capture(`"(.+)"`, 1); ok {
		return cap, nil
	}
	if mdLine.LookLike(`^\d+$`) {
		// ALL numbers in JSON considered a floating point.
		v, err := strconv.ParseFloat(mdLine.Value(), 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if mdLine.LookLike(`^\d+\.\d+$`) {
		v, err := strconv.ParseFloat(mdLine.Value(), 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if mdLine.LookLike("true") {
		return true, nil
	}
	if mdLine.LookLike("false") {
		return false, nil
	}
	if mdLine.LookLike("null") {
		return nil, nil
	}
	return nil, errors.New("unexpected value")
}

func New(op string) (mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator], error) {
	operatorMap := map[string]func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator]{
		"==": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&EqAssertion{})
		},
		"!=": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&NeAssertion{})
		},
		">": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&GtAssertion{})
		},
		">=": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&GteAssertion{})
		},
		"<": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&LtAssertion{})
		},
		"<=": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&LteAssertion{})
		},
		"is undefined": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg2[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&UndefinedAssertion{})
		},
		"is defined": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg2[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&DefinedAssertion{})
		},
	}
	if ao, ok := operatorMap[op]; ok {
		return ao(), nil
	}
	return mo.NewEither3Arg1[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](nil), errors.New("unrecognized assertion")
}

func (a *Assertion) SimilarTo(anotherAssertion *Assertion) bool {
	if anotherAssertion == nil {
		return false
	}

	if a.Operator.IsArg1() && anotherAssertion.Operator.IsArg1() {
		return a.Selector == anotherAssertion.Selector &&
			reflect.TypeOf(a.Value) == reflect.TypeOf(a.Value) &&
			a.Value == anotherAssertion.Value
	}

	if a.Operator.IsArg2() && anotherAssertion.Operator.IsArg2() {
		return a.Selector == anotherAssertion.Selector &&
			reflect.TypeOf(a.Value) == reflect.TypeOf(a.Value) &&
			a.Value == anotherAssertion.Value &&
			a.Operator.MustArg2().Name() == anotherAssertion.Operator.MustArg2().Name()
	}

	if a.Operator.IsArg3() && anotherAssertion.Operator.IsArg3() {
		return a.Selector == anotherAssertion.Selector &&
			reflect.TypeOf(a.Value) == reflect.TypeOf(a.Value) &&
			a.Value == anotherAssertion.Value &&
			a.Operator.MustArg3().Name() == anotherAssertion.Operator.MustArg3().Name()
	}

	return false
}

func (a *Assertion) String() string {
	if a.Operator.IsArg1() {
		return "undefined operator"
	}
	// unary
	if a.Operator.IsArg2() {
		return fmt.Sprintf("%s %s", a.Selector, a.Operator.MustArg2().Name())
	}
	// binary
	if a.Operator.IsArg3() {
		return fmt.Sprintf("%s %s %v", a.Selector, a.Operator.MustArg3().Name(), a.Value)
	}
	return "unexpected operator"
}
