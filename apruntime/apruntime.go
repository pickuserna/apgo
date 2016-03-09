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
	sum := reflect.ValueOf(x).Int() + reflect.ValueOf(y).Int()
	// Since this is a well-formed operation, the two types must be the
	// same, so convert to that type.
	// TODO: Handle other types, like floats.
	return reflect.ValueOf(sum).Convert(reflect.TypeOf(x)).Interface()
}

func sub(x interface{}, y interface{}) interface{} {
	sum := reflect.ValueOf(x).Int() - reflect.ValueOf(y).Int()
	// TODO: Handle other types.
	return reflect.ValueOf(sum).Convert(reflect.TypeOf(x)).Interface()
}

func mul(x interface{}, y interface{}) interface{} {
	sum := reflect.ValueOf(x).Int() * reflect.ValueOf(y).Int()
	// TODO: Handle other types.
	return reflect.ValueOf(sum).Convert(reflect.TypeOf(x)).Interface()
}

func quo(x interface{}, y interface{}) interface{} {
	sum := reflect.ValueOf(x).Int() / reflect.ValueOf(y).Int()
	// TODO: Handle other types.
	return reflect.ValueOf(sum).Convert(reflect.TypeOf(x)).Interface()
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

func neq(x interface{}, y interface{}) interface{} {
	return x != y
}

func leq(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() <= reflect.ValueOf(y).Int()
}

func geq(x interface{}, y interface{}) interface{} {
	// TODO: Handle other types.
	return reflect.ValueOf(x).Int() >= reflect.ValueOf(y).Int()
}


var BinaryOperators = map[token.Token]reflect.Value{
	token.ADD: reflect.ValueOf(add),
	token.SUB: reflect.ValueOf(sub),
	token.MUL: reflect.ValueOf(mul),
	token.QUO: reflect.ValueOf(quo),
	token.LSS: reflect.ValueOf(less),
	token.GTR: reflect.ValueOf(greater),
	token.LOR: reflect.ValueOf(lor),
	token.EQL: reflect.ValueOf(equal),
	token.NEQ: reflect.ValueOf(neq),
	token.LEQ: reflect.ValueOf(leq),
	token.GEQ: reflect.ValueOf(geq),
}

var AssignBinaryOperators = map[token.Token]reflect.Value{
	token.ADD_ASSIGN: reflect.ValueOf(add),
	token.SUB_ASSIGN: reflect.ValueOf(sub),
	token.MUL_ASSIGN: reflect.ValueOf(mul),
	token.QUO_ASSIGN: reflect.ValueOf(quo),
}

var IncDecOperators = map[token.Token]reflect.Value{
	token.INC: reflect.ValueOf(add),
	token.DEC: reflect.ValueOf(sub),
}

var FmtPackage = &NativePackage{
	Name: "fmt",
	Funcs: map[string]interface{} {
		"Print": fmt.Print,
		"Println": fmt.Println,
		"Sprint": fmt.Sprint,
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