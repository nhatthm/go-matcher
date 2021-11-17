package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"
)

// Matcher determines if the actual matches the expectation.
type Matcher interface {
	Match(actual interface{}) (bool, error)
	Expected() string
}

var _ Matcher = (*ExactMatcher)(nil)

// ExactMatcher matches by exact string.
type ExactMatcher struct {
	expected interface{}
}

// Expected returns the expectation.
func (m ExactMatcher) Expected() string {
	if v := strVal(m.expected); v != nil {
		return *v
	}

	return fmt.Sprintf("%+v", m.expected)
}

// Match determines if the actual is expected.
func (m ExactMatcher) Match(actual interface{}) (bool, error) {
	return assert.ObjectsAreEqual(m.expected, actual), nil
}

var _ Matcher = (*JSONMatcher)(nil)

// JSONMatcher matches by json with <ignore-diff> support.
type JSONMatcher struct {
	expected string
}

// Expected returns the expectation.
func (m JSONMatcher) Expected() string {
	return m.expected
}

// Match determines if the actual is expected.
func (m JSONMatcher) Match(actual interface{}) (bool, error) {
	actualBytes, err := jsonVal(actual)
	if err != nil {
		return false, err
	}

	return assertjson.FailNotEqual([]byte(m.expected), actualBytes) == nil, nil
}

var _ Matcher = (*RegexMatcher)(nil)

// RegexMatcher matches by regex.
type RegexMatcher struct {
	regexp *regexp.Regexp
}

// Expected returns the expectation.
func (m RegexMatcher) Expected() string {
	return m.regexp.String()
}

// Match determines if the actual is expected.
func (m RegexMatcher) Match(actual interface{}) (bool, error) {
	if v := strVal(actual); v != nil {
		return m.regexp.MatchString(*v), nil
	}

	return false, nil
}

var _ Matcher = (*LenMatcher)(nil)

// LenMatcher matches by the length of the value.
type LenMatcher struct {
	expected int
}

// Match determines if the actual is expected.
func (m LenMatcher) Match(actual interface{}) (_ bool, err error) {
	if actual == nil {
		return false, nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(recovered(r)) // nolint: goerr113
		}
	}()

	val := reflect.ValueOf(actual)

	if val.Type().Kind() == reflect.Ptr {
		return m.Match(val.Elem().Interface())
	}

	return val.Len() == m.expected, nil
}

// Expected returns the expectation.
func (m LenMatcher) Expected() string {
	return fmt.Sprintf("len is %d", m.expected)
}

var _ Matcher = (*EmptyMatcher)(nil)

// EmptyMatcher checks whether the value is empty.
type EmptyMatcher struct{}

// Match determines if the actual is expected.
func (EmptyMatcher) Match(actual interface{}) (bool, error) {
	return isEmpty(actual), nil
}

// Expected returns the expectation.
func (EmptyMatcher) Expected() string {
	return "is empty"
}

var _ Matcher = (*NotEmptyMatcher)(nil)

// NotEmptyMatcher checks whether the value is not empty.
type NotEmptyMatcher struct{}

// Match determines if the actual is expected.
func (NotEmptyMatcher) Match(actual interface{}) (bool, error) {
	return !isEmpty(actual), nil
}

// Expected returns the expectation.
func (NotEmptyMatcher) Expected() string {
	return "is not empty"
}

var _ Matcher = (*Callback)(nil)

// Callback matches by calling a function.
type Callback func() Matcher

// Expected returns the expectation.
func (m Callback) Expected() string {
	return m().Expected()
}

// Match determines if the actual is expected.
func (m Callback) Match(actual interface{}) (bool, error) {
	return m().Match(actual)
}

// Matcher returns the matcher.
func (m Callback) Matcher() Matcher {
	return m()
}

// Exact matches two objects by their exact values.
func Exact(expected interface{}) ExactMatcher {
	return ExactMatcher{expected: expected}
}

// Exactf matches two strings by the formatted expectation.
func Exactf(expected string, args ...interface{}) ExactMatcher {
	return ExactMatcher{expected: fmt.Sprintf(expected, args...)}
}

// JSON matches two json strings with <ignore-diff> support.
func JSON(expected interface{}) JSONMatcher {
	ex, err := jsonVal(expected)
	if err != nil {
		panic(err)
	}

	return JSONMatcher{expected: string(ex)}
}

// RegexPattern matches two strings by using regex.
func RegexPattern(pattern string) RegexMatcher {
	return RegexMatcher{regexp: regexp.MustCompile(pattern)}
}

// Regex matches two strings by using regex.
func Regex(regexp *regexp.Regexp) RegexMatcher {
	return RegexMatcher{regexp: regexp}
}

// Len matches by the length of the value.
func Len(expected int) LenMatcher {
	return LenMatcher{expected: expected}
}

// IsEmpty checks whether the value is empty.
func IsEmpty() EmptyMatcher {
	return EmptyMatcher{}
}

// IsNotEmpty checks whether the value is not empty.
func IsNotEmpty() NotEmptyMatcher {
	return NotEmptyMatcher{}
}

func match(v interface{}) Matcher {
	switch val := v.(type) {
	case Matcher:
		return val

	case func() Matcher:
		return Callback(val)

	case *regexp.Regexp:
		return Regex(val)

	case fmt.Stringer:
		return Exact(val.String())
	}

	return Exact(v)
}

// Match returns a matcher according to its type.
func Match(v interface{}) Matcher {
	return match(v)
}
