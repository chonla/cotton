package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/value"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessThanAssertionWithSameDataTypeAndSuccess(t *testing.T) {
	expected := float64(10)
	actual := value.New(float64(9))

	op := assertion.LtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestLessThanAssertionWithSameDataTypeAndFail(t *testing.T) {
	expected := float64(10)
	actual := value.New(float64(10))

	op := assertion.LtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("10 is expected to be less than 10, but not"), err)
	assert.False(t, result)
}

func TestLessThanAssertionWithInvalidActualDataType(t *testing.T) {
	expected := float64(8)
	actual := value.New("10")

	op := assertion.LtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 10 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestLessThanAssertionWithInvalidExpectedDataType(t *testing.T) {
	expected := "8"
	actual := value.New(float64(10))

	op := assertion.LtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}

func TestLessThanAssertionWithInvalidExpectedAndActualDataType(t *testing.T) {
	expected := "8"
	actual := value.New("10")

	op := assertion.LtAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 8 is expected to be number, but string"), err)
	assert.False(t, result)
}
