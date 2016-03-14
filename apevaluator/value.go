package apevaluator

import "fmt"

type Value interface {
	// Make an attempt to convert this to a native value, for example to
	// pass it to native code. Note that not all interpreted values are
	// possible to represent as native values.
	AsNative() interface{}
}

type NativeValue struct {
	val interface{}
}

func (nv *NativeValue) AsNative() interface{} {
	return nv.val
}

func (nv *NativeValue) String() string {
	return fmt.Sprint("NativeValue{", nv.val, "}")
}

type InterpretedStruct struct {
	// This is the concrete type of this struct instance.
	TypeName string
	Values map[string]Value
}

func (is *InterpretedStruct) Copy() *InterpretedStruct {
	newValues := make(map[string]Value)
	for key, value := range is.Values {
		newValues[key] = value
	}
	return &InterpretedStruct{
		is.TypeName,
		newValues,
	}
}