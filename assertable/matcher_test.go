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

func TestMatcherWithRegularExpressionShouldFailure(t *testing.T) {
	m := NewMatcher("/pattern/")
	assert.False(t, m.Match("my cotton in the web"))
}

func TestToStringShouldShowItIsRegexIfPatternIsRegularExpression(t *testing.T) {
	m := NewMatcher("/pattern/")
	result := m.String()
	assert.Equal(t, "Regex(/pattern/)", result)
}

func TestToStringShouldShowItIsRegexIfPatternIsNotRegularExpression(t *testing.T) {
	m := NewMatcher("/pattern/a")
	result := m.String()
	assert.Equal(t, "/pattern/a", result)
}

func TestMatchStringShouldDoExactMatch(t *testing.T) {
	m := NewMatcher("/pattern/a")
	result := m.Match("/pattern/a")
	assert.True(t, result)
}

func TestMatchStringShouldDoRegularExpressionMatch(t *testing.T) {
	m := NewMatcher("/^pattern/")
	result := m.Match("pattern is regular expression")
	assert.True(t, result)
}
