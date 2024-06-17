package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/response"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUndefinedAssertionWithUndefinedValue(t *testing.T) {
	actual := &response.DataValue{
		Value:       nil,
		TypeName:    "",
		IsUndefined: true,
	}

	op := assertion.UndefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestUndefinedAssertionWithNotUndefined(t *testing.T) {
	actual := &response.DataValue{
		Value:       4,
		TypeName:    "int",
		IsUndefined: false,
	}

	op := assertion.UndefinedAssertion{}

	result, err := op.Assert(actual)

	assert.Equal(t, errors.New("expecting undefined but got 4"), err)
	assert.False(t, result)
}
