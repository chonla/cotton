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

	expectedAfterRequest, _ := reqParser.Parse(`PATCH https://httpbin.org/patch HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 19

secret=updatedValue`)

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

	expected := &testcase.TestCase{
		Title:       "Test GET on httpbin.org",
		Description: "Test getting data from httpbin.org using multiple http requests.",
		Setups:      expectedSetups,
		Teardowns:   expectedTeardowns,
		Request:     expectedRequest,
	}

	assert.NoError(t, err)
	assert.True(t, expected.SimilarTo(result))
}
