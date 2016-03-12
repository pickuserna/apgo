package apevaluator

import (
	"github.com/alangpierce/apgo/apast"
)

type Context struct {
	Locals map[string]interface{}
	PackageValues map[string]interface{}
	// Slice of return values, or nil if the function hasn't returned yet.
	// This is used both for the values themselves and to communicate
	// control flow. For example, a function returning nothing should have
	// returnValues set to the empty slice upon returning, which signals to
	// other code that we want to finish the function now.
	returnValues []interface{}
	shouldBreak bool
}

func NewContext(pack *apast.Package) *Context {
	packageValues := make(map[string]interface{})
	for name, funcAst := range pack.Funcs {
		// Eager-bind the funcAst value so it points to the right loop
		// variable.
		packageValues[name] = func(funcAst *apast.FuncDecl) interface{} {
			return func(args ...interface{}) interface{} {
				result := EvaluateFunc(pack, funcAst, args...)
				if result == nil {
					return nil
				} else {
					return result[0]
				}
			}
		}(funcAst)
	}
	return &Context{
		Locals: make(map[string]interface{}),
		PackageValues: packageValues,
	}
}

func (ctx *Context) resolveValue(name string) ExprResult {
	if _, ok := ctx.Locals[name]; ok {
		return &VariableLValue{
			ctx.Locals,
			name,
		}
	} else if _, ok := ctx.PackageValues[name]; ok {
		return &VariableLValue{
			ctx.PackageValues,
			name,
		}
	} else {
		// If we didn't find anything, then create it as a local
		// variable.
		// TODO: Maybe we need to init to a zero value?
		ctx.Locals[name] = nil
		return &VariableLValue{
			ctx.Locals,
			name,
		}
	}
}

func (ctx *Context) isNameValid(name string) bool {
	if _, ok := ctx.Locals[name]; ok {
		return true
	} else if _, ok := ctx.PackageValues[name]; ok {
		return true
	}
	return false
}

func (ctx *Context) assignValue(name string, value interface{}) {
	ctx.Locals[name] = value
}