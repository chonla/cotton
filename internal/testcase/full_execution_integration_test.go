//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/console"
	"cotton/internal/executable"
	"cotton/internal/reader"
	"cotton/internal/request"
	"cotton/internal/testcase"
	"os"
	"testing"

	"github.com/chonla/httpreqparser"
	"github.com/stretchr/testify/assert"
)

func TestGetDataFromHttpBin(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}
	reader := reader.New(os.ReadFile)
	reqParser := httpreqparser.New()
	parser := testcase.NewParser(config, reader, reqParser)

	tc, err := parser.FromMarkdownFile("<rootDir>/etc/examples/httpbin.org/get.md")

	req, _ := reqParser.Parse(`GET https://httpbin.org/get?key1=value1&key2=value2 HTTP/1.1`)
	expectedRequest, _ := request.New(req)

	beforeReq, _ := reqParser.Parse(`POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

secret=thisIsASecretValue`)
	expectedBeforeRequest, _ := request.New(beforeReq)

	afterReq, _ := reqParser.Parse(`PATCH https://httpbin.org/patch HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 19

secret=updatedValue`)
	expectedAfterRequest, _ := request.New(afterReq)

	expectedBeforeCaptures := []*capture.Capture{
		{
			Name:    "secret",
			Locator: "$.form.secret",
		},
	}

	expectedSetups := []*executable.Executable{
		{
			Title:    "Post some data to host",
			Request:  expectedBeforeRequest,
			Captures: expectedBeforeCaptures,
		},
	}

	expectedTeardowns := []*executable.Executable{
		{
			Title:   "Patch some data to host",
			Request: expectedAfterRequest,
		},
	}

	expectedAssertions := []*assertion.Assertion{
		{
			Selector: "Body.args.key1",
			Value:    "value1",
			Operator: &assertion.EqualAssertion{},
		},
		{
			Selector: "Body.args.key2",
			Value:    "value2",
			Operator: &assertion.EqualAssertion{},
		},
	}

	expectedAssertionResults := []testcase.AssertionResult{
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

	expected := &testcase.TestCase{
		Title:       "Test GET on httpbin.org",
		Description: "Test getting data from httpbin.org using multiple http requests.",
		Setups:      expectedSetups,
		Teardowns:   expectedTeardowns,
		Request:     expectedRequest,
		Assertions:  expectedAssertions,
	}

	expectedTestResult := &testcase.TestResult{
		Title:      "Test GET on httpbin.org",
		Passed:     true,
		Assertions: expectedAssertionResults,
	}

	logger := console.NewNilConsole()
	result := tc.Execute(logger)

	assert.NoError(t, err)
	assert.True(t, expected.SimilarTo(tc))
	assert.Equal(t, expectedTestResult, result)
}

func TestGetDataFromHttpBinWithThreeTilkdedCodeBlock(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}
	reader := reader.New(os.ReadFile)
	reqParser := httpreqparser.New()
	parser := testcase.NewParser(config, reader, reqParser)

	tc, err := parser.FromMarkdownFile("<rootDir>/etc/examples/httpbin.org/3tildes.md")

	req, _ := reqParser.Parse(`GET https://httpbin.org/get?key1=value1&key2=value2 HTTP/1.1`)
	expectedRequest, _ := request.New(req)

	expectedAssertions := []*assertion.Assertion{
		{
			Selector: "Body.args.key1",
			Value:    "value1",
			Operator: &assertion.EqualAssertion{},
		},
		{
			Selector: "Body.args.key2",
			Value:    "value2",
			Operator: &assertion.EqualAssertion{},
		},
	}

	expectedAssertionResults := []testcase.AssertionResult{
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

	expected := &testcase.TestCase{
		Title:       "Test GET on httpbin.org with three-tilded code block",
		Description: "Test getting data from httpbin.org using multiple http requests.",
		Request:     expectedRequest,
		Assertions:  expectedAssertions,
	}

	expectedTestResult := &testcase.TestResult{
		Title:      "Test GET on httpbin.org with three-tilded code block",
		Passed:     true,
		Assertions: expectedAssertionResults,
	}

	logger := console.NewNilConsole()
	result := tc.Execute(logger)

	assert.NoError(t, err)
	assert.True(t, expected.SimilarTo(tc))
	assert.Equal(t, expectedTestResult, result)
}
