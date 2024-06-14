package assertion_test

import (
	"cotton/internal/assertion"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreaterThanAssertionWithSameDataTypeAndSuccess(t *testing.T) {
	expected := float64(10)
	actual := float64(11)

	op := assertion.GtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestGreaterThanAssertionWithSameDataTypeAndFail(t *testing.T) {
	expected := float64(10)
	actual := float64(10)

	op := assertion.GtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.False(t, result)
}

func TestGreaterThanAssertionWithInvalidActualDataType(t *testing.T) {
	expected := float64(8)
	actual := "10"

	op := assertion.GtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 10 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestGreaterThanAssertionWithInvalidExpectedDataType(t *testing.T) {
	expected := "8"
	actual := float64(10)

	op := assertion.GtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestGreaterThanAssertionWithInvalidExpectedAndActualDataType(t *testing.T) {
	expected := "8"
	actual := "10"

	op := assertion.GtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}
