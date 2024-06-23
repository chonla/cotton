//go:build integration
// +build integration

package executable_test

import (
	"cotton/internal/capture"
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteExecutableMarkdownFile(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir + "/../..",
	}

	parserOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader.New(os.ReadFile),
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}
	executableOptions := &executable.ExecutableOptions{
		RequestParser: &httphelper.HTTPRequestParser{},
		Logger:        logger.NewNilLogger(logger.Compact),
	}

	parser := executable.NewParser(parserOptions)
	expectedExecutable := executable.New("Untitled", `GET /get-info HTTP/1.1
Host: localhost`, executableOptions)
	expectedExecutable.AddCapture(capture.New("readiness", "$.readiness"))
	expectedExecutable.AddCapture(capture.New("version", "$.version"))

	result, err := parser.FromMarkdownFile("<rootDir>/etc/examples/executable_before.md")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
}
