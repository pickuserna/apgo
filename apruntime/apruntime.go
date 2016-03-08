// The apruntime package contains all base operations.
package apruntime

import (
	"reflect"
	"go/token"
)

type Context map[string]reflect.Value
type BuiltinFunc func (ctx Context, args []reflect.Value) reflect.Value

func add(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() + reflect.ValueOf(y).Int()
}

var BinaryOperators = map[token.Token]reflect.Value{
	token.ADD: reflect.ValueOf(add),
}