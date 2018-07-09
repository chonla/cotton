package main

import (
	"fmt"
	"os"

	"github.com/chonla/yas/parser"
	"github.com/kr/pretty"
)

func main() {
	parser := parser.NewParser()

	ts, e := parser.Parse("login-should-success-2.md")
	if e != nil {
		fmt.Printf("%s\n", e.Error())
		os.Exit(1)
	}
	pretty.Println(ts)
}
