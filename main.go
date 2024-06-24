package main

import (
	"cotton/internal/config"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/testcase"
	"flag"
	"fmt"
	"os"
)

const Version = "1.0.0"

func main() {
	var debug bool
	var compact bool
	var help bool
	var ver bool
	var insecure bool

	flag.Usage = usage
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&compact, "c", false, "compact mode")
	flag.BoolVar(&ver, "v", false, "display cotton version")
	flag.BoolVar(&help, "h", false, "display this help")
	flag.BoolVar(&insecure, "i", false, "disable certificate verification")
	flag.Parse()

	rootDir, _ := os.Getwd()
	config := &config.Config{
		RootDir: rootDir,
	}

	testDir := flag.Arg(0)
	if testDir == "" {
		testDir = rootDir
	}

	level := logger.Verbose
	if compact {
		level = logger.Compact
	}
	if debug {
		level = logger.Debug
	}

	if ver {
		fmt.Printf("cotton %s\n", Version)
		os.Exit(0)
	}

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

	ts, err := testcase.NewTestsuite(testDir, tcOptions)
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

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), `Usage of cotton:

  cotton [-d] [-c] <testpath|testdir>
  cotton -v
  cotton --help

`)
	flag.PrintDefaults()
}
