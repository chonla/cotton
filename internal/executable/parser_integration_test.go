//go:build integration
// +build integration

package executable_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/reader"
	"cotton/internal/request"
	"os"
	"testing"

	"github.com/chonla/httpreqparser"
	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteExecutableMarkdownFile(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}

	reader := reader.New(os.ReadFile)
	reqParser := httpreqparser.New()
	parser := executable.NewParser(config, reader, reqParser)

	result, err := parser.FromMarkdownFile("<rootDir>/etc/examples/executable_before.md")

	req, _ := reqParser.Parse(`GET /get-info HTTP/1.1
Host: localhost`)
	expectedRequest, _ := request.New(req)
	expectedCaptures := []*capture.Capture{
		{
			Name:    "readiness",
			Locator: "$.readiness",
		},
		{
			Name:    "version",
			Locator: "$.version",
		},
	}
	expectedExecutable := &executable.Executable{
		Request:  expectedRequest,
		Captures: expectedCaptures,
	}
	assert.NoError(t, err)
	assert.True(t, expectedExecutable.SimilarTo(result))
}
