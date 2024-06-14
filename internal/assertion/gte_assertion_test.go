package assertion_test

import (
	"cotton/internal/assertion"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreaterThanOrEqualAssertionWithSameDataTypeAndSuccessGreaterThanCase(t *testing.T) {
	expected := float64(10)
	actual := float64(11)

	op := assertion.GteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestGreaterThanOrEqualAssertionWithSameDataTypeAndSuccessEqualCase(t *testing.T) {
	expected := float64(10)
	actual := float64(10)

	op := assertion.GteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestGreaterThanOrEqualAssertionWithSameDataTypeAndFail(t *testing.T) {
	expected := float64(10)
	actual := float64(9)

	op := assertion.GteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("9 is expected to be greater than or equal to 10, but not"), err)
	assert.False(t, result)
}

func TestGreaterThanOrEqualAssertionWithInvalidActualDataType(t *testing.T) {
	expected := float64(8)
	actual := "10"

	op := assertion.GteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 10 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestGreaterThanOrEqualAssertionWithInvalidExpectedDataType(t *testing.T) {
	expected := "8"
	actual := float64(10)

	op := assertion.GteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestGreaterThanOrEqualAssertionWithInvalidExpectedAndActualDataType(t *testing.T) {
	expected := "8"
	actual := "10"

	op := assertion.GteAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}
