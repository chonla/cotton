# Code Review: Cotton Project

## Executive Summary

Cotton is a well-structured markdown-based API test runner with a clean modular architecture. The project successfully separates CLI concerns from core business logic, resulting in high cohesion and low coupling. However, there are several areas for improvement in naming conventions, error handling, code organization, and documentation.

**Overall Assessment**: ⭐⭐⭐⭐ (4/5)

---

## Status Update (Post-Review Changes)

The following improvements have been implemented based on this review:

- Added HTTP timeouts with CLI `--timeout <ms>` (validated 1..60000; capped at 60s); enforced in `http.Client`.
- Introduced typed errors: `ParseError` vs `RunError`; CLI returns non-zero on any failure.
- Strict variable substitution: missing `{{ var }}` now fails fast (ParseError).
- Cycle detection in setup/teardown links to prevent recursion loops.
- Hoisted parser regexes to package scope to avoid recompilation.
- Deterministic integration tests added using `httptest`; CLI exit code behavior tested.

Remaining items and recommendations below are still applicable unless marked as addressed.

## 1. Code Structure & Modularity

### ✅ Strengths

1. **Clear Package Separation**
   - Excellent CLI/Engine separation following clean architecture principles
   - CLI has zero dependencies on internal implementation details
   - Engine is completely reusable as a library

2. **Well-Defined Responsibilities**
   - Each struct has a single, clear purpose (Parser, Executor, Evaluator, Runner, Formatter)
   - Clear data flow: Parse → Resolve Variables → Execute → Evaluate → Report

3. **Dependency Flow**
   - Unidirectional: `main.go` → `cli` → `engine`
   - No circular dependencies
   - Proper use of public/private (exported/unexported) types

### ⚠️ Issues

1. **Missing Interface Definitions**
   - `HTTPExecutor` is tightly coupled to `http.Client`
   - No interface for `HTTPExecutor`, making testing difficult
   - **Recommendation**: Define an `HTTPClient` interface for better testability

   ```go
   type HTTPClient interface {
       Do(req *http.Request) (*http.Response, error)
   }
   ```

2. **Parser Doesn't Use Interfaces**
   - `MarkdownParser` directly reads files, no file abstraction
   - **Recommendation**: Accept `io.Reader` for better composability

   ```go
   type MarkdownParser struct {
       reader io.Reader
       lines  []string
   }
   
   func NewMarkdownParserFromReader(reader io.Reader) (*MarkdownParser, error)
   ```

3. **Missing Error Handling Strategy**
   - Some errors are logged to stdout (`fmt.Printf`) instead of stderr
   - Inconsistent error wrapping
   - **Recommendation**: Use `os.Stderr` for errors, consistent error wrapping

4. **No Configuration Management**
   - CLI flags are parsed inline in `Run()` method
   - No configuration object to pass through layers
   - **Recommendation**: Create a `Config` struct to encapsulate settings

---

## 2. Naming Conventions

### ✅ Strengths

- Clear, descriptive names for functions and types
- Follows Go conventions (CamelCase for exported, camelCase for private)
- Action verbs used appropriately (New*, Execute*, Evaluate*)

### ⚠️ Issues

1. **Inconsistent Abbreviations**
   - `resp` vs `response` (mixed usage in executor.go)
   - `fs` (flag set) is cryptic - could be `flagSet`
   - **Recommendation**: Standardize across codebase

2. **Generic Type Names**
   - `TestContext` is vague - what does it contain? 
   - Better name: `SpecExecutionContext` or `TestRunContext`
   - `TestRunner` could be `SpecRunner` for clarity

3. **Misleading Variable Names**
   - In `runner.go`, `setupResult` can be misleading - it's the result of executing a setup file, not the setup data itself
   - Better: `setupExecution`, `setupFileResult`, or more specific naming

4. **Variable Abbreviations in Tests**
   - `tt` for test table is acceptable but `exp` for expectation isn't clear
   - Consider `testCase` and `expectation` for clarity

---

## 3. Consistency

### ✅ Strengths

- Consistent error handling pattern with error wrapping
- Consistent use of maps for storing variables
- Consistent test structure (table-driven tests)

### ⚠️ Issues

1. **Inconsistent Return Order**
   - `parse()` returns `([]string, map[string]string)` (links, variables)
   - Could benefit from consistent patterns (data first, error last is Go standard)

2. **Inconsistent Error Messages**
   - Some wrap with context: `fmt.Errorf("failed to execute spec '%s': %w", specFile, err)`
   - Some don't: `fmt.Errorf("failed to read response body: %w")`
   - **Recommendation**: Add context consistently

3. **Inconsistent Visibility of Helpers**
   - Some helpers are public (exported) when they should be private:
     - `resolvePath()` - used only internally in runner.go but exported
     - `WriteJSONReport()` - method on TestResults but could be function

