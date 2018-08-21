package parser

import (
	"testing"

	"github.com/chonla/cotton/assertable"
	ts "github.com/chonla/cotton/testsuite"
	"github.com/stretchr/testify/assert"
)

const backticks = "```"

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
			Expectations: []assertable.Row{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Captured:     map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
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
			Expectations: []assertable.Row{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Captured:     map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
		},
		&ts.TestCase{
			Name:         "Another Test Case Name",
			Method:       "GET",
			Path:         "/list",
			Headers:      map[string]string{},
			Expectations: []assertable.Row{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Captured:     map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
		},
	}, result)
}

func TestParsePostAction(t *testing.T) {
	data := `# Test Case Name

## POST /todos

` + backticks + `
{
	"title": "Text data"
}
` + backticks

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:         "Test Case Name",
			Method:       "POST",
			Path:         "/todos",
			RequestBody:  "{\n\"title\": \"Text data\"\n}",
			Headers:      map[string]string{},
			Expectations: []assertable.Row{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Captured:     map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
		},
	}, result)
}

func TestParseActionWithHeader(t *testing.T) {
	data := `# Test Case Name

## POST /todos

| Header | Value |
| - | - |
| Content-Type | application/json |
| Authorization | Bearer test |

` + backticks + `
{
	"title": "Text data"
}
` + backticks

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:        "Test Case Name",
			Method:      "POST",
			Path:        "/todos",
			RequestBody: "{\n\"title\": \"Text data\"\n}",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test",
			},
			Expectations: []assertable.Row{},
			Captures:     map[string]string{},
			Variables:    map[string]string{},
			Captured:     map[string]string{},
			Setups:       []*ts.Task{},
			Teardowns:    []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
		},
	}, result)
}

func TestParseActionWithExpectations(t *testing.T) {
	data := `# Test Case Name

## POST /todos

| Header | Value |
| - | - |
| Content-Type | application/json |
| Authorization | Bearer test |

` + backticks + `
{
	"title": "Text data"
}
` + backticks + `

## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Content-Type | application/json |
| Data.title | Some text |
`

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:        "Test Case Name",
			Method:      "POST",
			Path:        "/todos",
			RequestBody: "{\n\"title\": \"Text data\"\n}",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test",
			},
			Expectations: []assertable.Row{
				assertable.Row{
					Field:       "StatusCode",
					Expectation: "200",
				},
				assertable.Row{
					Field:       "Header.Content-Type",
					Expectation: "application/json",
				},
				assertable.Row{
					Field:       "Data.title",
					Expectation: "Some text",
				},
			},
			Captures:  map[string]string{},
			Variables: map[string]string{},
			Captured:  map[string]string{},
			Setups:    []*ts.Task{},
			Teardowns: []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
		},
	}, result)
}

func TestParseActionWithCaptures(t *testing.T) {
	data := `# Test Case Name

## POST /todos

| Header | Value |
| - | - |
| Content-Type | application/json |
| Authorization | Bearer test |

` + backticks + `
{
	"title": "Text data"
}
` + backticks + `

## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Content-Type | application/json |
| Data.title | Some text |

## Captures

| Name | Value |
| - | - |
| status-code | StatusCode |
| text-title | Data.title |
`

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:        "Test Case Name",
			Method:      "POST",
			Path:        "/todos",
			RequestBody: "{\n\"title\": \"Text data\"\n}",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test",
			},
			Expectations: []assertable.Row{
				assertable.Row{
					Field:       "StatusCode",
					Expectation: "200",
				},
				assertable.Row{
					Field:       "Header.Content-Type",
					Expectation: "application/json",
				},
				assertable.Row{
					Field:       "Data.title",
					Expectation: "Some text",
				},
			},
			Captures: map[string]string{
				"status-code": "StatusCode",
				"text-title":  "Data.title",
			},
			Captured:  map[string]string{},
			Variables: map[string]string{},
			Setups:    []*ts.Task{},
			Teardowns: []*ts.Task{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
		},
	}, result)
}

