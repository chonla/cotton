package logger

type NilLogger struct{}

func NewNilLogger() Logger {
	return &NilLogger{}
}

func (c *NilLogger) PrintTestCaseTitle(title string) error {
	return nil
}

func (c *NilLogger) PrintTestResult(passed bool) error {
	return nil
}

func (c *NilLogger) Print(args ...any) (int, error) {
	return 0, nil
}

func (c *NilLogger) Println(args ...any) (int, error) {
	return 0, nil
}

func (c *NilLogger) Printf(format string, args ...any) (int, error) {
	return 0, nil
}

func (c *NilLogger) Printfln(format string, args ...any) (int, error) {
	return 0, nil
}
