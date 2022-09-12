package matcher

import (
	"encoding/json"
	"reflect"
	"regexp"
)

func strVal(v any) *string {
	switch v := v.(type) {
	case string:
		return &v

	case []byte:
		return ptr(string(v))
	}

	return nil
}

func jsonVal(v any) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil

	case []byte:
		return v, nil
	}

	return json.Marshal(v)
}

func regexpVal(v any) *regexp.Regexp {
	switch v := v.(type) {
	case *regexp.Regexp:
		return v

	case regexp.Regexp:
		return &v

	case string:
		return regexp.MustCompile(v)
	}

	return nil
}

// isEmpty gets whether the specified object is considered empty or not.
// nolint: exhaustive
func isEmpty(v any) bool {
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

func ptr[T any](v T) *T {
	return &v
}
