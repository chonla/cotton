package kindof

import "reflect"

func Int(v interface{}) (int64, bool) {
	aType := reflect.TypeOf(v)
	isInt := false
	var vInt int64
	switch aType.Kind() {
	case reflect.Int:
		vInt = int64(v.(int))
		isInt = true
	case reflect.Int8:
		vInt = int64(v.(int8))
		isInt = true
	case reflect.Int16:
		vInt = int64(v.(int16))
		isInt = true
	case reflect.Int32:
		vInt = int64(v.(int32))
		isInt = true
	case reflect.Int64:
		vInt = int64(v.(int64))
		isInt = true
	}
	return vInt, isInt
}
