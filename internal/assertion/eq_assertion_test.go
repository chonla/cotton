package assertion_test

import (
	"cotton/internal/assertion"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualAssertionWithSameValue(t *testing.T) {
	actual := "10"
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestEqualAssertionWithDifferentType(t *testing.T) {
	actual := 3
	expected := "3"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("type of 3 is expected to be string but int"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithDifferentValue(t *testing.T) {
	actual := "3"
	expected := "10"

	op := assertion.EqAssertion{}

	result, err := op.Assert(expected, actual)

	assert.Equal(t, errors.New("expecting 10 but got 3"), err)
	assert.False(t, result)
}
