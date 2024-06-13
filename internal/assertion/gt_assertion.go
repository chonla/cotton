package assertion

import (
	"cotton/internal/kindof"
	"errors"
	"fmt"
	"reflect"
)

type GtAssertion struct {
}

func (a *GtAssertion) Name() string {
	return ">"
}

func (a *GtAssertion) Assert(actual, expected interface{}) (bool, error) {
	// actual > expect?
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		return false, fmt.Errorf("type of %v is expected to be %s but %s", actual, reflect.TypeOf(expected).Name(), reflect.TypeOf(actual).Name())
	}

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

	return false, errors.New("unexpected value comparison")
}
