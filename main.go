package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chonla/yas/parser"
)

func main() {
	parser := parser.NewParser()
	var url string
	var help bool

	flag.Usage = usage

	flag.StringVar(&url, "u", "http://localhost:8080/", "set base url")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.Parse()

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

	ts.Run()
}

func usage() {
	fmt.Println("Usage of yas:")
	fmt.Println()
	fmt.Println("  yas [-u <base-url>] <test-cases>")
	fmt.Println()
	fmt.Println("  test-cases can be a markdown file or a directory contain markdowns.")
	fmt.Println()
	flag.PrintDefaults()
}
