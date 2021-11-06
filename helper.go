package matcher

import "encoding/json"

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
