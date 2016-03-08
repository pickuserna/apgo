package main

import (
	"github.com/alangpierce/apgo/interpreter"
)

func main() {
	interp := interpreter.NewInterpreter()
	interp.LoadPackage("sample")
	interp.RunMain()
}