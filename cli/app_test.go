package cli

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// captureOutput executes fn while capturing stdout output.
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
	}()

	fn()
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestApp_VersionFlagHasPriority(t *testing.T) {
	app := NewApp()

	output := captureOutput(func() {
		code := app.Run([]string{"-v"})
		if code != 0 {
			t.Fatalf("expected exit code 0 for version flag, got %d", code)
		}
	})

	expected := "cotton " + Version + "\n"
	if output != expected {
		t.Fatalf("version output mismatch, got %q want %q", output, expected)
	}
}

func TestApp_VersionFlagIgnoresMissingArgs(t *testing.T) {
	app := NewApp()

	output := captureOutput(func() {
		code := app.Run([]string{"-v", "-d"})
		if code != 0 {
			t.Fatalf("expected exit code 0 when -v provided, got %d", code)
		}
	})

	expected := "cotton " + Version + "\n"
	if output != expected {
		t.Fatalf("version output mismatch with extra flags, got %q want %q", output, expected)
	}
}
