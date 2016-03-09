// The apruntime package contains all base operations.
package apruntime

import (
	"reflect"
	"go/token"
	"fmt"
)

type NativePackage struct {
	Name string
	Funcs map[string]interface{}
	Globals map[string]*interface{}
}

func add(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() + reflect.ValueOf(y).Int()
}

func greater(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() > reflect.ValueOf(y).Int()
}


var BinaryOperators = map[token.Token]reflect.Value{
	token.ADD: reflect.ValueOf(add),
	token.GTR: reflect.ValueOf(greater),
}

var AssignBinaryOperators = map[token.Token]reflect.Value{
	token.ADD_ASSIGN: reflect.ValueOf(add),
}

var FmtPackage = &NativePackage{
	Name: "fmt",
	Funcs: map[string]interface{} {
		"Print": fmt.Print,
		"Println": fmt.Println,
	},
	Globals: map[string]*interface{} {},
}