package assertion

import (
	"fmt"
	"reflect"
)

type EqualAssertion struct {
}

func (a *EqualAssertion) Name() string {
	return "=="
}

func (a *EqualAssertion) Assert(actual, expected interface{}) (bool, error) {
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		return false, fmt.Errorf("type of %v is expected to be %s but %s", actual, reflect.TypeOf(expected).Name(), reflect.TypeOf(actual).Name())
	}
	if !reflect.DeepEqual(actual, expected) {
		return false, fmt.Errorf("expecting %v but got %v", expected, actual)
	}
	return true, nil
}
