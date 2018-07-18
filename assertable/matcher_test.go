package assertable

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRegShouldReturnTrueIfStringIsSurroundedBySlashes(t *testing.T) {
	result := isRegExp("/pattern/")
	assert.Equal(t, true, result)
}

func TestIsRegShouldReturnFalseIfStringIsOnlyEndedBySlashes(t *testing.T) {
	result := isRegExp("a/pattern/")
	assert.Equal(t, false, result)
}

func TestIsRegShouldReturnFalseIfStringIsOnlyStartedWithSlashes(t *testing.T) {
	result := isRegExp("/pattern/a")
	assert.Equal(t, false, result)
}

func TestIsRegShouldReturnFalseIfStringHasNoSlashesAtBeginningAndTheEnd(t *testing.T) {
	result := isRegExp("a/pattern/a")
	assert.Equal(t, false, result)
}

func TestNewMatcherShouldReturnRegexMatcherIfPatternIsRegularExpression(t *testing.T) {
	m := NewMatcher("/pattern/")
	assert.Equal(t, &Matcher{
		reg:   regexp.MustCompile("pattern"),
		value: "/pattern/",
	}, m)
}

func TestNewMatcherShouldReturnStringMatcherIfPatternIsNotRegularExpression(t *testing.T) {
	m := NewMatcher("a/pattern/b")
	assert.Equal(t, &Matcher{
		reg:   nil,
		value: "a/pattern/b",
	}, m)
}

func TestMatcherWithRegularExpressionShouldSuccess(t *testing.T) {
	m := NewMatcher("/pattern/")
	assert.True(t, m.Match("my pattern in the room"))
}

func TestMatcherWithRegularExpressionShouldFailure(t *testing.T) {
	m := NewMatcher("/pattern/")
	assert.False(t, m.Match("my cotton in the web"))
}

func TestNewMatcherWithValueShouldCompareWithValue(t *testing.T) {
	m := NewMatcher("a/pattern/b")
	assert.True(t, m.Match("a/pattern/b"))
}
