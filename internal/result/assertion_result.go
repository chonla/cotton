package result

type AssertionResult struct {
	Title    string
	Passed   bool
	Expected string
	Actual   string
	Error    error
}
