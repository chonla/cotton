package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/response"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualAssertionWithSameValue(t *testing.T) {
	actual := &response.DataValue{
		Value:       "10",
		IsUndefined: false,
	}
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestEqualAssertionWithDifferentType(t *testing.T) {
	actual := &response.DataValue{
		Value:       3,
		TypeName:    "int",
		IsUndefined: false,
	}
	expected := "3"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 3 is expected to be string but int"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithDifferentValue(t *testing.T) {
	actual := &response.DataValue{
		Value:       "3",
		TypeName:    "string",
		IsUndefined: false,
	}
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting 10 but got 3"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithUndefinedValue(t *testing.T) {
	actual := &response.DataValue{
		Value:       nil,
		TypeName:    "",
		IsUndefined: true,
	}
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("unexpected undefined value"), err)
	assert.False(t, result)
}
