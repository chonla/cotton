# Refactoring Summary: CLI-Engine Separation

## What Changed

The Cotton project has been successfully refactored to separate CLI concerns from the core engine, resulting in a clean, modular architecture.

### Before: Monolithic Structure
```
cotton/
├── main.go (CLI + startup logic)
├── runner.go (Test orchestration)
├── parser.go (Markdown parsing)
├── executor.go (HTTP execution)
└── evaluator.go (Expectation evaluation)
```

All files were in the `main` package with mixed responsibilities.

### After: Modular Structure
```
cotton/
├── main.go (Minimal entry point)
├── cli/
│   ├── app.go (CLI application & flag parsing)
│   └── formatter.go (Output presentation)
└── engine/
    ├── runner.go (Test orchestration)
    ├── parser.go (Markdown parsing)
    ├── executor.go (HTTP execution)
  ├── evaluator.go (Expectation evaluation)
  └── errors.go   (Typed errors: ParseError, RunError)
```

Two packages with clear responsibilities and well-defined boundaries.

## Key Improvements

### 1. **Separation of Concerns**
- **CLI Package**: User interaction, input parsing, output formatting
- **Engine Package**: Business logic, completely independent of CLI

### 2. **Modularity**
- Engine can be used as a library in other tools (web UI, API servers, IDE plugins)
- CLI can be replaced without touching engine logic
- New interfaces can be added without modifying existing code

### 3. **Cohesion & Coupling**
- **High Cohesion**: Each file/package has a single, clear purpose
- **Low Coupling**: Engine has zero dependencies on CLI; only CLI depends on engine

### 4. **Testability**
- Engine logic can be tested without requiring CLI setup
- Output formatting can be tested independently
- Integration tests can work with the engine directly

## Package Details

### CLI Package (`cli/`)
- **app.go**: Main CLI coordinator
  - Flag parsing (-spec)
  - Error handling
  - Exit code management
  - Orchestrates runner invocation

- **formatter.go**: Output presenter
  - Formats test results
  - Handles console output
  - Can be extended for different output formats (JSON, XML, etc.)

### Engine Package (`engine/`)
- **runner.go**: Test orchestrator
  - Executes setup specs
  - Manages variable context
  - Executes main request
  - Evaluates expectations
  - Executes teardown specs
  - Detects cycles in setup/teardown chains
  - Supports configurable HTTP timeout via `SetTimeout`

- **parser.go**: Markdown parser
  - Reads spec files
  - Extracts sections (setup, request, expectations, teardown)
  - Supports both triple backticks and triple tildes for code fences; language annotations are allowed (e.g., ```http, ~~~json)
  - Builds structured spec objects
  - Uses precompiled regexes for links, variables, and expectations

- **executor.go**: HTTP executor
  - Substitutes variables in requests (strict: missing variables fail)
  - Parses HTTP request format
  - Executes requests with `http.Client.Timeout`
  - Manages response bodies

- **evaluator.go**: Expectation evaluator
  - Validates JSON responses
  - Compares values with operators
  - Supports multiple comparison types
  - Uses GJSON; JSON paths in spec must start with `$`

## Usage Remains Unchanged

The command-line interface remains simple, with additional safety flags:

```bash
./cotton [file|dir|glob ...] --timeout 30000 -f   # fail-fast with 30s timeout
./cotton examples/*.spec.md -q                    # quiet mode (exit code only)
```

- `--timeout <ms>`: Valid values 1..60000; capped at 60000ms.
- Missing template variables cause a parse-time failure.
- Any failed spec returns a non-zero exit code.

## Benefits for Future Development

1. **Add Web UI**: Import `engine` package, add web handlers
2. **Add REST API**: Use `engine` as service backend
3. **Add Programmatic API**: Users can import `engine` directly
4. **Multiple Output Formats**: Create JSON/XML formatters alongside existing text formatter
5. **IDE Integration**: Tools can embed engine without CLI overhead
6. **Testing**: Unit test engine and CLI independently

## Documentation

See [ARCHITECTURE.md](ARCHITECTURE.md) for:
- Detailed package structure
- Design principles
- Component responsibilities
- Extension examples
- Import guidelines

## Migration Verification

✅ All functionality preserved  
✅ Binary builds successfully  
✅ Command-line interface compatible (adds `--timeout`, stricter validation)  
✅ No breaking changes to behavior

## Refactoring Guidelines

- **Preserve boundaries**: Keep `engine` independent from `cli`.
- **Add tests first**: When changing parsing/evaluation rules, write tests in `engine/` before refactoring.
- **Respect spec**: Ensure JSON paths start with `$` and code fences include both backticks and tildes.
- **Small steps**: Prefer incremental refactors with clear responsibilities over large rewrites.
- **Shared context**: Maintain variable sharing across setup/test/teardown; avoid hidden side effects.
