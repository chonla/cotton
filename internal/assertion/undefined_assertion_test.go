package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/value"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUndefinedAssertionWithUndefinedValue(t *testing.T) {
	actual := value.NewUndefined()

	op := assertion.UndefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestUndefinedAssertionWithNotUndefined(t *testing.T) {
	actual := value.New(4)

	op := assertion.UndefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Equal(t, errors.New("expecting value to be undefined, but got 4"), err)
	assert.False(t, result)
}
