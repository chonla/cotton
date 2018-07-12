package testsuite

// TestSuite holds a test suite
type TestSuite struct {
	Name      string
	BaseURL   string
	TestCases []*TestCase
	Stat      TestStat
}

// Run executes test case
func (ts *TestSuite) Run() {
	for _, tc := range ts.TestCases {
		if len(tc.Expectations) > 0 {
			tc.BaseURL = ts.BaseURL
			ts.Stat.Total++
			e := tc.Run()
			if e == nil {
				ts.Stat.Success++
			}
		}
	}
}
