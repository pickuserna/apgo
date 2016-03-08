// The apruntime package contains all base operations.
package apruntime

import (
	"reflect"
	"go/token"
)

type Context map[string]reflect.Value
type BuiltinFunc func (ctx Context, args []reflect.Value) reflect.Value

func add(ctx Context, args []reflect.Value) reflect.Value {
	return reflect.ValueOf(args[0].Int() + args[1].Int())
}

var BinaryOperators = map[token.Token]BuiltinFunc{
	token.ADD: add,
}