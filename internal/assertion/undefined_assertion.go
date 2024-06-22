package assertion

import (
	"cotton/internal/value"
	"fmt"
)

type UndefinedAssertion struct {
}

func (a *UndefinedAssertion) Name() string {
	return "is undefined"
}

func (a *UndefinedAssertion) Assert(actual *value.Value) (bool, error) {
	if actual.IsUndefined() {
		return true, nil
	}
	return false, fmt.Errorf("expecting value to be undefined, but got %v", actual)
}
