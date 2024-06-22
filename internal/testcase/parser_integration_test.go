//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/assertion"
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteTestcaseMarkdownFile(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}

	executableParserOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader.New(os.ReadFile),
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(false),
	}
	executableParser := executable.NewParser(executableParserOptions)

	parserOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader.New(os.ReadFile),
		RequestParser:    &httphelper.HTTPRequestParser{},
		Logger:           logger.NewNilLogger(false),
		ExecutableParser: executableParser,
	}
	parser := testcase.NewParser(parserOptions)

	testcaseOptions := &testcase.TestCaseOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(false),
	}
	result, err := parser.FromMarkdownFile("<rootDir>/etc/examples/testcase.md")

	executableOptions := &executable.ExecutableOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(false),
	}
	expectedSetup := executable.New("Link before the test will be executed before executing test", `GET /get-info HTTP/1.1
Host: localhost`, executableOptions)
	expectedSetup.AddCapture(capture.New("readiness", "$.readiness"))
	expectedSetup.AddCapture(capture.New("version", "$.version"))

	expectedTeardown := executable.New("Link after the test will be executed after executing test", `GET /time-taken HTTP/1.1
Host: localhost`, executableOptions)
	expectedTeardown.AddCapture(capture.New("time", "$.millisec"))

	expectedTestcase := testcase.New("This is title of test case written with ATX Heading 1", "The test case is described by providing paragraphs right after the test case title.\n\nThe description of test case can be single or multiple lines.\n\nCotton will consider only the first ATX Heading 1 as the test title.", `POST /some-path HTTP/1.1
Host: localhost

{
    "login": "login_name"
}`, testcaseOptions)
	expectedTestcase.AddSetup(expectedSetup)
	expectedTestcase.AddTeardown(expectedTeardown)

	eqOp, _ := assertion.NewOp("==")
	expectedTestcase.AddAssertion(assertion.New("$.form.result", eqOp, "success"))
	expectedTestcase.AddAssertion(assertion.New("$.form.result.length", eqOp, float64(1)))

	assert.NoError(t, err)
	assert.Equal(t, expectedTestcase, result)
}
