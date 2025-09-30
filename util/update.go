package util

import "reflect"

func ApplyUpdates[T any](original *T, updates T) {
	// Use reflection to iterate over the fields of the struct
	origValue := reflect.ValueOf(original).Elem()
	updateValue := reflect.ValueOf(updates)

	for i := 0; i < origValue.NumField(); i++ {
		field := origValue.Field(i)
		updateField := updateValue.Field(i)

		// Check if the update field is non-zero (not the zero value for its type)
		if !isZeroValue(updateField) {
			field.Set(updateField)
		}
	}
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Struct:
		// For structs, check if all fields are zero values
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// For other types, consider them non-zero by default
		return false
	}
}
