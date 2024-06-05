//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/reader"
	"cotton/internal/request"
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

	assert.NoError(t, err)
	assert.Equal(t, "This is title of test case written with ATX Heading 1", result.Title)
	assert.Equal(t, "The test case is described by providing paragraphs right after the test case title.\n\nThe description of test case can be single or multiple lines.\n\nCotton will consider only the first ATX Heading 1 as the test title.", result.Description)
	assert.Equal(t, []*executable.Executable{
		{
			Title:    "Link before the test will be executed before executing test",
			Request:  expectedRequestInBefore,
			Captures: expectedCapturesInBefore,
		},
	}, result.Setups)
	assert.True(t, request.Similar(expectedRequest, result.Request))
	// assert.Equal(t, &testcase.TestCase{
	// 	Title:       "This is title of test case written with ATX Heading 1",
	// 	Description: "The test case is described by providing paragraphs right after the test case title.\n\nThe description of test case can be single or multiple lines.\n\nCotton will consider only the first ATX Heading 1 as the test title.",
	// 	Request:     expectedRequest,
	// 	Setups: []*executable.Executable{
	// 		{
	// 			Title:    "Link before the test will be executed before executing test",
	// 			Request:  expectedRequestInBefore,
	// 			Captures: expectedCapturesInBefore,
	// 		},
	// 	},
	// }, result)
}
