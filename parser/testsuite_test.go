package parser

import (
	"testing"

	ts "github.com/chonla/yas/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestParseTestSuiteFileName(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteFileName("login-should-success.md")
	assert.Equal(t, "Login Should Success", result)
}

func TestParseTestSuiteFileNameWithoutExtension(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteFileName("login-should-success")
	assert.Equal(t, "Login Should Success", result)
}

func TestParseTestSuiteFileNameWithNotMdExtension(t *testing.T) {
	p := NewParser()
	result := p.parseTestSuiteFileName("login-should-success.txt")
	assert.Equal(t, "Login Should Success Txt", result)
}

func TestParseTestFile(t *testing.T) {
	readFileFn = func(string) ([]byte, error) {
		return []byte(`# Login Should Return Token Within

	## POST /login
	
	` + "```" + `
	{
		"login": "admin",
		"pwd": "admin"
	}
	` + "```" + `
	
	## Expectation
	
	| Assert | Expected |
	| - | - |
	| HEADER.content-type | application/json |
	| DATA.token | /.+/ |
`), nil
	}
	p := NewParser()
	result, _ := p.parseTestSuiteFile("test")

	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:        "Login Should Return Token Within",
			Method:      "POST",
			Path:        "/login",
			RequestBody: "{\n\"login\": \"admin\",\n\"pwd\": \"admin\"\n}",
			ContentType: "application/json",
			Expectations: []ts.Expectation{
				ts.Expectation{
					Key:   "HEADER.content-type",
					Value: "application/json",
				},
				ts.Expectation{
					Key:   "DATA.token",
					Value: "/.+/",
				},
			},
		},
	}, result)
}

func TestIsTitle(t *testing.T) {
	result, matchResult := isTitle("# Title is here")

	assert.True(t, matchResult)
	assert.NotEmpty(t, result)
}
