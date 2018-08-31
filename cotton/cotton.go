package cotton

import (
	"fmt"
	"os"
	"strings"

	"github.com/chonla/cotton/testsuite"

	"github.com/chonla/cotton/parser"
	"github.com/fatih/color"
)

var statFile = os.Stat

// Cotton represents cotton.
type Cotton struct {
	path   string
	parser parser.Interface
	Config
}

// Config holds cotton configuration
type Config struct {
	BaseURL   string
	Insecure  bool
	Verbose   bool
	Variables []string
}

// NewCotton creates a new cotton.
// if path does not exist, error is returned.
func NewCotton(path string, conf Config) (*Cotton, error) {
	if _, err := statFile(path); err != nil {
		return nil, err
	}

	return &Cotton{
		path:   path,
		parser: nil,
		Config: conf,
	}, nil
}

// SetParser to set markdown parser
func (c *Cotton) SetParser(p parser.Interface) {
	c.parser = p
}

// Run executes testsuites and return exit code
func (c *Cotton) Run() (testsuite.TestStat, int) {
	ts, e := c.parser.Parse(c.path)
	if e != nil {
		fmt.Printf("%s\n", e.Error())
		return testsuite.TestStat{}, 1
	}

	if len(c.Variables) > 0 {
		preVars := map[string]string{}
		for _, v := range c.Variables {
			s := strings.SplitN(v, "=", 2)
			if len(s) == 2 {
				preVars[s[0]] = s[1]
			}
		}

		if c.Verbose && len(preVars) > 0 {
			blue := color.New(color.FgBlue).SprintFunc()

			fmt.Printf("Injected variables:\n")
			for k := range preVars {
				fmt.Printf("* %s\n", blue(k))
			}
		}
		ts.SetVariables(preVars)
	}

	ts.SetBaseURL(c.BaseURL)
	ts.SetConfig(&testsuite.Config{
		Insecure: c.Insecure,
		Detail:   c.Verbose,
	})

	ts.Run()
	exitCode := ts.Summary()
	return ts.Stat(), exitCode
}
