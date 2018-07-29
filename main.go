package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chonla/cotton/parser"
	"github.com/chonla/cotton/testsuite"
)

// VERSION of cotton
const VERSION = "0.1.22"

func main() {
	parser := parser.NewParser()
	var url string
	var help bool
	var ver bool
	var insecure bool
	var detail bool

	flag.Usage = usage

	flag.StringVar(&url, "u", "http://localhost:8080", "set base url")
	flag.BoolVar(&detail, "d", false, "detail mode -- to dump test detail")
	flag.BoolVar(&insecure, "i", false, "insecure mode -- to disable certificate verification")
	flag.BoolVar(&ver, "v", false, "show cotton version")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.Parse()

	if ver {
		fmt.Printf("cotton %s\n", VERSION)
		os.Exit(0)
	}

	testpath := flag.Arg(0)
	if testpath == "" {
		flag.Usage()
		os.Exit(1)
	}

	ts, e := parser.Parse(testpath)
	if e != nil {
		fmt.Printf("%s\n", e.Error())
		os.Exit(1)
	}
	ts.BaseURL = url
	ts.Config = &testsuite.Config{
		Insecure: insecure,
		Detail:   detail,
	}

	ts.Run()
	exitCode := ts.Summary()
	os.Exit(exitCode)
}

func usage() {
	fmt.Println("Usage of cotton:")
	fmt.Println()
	fmt.Println("  cotton [-u <base-url>] <test-cases>")
	fmt.Println()
	fmt.Println("  test-cases can be a markdown file or a directory contain markdowns.")
	fmt.Println()
	flag.PrintDefaults()
}
