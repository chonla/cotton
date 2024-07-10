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
	"cotton/internal/variable"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingCompleteExecutableMarkdownFile(t *testing.T) {
	curdir, _ := os.Getwd()
	config := &config.Config{
		BaseDir: curdir + "/../..",
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
	expectedExecutable := executable.New("Untitled", `POST https://fakestoreapi.com/auth/login HTTP/1.1
Content-Type: application/json
Content-Length: 43

{"username":"{{username}}","password":"{{password}}"}`, executableOptions)
	expectedExecutable.AddCapture(capture.New("access_token", "Body.token"))
	expectedExecutable.AddVariable(&variable.Variable{
		Name:  "username",
		Value: "mor_2314",
	})
	expectedExecutable.AddVariable(&variable.Variable{
		Name:  "password",
		Value: "83r5^_",
	})

	result, err := parser.FromMarkdownFile("../../etc/examples/fakestoreapi.com/executables/auth.md")

	assert.NoError(t, err)
	assert.Equal(t, expectedExecutable, result)
}
