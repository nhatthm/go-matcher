package matcher_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/go-matcher"
)

func TestExactMatch_Expected(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		input    interface{}
		expected string
	}{
		{
			scenario: "string",
			input:    "foobar",
			expected: "foobar",
		},
		{
			scenario: "not a string",
			input:    42,
			expected: "42",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.Exact(tc.input)

			assert.Equal(t, tc.expected, m.Expected())
		})
	}
}

func TestExactMatch_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		actual   string
		expected bool
	}{
		{
			scenario: "match",
			actual:   "value",
			expected: true,
		},
		{
			scenario: "no match",
			actual:   "mismatch",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.Exact("value")
			result, err := m.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			assert.NoError(t, err)
		})
	}
}

func TestExactfMatch_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		format   string
		args     []interface{}
		actual   string
		expected bool
	}{
		{
			scenario: "match",
			format:   "Bearer %s",
			args:     []interface{}{"token"},
			actual:   "Bearer token",
			expected: true,
		},
		{
			scenario: "no match",
			format:   "Bearer %s",
			args:     []interface{}{"token"},
			actual:   "Bearer unknown",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.Exactf(tc.format, tc.args...)
			result, err := m.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			assert.NoError(t, err)
		})
	}
}

func TestJSONMatch_Panic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		matcher.JSON(make(chan error))
	})
}

func TestJSONMatch_Expected(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		input    interface{}
		expected string
	}{
		{
			scenario: "string",
			input:    "foobar",
			expected: "foobar",
		},
		{
			scenario: "not a string",
			input:    42,
			expected: "42",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.JSON(tc.input)

			assert.Equal(t, tc.expected, m.Expected())
		})
	}
}

func TestJSONMatch_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		json     string
		actual   string
		expected bool
	}{
		{
			scenario: "match",
			json: `{
	"username": "user"
}`,
			actual:   `{"username": "user"}`,
			expected: true,
		},
		{
			scenario: "match with <ignore-diff>",
			json:     `{"username": "<ignore-diff>"}`,
			actual:   `{"username": "user"}`,
			expected: true,
		},
		{
			scenario: "no match",
			json:     "{}",
			actual:   "[]",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.JSON(tc.json)
			result, err := m.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			assert.NoError(t, err)
		})
	}
}

func TestJSONMatch_Match_Error(t *testing.T) {
	t.Parallel()

	m := matcher.JSON(`{}`)
	result, err := m.Match(make(chan error))

	assert.False(t, result)
	assert.EqualError(t, err, `json: unsupported type: chan error`)
}

func TestRegexMatch_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.Regex(regexp.MustCompile(".*"))
	expected := ".*"

	assert.Equal(t, expected, m.Expected())
}

func TestLenMatcher_Match_NoError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		value    interface{}
		expected bool
	}{
		{
			scenario: "empty string",
			value:    "",
		},
		{
			scenario: "string len mismatched",
			value:    "foob",
		},
		{
			scenario: "string len matched",
			value:    "foo",
			expected: true,
		},
		{
			scenario: "empty slice",
			value:    []int{},
		},
		{
			scenario: "slice len mismatched",
			value:    []int{1, 2},
		},
		{
			scenario: "slice len matched",
			value:    []int{1, 2, 3},
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.Len(3)
			actual, err := m.Match(tc.value)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestLenMatcher_Match_Error(t *testing.T) {
	t.Parallel()

	m := matcher.Len(3)
	actual, err := m.Match(42)

	expected := `reflect: call of reflect.Value.Len on int Value`

	assert.False(t, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expected)
}

func TestLenMatcher_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.Len(5)
	expected := "len is 5"

	assert.Equal(t, expected, m.Expected())
}

func TestEmptyMatcher_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		value    string
		expected bool
	}{
		{
			scenario: "empty",
			expected: true,
		},
		{
			scenario: "not empty",
			value:    "foobar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.IsEmpty()
			actual, err := m.Match(tc.value)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestEmptyMatcher_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.IsEmpty()
	expected := "is empty"

	assert.Equal(t, expected, m.Expected())
}

func TestNotEmptyMatcher_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		value    string
		expected bool
	}{
		{
			scenario: "empty",
			expected: true,
		},
		{
			scenario: "not empty",
			value:    "foobar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.IsNotEmpty()
			actual, err := m.Match(tc.value)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestNotEmptyMatcher_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.IsNotEmpty()
	expected := "is not empty"

	assert.Equal(t, expected, m.Expected())
}

func TestRegexMatch_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		matcher  matcher.RegexMatcher
		actual   interface{}
		expected bool
	}{
		{
			scenario: "match with regexp",
			matcher:  matcher.Regex(regexp.MustCompile(".*")),
			actual:   `hello`,
			expected: true,
		},
		{
			scenario: "match with regexp pattern",
			matcher:  matcher.RegexPattern(".*"),
			actual:   `hello`,
			expected: true,
		},
		{
			scenario: "no match",
			matcher:  matcher.RegexPattern("^[0-9]+$"),
			actual:   "mismatch",
		},
		{
			scenario: "not a string",
			matcher:  matcher.Regex(nil),
			actual:   nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := tc.matcher.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			assert.NoError(t, err)
		})
	}
}

func TestCallback(t *testing.T) {
	t.Parallel()

	m := matcher.Callback(func() matcher.Matcher {
		return matcher.Exact("expected")
	})

	assert.Equal(t, matcher.Exact("expected"), m.Matcher())
}

func TestMatch(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		value    interface{}
		expected matcher.Matcher
	}{
		{
			scenario: "matcher",
			value:    matcher.Exact("expected"),
			expected: matcher.Exact("expected"),
		},
		{
			scenario: "[]byte",
			value:    []byte("expected"),
			expected: matcher.Exact([]byte("expected")),
		},
		{
			scenario: "string",
			value:    "expected",
			expected: matcher.Exact("expected"),
		},
		{
			scenario: "int",
			value:    42,
			expected: matcher.Exact(42),
		},
		{
			scenario: "regexp",
			value:    regexp.MustCompile(".*"),
			expected: matcher.RegexPattern(".*"),
		},
		{
			scenario: "fmt.Stringer",
			value:    time.UTC,
			expected: matcher.Exact("UTC"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, matcher.Match(tc.value))
		})
	}
}

func TestMatch_Callback(t *testing.T) {
	t.Parallel()

	m := matcher.Match(func() matcher.Matcher {
		return matcher.Exact("expected")
	})

	assert.Equal(t, "expected", m.Expected())

	result, err := m.Match("expected")

	assert.True(t, result)
	assert.NoError(t, err)

	result, err = m.Match("mismatched")

	assert.False(t, result)
	assert.NoError(t, err)
}
