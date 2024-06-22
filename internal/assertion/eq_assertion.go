package assertion

import (
	"cotton/internal/value"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

type EqAssertion struct {
}

func (a *EqAssertion) Name() string {
	return "=="
}

func (a *EqAssertion) Assert(expected interface{}, actual *value.Value) (bool, error) {
	if actual.IsUndefined() {
		return false, errors.New("unexpected undefined value")
	}
	expectedType := reflect.TypeOf(expected)
	if expectedType != nil && expectedType.String() == "*regexp.Regexp" {
		// value matching with regexp must be string
		if actual.Type() == "string" {
			result := expected.(*regexp.Regexp).MatchString(actual.Value().(string))
			if !result {
				return false, fmt.Errorf("expecting value matching /%s/, but got %v", expected.(*regexp.Regexp).String(), actual)
			}
		} else {
			return false, fmt.Errorf("type of %v is expected to be string but %s", actual, actual.Type())
		}
	} else {
		if reflect.TypeOf(actual.Value()) != expectedType {
			return false, fmt.Errorf("type of %v is expected to be %s but %s", actual, expectedType.Name(), actual.Type())
		}
		if !reflect.DeepEqual(actual.Value(), expected) {
			return false, fmt.Errorf("expecting %v but got %v", expected, actual)
		}
	}
	return true, nil
}
