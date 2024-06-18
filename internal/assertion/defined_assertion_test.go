package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/response"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefinedAssertionWithDefinedValue(t *testing.T) {
	actual := &response.DataValue{
		Value:       4,
		TypeName:    "int",
		IsUndefined: false,
	}

	op := assertion.DefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDefinedAssertionWithDefinedValueOnNullValue(t *testing.T) {
	actual := &response.DataValue{
		Value:       nil,
		TypeName:    "",
		IsUndefined: false,
	}

	op := assertion.DefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestDefinedAssertionWithUndefined(t *testing.T) {
	actual := &response.DataValue{
		Value:       nil,
		TypeName:    "",
		IsUndefined: true,
	}

	op := assertion.DefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Equal(t, errors.New("expecting value to be defined, but not"), err)
	assert.False(t, result)
}
