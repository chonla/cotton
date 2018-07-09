package testsuite

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
		ts.Stat.Total++
		e := tc.Run()
		if e == nil {
			ts.Stat.Success++
		}
	}
}
