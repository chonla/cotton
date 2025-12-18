# Project Architecture

## Overview

The Cotton project has been refactored to separate concerns between the CLI and the core engine, following principles of high cohesion and low coupling.

## Package Structure

```
cotton/
├── main.go                    # Entry point - minimal orchestration
├── go.mod / go.sum            # Dependencies
├── cli/                       # CLI Layer
│   ├── app.go                # Application entry point, flag parsing, coordination, timeout wiring
│   └── formatter.go          # Output formatting and presentation
└── engine/                   # Core Engine (business logic)
  ├── parser.go             # Markdown spec parsing (precompiled regexes)
  ├── executor.go           # HTTP request execution (validated substitution + timeout)
  ├── evaluator.go          # Expectation evaluation
  ├── runner.go             # Test orchestration (cycle detection, fail-fast, timeouts)
  └── errors.go             # Typed errors: ParseError, RunError
```

## Design Principles

### 1. **Separation of Concerns**

- **CLI Package**: Responsible for user interaction, argument parsing, and output presentation
- **Engine Package**: Responsible for all test execution logic, completely independent of how it's invoked

### 2. **High Cohesion**

Each package has a clear, single responsibility:

- `cli.App`: Handles command-line orchestration
- `cli.ResultFormatter`: Handles output presentation
- `engine.TestRunner`: Orchestrates test execution
- `engine.MarkdownParser`: Parses markdown specifications
- `engine.HTTPExecutor`: Executes HTTP requests
- `engine.Evaluator`: Evaluates test assertions
- `engine.ParseError` / `engine.RunError`: Distinguish parse-time vs run-time failures

### 3. **Low Coupling**

- The `engine` package has **zero dependencies** on the `cli` package
- `cli` depends on `engine` through well-defined public interfaces
- Engine types (`TestResult`, `ParsedSpec`, `Expectation`) are exported without requiring CLI knowledge
- The engine can be reused by other interfaces (web UI, API, programmatic calls, etc.)

### 4. **Dependency Flow**

```
main.go
  ↓
cli.App (flags, --timeout) ──→ cli.ResultFormatter
  ↓
engine.*
  ├── engine.TestRunner (SetFailFast, SetTimeout)
  ├── engine.MarkdownParser (precompiled regex)
  ├── engine.HTTPExecutor (validated substitution + http.Client Timeout)
  ├── engine.Evaluator
  └── engine.Errors (ParseError / RunError)
```

## File Descriptions

### CLI Layer

#### `cli/app.go`
- **Type**: `App`
- **Responsibility**: Main application coordinator
  - Parses command-line flags (debug/quiet/monochrome/fail-fast/version)
  - Validates and applies `--timeout <ms>` (1..60000 ms)
  - Optional reporters: `--json-report`, `--html-report` (HTML reserved)
  - Discovers spec files from file, directory, or glob
  - Resolves file paths to absolute paths
  - Orchestrates runner invocation
  - Handles exit codes
  - **Key Methods**:
    - `NewApp()`: Creates new CLI application
    - `Run(args []string) int`: Executes CLI with arguments, returns exit code

#### `cli/formatter.go`
- **Type**: `ResultFormatter`
- **Responsibility**: Formats and presents test results
  - Prints test titles, pass/fail status
  - Formats failures and errors for readability
  - Handles console output styling
  - **Key Methods**:
    - `NewResultFormatter()`: Creates formatter
    - `PrintResult(result *engine.TestResult)`: Formats and prints results

### Engine Layer

#### `engine/runner.go`
- **Types**: `TestRunner`, `TestContext`, `TestResult`
- **Responsibility**: Orchestrates the complete test execution pipeline
  - Executes setup specs sequentially
  - Resolves and manages variables across specs
  - Executes the main test request
  - Evaluates expectations
  - Executes teardown specs
  - Manages shared context state
  - Detects and prevents recursive setup/teardown cycles
  - Propagates failures to callers and aggregates summary
  - **Key Methods**:
    - `NewTestRunner()`: Creates runner
    - `SetTimeout(d time.Duration)`: Sets HTTP timeout (capped at 60s)
    - `RunSpec(specFile string)`: Entry point for running a spec
    - `executeSpec()`: Recursive execution of individual specs

