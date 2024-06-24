package main

import (
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"fmt"
	"os"
)

func main() {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir,
	}

	level := logger.Debug
	log := logger.NewTerminalLogger(level)
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

	ts, err := testcase.NewTestsuite("<rootDir>/etc/examples", tcOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = ts.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
