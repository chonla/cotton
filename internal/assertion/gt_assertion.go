package assertion

import (
	"cotton/internal/kindof"
	"fmt"
	"reflect"
)

type GtAssertion struct {
}

func (a *GtAssertion) Name() string {
	return ">"
}

func (a *GtAssertion) Assert(expected, actual interface{}) (bool, error) {
	// actual > expect?
	// greater than works only on numerical data type

	// try ints first
	actualInt, isActualInt := kindof.Int(actual)
	expectedInt, isExpectedInt := kindof.Int(expected)
	if isActualInt && isExpectedInt {
		return actualInt > expectedInt, nil
	}

	// try float
	actualFloat, isActualFloat := kindof.Float(actual)
	expectedFloat, isExpectedFloat := kindof.Float(expected)
	if isActualFloat && isExpectedFloat {
		return actualFloat > expectedFloat, nil
	}

	if !isExpectedFloat && !isExpectedInt {
		return false, fmt.Errorf("type of %v is expected to be number, but %v", expected, reflect.TypeOf(expected).Name())
	}

	return false, fmt.Errorf("type of %v is expected to be number, but %v", actual, reflect.TypeOf(actual).Name())
}
