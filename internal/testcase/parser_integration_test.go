//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"os"
	"testing"

	"github.com/chonla/httpreqparser"
	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteTestcaseMarkdownFile(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}

	reader := reader.New(os.ReadFile)
	reqParser := httpreqparser.New()
	parser := testcase.NewParser(config, reader, reqParser)

	result, err := parser.FromMarkdownFile("<rootDir>/etc/examples/testcase.md")

	expectedRequestInBefore, _ := reqParser.Parse(`GET /get-info HTTP/1.1
Host: localhost`)
	expectedRequestInAfter, _ := reqParser.Parse(`GET /time-taken HTTP/1.1
Host: localhost`)
	expectedRequest, _ := reqParser.Parse(`POST /some-path HTTP/1.1
Host: localhost

{
    "login": "login_name"
}`)
	expectedCapturesInBefore := []*capture.Capture{
		{
			Name:    "readiness",
			Locator: "$.readiness",
		},
		{
			Name:    "version",
			Locator: "$.version",
		},
	}
	expectedCapturesInAfter := []*capture.Capture{
		{
			Name:    "time",
			Locator: "$.millisec",
		},
	}
	expectedSetups := []*executable.Executable{
		{
			Title:    "Link before the test will be executed before executing test",
			Request:  expectedRequestInBefore,
			Captures: expectedCapturesInBefore,
		},
	}
	expectedTeardowns := []*executable.Executable{
		{
			Title:    "Link after the test will be executed after executing test",
			Request:  expectedRequestInAfter,
			Captures: expectedCapturesInAfter,
		},
	}
	expected := &testcase.TestCase{
		Title:       "This is title of test case written with ATX Heading 1",
		Description: "The test case is described by providing paragraphs right after the test case title.\n\nThe description of test case can be single or multiple lines.\n\nCotton will consider only the first ATX Heading 1 as the test title.",
		Setups:      expectedSetups,
		Teardowns:   expectedTeardowns,
		Request:     expectedRequest,
	}

	assert.NoError(t, err)
	assert.True(t, expected.SimilarTo(result))
}
