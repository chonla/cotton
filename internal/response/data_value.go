package response

import "fmt"

type DataValue struct {
	Value       interface{}
	TypeName    string
	IsUndefined bool
}

func (d *DataValue) String() string {
	if d.IsUndefined {
		return "undefined"
	}
	return fmt.Sprintf("%v", d.Value)
}
