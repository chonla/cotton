package assertable

import (
	"regexp"
	"testing"

	"github.com/chonla/cotton/response"
	"github.com/stretchr/testify/assert"
)

func TestIsRegShouldReturnTrueIfStringIsSurroundedBySlashes(t *testing.T) {
	result := isRegExp("/pattern/")
	assert.Equal(t, true, result)
}

func TestIsRegShouldReturnFalseIfStringIsOnlyEndedBySlash(t *testing.T) {
	result := isRegExp("a/pattern/")
	assert.Equal(t, false, result)
}

func TestIsRegShouldReturnFalseIfStringIsOnlyStartedWithSlash(t *testing.T) {
	result := isRegExp("/pattern/a")
	assert.Equal(t, false, result)
}

func TestIsRegShouldReturnFalseIfStringHasNoSlashAtBeginningAndTheEnd(t *testing.T) {
	result := isRegExp("a/pattern/a")
	assert.Equal(t, false, result)
}

func TestNewMatcherShouldReturnRegexMatcherIfPatternIsRegularExpression(t *testing.T) {
	m := NewMatcher("key", "/pattern/")
	assert.Equal(t, &Matcher{
		reg:     regexp.MustCompile("pattern"),
		value:   "/pattern/",
		key:     "key",
		builtIn: false,
	}, m)
}

func TestNewMatcherShouldReturnStringMatcherIfPatternIsNotRegularExpression(t *testing.T) {
	m := NewMatcher("key", "a/pattern/b")
	assert.Equal(t, &Matcher{
		reg:     nil,
		value:   "a/pattern/b",
		key:     "key",
		builtIn: false,
	}, m)
}

func TestMatcherWithRegularExpressionShouldFailure(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2] }"

	response := &response.Response{
		Proto:      "http",
		Status:     "200 OK",
		StatusCode: 200,
		Header: map[string][]string{
			"content-type": []string{
				"application/json; charset=utf-8",
			},
		},
		Body: jsonString,
	}

	assertable := NewAssertable(response)

	m := NewMatcher("header.content-type", "/pattern/")
	r, _ := m.Match(assertable)
	assert.False(t, r)
}

func TestToStringShouldShowItIsRegexIfPatternIsRegularExpression(t *testing.T) {
	m := NewMatcher("key", "/pattern/")
	result := m.String()
	assert.Equal(t, "Regex(/pattern/)", result)
}

func TestToStringShouldShowItIsRegexIfPatternIsNotRegularExpression(t *testing.T) {
	m := NewMatcher("key", "/pattern/a")
	result := m.String()
	assert.Equal(t, "/pattern/a", result)
}

func TestMatchStringShouldDoExactMatch(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2] }"

	response := &response.Response{
		Proto:      "http",
		Status:     "200 OK",
		StatusCode: 200,
		Header: map[string][]string{
			"content-type": []string{
				"application/json; charset=utf-8",
			},
		},
		Body: jsonString,
	}

	assertable := NewAssertable(response)

	m := NewMatcher("header.content-type", "application/json; charset=utf-8")
	r, _ := m.Match(assertable)
	assert.True(t, r)
}

func TestMatchStringShouldDoRegularExpressionMatch(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2] }"

	response := &response.Response{
		Proto:      "http",
		Status:     "200 OK",
		StatusCode: 200,
		Header: map[string][]string{
			"content-type": []string{
				"application/json; charset=utf-8",
			},
		},
		Body: jsonString,
	}

	assertable := NewAssertable(response)

	m := NewMatcher("header.content-type", "/^application/json/")
	r, _ := m.Match(assertable)
	assert.True(t, r)
}

func TestIsBuiltInShouldReturnTrueIfStringIsSurroundedByStars(t *testing.T) {
	result := isBuiltIn("*keyword*")
	assert.Equal(t, true, result)
}

func TestIsBuiltInShouldReturnFalseIfStringIsOnlyEndedByStar(t *testing.T) {
	result := isBuiltIn("a*keyword*")
	assert.Equal(t, false, result)
}

func TestIsBuiltInShouldReturnFalseIfStringIsOnlyStartedWithStar(t *testing.T) {
	result := isBuiltIn("*keyword*a")
	assert.Equal(t, false, result)
}

func TestIsBuiltInShouldReturnFalseIfStringHasNoStarAtBeginningAndTheEnd(t *testing.T) {
	result := isBuiltIn("a*keyword*a")
	assert.Equal(t, false, result)
}
