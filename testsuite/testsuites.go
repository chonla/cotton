package testsuite

import (
	"fmt"

	"github.com/fatih/color"
)

// TestSuites is several test suites
type TestSuites struct {
	Suites    []*TestSuite
	BaseURL   string
	Config    *Config
	Stat      TestStat
	Variables map[string]string
}

// Run executes test suite
func (ts *TestSuites) Run() {
	for _, suite := range ts.Suites {
		suite.BaseURL = ts.BaseURL
		suite.Config = ts.Config
		suite.Variables = ts.Variables
		suite.Run()
		ts.Stat.Total += suite.Stat.Total
		ts.Stat.Success += suite.Stat.Success
	}
}

// Summary prints test summary
func (ts *TestSuites) Summary() int {
	if ts.Stat.Total > 0 {
		magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

		fmt.Printf("%s\n", magenta("----"))
		fmt.Printf("Tests executed: ")
		color.White("%d", ts.Stat.Total)
		fmt.Printf("Tests passed: ")
		color.Green("%d (%0.2f%%)", ts.Stat.Success, float64(ts.Stat.Success*100)/float64(ts.Stat.Total))
		fmt.Printf("Tests failed: ")
		color.Red("%d (%0.2f%%)", ts.Stat.Total-ts.Stat.Success, float64((ts.Stat.Total-ts.Stat.Success)*100)/float64(ts.Stat.Total))
		if ts.Stat.Total == ts.Stat.Success {
			return 0
		}
		return 1
	}
	fmt.Printf("No tests executed.")
	return 0
}
