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
	m := NewMatcher("key", "/pattern/", nil)
	assert.Equal(t, &Matcher{
		reg:     regexp.MustCompile("pattern"),
		value:   "/pattern/",
		key:     "key",
		builtIn: false,
	}, m)
}

func TestNewMatcherShouldReturnStringMatcherIfPatternIsNotRegularExpression(t *testing.T) {
	m := NewMatcher("key", "a/pattern/b", nil)
	assert.Equal(t, &Matcher{
		reg:         nil,
		value:       "a/pattern/b",
		actualValue: "a/pattern/b",
		key:         "key",
		builtIn:     false,
	}, m)
}

func TestNewMatcherShouldReturnBuiltInMatcherIfPatternIsBuiltIn(t *testing.T) {
	m := NewMatcher("key", "*should exist*", nil)
	assert.Equal(t, &Matcher{
		reg:     nil,
		value:   "should exist",
		key:     "key",
		builtIn: true,
	}, m)
}

func TestNewMatcherShouldReturnStringMatcherIfPatternIsNotBuiltIn(t *testing.T) {
	m := NewMatcher("key", "should exist", nil)
	assert.Equal(t, &Matcher{
		reg:         nil,
		value:       "should exist",
		actualValue: "should exist",
		key:         "key",
		builtIn:     false,
	}, m)
}

func TestNewMatcherShouldReturnStringMatcherWithActualValueIfPatternHasVariables(t *testing.T) {
	m := NewMatcher("key", "{var1},{var2}", map[string]string{"var1": "value1", "var2": "value2"})
	assert.Equal(t, &Matcher{
		reg:         nil,
		value:       "{var1},{var2}",
		actualValue: "value1,value2",
		key:         "key",
		builtIn:     false,
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

	m := NewMatcher("header.content-type", "/pattern/", nil)
	r, _ := m.Match(assertable)
	assert.False(t, r)
}

func TestToStringShouldShowItIsRegexIfPatternIsRegularExpression(t *testing.T) {
	m := NewMatcher("key", "/pattern/", nil)
	result := m.String()
	assert.Equal(t, "with Regex(/pattern/)", result)
}

func TestToStringShouldShowBuiltinKeywordIfItIsBuiltIn(t *testing.T) {
	m := NewMatcher("key", "*pattern*", nil)
	result := m.String()
	assert.Equal(t, "pattern", result)
}

func TestToStringShouldShowItIsRegexIfPatternIsNotRegularExpression(t *testing.T) {
	m := NewMatcher("key", "/pattern/a", nil)
	result := m.String()
	assert.Equal(t, "with /pattern/a", result)
}

func TestToStringShouldShowActualValueIfPatternHasVariables(t *testing.T) {
	m := NewMatcher("key", "{var1},{var2}", map[string]string{"var1": "value1", "var2": "value2"})
	result := m.String()
	assert.Equal(t, "with value1,value2 as {var1},{var2}", result)
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

	m := NewMatcher("header.content-type", "application/json; charset=utf-8", nil)
	r, _ := m.Match(assertable)
	assert.True(t, r)
}
func TestMatchStringShouldDoExactMatchToActualValue(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"joined\": \"0-1-2\" }"

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
	assertable.Variables["first"] = "0"
	assertable.Variables["second"] = "1"
	assertable.Variables["third"] = "2"

	m := NewMatcher("data.joined", "{first}-{second}-{third}", assertable.Variables)
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

	m := NewMatcher("header.content-type", "/^application/json/", nil)
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

func TestIsBuiltInShouldReturnTrueIfStringIsSurroundedByUndersocres(t *testing.T) {
	result := isBuiltIn("_keyword_")
	assert.Equal(t, true, result)
}

func TestIsBuiltInShouldReturnFalseIfStringIsOnlyEndedByUnderscore(t *testing.T) {
	result := isBuiltIn("a_keyword_")
	assert.Equal(t, false, result)
}

func TestIsBuiltInShouldReturnFalseIfStringIsOnlyStartedWithUnderscore(t *testing.T) {
	result := isBuiltIn("_keyword_a")
	assert.Equal(t, false, result)
}

func TestIsBuiltInShouldReturnFalseIfStringHasNoUnderscoreAtBeginningAndTheEnd(t *testing.T) {
	result := isBuiltIn("a_keyword_a")
	assert.Equal(t, false, result)
}

func TestIsBuiltInShouldReturnFalseIfStringStartedWithStarAndEndedWithUnderscore(t *testing.T) {
	result := isBuiltIn("*keyword_")
	assert.Equal(t, false, result)
}

func TestIsBuiltInShouldReturnFalseIfStringStartedWithUnderscoreAndEndedWithStar(t *testing.T) {
	result := isBuiltIn("_keyword*")
	assert.Equal(t, false, result)
}

func TestBuiltInShouldExist(t *testing.T) {
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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should exist*",
		},
		{
			name:  "underscore",
			input: "_should exist_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.list", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldExistOnNullValue(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"target\": null }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should exist*",
		},
		{
			name:  "underscore",
			input: "_should exist_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.target", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldNotExist(t *testing.T) {
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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should not exist*",
		},
		{
			name:  "underscore",
			input: "_should not exist_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.no-a-list", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldBeNull(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": null }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should be null*",
		},
		{
			name:  "underscore",
			input: "_should be null_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldNotBeNull(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": \"element\" }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should not be null*",
		},
		{
			name:  "underscore",
			input: "_should not be null_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldBeNullOnStringOfNull(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": \"null\" }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should not be null*",
		},
		{
			name:  "underscore",
			input: "_should not be null_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldBeTrue(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": true }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should be true*",
		},
		{
			name:  "underscore",
			input: "_should be true_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldBeFalse(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": false }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should be false*",
		},
		{
			name:  "underscore",
			input: "_should be false_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, _ := m.Match(assertable)
			assert.True(t, r)
		})
	}
}

func TestBuiltInShouldBeTrueOnNonBooleanValue(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": \"true\" }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should be true*",
		},
		{
			name:  "underscore",
			input: "_should be true_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, e := m.Match(assertable)
			assert.False(t, r)
			assert.NotNil(t, e)
		})
	}
}

func TestBuiltInShouldBeFalseOnNonBooleanValue(t *testing.T) {
	jsonString := "{ \"data\": \"ok\", \"list\": [0, 1, 2], \"item\": \"false\" }"

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

	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "star",
			input: "*should be false*",
		},
		{
			name:  "underscore",
			input: "_should be false_",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			m := NewMatcher("data.item", v.input, nil)
			r, e := m.Match(assertable)
			assert.False(t, r)
			assert.NotNil(t, e)
		})
	}
}
