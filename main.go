package main

import (
	"lim-lang/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
