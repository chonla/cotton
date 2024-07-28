package main

import (
	"cotton/internal/clock"
	"cotton/internal/config"
	"cotton/internal/directory"
	"cotton/internal/executable"
	"cotton/internal/httphelper"
	"cotton/internal/logger"
	"cotton/internal/reader"
	"cotton/internal/reporter"
	"cotton/internal/testcase"
	"flag"
	"fmt"
	"os"
)

var Version = "1.0.0"

func main() {
	var debug bool
	var detailedDebug bool
	var compact bool
	var help bool
	var ver bool
	var insecure bool
	var stopWhenFailed bool
	var customBaseDir string
	var reporterType string

	flag.Usage = usage
	flag.BoolVar(&compact, "c", false, "compact mode")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.BoolVar(&detailedDebug, "p", false, "paranoid mode")
	flag.BoolVar(&insecure, "i", false, "disable certificate verification")
	flag.BoolVar(&stopWhenFailed, "s", false, "stop when test failed")
	flag.BoolVar(&ver, "v", false, "display cotton version")
	flag.BoolVar(&help, "h", false, "display this help")
	flag.StringVar(&customBaseDir, "b", "", "set base directory path")
	flag.StringVar(&reporterType, "r", "", "set reporter type, html or ctrf")
	flag.Parse()

	testPath := flag.Arg(0)
	if testPath == "" {
		testPath = "./"
	}

	dir := directory.New()
	testBaseDir, err := dir.DirectoryOf(testPath)

	baseDir := ""
	if customBaseDir == "" {
		if err == nil {
			baseDir = testBaseDir
		} else {
			baseDir, _ = os.Getwd()
		}
	} else {
		baseDir = customBaseDir
	}
	config := &config.Config{
		BaseDir: baseDir,
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
		Configurator:    config,
		FileReader:      reader,
		RequestParser:   reqParser,
		Logger:          log,
		ClockWrapper:    clockWrapper,
		InsecureRequest: insecure,
	}
	exParser := executable.NewParser(exOptions)

	tcOptions := &testcase.ParserOptions{
		Configurator:     config,
		FileReader:       reader,
		RequestParser:    reqParser,
		Logger:           log,
		ExecutableParser: exParser,
		ClockWrapper:     clockWrapper,
		InsecureRequest:  insecure,
	}

	options := &testcase.TestsuiteOptions{
		TestcaseParserOption: tcOptions,
		StopWhenFailed:       stopWhenFailed,
		Logger:               log,
		ClockWrapper:         clockWrapper,
		AppVersion:           Version,
		TestReporter:         reporter.NewReporter(reporter.ReporterType(reporterType)),
	}

	ts, err := testcase.NewTestsuite(testPath, options)
	if err != nil {
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

  cotton [-d] [-c] [-p] [-i] [-b <basedir>] [-r <reporttype>] <testpath|testdir>
  cotton -v
  cotton --help

`)
	flag.PrintDefaults()
}
