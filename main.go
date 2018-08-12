package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chonla/cotton/parser"
	"github.com/chonla/cotton/testsuite"
	"github.com/fatih/color"
)

// VERSION of cotton
const VERSION = "0.1.28"

// Vars are injected variables from command line
type Vars []string

// Set to set vars, can be multiple times
func (v *Vars) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func (v *Vars) String() string {
	return fmt.Sprint(*v)
}

func main() {
	parser := parser.NewParser()
	var url string
	var help bool
	var ver bool
	var insecure bool
	var detail bool
	var vars Vars

	flag.Usage = usage

	flag.StringVar(&url, "u", "http://localhost:8080", "set base url")
	flag.BoolVar(&detail, "d", false, "detail mode -- to dump test detail")
	flag.BoolVar(&insecure, "i", false, "insecure mode -- to disable certificate verification")
	flag.BoolVar(&ver, "v", false, "show cotton version")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.Var(&vars, "p", "to inject predefined in variable-name=variable-value format")
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

	if len(vars) > 0 {
		preVars := map[string]string{}
		for _, v := range vars {
			s := strings.SplitN(v, "=", 2)
			if len(s) == 2 {
				preVars[s[0]] = s[1]
			}
		}

		if detail && len(preVars) > 0 {
			blue := color.New(color.FgBlue).SprintFunc()

			fmt.Printf("Injected variables:\n")
			for k := range preVars {
				fmt.Printf("* %s\n", blue(k))
			}
		}
		ts.Variables = preVars
	}

	ts.Run()
	exitCode := ts.Summary()
	os.Exit(exitCode)
}

func usage() {
	fmt.Println("Usage of cotton:")
	fmt.Println()
	fmt.Println("  cotton [-u <base-url>] [-i] [-d] [-p name1=value1] [-p name2=value2] ... <test-cases>")
	fmt.Println()
	fmt.Println("  test-cases can be a markdown file or a directory contain markdowns.")
	fmt.Println()
	flag.PrintDefaults()
}
