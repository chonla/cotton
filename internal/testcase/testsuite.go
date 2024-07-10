package testcase

import (
	"cotton/internal/clock"
	"cotton/internal/directory"
	"cotton/internal/logger"
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/variable"
	"errors"
	"fmt"
	"strings"
)

type TestsuiteOptions struct {
	StopWhenFailed       bool
	Logger               logger.Logger
	TestcaseParserOption *ParserOptions
	ClockWrapper         clock.ClockWrapper
}

type Testsuite struct {
	testcases []*Testcase

	options *TestsuiteOptions
}

func NewTestsuite(path string, options *TestsuiteOptions) (*Testsuite, error) {
	dir := directory.New()
	files, err := dir.ListFiles(path, "md")
	if err != nil {
		return nil, err
	}

	options.Logger.PrintDebugMessage(fmt.Sprintf("Test path: %s", path))
	options.Logger.PrintDebugMessage(fmt.Sprintf("\nFile(s) scanned:\n- %s", strings.Join(files, "\n- ")))

	testcases := []*Testcase{}

	parser := NewParser(options.TestcaseParserOption)
	for _, file := range files {
		tc, err := parser.FromMarkdownFile(file)
		if err == nil && len(tc.assertions) > 0 {
			testcases = append(testcases, tc)
		} else {
			if err != nil {
				options.Logger.PrintError(file, err)
				if options.StopWhenFailed {
					return nil, err
				}
			}
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
		PassedCount:     0,
		ExecutionsCount: 0,
		FailedCount:     0,
		SkippedCount:    0,
		TestCount:       len(ts.testcases),
		EllapsedTime:    nil,
	}

	watch := stopwatch.New(ts.options.ClockWrapper)
	watch.Start()

	for index, tc := range ts.testcases {
		section := fmt.Sprintf("testcase %d/%d", index+1, testsuiteResult.TestCount)
		ts.options.Logger.PrintSectionedMessage(section, tc.Title())

		result := tc.Execute(initialVars)
		if result.Passed {
			testsuiteResult.PassedCount += 1
		} else {
			testsuiteResult.FailedCount += 1
		}
		testsuiteResult.ExecutionsCount += 1
		ts.options.Logger.PrintInlineTestResult(result.Passed)
		ts.options.Logger.PrintInlineTimeEllapsed(result.EllapsedTime)
		ts.options.Logger.PrintSectionTitle("result")
		ts.options.Logger.PrintTestResult(result.Passed)
		ts.options.Logger.PrintSectionTitle("ellapsed")
		ts.options.Logger.PrintTimeEllapsed(result.EllapsedTime)
		if !result.Passed && ts.options.StopWhenFailed {
			break
		}
	}
	testsuiteResult.SkippedCount = testsuiteResult.TestCount - testsuiteResult.ExecutionsCount
	testsuiteResult.EllapsedTime = watch.Stop()
	ts.options.Logger.PrintTestsuiteResult(testsuiteResult)

	return testsuiteResult, nil
}
