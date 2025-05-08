package matcher

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/swaggest/assertjson"

	"go.nhat.io/matcher/v3/format"
)

// Any returns a matcher that matches any value.
var Any = Func("is anything", func(any) (bool, error) {
	return true, nil
})

// Matcher determines if the actual matches the expectation.
//
//go:generate mockery --name Matcher --output mock --outpkg mock --filename matcher.go
type Matcher interface {
	Match(actual any) (bool, error)
	Expected() string
}

var _ Matcher = (*equalMatcher)(nil)

// equalMatcher matches by equal string.
type equalMatcher struct {
	expected any
}

// Expected returns the expectation.
func (m equalMatcher) Expected() string {
	if v := strVal(m.expected); v != nil {
		return *v
	}

	return fmt.Sprintf("%+v", m.expected)
}

// Match determines if the actual is expected.
func (m equalMatcher) Match(actual any) (bool, error) {
	return assert.ObjectsAreEqual(m.expected, actual), nil
}

func (m equalMatcher) Format(s fmt.State, r rune) {
	format.Format(s, r, m.expected)
}

var _ Matcher = (*jsonMatcher)(nil)

// jsonMatcher matches by json with <ignore-diff> support.
type jsonMatcher struct {
	expected string
}

// Expected returns the expectation.
func (m jsonMatcher) Expected() string {
	return m.expected
}

// Match determines if the actual is expected.
func (m jsonMatcher) Match(actual any) (bool, error) {
	actualBytes, err := jsonVal(actual)
	if err != nil {
		return false, err
	}

	return assertjson.FailNotEqual([]byte(m.expected), actualBytes) == nil, nil
}

func (m jsonMatcher) Format(s fmt.State, r rune) {
	format.Format(s, r, m.expected)
}

var _ Matcher = (*regexMatcher)(nil)

// regexMatcher matches by regex.
type regexMatcher struct {
	regexp *regexp.Regexp
}

// Expected returns the expectation.
func (m regexMatcher) Expected() string {
	return m.regexp.String()
}

// Match determines if the actual is expected.
func (m regexMatcher) Match(actual any) (bool, error) {
	if v := strVal(actual); v != nil {
		return m.regexp.MatchString(*v), nil
	}

	return false, nil
}

func (m regexMatcher) Format(s fmt.State, r rune) {
	format.Format(s, r, m.regexp)
}

var _ Matcher = (*typeMatcher)(nil)

// typeMatcher is a .typeMatcher.
type typeMatcher struct {
	typeOf reflect.Type
}

func (m typeMatcher) Match(actual any) (bool, error) {
	return reflect.DeepEqual(m.typeOf, reflect.TypeOf(actual)), nil
}

func (m typeMatcher) Expected() string {
	return "type is " + m.typeOf.String()
}

func (m typeMatcher) Format(s fmt.State, _ rune) {
	_, _ = fmt.Fprintf(s, "<type is %s>", m.typeOf.String()) //nolint: errcheck
}

var _ Matcher = (*lenMatcher)(nil)

// lenMatcher matches by the length of the value.
type lenMatcher struct {
	expected int
}

// Match determines if the actual is expected.
func (m lenMatcher) Match(actual any) (_ bool, err error) {
	if actual == nil {
		return false, nil
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(recovered(r)) //nolint: err113
		}
	}()

	val := reflect.ValueOf(actual)

	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	return val.Len() == m.expected, nil
}

// Expected returns the expectation.
func (m lenMatcher) Expected() string {
	return fmt.Sprintf("len is %d", m.expected)
}

func (m lenMatcher) Format(s fmt.State, _ rune) {
	_, _ = fmt.Fprintf(s, "<len is %d>", m.expected) //nolint: errcheck
}

var _ Matcher = (*emptyMatcher)(nil)

// emptyMatcher checks whether the value is empty.
type emptyMatcher struct{}

// Match determines if the actual is expected.
func (emptyMatcher) Match(actual any) (bool, error) {
	return isEmpty(actual), nil
}

// Expected returns the expectation.
func (emptyMatcher) Expected() string {
	return "is empty"
}

func (emptyMatcher) Format(s fmt.State, _ rune) {
	_, _ = s.Write([]byte("<is empty>")) //nolint: errcheck
}

var _ Matcher = (*notEmptyMatcher)(nil)

// notEmptyMatcher checks whether the value is not empty.
type notEmptyMatcher struct{}

// Match determines if the actual is expected.
func (notEmptyMatcher) Match(actual any) (bool, error) {
	return !isEmpty(actual), nil
}

// Expected returns the expectation.
func (notEmptyMatcher) Expected() string {
	return "is not empty"
}

func (notEmptyMatcher) Format(s fmt.State, _ rune) {
	_, _ = s.Write([]byte("<is not empty>")) //nolint: errcheck
}

var _ Matcher = (*funcMatcher)(nil)

// funcMatcher checks by calling a function.
type funcMatcher struct {
	expected string
	match    func(actual any) (bool, error)
}

// Match determines if the actual is expected.
func (f funcMatcher) Match(actual any) (bool, error) {
	return f.match(actual)
}

// Expected returns the expectation.
func (f funcMatcher) Expected() string {
	return f.expected
}

func (f funcMatcher) Format(s fmt.State, _ rune) {
	_, _ = fmt.Fprintf(s, "<%s>", f.expected) //nolint: errcheck
}

