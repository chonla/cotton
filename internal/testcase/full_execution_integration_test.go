//go:build integration
// +build integration

package testcase_test

import (
	"cotton/internal/config"
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

	// 	expectedBeforeRequest, _ := reqParser.Parse(`POST https://httpbin.org/post HTTP/1.1
	// Content-Type: application/x-www-form-urlencoded
	// Content-Length: 25

	// secret=thisIsASecretValue`)

	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, result.Request)
	// assert.Equal(t, &testcase.TestCase{
	// 	Title:       "Test GET on httpbin.org",
	// 	Description: "Test getting data from httpbin.org using multiple http requests.",
	// 	Request:     expectedRequest,
	// 	Setups: []*executable.Executable{
	// 		{
	// 			Title:   "Post some data to host",
	// 			Request: expectedBeforeRequest,
	// 			Captures: []*capture.Capture{
	// 				{
	// 					Name:    "secret",
	// 					Locator: "$.form.secret",
	// 				},
	// 			},
	// 		},
	// 	},
	// }, result)
}
