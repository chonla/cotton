package assertion_test

import (
	"cotton/internal/assertion"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualAssertionWithDifferentType(t *testing.T) {
	v1 := 3
	v2 := "3"

	op := assertion.EqualAssertion{}

	result, err := op.Assert(v1, v2)

	assert.Equal(t, errors.New("3 (int) and 3 (string) have different type"), err)
	assert.False(t, result)
}
