package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/value"
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEqualAssertionWithSameValue(t *testing.T) {
	actual := value.New("10")
	expected := "10"

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting 10 to be not equal to 10, but it is"), err)
	assert.False(t, result)
}

func TestNotEqualAssertionWithDifferentType(t *testing.T) {
	actual := value.New(3)
	expected := "3"

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 3 is expected to be string but int"), err)
	assert.False(t, result)
}

func TestNotEqualAssertionWithDifferentValue(t *testing.T) {
	actual := value.New("3")
	expected := "10"

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestNeAssertionWithRegexUnmatchValue(t *testing.T) {
	actual := value.New("108271X")
	expected, _ := regexp.Compile("^11")

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestNeAssertionWithRegexMatchValue(t *testing.T) {
	actual := value.New("108271X")
	expected, _ := regexp.Compile("^10")

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting value not matching /^10/, but got 108271X"), err)
	assert.False(t, result)
}

func TestNeAssertionWithRegexAgainstNonString(t *testing.T) {
	actual := value.New(10827)
	expected, _ := regexp.Compile("^10")

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 10827 is expected to be string but int"), err)
	assert.False(t, result)
}
