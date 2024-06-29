package variable

import (
	"errors"
	"fmt"
)

type Variables struct {
	names        []string
	values       map[string]interface{}
	stringValues map[string]string
}

func New() *Variables {
	return &Variables{
		names:        []string{},
		values:       map[string]interface{}{},
		stringValues: map[string]string{},
	}
}

func (v *Variables) ValueOf(name string) (interface{}, error) {
	if val, ok := v.values[name]; ok {
		return val, nil
	}
	return nil, errors.New("value not found")
}

func (v *Variables) nameExists(name string) bool {
	for _, existingName := range v.names {
		if existingName == name {
			return true
		}
	}
	return false
}

func (v *Variables) Add(variable *Variable) {
	v.Set(variable.Name, variable.Value)
}

func (v *Variables) Set(name string, value interface{}) {
	if !v.nameExists(name) {
		v.names = append(v.names, name)
	}

	v.values[name] = value
	v.stringValues[name] = fmt.Sprintf("%v", value)
}

func (v *Variables) Names() []string {
	return v.names
}

func (v *Variables) MergeWith(anotherVars *Variables) *Variables {
	newVars := v.Clone()
	anotherVarsNames := anotherVars.names
	for _, name := range anotherVarsNames {
		value, err := anotherVars.ValueOf(name)
		if err == nil {
			newVars.Set(name, value)
		}
	}
	return newVars
}

func (v *Variables) Clone() *Variables {
	newVars := New()
	for _, name := range v.names {
		value, err := v.ValueOf(name)
		if err == nil {
			newVars.Set(name, value)
		}
	}
	return newVars
}

func (v *Variables) ToStringMap() map[string]interface{} {
	m := map[string]interface{}{}

	for k, v := range v.stringValues {
		m[k] = v
	}
	return m
}
