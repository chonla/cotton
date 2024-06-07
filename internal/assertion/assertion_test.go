package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/line"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEqualAssertionWithInteger(t *testing.T) {
	mdLine := line.Line("* `$.form.value`==`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    30,
		Operator: &assertion.EqualAssertion{},
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseEqualAssertionWithString(t *testing.T) {
	mdLine := line.Line("* `$.form.value`==`\"30\"`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    "30",
		Operator: &assertion.EqualAssertion{},
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}
