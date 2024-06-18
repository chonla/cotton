package logger

type Logger interface {
	PrintTestCaseTitle(title string) error
	PrintTestResult(passed bool) error

	Print(args ...any) (int, error)
	Println(args ...any) (int, error)
	Printf(format string, args ...any) (int, error)
	Printfln(format string, args ...any) (int, error)
}
