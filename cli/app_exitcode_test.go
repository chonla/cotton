package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeTempSpec(t *testing.T, dir, name, body, expectation string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	content := fmt.Sprintf(`# Test

`+"```"+`
%s
`+"```"+`

- `+"`$.status`"+` == `+"`%s`"+`
`, body, expectation)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed writing spec: %v", err)
	}
	return path
}

func TestApp_NonZeroExitOnFailure(t *testing.T) {
	// Server returns status=bad to violate expectation
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "bad"})
	}))
	defer srv.Close()

	tmp := t.TempDir()
	req := fmt.Sprintf("GET %s/status HTTP/1.1", srv.URL)
	spec := writeTempSpec(t, tmp, "fail.spec.md", req, "\"ok\"")

	app := NewApp()
	code := app.Run([]string{"-q", spec})
	if code == 0 {
		t.Fatalf("expected non-zero exit for failing spec, got %d", code)
	}
}

func TestApp_ZeroExitOnSuccess(t *testing.T) {
	// Server returns status=ok to satisfy expectation
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	}))
	defer srv.Close()

	tmp := t.TempDir()
	req := fmt.Sprintf("GET %s/status HTTP/1.1", srv.URL)
	spec := writeTempSpec(t, tmp, "pass.spec.md", req, "\"ok\"")

	app := NewApp()
	code := app.Run([]string{"-q", spec})
	if code != 0 {
		t.Fatalf("expected zero exit for passing spec, got %d", code)
	}
}
