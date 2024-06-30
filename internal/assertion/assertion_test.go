package assertion_test

import (
	"cotton/internal/assertion"
	"cotton/internal/line"
	"regexp"
	"testing"

	"github.com/samber/mo"
	"github.com/stretchr/testify/assert"
)

func TestParseBinaryAssertionFromUnorderedListAsterisk(t *testing.T) {
	mdLine := line.Line("* `$.form.value`==`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseBinaryAssertionFromUnorderedListMinus(t *testing.T) {
	mdLine := line.Line("- `$.form.value`==`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseBinaryAssertionFromUnorderedListPlus(t *testing.T) {
	mdLine := line.Line("+ `$.form.value`==`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseBinaryAssertionFromOrderedList(t *testing.T) {
	mdLine := line.Line("3. `$.form.value`==`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseStringExpectaion(t *testing.T) {
	mdLine := line.Line("3. `$.form.value`==`\"30\"`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    "30",
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseRegularExpressionExpectaion(t *testing.T) {
	mdLine := line.Line("3. `$.form.value`==/^30$/")

	op, _ := assertion.NewRegexOp("==")
	regVal, _ := regexp.Compile("^30$")
	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    regVal,
		Operator: op,
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseBinaryAssertionWithIntegerAndGreaterAssertionOperator(t *testing.T) {
	mdLine := line.Line("* `$.form.value`>`30`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    float64(30),
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.GtAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseBinaryAssertionWithStringAndGreaterAssertionOperator(t *testing.T) {
	mdLine := line.Line("* `$.form.value`==`\"30\"`")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    "30",
		Operator: mo.NewEither3Arg3[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.EqAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseUnaryAssertionWithUnorderedListAsterisk(t *testing.T) {
	mdLine := line.Line("* `$.form.value` is undefined")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    nil,
		Operator: mo.NewEither3Arg2[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.UndefinedAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseUnaryAssertionWithUnorderedListPlus(t *testing.T) {
	mdLine := line.Line("+ `$.form.value` is undefined")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    nil,
		Operator: mo.NewEither3Arg2[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.UndefinedAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseUnaryAssertionWithUnorderedListMinus(t *testing.T) {
	mdLine := line.Line("- `$.form.value` is undefined")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    nil,
		Operator: mo.NewEither3Arg2[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.UndefinedAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}

func TestParseUnaryAssertionFromOrderedList(t *testing.T) {
	mdLine := line.Line("3. `$.form.value` is undefined")

	expected := &assertion.Assertion{
		Selector: "$.form.value",
		Value:    nil,
		Operator: mo.NewEither3Arg2[assertion.UndefinedOperator, assertion.UnaryAssertionOperator, assertion.BinaryAssertionOperator](&assertion.UndefinedAssertion{}),
	}

	result, ok := assertion.Try(mdLine)

	assert.True(t, ok)
	assert.True(t, expected.SimilarTo(result))
}
