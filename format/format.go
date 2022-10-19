package format

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
)

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, args ...any) string {
	carriers := make([]any, len(args))

	for i, arg := range args {
		if _, ok := arg.(fmt.Formatter); !ok {
			carriers[i] = carry(arg)
		}
	}

	return fmt.Sprintf(format, carriers...)
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, args ...any) (int, error) {
	carriers := make([]any, len(args))

	for i, arg := range args {
		if _, ok := arg.(fmt.Formatter); !ok {
			carriers[i] = carry(arg)
		}
	}

	return fmt.Fprintf(w, format, carriers...)
}

type carrier func() any

func (c carrier) Format(s fmt.State, r rune) {
	Format(s, r, c())
}

func carry(v any) carrier {
	return func() any {
		return v
	}
}

// Format formats the value according to the format specifier.
func Format(s fmt.State, r rune, value any) {
	hasPlus, hasSharp, r, converted := configureFormatValue(s, r, value)

	switch r {
	case 's':
		formatValueWithoutType(s, hasPlus, hasSharp, converted)

	case 'v':
		formatValueWithType(s, hasPlus, hasSharp, value, converted)

	case 'q':
		formatString(s, hasSharp, value, converted.(string))
	}
}

func configureFormatValue(s fmt.State, r rune, value any) (hasPlus bool, hasSharp bool, convertedRune rune, convertedValue any) { //nolint: nonamedreturns
	hasPlus, hasSharp = stateFlags(s)
	convertedRune = r
	convertedValue = value

	switch val := convertedValue.(type) {
	case json.RawMessage:
		convertedValue = string(val)
		hasSharp = false

		if convertedRune == 'q' {
			convertedRune = 's'
		}

	case *regexp.Regexp:
		convertedValue = val.String()

	case string:

	default:
		if convertedRune == 'q' {
			convertedRune = 'v'
		}
	}

	return hasPlus, hasSharp, convertedRune, convertedValue
}

func formatValueWithoutType(w io.Writer, hasPlus, hasSharp bool, v any) {
	switch {
	case hasPlus:
		fprintf(w, "%+v", v)
	case hasSharp:
		fprintf(w, "%#v", v)
	default:
		fprintf(w, "%v", v)
	}
}

func formatValueWithType(w io.Writer, hasPlus, hasSharp bool, original any, converted any) {
	switch typeOf := reflect.TypeOf(converted); {
	case hasPlus:
		fprintf(w, "%T(%+v)", original, converted)
	case hasSharp && hasTypeInOutput(typeOf):
		fprintf(w, "%#v", converted)
	case hasSharp:
		fprintf(w, "%T(%#v)", original, converted)
	default:
		fprintf(w, "%T(%v)", original, converted)
	}
}

func formatString(w io.Writer, hasSharp bool, original any, converted string) {
	if hasSharp {
		fprintf(w, "%T(%#v)", original, converted)
	} else {
		fprintf(w, "%q", converted)
	}
}

func stateFlags(s fmt.State) (hasPlus bool, hasSharp bool) { //nolint: nonamedreturns
	return s.Flag('+'), s.Flag('#')
}

func hasTypeInOutput(t reflect.Type) bool {
	switch t.Kind() { //nolint: exhaustive
	case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return true
	default:
		return false
	}
}

func fprintf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...) //nolint: errcheck
}
