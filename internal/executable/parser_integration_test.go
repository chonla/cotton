//go:build integration
// +build integration

package executable_test

import (
	"bufio"
	"cotton/internal/capture"
	"cotton/internal/executable"
	"cotton/internal/reader"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteMarkdownFile(t *testing.T) {
	reader := reader.New(os.ReadFile)
	parser := executable.NewParser(reader)

	curdir, _ := os.Getwd()
	result, err := parser.FromMarkdownFile(curdir + "/../../etc/examples/executable.md")

	expectedRequest, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(`GET /get-info HTTP/1.0
Host: http://localhost`)))
	expectedCaptures := []*capture.Captured{
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
	assert.Equal(t, &executable.Executable{
		Request:  expectedRequest,
		Captures: expectedCaptures,
	}, result)
}