var _ Matcher = (*Callback)(nil)

// Callback matches by calling a function.
type Callback func() Matcher

// Expected returns the expectation.
func (m Callback) Expected() string {
	return m().Expected()
}

// Match determines if the actual is expected.
func (m Callback) Match(actual any) (bool, error) {
	return m().Match(actual)
}

// Matcher returns the matcher.
func (m Callback) Matcher() Matcher {
	return m()
}

// Equal matches two objects.
func Equal(expected any) Matcher {
	return equalMatcher{expected: expected}
}

// Equalf matches two strings by the formatted expectation.
func Equalf(expected string, args ...any) Matcher {
	return equalMatcher{expected: fmt.Sprintf(expected, args...)}
}

// JSON matches two json strings with <ignore-diff> support.
func JSON(expected any) Matcher {
	ex, err := jsonVal(expected)
	if err != nil {
		panic(err)
	}

	return jsonMatcher{expected: string(ex)}
}

// Regex matches two strings by using regex.
func Regex[T ~string | *regexp.Regexp | regexp.Regexp](regexp T) Matcher {
	return regexMatcher{regexp: regexpVal(regexp)}
}

// Wildcard matches two strings by using equal or regex with wildcard support.
func Wildcard[T ~string](pattern T) Matcher {
	parts := strings.Split(string(pattern), "*")

	if len(parts) == 1 {
		return Equal(pattern)
	}

	var patternBuilder strings.Builder

	for i, part := range parts {
		if i > 0 {
			patternBuilder.WriteString(".*")
		}

		patternBuilder.WriteString(regexp.QuoteMeta(part))
	}

	return Regex(regexp.MustCompile("^" + patternBuilder.String() + "$"))
}

// IsType matches two types.
func IsType[T any]() Matcher {
	var t *T

	return typeMatcher{typeOf: reflect.TypeOf(t).Elem()}
}

// SameTypeAs matches two types.
func SameTypeAs(expected any) Matcher {
	return typeMatcher{typeOf: reflect.TypeOf(expected)}
}

// Len matches by the length of the value.
func Len[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](expected T) Matcher {
	return lenMatcher{expected: int(expected)}
}

// IsEmpty checks whether the value is empty.
func IsEmpty() Matcher {
	return emptyMatcher{}
}

// IsNotEmpty checks whether the value is not empty.
func IsNotEmpty() Matcher {
	return notEmptyMatcher{}
}

// Func matches by calling a function.
func Func(expected string, match func(actual any) (bool, error)) Matcher {
	return funcMatcher{expected: expected, match: match}
}

// Match returns a matcher according to its type.
func Match(v any) Matcher {
	switch val := v.(type) {
	case Matcher:
		return val

	case func() Matcher:
		return Callback(val)

	case regexp.Regexp, *regexp.Regexp:
		return Regex(regexpVal(val))

	case fmt.Stringer:
		return Equal(val.String())
	}

	return Equal(v)
}

const (
	logicalOperatorAnd = "and"
	logicalOperatorOr  = "or"
)

type logicalOperator string

type binaryLogicalMatcher struct {
	matchers []Matcher
	operator logicalOperator
	nested   bool
}

func (m *binaryLogicalMatcher) Expected() string {
	expected := make([]string, len(m.matchers))

	for i, matcher := range m.matchers {
		expected[i] = matcher.Expected()
	}

	if len(m.matchers) == 1 {
		return expected[0]
	}

	result := strings.Join(expected, " "+string(m.operator)+" ")

	if m.nested {
		result = "(" + result + ")"
	}

	return result
}

type orLogicalMatcher struct {
	*binaryLogicalMatcher
}

func (m *orLogicalMatcher) Match(actual any) (bool, error) {
	for _, matcher := range m.matchers {
		if ok, err := matcher.Match(actual); err != nil {
			return false, err
		} else if ok {
			return true, nil
		}
	}

	return false, nil
}

// Or returns a matcher that matches if any of the matchers match.
func Or(matchers ...any) Matcher {
	return &orLogicalMatcher{
		binaryLogicalMatcher: &binaryLogicalMatcher{
			matchers: makeNestableMatchers(matchers...),
			operator: logicalOperatorOr,
		},
	}
}

type andLogicalMatcher struct {
	*binaryLogicalMatcher
}

func (m *andLogicalMatcher) Match(actual any) (bool, error) {
	for _, matcher := range m.matchers {
		if ok, err := matcher.Match(actual); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		}
	}

	return true, nil
}

// And returns a matcher that matches if all of the matchers match.
func And(matchers ...any) Matcher {
	return &andLogicalMatcher{
		binaryLogicalMatcher: &binaryLogicalMatcher{
			matchers: makeNestableMatchers(matchers...),
			operator: logicalOperatorAnd,
		},
	}
}

func makeNestableMatchers(v ...any) []Matcher {
	matchers := make([]Matcher, len(v))

	for i, matcher := range v {
		matchers[i] = makeNestedMatcher(matcher)
	}

	return matchers
}

func makeNestedMatcher(v any) Matcher {
	m := Match(v)

	switch m := m.(type) {
	case *orLogicalMatcher:
		m.nested = true
	case *andLogicalMatcher:
		m.nested = true
	}

	return m
}
