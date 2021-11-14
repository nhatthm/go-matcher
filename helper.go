package matcher

import (
	"encoding/json"
	"reflect"
)

func strVal(v interface{}) *string {
	switch v := v.(type) {
	case string:
		return &v

	case []byte:
		return strPtr(string(v))
	}

	return nil
}

func strPtr(v string) *string {
	return &v
}

func jsonVal(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil

	case []byte:
		return v, nil
	}

	return json.Marshal(v)
}

// isEmpty gets whether the specified object is considered empty or not.
// nolint: exhaustive
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() == 0

	case reflect.Ptr:
		if val.IsNil() {
			return true
		}

		return isEmpty(val.Elem().Interface())
	}

	zero := reflect.Zero(val.Type())

	return reflect.DeepEqual(v, zero.Interface())
}