4. **Mixed Time Representation**
   - Using `.UnixMilli()` in some places, milliseconds manually calculated in others
   - Inconsistent elapsed time calculation logic

---

## 4. Code Smell & Best Practices

### ⚠️ Code Smells Detected

1. **Long Methods**
   - `executeSpec()` in runner.go (150+ lines) does too much:
     - Parses spec
     - Executes setup
     - Resolves variables
     - Executes request
     - Evaluates expectations
     - Executes teardown
   - **Fix**: Extract setup/teardown execution to separate methods

2. **Tight Coupling in Variable Resolution**
   ```go
   // In runner.go - mixed concerns
   if strings.HasPrefix(v, "$") {
       // This is a JSON path
       result := gjson.Get(lastBody, v)
   }
   ```
   - Parser and executor logic mixed with runner
   - **Fix**: Encapsulate in dedicated method

3. **Regex Compilation in Loop**
   - `extractBeforeCodeFence()` and `extractAfterCodeFence()` compile regexes on every call
   - **Fix**: Compile at package level or struct initialization

   ```go
   // Current (inefficient)
   varRegex := regexp.MustCompile("`([^`]+)`\\s*=\\s*`([^`]+)`")
   
   // Should be
   var varRegex = regexp.MustCompile("`([^`]+)`\\s*=\\s*`([^`]+)`")
   ```

4. **Magic Strings & Numbers**
   - HTTP version "HTTP/1.1" hardcoded in tests
   - Status codes not defined as constants
   - Flag names repeated as strings

5. **Nil Checks Scattered Throughout**
   - Multiple `if result == nil` checks in runner.go
   - Could use builder pattern or defensive initialization

6. **Variable Shadowing Risk**
   ```go
   for _, setupFile := range spec.SetupLinks {
       setupPath := resolvePath(specFile, setupFile) // 'setupFile' reused
   ```

### ✅ Best Practices Followed

- Table-driven tests (parser_test.go, executor_test.go)
- Defer for resource cleanup
- Error wrapping with context
- Constants for magic values (color codes in formatter)
- Unidirectional dependency flow

### ⚠️ Best Practices Not Followed

1. **Missing Documentation Comments**
   - Public types lack descriptions
   - Complex logic needs explanation (e.g., variable resolution strategy)
   - **Example Missing Docs**:
     ```go
     // UpdateFromResponse updates context variables from the latest HTTP response
     func (tr *TestRunner) updateVariablesFromResponse(...)
     ```

2. **No Input Validation**
   - `Run()` doesn't validate file existence before discovery
   - `parseRequest()` accepts any method string without validation
   - JSON path extraction doesn't validate path format

3. **Inconsistent Logging**
   - Mix of `fmt.Printf` and `log.Fatalf`
   - Debug output goes to stdout instead of logger
   - No structured logging

4. **No Context Support**
   - HTTP requests don't use `context.Context` for timeout management
   - Can't cancel long-running requests
   - **Recommendation**: Add context support for HTTP client

5. **Error Recovery Missing**
   - Setup failures cause entire test to fail (might be intentional but not configurable)
   - No retry logic or failure strategies
   - No panic recovery

---

## 5. Specific Issues & Recommendations

### Critical

1. **HTTP Request Parser Is Fragile**
   - Assumes exact "METHOD PATH HTTP/1.1" format
   - Doesn't handle HTTPS URLs without scheme
   - No support for absolute URLs (only paths)
   - **Fix**: Improve URL handling, validate against RFC 7230

2. **Variable Substitution Can Cause Infinite Loops**
   - If variable references itself: `var1 = "{{ var1 }}"`, will substitute infinitely
   - **Fix**: Track variable dependencies, detect cycles

