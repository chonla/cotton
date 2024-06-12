package console

import (
	"cotton/internal/line"

	"fmt"
)

type Terminal struct{}

func NewTerminal() Console {
	return &Terminal{}
}

func (c *Terminal) Print(args ...any) (int, error) {
	return fmt.Print(args...)
}

func (c *Terminal) Println(args ...any) (int, error) {
	return fmt.Println(args...)
}

func (c *Terminal) Printf(format string, args ...any) (int, error) {
	return fmt.Printf(format, args...)
}

func (c *Terminal) Printfln(format string, args ...any) (int, error) {
	return fmt.Printf(fmt.Sprintf("%s%s", format, line.DetectLineSeparator()), args...)
}
