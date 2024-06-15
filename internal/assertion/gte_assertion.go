package assertion

import (
	"cotton/internal/kindof"
	"fmt"
	"reflect"
)

type GteAssertion struct {
}

func (a *GteAssertion) Name() string {
	return ">="
}

func (a *GteAssertion) Assert(expected, actual interface{}) (bool, error) {
	// actual >= expect? : greater than or equal to works only on numerical data type

	// try ints first
	actualInt, isActualInt := kindof.Int(actual)
	expectedInt, isExpectedInt := kindof.Int(expected)
	if isActualInt && isExpectedInt {
		if actualInt >= expectedInt {
			return true, nil
		}
		return false, fmt.Errorf("%v is expected to be greater than or equal to %v, but not", actualInt, expectedInt)
	}

	// try float
	actualFloat, isActualFloat := kindof.Float(actual)
	expectedFloat, isExpectedFloat := kindof.Float(expected)
	if isActualFloat && isExpectedFloat {
		if actualFloat >= expectedFloat {
			return true, nil
		}
		return false, fmt.Errorf("%v is expected to be greater than or equal to %v, but not", actualFloat, expectedFloat)
	}

	if !isExpectedFloat && !isExpectedInt {
		return false, fmt.Errorf("type of %v is expected to be number, but %v", expected, reflect.TypeOf(expected).Name())
	}

	return false, fmt.Errorf("type of %v is expected to be number, but %v", actual, reflect.TypeOf(actual).Name())
}
