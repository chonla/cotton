//go:build integration
// +build integration

package testcase_test

import (
	"bufio"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteMarkdownFile(t *testing.T) {
	reader := reader.New(os.ReadFile)
	parser := testcase.NewParser(reader)

	curdir, _ := os.Getwd()
	result, err := parser.FromMarkdownFile(curdir + "/../../etc/examples/testcase.md")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`POST http://localhost/some-path HTTP/1.1

{
	"login": "login_name"
}`)))
	expectedBeforeRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`GET /get-info HTTP/1.0
Host: http://localhost`)))

	assert.NoError(t, err)
	assert.Equal(t, &testcase.TestCase{
		Title:       "This is title of test case written with ATX Heading 1",
		Description: "The test case is described by providing paragraphs right after the test case title.\n\nThe description of test case can be single or multiple lines.\n\nCotton will consider only the first ATX Heading 1 as the test title.",
		Request:     expectedRequest,
		Setups: []*executable.Executable{
			{
				Title:   "Link before the test will be executed before executing test",
				Request: expectedBeforeRequest,
				Captures: []*capture.Captured{
					{
						Name:    "readiness",
						Locator: "$.readiness",
					},
					{
						Name:    "version",
						Locator: "$.version",
					},
				},
			},
		},
	}, result)
}
