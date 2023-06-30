package utilities

import "reflect"

func StructIsEmpty(data interface{}) bool {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Struct {
		return true
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.String, reflect.Slice:
			if field.Len() > 0 {
				return false
			}
		case reflect.Struct:
			if !StructIsEmpty(field.Interface()) {
				return false
			}
		}
	}

	return true
}
