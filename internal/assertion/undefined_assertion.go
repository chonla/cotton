package assertion

import (
	"cotton/internal/response"
	"fmt"
)

type UndefinedAssertion struct {
}

func (a *UndefinedAssertion) Name() string {
	return "is undefined"
}

func (a *UndefinedAssertion) Assert(actual *response.DataValue) (bool, error) {
	if actual.IsUndefined {
		return true, nil
	}
	return false, fmt.Errorf("expecting undefined but got %v", actual)
}
