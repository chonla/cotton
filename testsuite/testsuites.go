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
	stat      TestStat
	Variables map[string]string
}

// TestSuitesInterface is interface of TestSuites
type TestSuitesInterface interface {
	Run()
	Summary() int
	SetVariables(map[string]string)
	SetBaseURL(string)
	SetConfig(*Config)
	Stat() TestStat
}

// SetVariables to set variables to test suites
func (ts *TestSuites) SetVariables(v map[string]string) {
	ts.Variables = v
}

// SetBaseURL to set base url
func (ts *TestSuites) SetBaseURL(url string) {
	ts.BaseURL = url
}

// SetConfig to set configuration
func (ts *TestSuites) SetConfig(conf *Config) {
	ts.Config = &Config{
		Insecure:       conf.Insecure,
		Detail:         conf.Detail,
		StopWhenFailed: conf.StopWhenFailed,
	}
}

// Stat returns test stat
func (ts *TestSuites) Stat() TestStat {
	return ts.stat
}

// Run executes test suite
func (ts *TestSuites) Run() {
	ts.stat = TestStat{
		Total:   0,
		Success: 0,
	}
	for _, suite := range ts.Suites {
		suite.Stat = TestStat{
			Total:   0,
			Success: 0,
		}
		suite.BaseURL = ts.BaseURL
		suite.Config = ts.Config
		suite.Variables = ts.Variables
		e := suite.Run()
		ts.stat.Total += suite.Stat.Total
		ts.stat.Success += suite.Stat.Success

		if e != nil && ts.Config.StopWhenFailed {
			return
		}
	}
}

// Summary prints test summary
func (ts *TestSuites) Summary() int {
	if ts.stat.Total > 0 {
		magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

		fmt.Printf("%s\n", magenta("----"))
		fmt.Printf("Tests executed: ")
		color.White("%d", ts.stat.Total)
		fmt.Printf("Tests passed: ")
		color.Green("%d (%0.2f%%)", ts.stat.Success, float64(ts.stat.Success*100)/float64(ts.stat.Total))
		fmt.Printf("Tests failed: ")
		color.Red("%d (%0.2f%%)", ts.stat.Total-ts.stat.Success, float64((ts.stat.Total-ts.stat.Success)*100)/float64(ts.stat.Total))
		if ts.stat.Total == ts.stat.Success {
			return 0
		}
		return 1
	}
	fmt.Printf("No tests executed.")
	return 0
}