3. **Regex Patterns Not Escaped Properly**
   - Pattern: `` `([^`]+)` `` could miss edge cases with nested backticks
   - **Test**: Check with backticks in variable names

### Major

1. **No Timeout on HTTP Requests**
   - Client has no timeout configured
   - Long-running requests will block indefinitely
   - **Fix**: Add configurable timeout

   ```go
   client := &http.Client{
       Timeout: 30 * time.Second,
   }
   ```

2. **Variable Type Information Lost**
   - Everything stored as strings, numeric comparisons assume string parsing succeeds
   - Could fail silently with wrong output
   - **Fix**: Preserve JSON types or document limitations

3. **Path Resolution Could Have Security Issues**
   - `resolvePath()` doesn't check for `../` directory traversal
   - **Fix**: Validate paths don't escape boundaries

   ```go
   if !strings.HasPrefix(absPath, baseDir) {
       return nil, fmt.Errorf("path escapes base directory")
   }
   ```

4. **No Support for File Encoding**
   - Assumes UTF-8, no BOM handling
   - Could fail on different encodings

### Minor

1. **Unused Exported Functions**
   - `resolvePath()` exported but used only internally
   - `WriteJSONReport()` is method but could be function
   - **Fix**: Make private or create proper API

2. **Test Files Missing Edge Cases**
   - No tests for concurrent execution
   - No tests for variable circular dependencies
   - No tests for large responses

3. **Empty Code Fence Handling**
   - Parser silently ignores specs with no code fence
   - No warning to user
   - **Fix**: Warn about specs without requests

---

## 6. Performance Considerations

### Current Issues

1. **Regex Recompilation**
   - Patterns compiled on every parse
   - **Impact**: Noticeable with many files
   - **Fix**: Cache compiled regexes

2. **Inefficient Variable Substitution**
   - `strings.ReplaceAll()` called for each variable
   - **Better**: Single pass with template engine
   - **Fix**: Use `text/template` or custom tokenizer

3. **Full Body Reads for Each Expectation**
   - Response body read twice: once for expectations, once for variable extraction
   - **Fix**: Parse once, reuse

4. **No Response Streaming**
   - Entire response loaded into memory
   - **Fix**: Consider streaming for large responses (breaking change)

---

## 7. Testing

### ✅ Strengths

- Good test coverage for core parsers
- Table-driven test style
- Clear test names

### ⚠️ Issues

1. **Test Coverage Gaps**
   - runner_test.go appears empty or minimal
   - No integration tests with actual HTTP calls
   - Evaluator tests missing
   - formatter tests missing

2. **Mock/Stub Limitations**
   - Tests don't mock HTTP client
   - Can't test HTTP error scenarios
   - **Fix**: Add HTTPClient interface and mock implementations

3. **Missing Negative Tests**
   - No tests for malformed markdown
   - No tests for invalid expectations
   - No tests for missing files in setup links

---

## 8. Documentation

### ✅ What's Good

- README is clear and comprehensive
- SPEC.md explains format well
- ARCHITECTURE.md documents design decisions

### ⚠️ What's Missing

1. **Code Comments**
   - Complex algorithms need explanation
   - Regex patterns need comments explaining groups
   - JSON path resolution needs documentation

2. **API Documentation**
   - No godoc comments on exported types
   - No usage examples in code comments
   - **Example Missing**:
     ```go
     // ParsedSpec represents a parsed Markdown API specification.
     // It contains the test case metadata, HTTP request details, 
     // and expectations to validate against the response.
     type ParsedSpec struct { ... }
     ```

3. **Troubleshooting Guide**
   - No guide for common errors
   - No FAQ for spec format questions

---

## 9. Security Considerations

### ⚠️ Issues

1. **No Input Validation**
   - File paths not validated for escape sequences
   - HTTP headers accepted without validation
   - Variable values not sanitized

2. **Potential Code Injection**
   - HTTP request parsing could be abused
   - Variables directly substituted without escaping

3. **No SSL/TLS Configuration**
   - No control over certificate validation
   - Potential for MITM attacks in untrusted networks

---

## 10. Summary of Recommendations

### Priority 1 (Critical)
- [x] Add HTTP request timeout configuration
- [ ] Add path escape validation in resolvePath()
- [ ] Implement HTTPClient interface for testability
- [ ] Add documentation comments to exported types

### Priority 2 (Important)
- [ ] Extract long methods (executeSpec > 150 lines)
- [x] Cache compiled regex patterns
- [ ] Add input validation for HTTP requests
- [ ] Implement context.Context support
- [ ] Standardize error output (use os.Stderr)
- [x] Add comprehensive integration tests

### Priority 3 (Nice to Have)
- [ ] Refactor variable substitution (consider template engine)
- [ ] Add performance benchmarks
- [ ] Implement proper logging (structured logging)
- [ ] Add configuration struct
- [ ] Support for more HTTP features (redirects, custom client config)
- [ ] Add profiling support for large test suites

### Priority 4 (Future)
- [ ] Parallel test execution support
- [ ] Response streaming for large payloads
- [ ] HTML report generation (mentioned but not implemented)
- [ ] Interactive test builder
- [ ] IDE integration plugins

---

## Conclusion

Cotton is a well-designed project with a solid foundation in clean architecture. The separation of concerns is exemplary for a Go project. The main areas needing attention are:

1. **Robustness**: Add validation, timeouts, and better error handling
2. **Maintainability**: Extract long methods, add documentation
3. **Testability**: Add interfaces for dependencies, improve test coverage
4. **Performance**: Cache compiled regexes, optimize variable handling

With the implemented improvements (timeouts, typed errors, cycle detection, regex caching, and deterministic tests), Cotton has significantly improved robustness. Addressing the remaining items will further strengthen maintainability, testability, and performance.

---

*Review updated: December 18, 2025*
