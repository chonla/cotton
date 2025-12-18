package engine

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// MarkdownParser parses markdown spec files
type MarkdownParser struct {
	content string
	lines   []string
}

// ParsedSpec represents the parsed markdown specification
type ParsedSpec struct {
	Title         string
	SetupLinks    []string
	Variables     map[string]string
	RequestBody   string
	Expectations  []Expectation
	TeardownLinks []string
	IsTestCase    bool
}

// Expectation represents an expectation rule
type Expectation struct {
	Field    string
	Operator string
	Value    string
}

var (
	linkRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	varRegex  = regexp.MustCompile("`([^`]+)`\\s*=\\s*`([^`]+)`")
	expRegex  = regexp.MustCompile("`([^`]+)`\\s*(==|!=|<|>|<=|>=|=~|!~)\\s*`([^`]*)`|`([^`]+)`\\s*(==|!=|<|>|<=|>=|=~|!~)\\s*(undefined|null)")
)

// NewMarkdownParser creates a new parser
func NewMarkdownParser(filePath string) (*MarkdownParser, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	return &MarkdownParser{
		content: string(content),
		lines:   lines,
	}, nil
}

// Parse parses the markdown spec
func (p *MarkdownParser) Parse() (*ParsedSpec, error) {
	spec := &ParsedSpec{
		Variables:     make(map[string]string),
		Expectations:  []Expectation{},
		SetupLinks:    []string{},
		TeardownLinks: []string{},
	}

	// Extract title (first heading)
	for _, line := range p.lines {
		if strings.HasPrefix(line, "# ") {
			spec.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}

	// Find code fence position (support ``` and ~~~, with optional language annotations)
	codeFenceStart := -1
	codeFenceEnd := -1
	fenceMarker := ""
	for i, line := range p.lines {
		// Detect start or end of a code fence; allow language annotations (e.g., ```http or ~~~json)
		isBacktickFence := strings.HasPrefix(line, "```")
		isTildeFence := strings.HasPrefix(line, "~~~")
		if isBacktickFence || isTildeFence {
			if codeFenceStart == -1 {
				codeFenceStart = i
				if isBacktickFence {
					fenceMarker = "```"
				} else {
					fenceMarker = "~~~"
				}
			} else if codeFenceEnd == -1 {
				// End fence must match the marker used at start
				if fenceMarker != "" && strings.HasPrefix(line, fenceMarker) {
					codeFenceEnd = i
					break
				}
			}
		}
	}

	if codeFenceStart != -1 && codeFenceEnd != -1 {
		// Extract content before code fence (setup, variables)
		spec.SetupLinks, spec.Variables = p.extractBeforeCodeFence(codeFenceStart)

		// Extract request body (content inside code fence)
		spec.RequestBody = p.extractCodeFenceContent(codeFenceStart, codeFenceEnd)

		// Extract content after code fence (expectations, teardown)
		spec.Expectations, spec.TeardownLinks = p.extractAfterCodeFence(codeFenceEnd)
	}

	// A spec is a "test case" if it has one or more expectations
	spec.IsTestCase = len(spec.Expectations) > 0

	return spec, nil
}

// extractBeforeCodeFence extracts setup links and variables from content before the code fence
func (p *MarkdownParser) extractBeforeCodeFence(codeFenceStart int) ([]string, map[string]string) {
	setupLinks := []string{}
	variables := make(map[string]string)

	for i := 0; i < codeFenceStart; i++ {
		line := strings.TrimSpace(p.lines[i])

		// Extract links [text](url)
		if matches := linkRegex.FindStringSubmatch(line); matches != nil {
			setupLinks = append(setupLinks, matches[2])
		}

		// Extract variables: `<name>` = `<value>`
		if matches := varRegex.FindStringSubmatch(line); matches != nil {
			variables[matches[1]] = matches[2]
		}
	}

	return setupLinks, variables
}

// extractCodeFenceContent extracts the content inside the code fence
func (p *MarkdownParser) extractCodeFenceContent(start, end int) string {
	// Content between the code fences (excluding the fence markers)
	if start+1 >= end {
		return ""
	}

	var content []string
	for i := start + 1; i < end; i++ {
		content = append(content, p.lines[i])
	}
	return strings.Join(content, "\n")
}

// extractAfterCodeFence extracts expectations and teardown links from content after the code fence
func (p *MarkdownParser) extractAfterCodeFence(codeFenceEnd int) ([]Expectation, []string) {
	expectations := []Expectation{}
	teardownLinks := []string{}

	for i := codeFenceEnd + 1; i < len(p.lines); i++ {
		line := strings.TrimSpace(p.lines[i])

		// Extract links [text](url)
		if matches := linkRegex.FindStringSubmatch(line); matches != nil {
			teardownLinks = append(teardownLinks, matches[2])
		}

		// Extract expectations: `<field>` <op> `<value>` or `<field>` <op> null/undefined
		// Pattern: backticks around field, operator, then backticks around value or null/undefined
		if matches := expRegex.FindStringSubmatch(line); matches != nil {
			field := matches[1]
			op := matches[2]
			value := matches[3]

			if field == "" {
				field = matches[4]
				op = matches[5]
				value = matches[6]
			}

			if field != "" {
				expectations = append(expectations, Expectation{
					Field:    field,
					Operator: op,
					Value:    value,
				})
			}
		}
	}

	return expectations, teardownLinks
}
