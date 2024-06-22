package value

import (
	"fmt"
	"reflect"
)

type Value struct {
	value       interface{}
	valueType   string
	isUndefined bool
}

func New(value interface{}) *Value {
	if value == nil {
		return &Value{
			value:       nil,
			valueType:   "unknown",
			isUndefined: false,
		}
	}
	typeName := reflect.TypeOf(value).Name()
	return &Value{
		value:       value,
		valueType:   typeName,
		isUndefined: false,
	}
}

func NewUndefined() *Value {
	return &Value{
		isUndefined: true,
	}
}

func (v *Value) IsNil() bool {
	return !v.isUndefined && v.value == nil
}

func (v *Value) IsUndefined() bool {
	return v.isUndefined
}

func (v *Value) Value() interface{} {
	return v.value
}

func (v *Value) Type() string {
	return v.valueType
}

func (v *Value) String() string {
	if v.isUndefined {
		return "undefined"
	}
	if v.value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v.value)
}
