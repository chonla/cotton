package console

type Console interface {
	Print(args ...any) (int, error)
	Println(args ...any) (int, error)
	Printf(format string, args ...any) (int, error)
	Printfln(format string, args ...any) (int, error)
}
