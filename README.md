# Cotton - API Spec Test Runner

A modular test runner that executes API specification documents written in Markdown format. Designed for high cohesion and low coupling.

## Architecture

Cotton follows a clean, modular architecture with clear CLI/engine separation:

- cli/
	- app.go: CLI, flags, discovery, timeout wiring, exit codes
	- formatter.go: Output formatting (stdout, JSON)
- engine/
	- parser.go: Parses markdown specs (precompiled regexes)
	- executor.go: Executes HTTP with validated substitution and timeouts
	- evaluator.go: Evaluates expectations
	- runner.go: Orchestrates execution (fail-fast, cycle detection, timeouts)
	- errors.go: Typed errors (ParseError, RunError)

## Features

- Parse API test specs from Markdown files
- Execute HTTP requests with variable substitution (`{{ name }}`) — missing variables fail fast
- Evaluate expectations using GJSON JSON paths
- Sequential setup and teardown execution with cycle detection
- Shared variable context across specs
- Distinguish test cases (with expectations) vs executables (without expectations)
- Supports triple backtick and triple tilde code fences (optional language: ```http, ~~~json)
- Typed errors (ParseError vs RunError) and consistent exit codes
- Configurable HTTP timeout with sane defaults and caps

## Installation

```bash
go build
```

## Usage

Run specs by file, directory, or glob. Key flags shown below.

```bash
# Run a single spec with 30s timeout (default)
./cotton examples/example.spec.md --timeout 30000

# Run a directory and fail fast on first error
./cotton examples/ -f

# Run multiple globs quietly (exit code only)
./cotton "examples/*.spec.md" -q
```

- `--timeout <ms>`: Valid range 1..60000; defaults to 30000; values are capped at 60000.
- `-f, --fail-fast`: Stop on first failing spec.
- `-q, --quiet`: Suppress output (exit code still reflects success/failure).
- Other flags: `-d` (debug), `-m` (monochrome), `-v` (version), `--json-report`/`--html-report` (optional outputs; HTML reserved).

## Specification Format

### Basic Structure

```markdown
# Test Case Title

Variable declarations and setup links

\`\`\`
HTTP Request
\`\`\`

Expectations and teardown links
```

Code fence: Use either triple backticks or triple tildes. Language annotations (like `http`, `json`) are supported and ignored by the parser:

- Backticks: ```http ... ```
- Tildes: ~~~json ... ~~~

### Variables

Define variables before the request code fence:

- Immediate value: `` `product_id` = `3` ``
- JSON path (from response): `` `token` = `$.data.access_token` ``

Variables can be referenced in requests using `{{ variable_name }}`.
If a placeholder has no value at substitution time, Cotton raises a ParseError and aborts the spec.

### Expectations

List expectations after the request code fence:

- `` `$.data.name` == `"John"` ``
- `` `$.data.count` > `0` ``
- `` `$.data.records` != undefined ``
- `` `$.data.title` =~ `/^Executive/` ``

**Supported operators:**
- Numeric: `<`, `>`, `<=`, `>=`, `==`, `!=`
- String: `==`, `!=`
- Regex: `=~` (match), `!~` (not match)

**Value formats:**
- Strings: `"text"` (quoted)
- Numbers: `42`, `3.14`
- Regex: `/pattern/` (forward-slash delimited; e.g., `/^[A-Z]/`, `/[0-9]{3}/`)
- Special: `null`, `undefined`

### JSON Paths

All JSON paths must start with `$` and follow GJSON syntax:

- `$.field` - access field
- `$0.field` - array index (first element's field)
- `$.nested.field` - nested access
- `$.records.#.name` - all names in array

## Setup and Teardown

- **Setup**: Markdown links that appear before the first code fence are executed sequentially before the main request.
- **Teardown**: Markdown links that appear after the first code fence are executed sequentially after evaluation.
- **Variable Sharing**: Variables declared in setup, main test, and teardown share one context. Later declarations overwrite earlier ones.
- **Cycle Detection**: Recursive link cycles are detected and reported as RunError.

## Example Spec

```markdown
# Get Product

- `base_url` = `"https://api.example.com"`
- `sku` = `$.data.sku`
- [Login](./login.md)

```http
GET {{ base_url }}/products/42 HTTP/1.1
Accept: application/json
```

- `$.status` == `"ok"`
- `$.data.id` == `42`
- `$.data.name` != undefined
- `$.data.title` =~ `/^Product/`
- `$.data.sku` == `{{ sku }}`

[Cleanup](./cleanup.md)
```

## Testing

Run all tests:

```bash
go test -v ./...
```

Tests are deterministic and use local `httptest` servers; no external network access is required.

## Code Style Guidelines

- Consistent error handling using `fmt.Errorf` with context; ParseError vs RunError labeling
- Clear function names with descriptive receiver names
- Comments for exported types and functions
- Modular design with single responsibility per package
- Sequential execution of setup/teardown phases
