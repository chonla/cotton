package assertion

import (
	"cotton/internal/response"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

type NeAssertion struct {
}

func (a *NeAssertion) Name() string {
	return "!="
}

func (a *NeAssertion) Assert(expected interface{}, actual *response.DataValue) (bool, error) {
	if actual.IsUndefined {
		return false, errors.New("unexpected undefined value")
	}
	expectedType := reflect.TypeOf(expected)
	if expectedType != nil && expectedType.String() == "*regexp.Regexp" {
		// value matching with regexp must be string
		if actual.TypeName == "string" {
			result := expected.(*regexp.Regexp).MatchString(actual.Value.(string))
			if result {
				return false, fmt.Errorf("expecting value not matching /%s/, but got %v", expected.(*regexp.Regexp).String(), actual)
			}
		} else {
			return false, fmt.Errorf("type of %v is expected to be string but %s", actual, actual.TypeName)
		}
	} else {
		if reflect.TypeOf(actual.Value) != expectedType {
			return false, fmt.Errorf("type of %v is expected to be %s but %s", actual, expectedType.Name(), actual.TypeName)
		}
		if reflect.DeepEqual(actual.Value, expected) {
			return false, fmt.Errorf("expecting %v to be not equal to %v, but it is", actual, expected)
		}
	}
	return true, nil
}
