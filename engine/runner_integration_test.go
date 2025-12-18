package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunner_CycleDetection(t *testing.T) {
	// Create temporary spec files that reference each other
	tmpDir := t.TempDir()

	specA := filepath.Join(tmpDir, "a.spec.md")
	specB := filepath.Join(tmpDir, "b.spec.md")

	// a.spec.md links to b.spec.md as setup
	contentA := fmt.Sprintf(`# Test A

Setup: [B](b.spec.md)

` + "```" + `
GET http://example.com HTTP/1.1
` + "```" + `

- ` + "`$.status`" + ` == ` + "`\"ok\"`" + `
`)

	// b.spec.md links back to a.spec.md as setup (creating a cycle)
	contentB := fmt.Sprintf(`# Test B

Setup: [A](a.spec.md)

` + "```" + `
GET http://example.com HTTP/1.1
` + "```" + `

- ` + "`$.status`" + ` == ` + "`\"ok\"`" + `
`)

	if err := os.WriteFile(specA, []byte(contentA), 0644); err != nil {
		t.Fatalf("Failed to write spec A: %v", err)
	}
	if err := os.WriteFile(specB, []byte(contentB), 0644); err != nil {
		t.Fatalf("Failed to write spec B: %v", err)
	}

	runner := NewTestRunner()
	result, err := runner.RunSpec(specA)

	if err == nil {
		t.Fatal("Expected cycle detection error, got nil")
	}

	// Verify it's a RunError about cycles
	if _, ok := err.(*RunError); !ok {
		t.Errorf("Expected RunError for cycle, got %T: %v", err, err)
	}

	if !containsStrHelper(err.Error(), "cycle detected") {
		t.Errorf("Error should mention cycle, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result even on error, got nil")
	}

	if result.Passed {
		t.Error("Result should not be passed when cycle is detected")
	}
}

func TestRunner_NoCycleWithLinearSetup(t *testing.T) {
	tmpDir := t.TempDir()

	specSetup := filepath.Join(tmpDir, "setup.spec.md")
	specMain := filepath.Join(tmpDir, "main.spec.md")

	// Local HTTP server to avoid external dependency
	var baseURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/post":
			var payload map[string]any
			_ = json.NewDecoder(r.Body).Decode(&payload)
			_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "payload": payload})
		case "/get":
			_ = json.NewEncoder(w).Encode(map[string]any{"url": baseURL + "/get"})
		default:
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": "not found"})
		}
	}))
	defer srv.Close()
	baseURL = srv.URL

	// setup.spec.md posts to local server
	contentSetup := fmt.Sprintf(`# Setup

`+"```"+`
POST %s/post HTTP/1.1
Content-Type: application/json

{"setup":"done"}
`+"```"+`
`, baseURL)

	// main.spec.md depends on setup and calls local server
	contentMain := fmt.Sprintf(`# Main Test

Setup: [Setup](setup.spec.md)

`+"```"+`
GET %s/get HTTP/1.1
`+"```"+`

- `+"`$.url`"+` == `+"`\"%s/get\"`"+`
`, baseURL, baseURL)

	if err := os.WriteFile(specSetup, []byte(contentSetup), 0644); err != nil {
		t.Fatalf("Failed to write setup spec: %v", err)
	}
	if err := os.WriteFile(specMain, []byte(contentMain), 0644); err != nil {
		t.Fatalf("Failed to write main spec: %v", err)
	}

	runner := NewTestRunner()
	runner.SetTimeout(5000 * time.Millisecond)

	// This should complete without cycle errors (may fail due to network/expectations)
	result, err := runner.RunSpec(specMain)

	// We don't care if the test passes or fails, just that there's no cycle error
	if err != nil {
		if _, ok := err.(*RunError); ok && containsStrHelper(err.Error(), "cycle") {
			t.Errorf("Should not detect cycle in linear setup chain: %v", err)
		}
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if len(result.SetupTitles) == 0 {
		t.Error("Expected setup to be executed")
	}
}
