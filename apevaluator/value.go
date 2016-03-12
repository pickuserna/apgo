package apevaluator

import (
	"reflect"
	"fmt"
)

// ExprResult is what you get when evaluating an expression. It is a little more
// generic than just a value because sometimes it can implicitly be assignable
// and/or have a pointer associated with it.
type ExprResult interface {
	get() interface{}
	set(val interface{})
}

type RValue struct {
	val interface{}
}

func (rv *RValue) get() interface{} {
	return rv.val
}

func (rv *RValue) set(val interface{}) {
	panic(fmt.Sprint("Called set on RValue ", rv.val))
}

type VariableLValue struct {
	varMap map[string]interface{}
	name string
}

func (lv *VariableLValue) get() interface{} {
	return lv.varMap[lv.name]
}

func (lv *VariableLValue) set(val interface{}) {
	lv.varMap[lv.name] = val
}

type ReflectValLValue struct {
	val reflect.Value
}

func (lv *ReflectValLValue) get() interface{} {
	return lv.val.Interface()
}

func (lv *ReflectValLValue) set(val interface{}) {
	lv.val.Set(reflect.ValueOf(val))
}