package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chonla/cotton/cotton"
	"github.com/chonla/cotton/parser"
	"github.com/chonla/cotton/testsuite"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// VERSION of cotton
const VERSION = "0.4.0"

// Vars are injected variables from command line
type Vars []string

// watcher is file changes watcher
var watcher *fsnotify.Watcher

// Set to set vars, can be multiple times
func (v *Vars) Set(value string) error {
	*v = append(*v, value)
	return nil
}

func (v *Vars) String() string {
	return fmt.Sprint(*v)
}

func main() {
	// parser := parser.NewParser()
	var url string
	var help bool
	var ver bool
	var insecure bool
	var detail bool
	var watch bool
	var stopWhenFailed bool
	var vars Vars

	flag.Usage = usage

	flag.StringVar(&url, "u", "http://localhost:8080", "set base url")
	flag.BoolVar(&detail, "d", false, "detail mode -- to dump test detail")
	flag.BoolVar(&insecure, "i", false, "insecure mode -- to disable certificate verification")
	flag.BoolVar(&watch, "w", false, "watch mode -- to auto-rerun when files are changed")
	flag.BoolVar(&stopWhenFailed, "s", false, "panic mode -- to stop when failed")
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

	c, e := cotton.NewCotton(testpath, cotton.Config{
		BaseURL:        url,
		Insecure:       insecure,
		Verbose:        detail,
		Variables:      vars,
		StopWhenFailed: stopWhenFailed,
	})
	if e != nil {
		fmt.Printf("%s\n", e.Error())
		os.Exit(1)
	}
	c.SetParser(parser.NewParser())

	_, exitCode := c.Run()

	if watch {
		watcher, _ = fsnotify.NewWatcher()
		defer watcher.Close()

		// register all files/directories to be watched
		if err := filepath.Walk(testpath, watchDir); err != nil {
			fmt.Println("ERROR", err)
		}

		done := make(chan bool)

		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("\n%s\n", yellow("(Watching...)"))

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						fmt.Printf("\n%s\n\n", yellow("(File changes detected. Rerun tests.)"))
						_, exitCode = c.Run()
						fmt.Printf("\n%s\n", yellow("(Watching...)"))
					}
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()

		<-done
	} else {
		os.Exit(exitCode)
	}
}

func usage() {
	fmt.Println("Usage of cotton:")
	fmt.Println()
	fmt.Println("  cotton [-u <base-url>] [-i] [-d] [-s] [-p name1=value1] [-p name2=value2] ... <test-cases>")
	fmt.Println()
	fmt.Println("  test-cases can be a markdown file or a directory contain markdowns.")
	fmt.Println()
	flag.PrintDefaults()
}

func dispatchTests(ts testsuite.TestSuitesInterface, vars Vars, detail bool) int {
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
		ts.SetVariables(preVars)
	}

	ts.Run()
	exitCode := ts.Summary()
	return exitCode
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
