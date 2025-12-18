package engine

import (
	"strings"
	"testing"
)

func TestMarkdownParser_ExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Simple title",
			content:  "# Test Title\n\nSome content",
			expected: "Test Title",
		},
		{
			name:     "Title with special characters",
			content:  "# API Test: GET /products\n\nContent",
			expected: "API Test: GET /products",
		},
		{
			name:     "No title",
			content:  "Some content without heading",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MarkdownParser{
				content: tt.content,
				lines:   strings.Split(tt.content, "\n"),
			}
			spec, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if spec.Title != tt.expected {
				t.Errorf("Title = %q, want %q", spec.Title, tt.expected)
			}
		})
	}
}

func TestMarkdownParser_ExtractVariables(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name:    "Single variable",
			content: "`var_name` = `value`\n\n```\nGET http://example.com HTTP/1.1\n```",
			expected: map[string]string{
				"var_name": "value",
			},
		},
		{
			name:    "Multiple variables",
			content: "`id` = `42`\n`name` = `\"John\"`\n\n```\nGET http://example.com HTTP/1.1\n```",
			expected: map[string]string{
				"id":   "42",
				"name": "\"John\"",
			},
		},
		{
			name:    "JSON path variable",
			content: "`token` = `$.auth.token`\n\n```\nGET http://example.com HTTP/1.1\n```",
			expected: map[string]string{
				"token": "$.auth.token",
			},
		},
		{
			name:     "No variables",
			content:  "```\nGET http://example.com HTTP/1.1\n```",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MarkdownParser{
				content: tt.content,
				lines:   strings.Split(tt.content, "\n"),
			}
			spec, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(spec.Variables) != len(tt.expected) {
				t.Errorf("Variables count = %d, want %d", len(spec.Variables), len(tt.expected))
			}
			for k, v := range tt.expected {
				if got, ok := spec.Variables[k]; !ok || got != v {
					t.Errorf("Variable %q = %q, want %q", k, got, v)
				}
			}
		})
	}
}

func TestMarkdownParser_ExtractRequestBody(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Simple GET request",
			content:  "# Test\n\n```\nGET http://example.com HTTP/1.1\n```",
			expected: "GET http://example.com HTTP/1.1",
		},
		{
			name:     "POST with headers and body",
			content:  "# Test\n\n```\nPOST http://example.com HTTP/1.1\nContent-Type: application/json\n\n{\"key\": \"value\"}\n```",
			expected: "POST http://example.com HTTP/1.1\nContent-Type: application/json\n\n{\"key\": \"value\"}",
		},
		{
			name:     "Tilde code fence",
			content:  "# Test\n\n~~~\nGET http://example.com HTTP/1.1\n~~~",
			expected: "GET http://example.com HTTP/1.1",
		},
		{
			name:     "Backtick fence with language",
			content:  "# Test\n\n```http\nGET http://example.com HTTP/1.1\n```",
			expected: "GET http://example.com HTTP/1.1",
		},
		{
			name:     "Tilde fence with language",
			content:  "# Test\n\n~~~http\nPOST http://example.com HTTP/1.1\nContent-Type: text/plain\n\nbody\n~~~",
			expected: "POST http://example.com HTTP/1.1\nContent-Type: text/plain\n\nbody",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MarkdownParser{
				content: tt.content,
				lines:   strings.Split(tt.content, "\n"),
			}
			spec, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if strings.TrimSpace(spec.RequestBody) != strings.TrimSpace(tt.expected) {
				t.Errorf("RequestBody = %q, want %q", spec.RequestBody, tt.expected)
			}
		})
	}
}

func TestMarkdownParser_ExtractExpectations(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Expectation
	}{
		{
			name:    "Single expectation",
			content: "```\nGET http://example.com HTTP/1.1\n```\n\n- `$.status` == `\"ok\"`",
			expected: []Expectation{
				{Field: "$.status", Operator: "==", Value: "\"ok\""},
			},
		},
		{
			name:    "Multiple expectations",
			content: "```\nGET http://example.com HTTP/1.1\n```\n\n- `$.id` == `42`\n- `$.name` != undefined",
			expected: []Expectation{
				{Field: "$.id", Operator: "==", Value: "42"},
				{Field: "$.name", Operator: "!=", Value: "undefined"},
			},
		},
		{
			name:    "Comparison operators",
			content: "```\nGET http://example.com HTTP/1.1\n```\n\n- `$.count` > `0`\n- `$.score` <= `100`",
			expected: []Expectation{
				{Field: "$.count", Operator: ">", Value: "0"},
				{Field: "$.score", Operator: "<=", Value: "100"},
			},
		},
		{
			name:    "Regex operators",
			content: "```\nGET http://example.com HTTP/1.1\n```\n\n- `$.title` =~ `/^Exec/`\n- `$.department` !~ `/^Sales/`",
			expected: []Expectation{
				{Field: "$.title", Operator: "=~", Value: "/^Exec/"},
				{Field: "$.department", Operator: "!~", Value: "/^Sales/"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MarkdownParser{
				content: tt.content,
				lines:   strings.Split(tt.content, "\n"),
			}
			spec, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(spec.Expectations) != len(tt.expected) {
				t.Errorf("Expectations count = %d, want %d", len(spec.Expectations), len(tt.expected))
			}
			for i, exp := range tt.expected {
				if i >= len(spec.Expectations) {
					break
				}
				if spec.Expectations[i].Field != exp.Field || spec.Expectations[i].Operator != exp.Operator || spec.Expectations[i].Value != exp.Value {
					t.Errorf("Expectation %d = %+v, want %+v", i, spec.Expectations[i], exp)
				}
			}
		})
	}
}

func TestMarkdownParser_IsTestCase(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "Test case with expectations",
			content:  "# Test\n\n```\nGET http://example.com HTTP/1.1\n```\n\n- `$.status` == `\"ok\"`",
			expected: true,
		},
		{
			name:     "Executable without expectations",
			content:  "# Setup\n\n```\nPOST http://example.com HTTP/1.1\n```",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &MarkdownParser{
				content: tt.content,
				lines:   strings.Split(tt.content, "\n"),
			}
			spec, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if spec.IsTestCase != tt.expected {
				t.Errorf("IsTestCase = %v, want %v", spec.IsTestCase, tt.expected)
			}
		})
	}
}
