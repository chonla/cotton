package main

import (
	"cotton/internal/config"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"fmt"
	"os"

	"github.com/chonla/httpreqparser"
)

func main() {
	curdir, _ := os.Getwd()
	config := &config.Config{
		RootDir: curdir,
	}
	reader := reader.New(os.ReadFile)
	reqParser := httpreqparser.New()
	parser := testcase.NewParser(config, reader, reqParser)

	tc, err := parser.FromMarkdownFile("<rootDir>/etc/examples/opensource.org/get_copyleft.md")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log := logger.NewTerminalLogger()
	result := tc.Execute(log)
	log.PrintTestResult(result.Passed)
}
