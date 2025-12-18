package engine

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestEvaluator_EvaluateExpectations_Success(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		expectations []Expectation
		wantFailures int
	}{
		{
			name:         "String equality",
			responseBody: `{"status":"ok","message":"success"}`,
			expectations: []Expectation{
				{Field: "status", Operator: "==", Value: "\"ok\""},
			},
			wantFailures: 0,
		},
		{
			name:         "Numeric comparison",
			responseBody: `{"count":42,"limit":100}`,
			expectations: []Expectation{
				{Field: "count", Operator: ">", Value: "0"},
				{Field: "count", Operator: "<=", Value: "100"},
				{Field: "limit", Operator: ">=", Value: "50"},
			},
			wantFailures: 0,
		},
		{
			name:         "Null check",
			responseBody: `{"status":null,"data":{"value":123}}`,
			expectations: []Expectation{
				{Field: "data.value", Operator: "==", Value: "123"},
			},
			wantFailures: 0,
		},
		{
			name:         "Regex match",
			responseBody: `{"title":"Executive Summary"}`,
			expectations: []Expectation{
				{Field: "title", Operator: "=~", Value: "/^Executive/"},
			},
			wantFailures: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Body: io.NopCloser(strings.NewReader(tt.responseBody)),
			}
			defer resp.Body.Close()

			evaluator := NewEvaluator()
			failures, err := evaluator.EvaluateExpectations(resp, tt.expectations)

			if err != nil {
				t.Errorf("EvaluateExpectations() error = %v", err)
			}
			if len(failures) != tt.wantFailures {
				t.Errorf("Failures count = %d, want %d. Failures: %v", len(failures), tt.wantFailures, failures)
			}
		})
	}
}

func TestEvaluator_EvaluateExpectations_Failure(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		expectations []Expectation
		wantFailures int
	}{
		{
			name:         "String inequality",
			responseBody: `{"status":"error"}`,
			expectations: []Expectation{
				{Field: "status", Operator: "==", Value: "\"ok\""},
			},
			wantFailures: 1,
		},
		{
			name:         "Numeric comparison fails",
			responseBody: `{"count":5}`,
			expectations: []Expectation{
				{Field: "count", Operator: ">", Value: "10"},
			},
			wantFailures: 1,
		},
		{
			name:         "Multiple failures",
			responseBody: `{"status":"error","count":0}`,
			expectations: []Expectation{
				{Field: "status", Operator: "==", Value: "\"ok\""},
				{Field: "count", Operator: ">", Value: "0"},
			},
			wantFailures: 2,
		},
		{
			name:         "Regex non-match",
			responseBody: `{"title":"Associate"}`,
			expectations: []Expectation{
				{Field: "title", Operator: "=~", Value: "/^Executive/"},
			},
			wantFailures: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Body: io.NopCloser(strings.NewReader(tt.responseBody)),
			}
			defer resp.Body.Close()

			evaluator := NewEvaluator()
			failures, err := evaluator.EvaluateExpectations(resp, tt.expectations)

			if err != nil {
				t.Errorf("EvaluateExpectations() error = %v", err)
			}
			if len(failures) != tt.wantFailures {
				t.Errorf("Failures count = %d, want %d. Failures: %v", len(failures), tt.wantFailures, failures)
			}
		})
	}
}

func TestEvaluator_EqualsValue(t *testing.T) {
	tests := []struct {
		name  string
		json  string
		path  string
		value string
		want  bool
	}{
		{
			name:  "String equality",
			json:  `{"name":"John"}`,
			path:  "name",
			value: "\"John\"",
			want:  true,
		},
		{
			name:  "Number equality",
			json:  `{"age":30}`,
			path:  "age",
			value: "30",
			want:  true,
		},
		{
			name:  "Boolean true",
			json:  `{"active":true}`,
			path:  "active",
			value: "true",
			want:  true,
		},
		{
			name:  "Boolean false",
			json:  `{"active":false}`,
			path:  "active",
			value: "false",
			want:  true,
		},
		{
			name:  "String inequality",
			json:  `{"name":"John"}`,
			path:  "name",
			value: "\"Jane\"",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Body: io.NopCloser(strings.NewReader(tt.json)),
			}
			defer resp.Body.Close()

			evaluator := NewEvaluator()
			failures, err := evaluator.EvaluateExpectations(resp, []Expectation{
				{Field: tt.path, Operator: "==", Value: tt.value},
			})

			if err != nil {
				t.Errorf("EvaluateExpectations() error = %v", err)
			}

			hasFailure := len(failures) > 0
			if hasFailure != !tt.want {
				t.Errorf("Expected %v, got failure=%v", tt.want, hasFailure)
			}
		})
	}
}

func TestEvaluator_UnsupportedOperator(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(`{"value":10}`)),
	}
	defer resp.Body.Close()

	evaluator := NewEvaluator()
	failures, err := evaluator.EvaluateExpectations(resp, []Expectation{
		{Field: "value", Operator: "~~", Value: "10"},
	})

	if err != nil {
		t.Errorf("EvaluateExpectations() error = %v", err)
	}
	if len(failures) == 0 {
		t.Errorf("Expected failure for unsupported operator")
	}
}

func TestEvaluator_InvalidRegex(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(`{"title":"Chief"}`)),
	}
	defer resp.Body.Close()

	evaluator := NewEvaluator()
	failures, err := evaluator.EvaluateExpectations(resp, []Expectation{
		{Field: "title", Operator: "=~", Value: "/[unclosed"},
	})

	if err != nil {
		t.Errorf("EvaluateExpectations() error = %v", err)
	}
	if len(failures) == 0 {
		t.Errorf("Expected failure for invalid regex pattern")
	}
}

func TestEvaluator_InvalidJSON(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(`{invalid json}`)),
	}
	defer resp.Body.Close()

	evaluator := NewEvaluator()
	_, err := evaluator.EvaluateExpectations(resp, []Expectation{
		{Field: "field", Operator: "==", Value: "value"},
	})

	if err == nil {
		t.Errorf("EvaluateExpectations() expected error for invalid JSON")
	}
}
