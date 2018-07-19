package parser

import (
	"testing"

	ts "github.com/chonla/cotton/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestParseSimpleAction(t *testing.T) {
	data := `# Test Case Name

## GET /todos
`

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:         "Test Case Name",
			Method:       "GET",
			Path:         "/todos",
			Headers:      map[string]string{},
			Expectations: map[string]string{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
		},
	}, result)
}

func TestParseMultipleSimpleAction(t *testing.T) {
	data := `# Test Case Name

## GET /todos

# Another Test Case Name

## GET /list
`

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:         "Test Case Name",
			Method:       "GET",
			Path:         "/todos",
			Headers:      map[string]string{},
			Expectations: map[string]string{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
		},
		&ts.TestCase{
			Name:         "Another Test Case Name",
			Method:       "GET",
			Path:         "/list",
			Headers:      map[string]string{},
			Expectations: map[string]string{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
		},
	}, result)
}

// func TestParseTestSuiteFileName(t *testing.T) {
// 	p := NewParser()
// 	result := p.parseTestSuiteFileName("login-should-success.md")
// 	assert.Equal(t, "Login Should Success", result)
// }

// func TestParseTestSuiteFileNameWithoutExtension(t *testing.T) {
// 	p := NewParser()
// 	result := p.parseTestSuiteFileName("login-should-success")
// 	assert.Equal(t, "Login Should Success", result)
// }

// func TestParseTestSuiteFileNameWithNotMdExtension(t *testing.T) {
// 	p := NewParser()
// 	result := p.parseTestSuiteFileName("login-should-success.txt")
// 	assert.Equal(t, "Login Should Success Txt", result)
// }

// func TestParseTestFile(t *testing.T) {
// 	readFileFn = func(string) ([]byte, error) {
// 		return []byte(`# Login Should Return Token Within

// 	## POST /login

// 	` + "```" + `
// 	{
// 		"login": "admin",
// 		"pwd": "admin"
// 	}
// 	` + "```" + `

// 	## Expectation

// 	| Assert | Expected |
// 	| - | - |
// 	| HEADER.content-type | application/json |
// 	| DATA.token | /.+/ |
// `), nil
// 	}
// 	p := NewParser()
// 	result, _ := p.parseTestSuiteFile("test")

// 	assert.Equal(t, []*ts.TestCase{
// 		&ts.TestCase{
// 			Name:        "Login Should Return Token Within",
// 			Method:      "POST",
// 			Path:        "/login",
// 			RequestBody: "{\n\"login\": \"admin\",\n\"pwd\": \"admin\"\n}",
// 			ContentType: "application/json",
// 			Headers:     map[string]string{},
// 			Expectations: map[string]string{
// 				"HEADER.content-type": "application/json",
// 				"DATA.token":          "/.+/",
// 			},
// 		},
// 	}, result)
// }

// func TestIsTitle(t *testing.T) {
// 	result, matchResult := isTitle("# Title is here")

// 	assert.True(t, matchResult)
// 	assert.NotEmpty(t, result)
// }
