package variable

import "errors"

type Variables struct {
	names  []string
	values map[string]interface{}
}

func New() *Variables {
	return &Variables{
		values: map[string]interface{}{},
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

func (v *Variables) Set(name string, value interface{}) {
	if !v.nameExists(name) {
		v.names = append(v.names, name)
	}

	v.values[name] = value
}

func (v *Variables) Names() []string {
	return v.names
}

func (v *Variables) MergeWith(anotherVars *Variables) *Variables {
	newVars := New()
	for k, v := range v.values {
		newVars.Set(k, v)
	}
	anotherVarsNames := anotherVars.names
	for _, name := range anotherVarsNames {
		value, err := anotherVars.ValueOf(name)
		if err == nil {
			newVars.Set(name, value)
		}
	}
	return newVars
}
