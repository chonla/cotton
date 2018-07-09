package testsuite

import (
	"fmt"

	"github.com/chonla/console"
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
		fmt.Printf("Tests executed: ")
		console.Printfln("%d", ts.Stat.Total, console.ColorBold+console.ColorWhite)
		fmt.Printf("Tests passed: ")
		console.Printfln("%d (%s)", ts.Stat.Success, fmt.Sprintf("%0.2f%%", float64(ts.Stat.Success*100)/float64(ts.Stat.Total)), console.ColorBold+console.ColorGreen)
		fmt.Printf("Tests failed: ")
		console.Printfln("%d (%s)", ts.Stat.Total-ts.Stat.Success, fmt.Sprintf("%0.2f%%", float64((ts.Stat.Total-ts.Stat.Success)*100)/float64(ts.Stat.Total)), console.ColorBold+console.ColorRed)
		if ts.Stat.Total == ts.Stat.Success {
			return 0
		}
		return 1
	}
	fmt.Printf("No tests executed.")
	return 0
}
