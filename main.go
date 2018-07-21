package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chonla/cotton/parser"
)

// VERSION of cotton
const VERSION = "0.1.7"

func main() {
	parser := parser.NewParser()
	var url string
	var help bool
	var ver bool
	var insecure bool

	flag.Usage = usage

	flag.StringVar(&url, "u", "http://localhost:8080", "set base url")
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
	ts.Insecure = insecure

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