func TestParseActionWithFullSections(t *testing.T) {
	data := `# Test Case Name

## POST /todos

| Header | Value |
| - | - |
| Content-Type | application/json |
| Authorization | Bearer test |

` + backticks + `
{
	"title": "Text data"
}
` + backticks + `

## Preconditions

* Do [Login](Login.md)
* And [Create Todo](CreateTodo.md) Item

## Expectations

| Assert | Expected |
| - | - |
| StatusCode | 200 |
| Header.Content-Type | application/json |
| Data.title | Some text |

## Captures

| Name | Value |
| - | - |
| status-code | StatusCode |
| text-title | Data.title |

## Finally

* [Delete Todo](DeleteTodo.md) Item
* And [Logout](Logout.md)
`

	readFileFn = func(filename string) ([]byte, error) {
		data := ""
		switch filename {

		case "/Login.md":
			data = `# Login

## POST /login

` + backticks + `
{
	"login": "admin",
	"password" : "password"
}
` + backticks + `

## Captures

| Name | Value |
| - | - |
| token | Data.token |
`

		case "/CreateTodo.md":
			data = `# Create ToDo

			## POST /todos
			
			` + backticks + `
			{
				"title": "something"
			}
			` + backticks + `
			
			## Captures
			
			| Name | Value |
			| - | - |
			| location | Header.Location |
			`

		case "/DeleteTodo.md":
			data = `# Delete ToDo

			## DELETE /todos/3
			`

		case "/Logout.md":
			data = `# Logout

			## GET /logout
			`
		}
		return []byte(data), nil
	}

	p := NewParser()
	result, e := p.ParseString(data, "")

	assert.Nil(t, e)
	assert.Equal(t, []*ts.TestCase{
		&ts.TestCase{
			Name:        "Test Case Name",
			Method:      "POST",
			Path:        "/todos",
			RequestBody: "{\n\"title\": \"Text data\"\n}",
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test",
			},
			Expectations: []assertable.Row{
				assertable.Row{
					Field:       "StatusCode",
					Expectation: "200",
				},
				assertable.Row{
					Field:       "Header.Content-Type",
					Expectation: "application/json",
				},
				assertable.Row{
					Field:       "Data.title",
					Expectation: "Some text",
				},
			},
			Captures: map[string]string{
				"status-code": "StatusCode",
				"text-title":  "Data.title",
			},
			Captured:  map[string]string{},
			Variables: map[string]string{},
			Config: &ts.Config{
				Insecure: false,
				Detail:   false,
			},
			Setups: []*ts.Task{
				&ts.Task{
					Name:        "Login",
					Method:      "POST",
					Path:        "/login",
					RequestBody: "{\n\"login\": \"admin\",\n\"password\" : \"password\"\n}",
					Headers:     map[string]string{},
					Captures: map[string]string{
						"token": "Data.token",
					},
					Variables: map[string]string{},
					Captured:  map[string]string{},
					Config: &ts.Config{
						Insecure: false,
						Detail:   false,
					},
				},
				&ts.Task{
					Name:        "Create ToDo",
					Method:      "POST",
					Path:        "/todos",
					RequestBody: "{\n\"title\": \"something\"\n}",
					Headers:     map[string]string{},
					Captures: map[string]string{
						"location": "Header.Location",
					},
					Variables: map[string]string{},
					Captured:  map[string]string{},
					Config: &ts.Config{
						Insecure: false,
						Detail:   false,
					},
				},
			},
			Teardowns: []*ts.Task{
				&ts.Task{
					Name:      "Delete ToDo",
					Method:    "DELETE",
					Path:      "/todos/3",
					Headers:   map[string]string{},
					Captures:  map[string]string{},
					Variables: map[string]string{},
					Captured:  map[string]string{},
					Config: &ts.Config{
						Insecure: false,
						Detail:   false,
					},
				},
				&ts.Task{
					Name:      "Logout",
					Method:    "GET",
					Path:      "/logout",
					Headers:   map[string]string{},
					Captures:  map[string]string{},
					Variables: map[string]string{},
					Captured:  map[string]string{},
					Config: &ts.Config{
						Insecure: false,
						Detail:   false,
					},
				},
			},
		},
	}, result)
}
