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

func (v *Variables) Set(name string, value interface{}) {
	v.values[name] = value
	keys := []string{}
	for k := range v.values {
		keys = append(keys, k)
	}
	v.names = keys
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