#### `engine/parser.go`
- **Types**: `MarkdownParser`, `ParsedSpec`, `Expectation`
- **Responsibility**: Parses markdown spec files
  - Reads markdown files
  - Extracts titles, setup links, variables
  - Extracts HTTP request bodies from code fences
  - Extracts expectations and teardown links
  - Supports triple backtick and triple tilde fences, with optional language annotations (e.g., ```http, ~~~json)
  - **Key Methods**:
    - `NewMarkdownParser(filePath string)`: Creates parser
    - `Parse()`: Parses and returns structured spec
    - Helper methods for extracting different sections

#### `engine/executor.go`
- **Type**: `HTTPExecutor`
- **Responsibility**: Executes HTTP requests and manages execution state
  - Maintains variable dictionary for substitution
  - Substitutes variables in request specs (validated; missing variables cause ParseError)
  - Parses HTTP requests from text format
  - Executes HTTP requests with `http.Client{Timeout: …}`
  - Stores response bodies for variable extraction
  - **Key Methods**:
    - `NewHTTPExecutor()` / `NewHTTPExecutorWithTimeout(d)`: Creates executor
    - `ExecuteRequest(requestSpec string)`: Executes HTTP request
    - `SetVariables(vars map[string]string)`: Sets substitution variables
    - `GetLastBody()`: Retrieves last response body

#### `engine/evaluator.go`
- **Type**: `Evaluator`
- **Responsibility**: Evaluates test expectations against responses
  - Validates JSON response structure
  - Evaluates field values against expected values
  - Supports multiple comparison operators (`==`, `!=`, `<`, `>`, `<=`, `>=`)
  - Handles JSON path extraction via GJSON
  - Enforces JSON path convention from spec: paths start with `$` and comply with GJSON syntax
  - **Key Methods**:
    - `NewEvaluator()`: Creates evaluator
    - `EvaluateExpectations(resp, expectations)`: Evaluates all expectations
    - Comparison methods: `lessThan()`, `greaterThan()`, etc.

### Entry Point

#### `main.go`
- **Responsibility**: Minimal entry point
  - Creates CLI application
  - Passes command-line arguments
  - Exits with appropriate code
  - Keeps dependency resolution minimal

## Extending the Project

### Adding a Web Interface

Since the `engine` package is independent of the CLI:

```go
package webui

import "github.com/chonla/cotton/engine"

func runSpecViaHTTP(w http.ResponseWriter, r *http.Request) {
    runner := engine.NewTestRunner()
    result, err := runner.RunSpec(r.URL.Query().Get("spec"))
    // Return result as JSON
}
```

### Programmatic Usage

```go
import "github.com/chonla/cotton/engine"

runner := engine.NewTestRunner()
result, err := runner.RunSpec("tests/api.spec.md")
if err != nil {
    // Handle error
}
if !result.Passed {
    // Handle test failure
}
```

### Adding Different Output Formats

Create new formatter implementations without touching the engine:

```go
package cli

type JSONFormatter struct{}
func (f *JSONFormatter) PrintResult(result *engine.TestResult) {
    // Marshal to JSON and print
}
```

## Benefits of This Architecture

1. **Testability**: Engine logic can be tested independently of CLI
2. **Reusability**: Engine can be imported and used by other tools
3. **Maintainability**: Clear separation makes code easier to understand and modify
4. **Extensibility**: New interfaces (web, API, etc.) can be added without modifying engine
5. **Scalability**: Components can be evolved independently
6. **Decoupling**: Changes to CLI output don't affect engine logic

## Import Guidelines

- ✅ `cli` package imports from `engine`
- ❌ `engine` package should NOT import from `cli`
- ✅ Other packages can import `engine` for programmatic usage
- ✅ `main` package only imports from `cli`

## Spec Conventions

The engine parses and executes Markdown specs with the following rules:

- **Title**: The first `#` heading is the test case title.
- **Request**: The first code fence contains the HTTP request. Supported fences: triple backticks and triple tildes, optionally with a language (e.g., ```http or ~~~json).
- **Setup Links**: Markdown links that appear before the request fence are executed sequentially as setup steps.
- **Teardown Links**: Markdown links that appear after the request fence are executed sequentially as teardown steps.
- **Variables**: Declared as bullets before the request fence using `` `name` = `value` ``. `value` may be an immediate value or a JSON path (prefixed with `$`) resolved from the latest response.
- **Variable Substitution**: Placeholders in requests must be in the form `{{ name }}`. Missing variables produce a ParseError and abort execution.
- **Expectations**: Listed as bullets after the request fence using `` `$.path` <op> `value` `` where `<op>` is one of `==`, `!=`, `<`, `>`, `<=`, `>=`, `=~`, `!~`.
  - Operators: Numeric (`<`, `>`, `<=`, `>=`, `==`, `!=`), String (`==`, `!=`), Regex (`=~` match, `!~` not match).
  - Value formats: Strings (`"text"`), Numbers (`42`), Regex (`/pattern/`), Special (`null`, `undefined`).
- **Executables vs Tests**: Specs without expectations are executables (e.g., setup scripts) and run only when linked from a test.
- **Variable Sharing**: Variables declared in setups, the test case, and teardowns share a common context. Later declarations overwrite earlier ones.
- **JSON Paths**: Must start with `$` and follow GJSON syntax.

## Error Model & Exit Codes

- Parse-time errors (e.g., invalid `--timeout`, missing template variables, spec read/parse failures) are surfaced as `ParseError` and cause a non-zero exit.
- Run-time errors (e.g., HTTP failures, evaluation errors, cycle detection) are surfaced as `RunError` and cause a non-zero exit.
- When any test fails, the CLI exits with code 1 after printing results; quiet mode (`-q`) prints nothing but still returns the appropriate code.

## Performance Notes

- Parser uses precompiled regexes for links, variables, and expectations to avoid recompilation per line.
- HTTP requests enforce a default timeout (30s) and maximum cap (60s) to avoid hangs.

## Execution Flow

1. **Parse**: `MarkdownParser.Parse()` reads the spec and produces `ParsedSpec` (title, setup/teardown links, variables, request, expectations).
2. **Setup**: `TestRunner` resolves and runs setup specs in order, sharing variables.
3. **Request**: `HTTPExecutor` substitutes variables, parses, and executes the request from the main spec.
4. **Evaluate**: `Evaluator` validates the response JSON and evaluates expectations via GJSON.
5. **Teardown**: `TestRunner` runs teardown specs in order, using the shared variable context.
