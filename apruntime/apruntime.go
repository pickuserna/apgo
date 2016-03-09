// The apruntime package contains all base operations.
package apruntime

import (
	"reflect"
	"go/token"
	"fmt"
	"time"
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

func sub(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() - reflect.ValueOf(y).Int()
}

func less(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() < reflect.ValueOf(y).Int()
}

func greater(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() > reflect.ValueOf(y).Int()
}

func lor(x interface{}, y interface{}) interface{} {
	// TODO: Short-circuit.
	return x.(bool) || y.(bool)
}

func equal(x interface{}, y interface{}) interface{} {
	return x == y
}


var BinaryOperators = map[token.Token]reflect.Value{
	token.ADD: reflect.ValueOf(add),
	token.SUB: reflect.ValueOf(sub),
	token.LSS: reflect.ValueOf(less),
	token.GTR: reflect.ValueOf(greater),
	token.LOR: reflect.ValueOf(lor),
	token.EQL: reflect.ValueOf(equal),
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

var TimePackage = &NativePackage{
	Name: "time",
	Funcs: map[string]interface{} {
		"Now": time.Now,
		"Since": time.Since,
	},
	Globals: map[string]*interface{} {},
}