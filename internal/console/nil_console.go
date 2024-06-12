package console

type NilConsole struct{}

func NewNilConsole() Console {
	return &NilConsole{}
}

func (c *NilConsole) Print(args ...any) (int, error) {
	return 0, nil
}

func (c *NilConsole) Println(args ...any) (int, error) {
	return 0, nil
}

func (c *NilConsole) Printf(format string, args ...any) (int, error) {
	return 0, nil
}

func (c *NilConsole) Printfln(format string, args ...any) (int, error) {
	return 0, nil
}
