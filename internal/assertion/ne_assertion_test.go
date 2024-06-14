package assertion_test

import (
	"cotton/internal/assertion"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEqualAssertionWithSameValue(t *testing.T) {
	actual := "10"
	expected := "10"

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting 10 to be not equal to 10, but it is"), err)
	assert.False(t, result)
}

func TestNotEqualAssertionWithDifferentType(t *testing.T) {
	actual := 3
	expected := "3"

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 3 is expected to be string but int"), err)
	assert.False(t, result)
}

func TestNotEqualAssertionWithDifferentValue(t *testing.T) {
	actual := "3"
	expected := "10"

	op := assertion.NeAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}
