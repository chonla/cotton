package assertion

import (
	"cotton/internal/response"
	"errors"
	"fmt"
	"reflect"
)

type EqAssertion struct {
}

func (a *EqAssertion) Name() string {
	return "=="
}

func (a *EqAssertion) Assert(expected interface{}, actual *response.DataValue) (bool, error) {
	if actual.IsUndefined {
		return false, errors.New("unexpected undefined value")
	}
	if reflect.TypeOf(actual.Value) != reflect.TypeOf(expected) {
		return false, fmt.Errorf("type of %v is expected to be %s but %s", actual, reflect.TypeOf(expected).Name(), actual.TypeName)
	}
	if !reflect.DeepEqual(actual.Value, expected) {
		return false, fmt.Errorf("expecting %v but got %v", expected, actual)
	}
	return true, nil
}
