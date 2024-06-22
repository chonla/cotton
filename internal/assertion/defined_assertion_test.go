package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/value"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefinedAssertionWithDefinedValue(t *testing.T) {
	actual := value.New(4)

	op := assertion.DefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDefinedAssertionWithDefinedValueOnNullValue(t *testing.T) {
	actual := value.New(nil)

	op := assertion.DefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDefinedAssertionWithUndefined(t *testing.T) {
	actual := value.NewUndefined()

	op := assertion.DefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Equal(t, errors.New("expecting value to be defined, but not"), err)
	assert.False(t, result)
}
