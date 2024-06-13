package kindof

import "reflect"

func Float(v interface{}) (float64, bool) {
	aType := reflect.TypeOf(v)
	isFloat := false
	var vFloat float64
	switch aType.Kind() {
	case reflect.Float32:
		vFloat = float64(v.(float32))
		isFloat = true
	case reflect.Float64:
		vFloat = float64(v.(float64))
		isFloat = true
	}
	return vFloat, isFloat
}
