//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/clock"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/result"
	"cotton/internal/stopwatch"
	"cotton/internal/testcase"
	"cotton/internal/variable"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDataFromHttpBin(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}

	mockTime, _ := time.Parse(time.RFC3339, "2024-06-26T15:27:05+07:00")

	mockClock := new(clock.MockClock)
	mockClock.On("Now").Return(mockTime)

	executableParserOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader.New(os.ReadFile),
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}
	executableParser := executable.NewParser(executableParserOptions)

	parserOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader.New(os.ReadFile),
		RequestParser:    &httphelper.HTTPRequestParser{},
		Logger:           logger.NewNilLogger(logger.Compact),
		ExecutableParser: executableParser,
		ClockWrapper:     mockClock,
	}

	executableOptions := &executable.ExecutableOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}

	testcaseOptions := &testcase.TestcaseOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
		ClockWrapper:  mockClock,
	}

	parser := testcase.NewParser(parserOptions)
	tc, err := parser.FromMarkdownFile("<rootDir>/etc/examples/httpbin.org/get.md")

	expectedSetup := executable.New("Post some data to host", `POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

secret=thisIsASecretValue`, executableOptions)
	expectedSetup.AddCapture(capture.New("secret", "$.form.secret"))

	expectedTeardown := executable.New("Patch some data to host", `PATCH https://httpbin.org/patch HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 19

secret=updatedValue`, executableOptions)

	eqOp, _ := assertion.NewOp("==")
	expectedTestcase := testcase.NewTestcase("Test GET on httpbin.org", "Test getting data from httpbin.org using multiple http requests.", "GET https://httpbin.org/get?key1=value1&key2=value2 HTTP/1.1", testcaseOptions)
	expectedTestcase.AddSetup(expectedSetup)
	expectedTestcase.AddTeardown(expectedTeardown)
	expectedTestcase.AddAssertion(assertion.New("Body.args.key1", eqOp, "value1"))
	expectedTestcase.AddAssertion(assertion.New("Body.args.key2", eqOp, "value2"))

	expectedAssertionResults := []*result.AssertionResult{
		{
			Title:    "Body.args.key1 == value1",
			Passed:   true,
			Actual:   "value1",
			Expected: "value1",
		},
		{
			Title:    "Body.args.key2 == value2",
			Passed:   true,
			Actual:   "value2",
			Expected: "value2",
		},
	}

	expectedTestResult := &result.TestResult{
		Title:        "Test GET on httpbin.org",
		Passed:       true,
		Assertions:   expectedAssertionResults,
		EllapsedTime: stopwatch.NewEllapsedTime(0),
	}

	initialVars := variable.New()
	result := tc.Execute(initialVars)

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, tc)
	assert.Equal(t, expectedTestResult, result)
}

func TestGetDataFromHttpBinWithThreeTildedCodeBlock(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}

	mockTime, _ := time.Parse(time.RFC3339, "2024-06-26T15:27:05+07:00")

	mockClock := new(clock.MockClock)
	mockClock.On("Now").Return(mockTime)

	executableParserOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader.New(os.ReadFile),
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}
	executableParser := executable.NewParser(executableParserOptions)

	parserOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader.New(os.ReadFile),
		RequestParser:    &httphelper.HTTPRequestParser{},
		Logger:           logger.NewNilLogger(logger.Compact),
		ExecutableParser: executableParser,
		ClockWrapper:     mockClock,
	}

	executableOptions := &executable.ExecutableOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}

	testcaseOptions := &testcase.TestcaseOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
		ClockWrapper:  mockClock,
	}

	parser := testcase.NewParser(parserOptions)
	tc, err := parser.FromMarkdownFile("<rootDir>/etc/examples/httpbin.org/3tildes.md")

	expectedSetup := executable.New("Post some data to host", `POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

secret=thisIsASecretValue`, executableOptions)
	expectedSetup.AddCapture(capture.New("secret", "$.form.secret"))

	expectedTeardown := executable.New("Patch some data to host", `PATCH https://httpbin.org/patch HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 19

secret=updatedValue`, executableOptions)

	eqOp, _ := assertion.NewOp("==")
	expectedTestcase := testcase.NewTestcase("Test GET on httpbin.org with three-tilded code block", "Test getting data from httpbin.org using multiple http requests.", "GET https://httpbin.org/get?key1=value1&key2=value2 HTTP/1.1", testcaseOptions)
	expectedTestcase.AddSetup(expectedSetup)
	expectedTestcase.AddTeardown(expectedTeardown)
	expectedTestcase.AddAssertion(assertion.New("Body.args.key1", eqOp, "value1"))
	expectedTestcase.AddAssertion(assertion.New("Body.args.key2", eqOp, "value2"))

	expectedAssertionResults := []*result.AssertionResult{
		{
			Title:    "Body.args.key1 == value1",
			Passed:   true,
			Actual:   "value1",
			Expected: "value1",
		},
		{
			Title:    "Body.args.key2 == value2",
			Passed:   true,
			Actual:   "value2",
			Expected: "value2",
		},
	}

	expectedTestResult := &result.TestResult{
		Title:        "Test GET on httpbin.org with three-tilded code block",
		Passed:       true,
		Assertions:   expectedAssertionResults,
		EllapsedTime: stopwatch.NewEllapsedTime(0),
	}

	initialVars := variable.New()
	result := tc.Execute(initialVars)

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, tc)
	assert.Equal(t, expectedTestResult, result)
}
