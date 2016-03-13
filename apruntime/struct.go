package apruntime

type InterpretedType string

type InterpretedStruct struct {
	// This is the concrete type of this struct instance.
	TypeName string
	Values map[string]interface{}
}

func (is *InterpretedStruct) Copy() *InterpretedStruct {
	newValues := make(map[string]interface{})
	for key, value := range is.Values {
		newValues[key] = value
	}
	return &InterpretedStruct{
		is.TypeName,
		newValues,
	}
}