package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/value"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualAssertionWithSameValue(t *testing.T) {
	actual := value.New("10")
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestEqualAssertionWithDifferentType(t *testing.T) {
	actual := value.New(3)
	expected := "3"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 3 is expected to be string but int"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithDifferentValue(t *testing.T) {
	actual := value.New("3")
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting 10 but got 3"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithUndefinedValue(t *testing.T) {
	actual := value.NewUndefined()
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("unexpected undefined value"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithRegexMatchValue(t *testing.T) {
	actual := value.New("108271X")
	expected, _ := regexp.Compile("^10")

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestEqualAssertionWithRegexUnmatchValue(t *testing.T) {
	actual := value.New("108271X")
	expected, _ := regexp.Compile("^11")

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting value matching /^11/, but got 108271X"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithRegexAgainstNonString(t *testing.T) {
	actual := value.New(10827)
	expected, _ := regexp.Compile("^10")

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 10827 is expected to be string but int"), err)
	assert.False(t, result)
}
