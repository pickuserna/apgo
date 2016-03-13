package apruntime

type InterpretedType string

type InterpretedStruct struct {
	// This is the concrete type of this struct instance.
	TypeName string
	Values map[string]interface{}
}
