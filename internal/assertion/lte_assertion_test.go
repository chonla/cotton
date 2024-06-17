package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/response"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessThanOrEqualAssertionWithSameDataTypeAndSuccessLessThanCase(t *testing.T) {
	expected := float64(10)
	actual := &response.DataValue{
		Value:       float64(9),
		TypeName:    "float64",
		IsUndefined: false,
	}

	op := assertion.LteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestLessThanOrEqualAssertionWithSameDataTypeAndSuccessEqualCase(t *testing.T) {
	expected := float64(10)
	actual := &response.DataValue{
		Value:       float64(10),
		TypeName:    "float64",
		IsUndefined: false,
	}

	op := assertion.LteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestLessThanOrEqualAssertionWithSameDataTypeAndFail(t *testing.T) {
	expected := float64(10)
	actual := &response.DataValue{
		Value:       float64(11),
		TypeName:    "float64",
		IsUndefined: false,
	}

	op := assertion.LteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("11 is expected to be less than or equal to 10, but not"), err)
	assert.False(t, result)
}

func TestLessThanOrEqualAssertionWithInvalidActualDataType(t *testing.T) {
	expected := float64(8)
	actual := &response.DataValue{
		Value:       "10",
		TypeName:    "string",
		IsUndefined: false,
	}

	op := assertion.LteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 10 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestLessThanOrEqualAssertionWithInvalidExpectedDataType(t *testing.T) {
	expected := "8"
	actual := &response.DataValue{
		Value:       float64(10),
		TypeName:    "float64",
		IsUndefined: false,
	}

	op := assertion.LteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestLessThanOrEqualAssertionWithInvalidExpectedAndActualDataType(t *testing.T) {
	expected := "8"
	actual := &response.DataValue{
		Value:       "10",
		TypeName:    "string",
		IsUndefined: false,
	}

	op := assertion.LteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}
