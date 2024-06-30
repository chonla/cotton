package assertion

import (
	"cotton/internal/line"
	"cotton/internal/value"
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/mo"
)

type UndefinedOperator interface {
	Name() string
}

type UnaryAssertionOperator interface {
	Assert(actual *value.Value) (bool, error)
	Name() string
}

type BinaryAssertionOperator interface {
	Assert(expected interface{}, actual *value.Value) (bool, error)
	Name() string
}

type Assertion struct {
	Selector string
	Value    interface{}
	Operator mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator]
}

func New(selector string, operator mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator], value interface{}) *Assertion {
	return &Assertion{
		Selector: selector,
		Operator: operator,
		Value:    value,
	}
}

func Try(mdLine line.Line) (*Assertion, bool) {
	// binary assertion operator
	if caps, ok := mdLine.CaptureAll("\\s*\\*\\s+`([^`]+)`\\s*(.+)\\s*`([^`]+)`"); ok {
		op, err := NewOp(line.Line(caps[2]).Trim().Value())
		if err == nil {
			value, err := line.Line(caps[3]).Trim().ReflectJSValue()
			if err == nil {
				return &Assertion{
					Selector: line.Line(caps[1]).Trim().Value(),
					Value:    value,
					Operator: op,
				}, true
			}
		}
	}
	// binary assertion operator for regex
	if caps, ok := mdLine.CaptureAll("\\s*\\*\\s+`([^`]+)`\\s*(.+)\\s*/(.+)/"); ok {
		op, err := NewRegexOp(line.Line(caps[2]).Trim().Value())
		if err == nil {
			value, err := line.Line(caps[3]).Trim().ReflectRegexValue()
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
		op, err := NewOp(line.Line(caps[2]).Trim().Value())
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

func NewRegexOp(op string) (mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator], error) {
	operatorMap := map[string]func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator]{
		"==": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&EqAssertion{})
		},
		"!=": func() mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator] {
			return mo.NewEither3Arg3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](&NeAssertion{})
		},
	}
	if ao, ok := operatorMap[op]; ok {
		return ao(), nil
	}
	return mo.NewEither3Arg1[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator](nil), errors.New("unrecognized assertion")
}

func NewOp(op string) (mo.Either3[UndefinedOperator, UnaryAssertionOperator, BinaryAssertionOperator], error) {
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

func (a *Assertion) Clone() *Assertion {
	return &Assertion{
		Selector: a.Selector,
		Operator: a.Operator,
		Value:    a.Value,
	}
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
		aType := reflect.TypeOf(a.Value)
		aValue := a.Value
		if aType != nil && aType.String() == "*regexp.Regexp" {
			aValue = fmt.Sprintf("/%v/", aValue)
		}
		return fmt.Sprintf("%s %s %v", a.Selector, a.Operator.MustArg3().Name(), aValue)
	}
	return "unexpected operator"
}
