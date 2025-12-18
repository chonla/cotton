package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chonla/cotton/engine"
)

// App represents the CLI application
type App struct {
	runner    *engine.TestRunner
	formatter *ResultFormatter
}

// NewApp creates a new CLI application
func NewApp() *App {
	return &App{
		runner:    engine.NewTestRunner(),
		formatter: NewResultFormatter(),
	}
}

// Run executes the CLI application
func (a *App) Run(args []string) int {
	fs := flag.NewFlagSet("cotton", flag.ContinueOnError)
	debug := fs.Bool("d", false, "Enable debug logging")
	quiet := fs.Bool("q", false, "Silence execution (exit code only)")
	monochrome := fs.Bool("m", false, "Monochrome logging")
	failFast := fs.Bool("f", false, "Fail fast (stop on first failure)")
	showVersion := fs.Bool("v", false, "Show version and exit")
	htmlReport := fs.String("html-report", "", "Generate HTML report to file")
	jsonReport := fs.String("json-report", "", "Generate JSON report to file")
	timeout := fs.Int("timeout", 30000, "HTTP timeout in milliseconds (max 60000)")

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Version flag has highest priority and bypasses other options
	if *showVersion {
		fmt.Printf("cotton %s\n", Version)
		return 0
	}

	// Validate timeout
	if *timeout <= 0 || *timeout > 60000 {
		if !*quiet {
			log.Printf("Parse error: timeout must be between 1 and 60000 milliseconds")
		}
		return 1
	}

	// Get positional argument (file, path, or pattern)
	if fs.NArg() == 0 {
		fmt.Println("Usage: cotton [options] <filename|path|file-spec>")
		fmt.Println("Options:")
		fmt.Println("  -d              Enable debug logging")
		fmt.Println("  -q              Silence execution (exit code only)")
		fmt.Println("  -m              Monochrome logging")
		fmt.Println("  -f              Fail fast (stop on first failure)")
		fmt.Println("  -v              Show version and exit")
		fmt.Println("  --timeout       HTTP timeout in milliseconds (default 30000, max 60000)")
		fmt.Println("  --html-report   Generate HTML report to file")
		fmt.Println("  --json-report   Generate JSON report to file")
		return 1
	}

	// Configure formatter and runner
	a.formatter.SetDebug(*debug)
	a.formatter.SetQuiet(*quiet)
	a.formatter.SetMonochrome(*monochrome)
	a.runner.SetFailFast(*failFast)
	a.runner.SetTimeout(time.Duration(*timeout) * time.Millisecond)

	// Discover spec files from all arguments
	var specFiles []string
	for i := 0; i < fs.NArg(); i++ {
		target := fs.Arg(i)
		files, err := a.discoverSpecFiles(target)
		if err != nil {
			if !*quiet {
				log.Printf("Warning: Failed to discover spec files from '%s': %v", target, err)
			}
			continue
		}
		specFiles = append(specFiles, files...)
	}

	if len(specFiles) == 0 {
		if !*quiet {
			log.Fatalf("No spec files found")
		}
		return 1
	}

	// Run all specs and collect results
	results, err := a.runner.RunSpecs(specFiles)
	if err != nil {
		if !*quiet {
			log.Fatalf("Failed to run specs: %v", err)
		}
		return 1
	}

	// Print results
	a.formatter.PrintResults(results)

	// Generate reports if requested
	if *jsonReport != "" {
		if err := a.generateJSONReport(*jsonReport, results); err != nil {
			if !*quiet {
				log.Printf("Failed to generate JSON report: %v", err)
			}
		}
	}

	if *htmlReport != "" {
		if !*quiet {
			log.Printf("HTML report generation not yet implemented")
		}
	}

	if results.Summary.Failed > 0 {
		return 1
	}

	return 0
}

// discoverSpecFiles discovers spec files from a file, directory, or pattern
func (a *App) discoverSpecFiles(target string) ([]string, error) {
	absPath, err := filepath.Abs(target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if it's a file
	info, err := os.Stat(absPath)
	if err == nil {
		if !info.IsDir() {
			// Single file
			return []string{absPath}, nil
		}
		// Directory - find all .spec.md files recursively
		return a.findSpecFilesInDir(absPath)
	}

	// Try as glob pattern
	matches, err := filepath.Glob(target)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	if len(matches) > 0 {
		// Convert to absolute paths and filter files only
		var files []string
		for _, match := range matches {
			absMatch, err := filepath.Abs(match)
			if err != nil {
				continue
			}
			info, err := os.Stat(absMatch)
			if err == nil && !info.IsDir() {
				files = append(files, absMatch)
			}
		}
		return files, nil
	}

	return nil, fmt.Errorf("no files found matching: %s", target)
}

// findSpecFilesInDir recursively finds all .spec.md files in a directory
func (a *App) findSpecFilesInDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".spec.md") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// generateJSONReport generates a JSON report file
func (a *App) generateJSONReport(filename string, results *engine.TestResults) error {
	return results.WriteJSONReport(filename)
}
