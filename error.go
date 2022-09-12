package matcher

import "fmt"

func recovered(v any) string {
	switch v := v.(type) {
	case error:
		return v.Error()

	case string:
		return v
	}

	return fmt.Sprintf("%+v", v)
}
