package main

import (
	"cotton/internal/config"
	"cotton/internal/console"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"fmt"
	"os"

	"github.com/chonla/httpreqparser"
	"github.com/kr/pretty"
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

	logger := console.NewTerminal()
	result := tc.Execute(logger)

	pretty.Println(result)
}
