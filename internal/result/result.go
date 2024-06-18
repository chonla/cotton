package result

type TestResult struct {
	Title      string
	Passed     bool
	Assertions []AssertionResult
	Error      error
}

type AssertionResult struct {
	Title    string
	Passed   bool
	Expected string
	Actual   string
	Error    error
}
