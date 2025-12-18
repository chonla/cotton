package engine

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// HTTPExecutor executes HTTP requests
type HTTPExecutor struct {
	client    *http.Client
	variables map[string]interface{}
	lastBody  string
}

const (
	defaultHTTPTimeout = 30 * time.Second
	maxHTTPTimeout     = 60 * time.Second
)

var placeholderRegexp = regexp.MustCompile(`{{\s*([a-zA-Z0-9_.-]+)\s*}}`)

// NewHTTPExecutor creates a new HTTP executor
func NewHTTPExecutor() *HTTPExecutor {
	return NewHTTPExecutorWithTimeout(defaultHTTPTimeout)
}

// NewHTTPExecutorWithTimeout creates a new HTTP executor with timeout
func NewHTTPExecutorWithTimeout(timeout time.Duration) *HTTPExecutor {
	// Enforce a sensible default and cap if callers pass invalid values
	if timeout <= 0 || timeout > maxHTTPTimeout {
		timeout = defaultHTTPTimeout
	}

	return &HTTPExecutor{
		client:    &http.Client{Timeout: timeout},
		variables: make(map[string]interface{}),
	}
}

// SetVariables sets multiple variables for substitution
func (e *HTTPExecutor) SetVariables(vars map[string]string) {
	for k, v := range vars {
		e.variables[k] = v
	}
}

// SetVariable sets a single variable
func (e *HTTPExecutor) SetVariable(key string, value interface{}) {
	e.variables[key] = value
}

// GetLastBody returns the last response body
func (e *HTTPExecutor) GetLastBody() string {
	return e.lastBody
}

// ExecuteRequest executes an HTTP request from the request spec
func (e *HTTPExecutor) ExecuteRequest(requestSpec string) (*http.Response, error) {
	// Substitute variables in the request spec
	substituted, err := e.substituteVariables(requestSpec)
	if err != nil {
		return nil, err
	}

	// Parse the request
	req, err := e.parseRequest(substituted)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	// Execute the request
	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Read and store the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))
	e.lastBody = string(body)

	return resp, nil
}

// substituteVariables replaces {{ varname }} with variable values and fails when a variable is missing
func (e *HTTPExecutor) substituteVariables(text string) (string, error) {
	matches := placeholderRegexp.FindAllStringSubmatchIndex(text, -1)
	if len(matches) == 0 {
		return text, nil
	}

	var builder strings.Builder
	last := 0

	for _, m := range matches {
		start, end := m[0], m[1]
		keyStart, keyEnd := m[2], m[3]
		key := text[keyStart:keyEnd]

		value, ok := e.variables[key]
		if !ok {
			return "", newParseError(fmt.Sprintf("missing variable %q", key))
		}

		builder.WriteString(text[last:start])
		builder.WriteString(fmt.Sprintf("%v", value))
		last = end
	}

	builder.WriteString(text[last:])
	return builder.String(), nil
}

// parseRequest parses an HTTP request from text format
func (e *HTTPExecutor) parseRequest(requestText string) (*http.Request, error) {
	lines := strings.Split(requestText, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	// Parse request line (METHOD PATH HTTP/1.1)
	requestLine := strings.Fields(strings.TrimSpace(lines[0]))
	if len(requestLine) < 2 {
		return nil, fmt.Errorf("invalid request line: %s", lines[0])
	}

	method := requestLine[0]
	path := requestLine[1]

	// Find headers and body separator
	headerEnd := 1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			headerEnd = i
			break
		}
	}

	// Parse headers
	headers := make(map[string]string)
	for i := 1; i < headerEnd; i++ {
		line := strings.TrimSpace(lines[i])
		if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Extract body
	var body string
	if headerEnd+1 < len(lines) {
		body = strings.Join(lines[headerEnd+1:], "\n")
	}

	// Create request
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
