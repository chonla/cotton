package assertion_test

import (
	"cotton/internal/assertion"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualAssertionWithSameValue(t *testing.T) {
	v1 := "10"
	v2 := "10"

	op := assertion.EqualAssertion{}

	result, err := op.Assert(v1, v2)

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestEqualAssertionWithDifferentType(t *testing.T) {
	v1 := 3
	v2 := "3"

	op := assertion.EqualAssertion{}

	result, err := op.Assert(v1, v2)

	assert.Equal(t, errors.New("type of 3 is expected to be string but int"), err)
	assert.False(t, result)
}

func TestEqualAssertionWithDifferentValue(t *testing.T) {
	v1 := "3"
	v2 := "10"

	op := assertion.EqualAssertion{}

	result, err := op.Assert(v1, v2)

	assert.Equal(t, errors.New("expecting 10 but got 3"), err)
	assert.False(t, result)
}
