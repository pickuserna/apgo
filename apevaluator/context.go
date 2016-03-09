package apevaluator

import (
	"reflect"
	"github.com/alangpierce/apgo/apast"
	"fmt"
)

type Context struct {
	Locals map[string]reflect.Value
	PackageValues map[string]reflect.Value
	// Slice of return values, or nil if the function hasn't returned yet.
	// This is used both for the values themselves and to communicate
	// control flow. For example, a function returning nothing should have
	// returnValues set to the empty slice upon returning, which signals to
	// other code that we want to finish the function now.
	returnValues []reflect.Value
}

func NewContext(pack *apast.Package) *Context {
	packageValues := make(map[string]reflect.Value)
	for name, funcAst := range pack.Funcs {
		// Eager-bind the funcAst value so it points to the right loop
		// variable.
		packageValues[name] = func(funcAst *apast.FuncDecl) reflect.Value {
			f := func(args ...interface{}) interface{} {
				result := EvaluateFunc(pack, funcAst, args...)
				if result == nil {
					return nil
				} else {
					return result[0].Interface()
				}
			}
			return reflect.ValueOf(f)
		}(funcAst)
	}
	return &Context{
		Locals: make(map[string]reflect.Value),
		PackageValues: packageValues,
	}
}

func (ctx *Context) resolveValue(name string) reflect.Value {
	if local, ok := ctx.Locals[name]; ok {
		return local
	} else if packageVal, ok := ctx.PackageValues[name]; ok {
		return packageVal
	} else {
		panic(fmt.Sprint("Variable not found: ", name))
	}
}

func (ctx *Context) assignValue(name string, value reflect.Value) {
	ctx.Locals[name] = value
}