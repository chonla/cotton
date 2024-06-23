package main

import (
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"cotton/internal/variable"
	"fmt"
	"os"
)

func main() {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir,
	}

	debug := true
	log := logger.NewTerminalLogger(debug)
	reader := reader.New(os.ReadFile)
	reqParser := &httphelper.HTTPRequestParser{}

	exOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader,
		RequestParser: reqParser,
		Logger:        log,
	}
	exParser := executable.NewParser(exOptions)

	tcOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader,
		RequestParser:    reqParser,
		Logger:           log,
		ExecutableParser: exParser,
	}
	parser := testcase.NewParser(tcOptions)

	tc, err := parser.FromMarkdownFile("<rootDir>/etc/examples/mixed/test.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initialVars := variable.New()

	result := tc.Execute(initialVars)
	log.PrintTestResult(result.Passed)
	log.PrintAssertionResults(result.Assertions)
}
