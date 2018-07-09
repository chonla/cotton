package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTestSuiteNameWithSimpleName(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("aaa")
	assert.Equal(t, "Aaa", result)
}

func TestParseTestSuiteNameWithMultipleWordWithDash(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("login-should-success")
	assert.Equal(t, "Login Should Success", result)
}

func TestParseTestSuiteNameWithMultipleWordWithUnderScore(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("login_should_success")
	assert.Equal(t, "Login Should Success", result)
}

func TestParseTestSuiteNameWithMultipleWordWithMixedDashAndUnderScore(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("login_should-success")
	assert.Equal(t, "Login Should Success", result)
}

func TestParseTestSuiteNameWithMultipleWordWithMixedCases(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("loginShouldSuccess")
	assert.Equal(t, "Login Should Success", result)
}

func TestParseTestSuiteNameWithMultipleWordWithAllCaps(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("LOGIN")
	assert.Equal(t, "Login", result)
}

func TestParseTestSuiteNameWithMultipleWordWithNumber(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteName("loginShouldSuccess1234")
	assert.Equal(t, "Login Should Success 1234", result)
}

func TestTokenizeTestSuiteNameWithSimpleName(t *testing.T) {
	p := NewParser()
	result := p.tokenizeTestSuiteName("login")
	assert.Equal(t, []string{"login"}, result)
}

func TestTokenizeTestSuiteNameWithMultipleWordWithDash(t *testing.T) {
	p := NewParser()
	result := p.tokenizeTestSuiteName("login-should-success")
	assert.Equal(t, []string{"login", "should", "success"}, result)
}

func TestTokenizeTestSuiteNameWithMultipleWordWithUnderScore(t *testing.T) {
	p := NewParser()
	result := p.tokenizeTestSuiteName("login_should_success")
	assert.Equal(t, []string{"login", "should", "success"}, result)
}

func TestTokenizeTestSuiteNameWithMultipleWordWithMixedCases(t *testing.T) {
	p := NewParser()
	result := p.tokenizeTestSuiteName("loginShouldSuccess")
	assert.Equal(t, []string{"login", "Should", "Success"}, result)
}

func TestTokenizeTestSuiteNameWithNumber(t *testing.T) {
	p := NewParser()
	result := p.tokenizeTestSuiteName("loginShouldSuccess1234")
	assert.Equal(t, []string{"login", "Should", "Success", "1234"}, result)
}
