package testcase

import (
	"cotton/internal/directory"
	"cotton/internal/result"
	"cotton/internal/variable"
	"errors"
	"fmt"
	"strings"
)

type Testsuite struct {
	testcases []*Testcase

	options *ParserOptions
}

func NewTestsuite(path string, options *ParserOptions) (*Testsuite, error) {
	actualPath := options.Configurator.ResolvePath(path)
	dir := directory.New()
	files, err := dir.ListFiles(actualPath, "md")
	if err != nil {
		return nil, err
	}

	options.Logger.PrintDebugMessage(fmt.Sprintf("Test path: %s", actualPath))
	options.Logger.PrintDebugMessage(fmt.Sprintf("\nFile(s) scanned:\n- %s", strings.Join(files, "\n- ")))

	testcases := []*Testcase{}

	parser := NewParser(options)
	for _, file := range files {
		tc, err := parser.FromMarkdownFile(file)
		if err == nil && len(tc.assertions) > 0 {
			testcases = append(testcases, tc)
		}
	}

	if len(testcases) == 0 {
		return nil, errors.New("no executable testcase found")
	}

	return &Testsuite{
		testcases: testcases,
		options:   options,
	}, nil
}

func (ts *Testsuite) Execute() (*result.TestsuiteResult, error) {
	initialVars := variable.New()

	testsuiteResult := &result.TestsuiteResult{
		PassedCount: 0,
		TestCount:   len(ts.testcases),
	}

	for index, tc := range ts.testcases {
		// ts.options.Logger.PrintTestcaseSequence(index+1, testsuiteResult.TestCount)
		section := fmt.Sprintf("testcase %d/%d", index+1, testsuiteResult.TestCount)
		ts.options.Logger.PrintSectionedMessage(section, tc.title)

		result := tc.Execute(initialVars)
		if result.Passed {
			testsuiteResult.PassedCount += 1
		}
		ts.options.Logger.PrintInlineTestResult(result.Passed)
		ts.options.Logger.PrintSectionTitle("result")
		ts.options.Logger.PrintTestResult(result.Passed)
	}
	ts.options.Logger.PrintTestsuiteResult(testsuiteResult)

	return testsuiteResult, nil
}
