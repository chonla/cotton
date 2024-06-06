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

func TestGetDataFromHttpBin(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}
	reader := reader.New(os.ReadFile)
	reqParser := httpreqparser.New()
	parser := testcase.NewParser(config, reader, reqParser)

	result, err := parser.FromMarkdownFile("<rootDir>/etc/examples/httpbin.org/get.md")

	expectedRequest, _ := reqParser.Parse(`GET https://httpbin.org/get?key1=value1 HTTP/1.1`)

	expectedBeforeRequest, _ := reqParser.Parse(`POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 25

secret=thisIsASecretValue`)

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

	assert.NoError(t, err)
	assert.Equal(t, "Test GET on httpbin.org", result.Title)
	assert.Equal(t, "Test getting data from httpbin.org using multiple http requests.", result.Description)
	assert.True(t, request.Similar(expectedRequest, result.Request))
	assert.Equal(t, len(expectedSetups), len(result.Setups))
	assert.True(t, expectedSetups[0].SimilarTo(result.Setups[0]))
}
