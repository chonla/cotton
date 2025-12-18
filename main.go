package main

import (
	"os"

	"github.com/chonla/cotton/cli"
)

func main() {
	app := cli.NewApp()
	exitCode := app.Run(os.Args[1:])
	os.Exit(exitCode)
}
