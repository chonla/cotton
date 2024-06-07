package assertion

import (
	"fmt"
	"reflect"
)

type EqualAssertion struct {
}

func (a *EqualAssertion) Name() string {
	return "eq"
}

func (a *EqualAssertion) Assert(actual, expected interface{}) (bool, error) {
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		return false, fmt.Errorf("%v (%s) and %v (%s) have different type", actual, reflect.TypeOf(actual).Name(), expected, reflect.TypeOf(expected).Name())
	}
	if reflect.ValueOf(actual) != reflect.ValueOf(expected) {
		return false, fmt.Errorf("%v and %v are not equal", actual, expected)
	}
	return true, nil
}
