package engine

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPExecutor_Timeout(t *testing.T) {
	// Create a slow server that takes 200ms to respond
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok"}`)
	}))
	defer slowServer.Close()

	// Create executor with 50ms timeout
	executor := NewHTTPExecutorWithTimeout(50 * time.Millisecond)

	requestSpec := fmt.Sprintf("GET %s HTTP/1.1\nHost: localhost", slowServer.URL)

	_, err := executor.ExecuteRequest(requestSpec)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Verify it's a timeout error
	if !isTimeoutError(err) {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestHTTPExecutor_MissingVariable(t *testing.T) {
	executor := NewHTTPExecutor()
	executor.SetVariable("known", "value")

	requestSpec := "GET /api/{{ unknown }} HTTP/1.1\nHost: localhost"

	_, err := executor.ExecuteRequest(requestSpec)
	if err == nil {
		t.Fatal("Expected parse error for missing variable, got nil")
	}

	// Verify it's a ParseError
	if _, ok := err.(*ParseError); !ok {
		t.Errorf("Expected ParseError, got %T: %v", err, err)
	}

	expectedMsg := `Parse error: missing variable "unknown"`
	if err.Error() != expectedMsg {
		t.Errorf("Error message = %q, want %q", err.Error(), expectedMsg)
	}
}

func TestHTTPExecutor_SuccessfulSubstitution(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/42" {
			t.Errorf("Expected path /api/42, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"id":42}`)
	}))
	defer server.Close()

	executor := NewHTTPExecutor()
	executor.SetVariable("id", "42")

	requestSpec := fmt.Sprintf("GET %s/api/{{ id }} HTTP/1.1\nHost: localhost", server.URL)

	resp, err := executor.ExecuteRequest(requestSpec)
	if err != nil {
		t.Fatalf("ExecuteRequest() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status = %d, want 200", resp.StatusCode)
	}

	body := executor.GetLastBody()
	if body != `{"id":42}`+"\n" {
		t.Errorf("Body = %q, want %q", body, `{"id":42}`+"\n")
	}
}

func isTimeoutError(err error) bool {
	// Check if error contains timeout-related message
	return err != nil && (err.Error() == "context deadline exceeded" ||
		err.Error() == "Client.Timeout exceeded while awaiting headers" ||
		containsStr(err.Error(), "timeout") ||
		containsStr(err.Error(), "deadline"))
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && containsStrHelper(s, substr)))
}

func containsStrHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
