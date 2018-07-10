package testsuite

import (
	"fmt"

	"github.com/fatih/color"
)

// TestSuite holds a test suite
type TestSuite struct {
	Name      string
	BaseURL   string
	TestCases []*TestCase
	Stat      TestStat
}

// TestStat contains test result
type TestStat struct {
	Total   int
	Success int
}

// Run executes test case
func (ts *TestSuite) Run() {
	for _, tc := range ts.TestCases {
		tc.BaseURL = ts.BaseURL
		ts.Stat.Total++
		e := tc.Run()
		if e == nil {
			ts.Stat.Success++
		}
	}
}

// Summary prints test summary
func (ts *TestSuite) Summary() int {
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
