package assertion

import (
	"fmt"
	"reflect"
)

type NeAssertion struct {
}

func (a *NeAssertion) Name() string {
	return "!="
}

func (a *NeAssertion) Assert(expected, actual interface{}) (bool, error) {
	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		return false, fmt.Errorf("type of %v is expected to be %s but %s", actual, reflect.TypeOf(expected).Name(), reflect.TypeOf(actual).Name())
	}
	if reflect.DeepEqual(actual, expected) {
		return false, fmt.Errorf("expecting %v to be not equal to %v, but it is", actual, expected)
	}
	return true, nil
}
