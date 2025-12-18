package engine

import (
	"testing"
)

func TestHTTPExecutor_SubstituteVariables(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		variables map[string]interface{}
		expected  string
	}{
		{
			name: "Single variable substitution",
			text: "GET /api/products/{{ id }} HTTP/1.1",
			variables: map[string]interface{}{
				"id": "42",
			},
			expected: "GET /api/products/42 HTTP/1.1",
		},
		{
			name: "Multiple variables",
			text: "POST {{ path }} HTTP/1.1\nAuthorization: {{ token }}",
			variables: map[string]interface{}{
				"path":  "/api/users",
				"token": "Bearer abc123",
			},
			expected: "POST /api/users HTTP/1.1\nAuthorization: Bearer abc123",
		},
		{
			name: "No variables to substitute",
			text: "GET /api/products HTTP/1.1",
			variables: map[string]interface{}{
				"id": "42",
			},
			expected: "GET /api/products HTTP/1.1",
		},
		{
			name: "Numeric variable substitution",
			text: "GET /api/products?limit={{ limit }}",
			variables: map[string]interface{}{
				"limit": 100,
			},
			expected: "GET /api/products?limit=100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewHTTPExecutor()
			for k, v := range tt.variables {
				executor.SetVariable(k, v)
			}
			result, err := executor.substituteVariables(tt.text)
			if err != nil {
				t.Fatalf("substituteVariables() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("substituteVariables() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHTTPExecutor_SetAndGetVariables(t *testing.T) {
	executor := NewHTTPExecutor()

	vars := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	executor.SetVariables(vars)

	if len(executor.variables) != len(vars) {
		t.Errorf("SetVariables() stored %d variables, want %d", len(executor.variables), len(vars))
	}

	for k, v := range vars {
		if got, ok := executor.variables[k]; !ok || got != v {
			t.Errorf("Variable %q = %q, want %q", k, got, v)
		}
	}
}

func TestHTTPExecutor_SetVariable(t *testing.T) {
	executor := NewHTTPExecutor()

	executor.SetVariable("name", "John")
	executor.SetVariable("age", 30)

	if val, ok := executor.variables["name"]; !ok || val != "John" {
		t.Errorf("SetVariable() failed for string: got %v", val)
	}
	if val, ok := executor.variables["age"]; !ok || val != 30 {
		t.Errorf("SetVariable() failed for int: got %v", val)
	}
}

func TestHTTPExecutor_ParseRequest(t *testing.T) {
	tests := []struct {
		name        string
		requestText string
		wantMethod  string
		wantURL     string
		wantHeader  string
		wantErr     bool
	}{
		{
			name:        "Simple GET request",
			requestText: "GET http://example.com/api HTTP/1.1",
			wantMethod:  "GET",
			wantURL:     "http://example.com/api",
			wantErr:     false,
		},
		{
			name:        "POST with headers",
			requestText: "POST http://example.com/api HTTP/1.1\nContent-Type: application/json\n\n{\"test\": \"data\"}",
			wantMethod:  "POST",
			wantURL:     "http://example.com/api",
			wantHeader:  "application/json",
			wantErr:     false,
		},
		{
			name:        "Invalid request line",
			requestText: "INVALID",
			wantErr:     true,
		},
		{
			name:        "Empty request",
			requestText: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewHTTPExecutor()
			req, err := executor.parseRequest(tt.requestText)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.Method != tt.wantMethod {
					t.Errorf("Method = %q, want %q", req.Method, tt.wantMethod)
				}
				if req.URL.String() != tt.wantURL {
					t.Errorf("URL = %q, want %q", req.URL.String(), tt.wantURL)
				}
				if tt.wantHeader != "" && req.Header.Get("Content-Type") != tt.wantHeader {
					t.Errorf("Header = %q, want %q", req.Header.Get("Content-Type"), tt.wantHeader)
				}
			}
		})
	}
}

func TestHTTPExecutor_LastBody(t *testing.T) {
	executor := NewHTTPExecutor()

	if executor.GetLastBody() != "" {
		t.Errorf("GetLastBody() initial = %q, want empty", executor.GetLastBody())
	}
}
