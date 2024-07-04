package main

import (
	"cotton/internal/clock"
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
	var detailedDebug bool
	var compact bool
	var help bool
	var ver bool
	var insecure bool
	var stopWhenFailed bool
	var customRootDir string

	flag.Usage = usage
	flag.BoolVar(&compact, "c", false, "compact mode")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&detailedDebug, "p", false, "paranoid mode")
	flag.BoolVar(&insecure, "i", false, "disable certificate verification")
	flag.BoolVar(&stopWhenFailed, "s", false, "stop when test failed")
	flag.BoolVar(&ver, "v", false, "display cotton version")
	flag.BoolVar(&help, "h", false, "display this help")
	flag.StringVar(&customRootDir, "r", "", "set rootDir path")
	flag.Parse()

	rootDir := ""
	if customRootDir == "" {
		rootDir, _ = os.Getwd()
	} else {
		rootDir = customRootDir
	}
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
	if detailedDebug {
		level = logger.DetailedDebug
	}

	if ver {
		fmt.Printf("cotton %s\n", Version)
		os.Exit(0)
	}

	if help {
		usage()
		os.Exit(0)
	}

	log := logger.NewTerminalLogger(level)
	reader := reader.New(os.ReadFile)
	reqParser := &httphelper.HTTPRequestParser{}
	clockWrapper := clock.New()

	exOptions := &executable.ParserOptions{
		Configurator:  config,
		FileReader:    reader,
		RequestParser: reqParser,
		Logger:        log,
		ClockWrapper:  clockWrapper,
	}
	exParser := executable.NewParser(exOptions)

	tcOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader,
		RequestParser:    reqParser,
		Logger:           log,
		ExecutableParser: exParser,
		ClockWrapper:     clockWrapper,
	}

	options := &testcase.TestsuiteOptions{
		TestcaseParserOption: tcOptions,
		StopWhenFailed:       stopWhenFailed,
		Logger:               log,
		ClockWrapper:         clockWrapper,
	}

	ts, err := testcase.NewTestsuite(testDir, options)
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

  cotton [-d] [-c] [-p] [-r <rootdir>] <testpath|testdir>
  cotton -v
  cotton --help

`)
	flag.PrintDefaults()
}
