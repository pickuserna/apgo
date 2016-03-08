package main

import (
	"github.com/alangpierce/apgo/interpreter"
	"github.com/alangpierce/apgo/apruntime"
)

func main() {
	interp := interpreter.NewInterpreter()
	interp.LoadNativePackage(apruntime.FmtPackage)
	interp.LoadPackage("sample")
	interp.RunMain()
}