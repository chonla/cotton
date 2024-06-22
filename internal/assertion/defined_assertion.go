package assertion

import (
	"cotton/internal/value"
	"errors"
)

type DefinedAssertion struct {
}

func (a *DefinedAssertion) Name() string {
	return "is defined"
}

func (a *DefinedAssertion) Assert(actual *value.Value) (bool, error) {
	if !actual.IsUndefined() {
		return true, nil
	}
	return false, errors.New("expecting value to be defined, but not")
}
