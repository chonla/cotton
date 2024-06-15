package assertion

import (
	"cotton/internal/kindof"
	"fmt"
	"reflect"
)

type LtAssertion struct {
}

func (a *LtAssertion) Name() string {
	return "<"
}

func (a *LtAssertion) Assert(expected, actual interface{}) (bool, error) {
	// actual < expect? : less than works only on numerical data type

	// try ints first
	actualInt, isActualInt := kindof.Int(actual)
	expectedInt, isExpectedInt := kindof.Int(expected)
	if isActualInt && isExpectedInt {
		if actualInt < expectedInt {
			return true, nil
		}
		return false, fmt.Errorf("%v is expected to be less than %v, but not", actualInt, expectedInt)
	}

	// try float
	actualFloat, isActualFloat := kindof.Float(actual)
	expectedFloat, isExpectedFloat := kindof.Float(expected)
	if isActualFloat && isExpectedFloat {
		if actualFloat < expectedFloat {
			return true, nil
		}
		return false, fmt.Errorf("%v is expected to be less than %v, but not", actualFloat, expectedFloat)
	}

	if !isExpectedFloat && !isExpectedInt {
		return false, fmt.Errorf("type of %v is expected to be number, but %v", expected, reflect.TypeOf(expected).Name())
	}

	return false, fmt.Errorf("type of %v is expected to be number, but %v", actual, reflect.TypeOf(actual).Name())
}
