package matcher_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/matcher/v3"
)

func TestAny(t *testing.T) {
	t.Parallel()

	t.Run("expected", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "is anything", matcher.Any.Expected())
	})

	t.Run("format", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "<is anything>", fmt.Sprintf("%#v", matcher.Any))
	})

	testCases := []struct {
		scenario string
		actual   any
	}{
		{
			scenario: "int",
			actual:   42,
		},
		{
			scenario: "string",
			actual:   "foobar",
		},
		{
			scenario: "struct",
			actual:   struct{}{},
		},
		{
			scenario: "nil",
			actual:   nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			matched, err := matcher.Any.Match(tc.actual)

			assert.True(t, matched)
			require.NoError(t, err)
		})
	}
}

func TestEqualMatcher_Format(t *testing.T) {
	t.Parallel()

	type data struct {
		Name string
	}

	testCases := []struct {
		scenario string
		format   string
		value    any
		expected string
	}{
		{
			scenario: "string - %T",
			format:   "%T",
			value:    "foobar",
			expected: "matcher.equalMatcher",
		},
		{
			scenario: "string - %s",
			format:   "%s",
			value:    "foobar",
			expected: "foobar",
		},
		{
			scenario: "string - %+s",
			format:   "%+s",
			value:    "foobar",
			expected: "foobar",
		},
		{
			scenario: "string - %#s",
			format:   "%#s",
			value:    "foobar",
			expected: `"foobar"`,
		},
		{
			scenario: "string - %v",
			format:   "%v",
			value:    "foobar",
			expected: "string(foobar)",
		},
		{
			scenario: "string - %+v",
			format:   "%+v",
			value:    "foobar",
			expected: "string(foobar)",
		},
		{
			scenario: "string - %#v",
			format:   "%#v",
			value:    "foobar",
			expected: `string("foobar")`,
		},
		{
			scenario: "string - %q",
			format:   "%q",
			value:    "foobar",
			expected: `"foobar"`,
		},
		{
			scenario: "string - %+q",
			format:   "%+q",
			value:    "foobar",
			expected: `"foobar"`,
		},
		{
			scenario: "string - %#q",
			format:   "%#q",
			value:    "foobar",
			expected: `string("foobar")`,
		},
		{
			scenario: "struct - %T",
			format:   "%T",
			value:    data{Name: "foobar"},
			expected: "matcher.equalMatcher",
		},
		{
			scenario: "struct - %s",
			format:   "%s",
			value:    data{Name: "foobar"},
			expected: "{foobar}",
		},
		{
			scenario: "struct - %+s",
			format:   "%+s",
			value:    data{Name: "foobar"},
			expected: "{Name:foobar}",
		},
		{
			scenario: "struct - %#s",
			format:   "%#s",
			value:    data{Name: "foobar"},
			expected: `matcher_test.data{Name:"foobar"}`,
		},
		{
			scenario: "struct - %v",
			format:   "%v",
			value:    data{Name: "foobar"},
			expected: "matcher_test.data({foobar})",
		},
		{
			scenario: "struct - %+v",
			format:   "%+v",
			value:    data{Name: "foobar"},
			expected: "matcher_test.data({Name:foobar})",
		},
		{
			scenario: "struct - %#v",
			format:   "%#v",
			value:    data{Name: "foobar"},
			expected: `matcher_test.data{Name:"foobar"}`,
		},
		{
			scenario: "struct - %q",
			format:   "%q",
			value:    data{Name: "foobar"},
			expected: `matcher_test.data({foobar})`,
		},
		{
			scenario: "struct - %+q",
			format:   "%+q",
			value:    data{Name: "foobar"},
			expected: `matcher_test.data({Name:foobar})`,
		},
		{
			scenario: "struct - %#q",
			format:   "%#q",
			value:    data{Name: "foobar"},
			expected: `matcher_test.data{Name:"foobar"}`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := fmt.Sprintf(tc.format, matcher.Equal(tc.value))

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestEqual_Expected(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		input    any
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

			m := matcher.Equal(tc.input)

			assert.Equal(t, tc.expected, m.Expected())
		})
	}
}

func TestEqual_Match(t *testing.T) {
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

			m := matcher.Equal("value")
			result, err := m.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			require.NoError(t, err)
		})
	}
}

func TestEqualf_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		format   string
		args     []any
		actual   string
		expected bool
	}{
		{
			scenario: "match",
			format:   "Bearer %s",
			args:     []any{"token"},
			actual:   "Bearer token",
			expected: true,
		},
		{
			scenario: "no match",
			format:   "Bearer %s",
			args:     []any{"token"},
			actual:   "Bearer unknown",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.Equalf(tc.format, tc.args...)
			result, err := m.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			require.NoError(t, err)
		})
	}
}

func TestJsonMatcher_Format(t *testing.T) {
	t.Parallel()

	const payload = `{"username": "user"}`

	testCases := []struct {
		scenario string
		format   string
		value    string
		expected string
	}{
		{
			scenario: "type - %T",
			format:   "%T",
			value:    payload,
			expected: "matcher.jsonMatcher",
		},
		{
			scenario: "string - %s",
			format:   "%s",
			value:    payload,
			expected: payload,
		},
		{
			scenario: "string - %+s",
			format:   "%+s",
			value:    payload,
			expected: payload,
		},
		{
			scenario: "string - %#s",
			format:   "%#s",
			value:    payload,
			expected: `"{\"username\": \"user\"}"`,
		},
		{
			scenario: "string - %v",
			format:   "%v",
			value:    payload,
			expected: `string({"username": "user"})`,
		},
		{
			scenario: "string - %+v",
			format:   "%+v",
			value:    payload,
			expected: `string({"username": "user"})`,
		},
		{
			scenario: "string - %#v",
			format:   "%#v",
			value:    payload,
			expected: `string("{\"username\": \"user\"}")`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := fmt.Sprintf(tc.format, matcher.JSON(tc.value))

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestJSON_Panic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		matcher.JSON(make(chan error))
	})
}

func TestJSON_Expected(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		input    any
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

func TestJSON_Match(t *testing.T) {
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
			require.NoError(t, err)
		})
	}
}

func TestJSON_Match_Error(t *testing.T) {
	t.Parallel()

	m := matcher.JSON(`{}`)
	result, err := m.Match(make(chan error))

	assert.False(t, result)
	require.EqualError(t, err, `json: unsupported type: chan error`)
}

func TestRegexMatcher_Format(t *testing.T) {
	t.Parallel()

	const pattern = `.*`

	testCases := []struct {
		scenario string
		format   string
		value    string
		expected string
	}{
		{
			scenario: "type - %T",
			format:   "%T",
			value:    pattern,
			expected: "matcher.regexMatcher",
		},
		{
			scenario: "string - %s",
			format:   "%s",
			value:    pattern,
			expected: pattern,
		},
		{
			scenario: "string - %+s",
			format:   "%+s",
			value:    pattern,
			expected: pattern,
		},
		{
			scenario: "string - %#s",
			format:   "%#s",
			value:    pattern,
			expected: `".*"`,
		},
		{
			scenario: "string - %v",
			format:   "%v",
			value:    pattern,
			expected: `*regexp.Regexp(.*)`,
		},
		{
			scenario: "string - %+v",
			format:   "%+v",
			value:    pattern,
			expected: `*regexp.Regexp(.*)`,
		},
		{
			scenario: "string - %#v",
			format:   "%#v",
			value:    pattern,
			expected: `*regexp.Regexp(".*")`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := fmt.Sprintf(tc.format, matcher.Regex(tc.value))

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestRegex_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.Regex(regexp.MustCompile(".*"))
	expected := ".*"

	assert.Equal(t, expected, m.Expected())
}

func TestTypeMatcher_Format(t *testing.T) {
	t.Parallel()

	m := matcher.IsType[string]()

	actual := fmt.Sprintf("%#v", m)
	expected := "<type is string>"

	assert.Equal(t, expected, actual)
}

func TestIsType_Match(t *testing.T) {
	t.Parallel()

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		m := matcher.IsType[bool]()

		actual, err := m.Match(true)

		assert.True(t, actual)
		require.NoError(t, err)

		actual, err = m.Match(1)

		assert.False(t, actual)
		require.NoError(t, err)
	})

	t.Run("*time.Time", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		m := matcher.IsType[*time.Time]()

		actual, err := m.Match(&now)

		assert.True(t, actual)
		require.NoError(t, err)

		actual, err = m.Match(now)

		assert.False(t, actual)
		require.NoError(t, err)
	})
}

func TestIsType_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.IsType[bool]()

	expected := `type is bool`

	assert.Equal(t, expected, m.Expected())
}

func TestSameTypeAs_Match(t *testing.T) {
	t.Parallel()

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		m := matcher.SameTypeAs(true)

		actual, err := m.Match(false)

		assert.True(t, actual)
		require.NoError(t, err)

		actual, err = m.Match(1)

		assert.False(t, actual)
		require.NoError(t, err)
	})

	t.Run("*time.Time", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		m := matcher.SameTypeAs(&time.Time{})

		actual, err := m.Match(&now)

		assert.True(t, actual)
		require.NoError(t, err)

		actual, err = m.Match(now)

		assert.False(t, actual)
		require.NoError(t, err)
	})
}

func TestSameTypeAs_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.SameTypeAs(true)

	expected := `type is bool`

	assert.Equal(t, expected, m.Expected())
}

func TestLen_Match_NoError(t *testing.T) {
	t.Parallel()

	str := "foo"

	testCases := []struct {
		scenario string
		value    any
		expected bool
	}{
		{
			scenario: "nil",
		},
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
			value:    str,
			expected: true,
		},
		{
			scenario: "string len ptr matched",
			value:    &str,
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
			scenario: "slice ptr len mismatched",
			value:    &[]int{1, 2},
		},
		{
			scenario: "slice len matched",
			value:    []int{1, 2, 3},
			expected: true,
		},
		{
			scenario: "slice ptr len matched",
			value:    &[]int{1, 2, 3},
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.Len(3)
			actual, err := m.Match(tc.value)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestLen_Match_Error(t *testing.T) {
	t.Parallel()

	m := matcher.Len(3)
	actual, err := m.Match(42)

	expected := `reflect: call of reflect.Value.Len on int Value`

	assert.False(t, actual)
	require.Error(t, err)
	require.EqualError(t, err, expected)
}

func TestLenMatcher_Format(t *testing.T) {
	t.Parallel()

	m := matcher.Len(10)

	actual := fmt.Sprintf("%#v", m)
	expected := "<len is 10>"

	assert.Equal(t, expected, actual)
}

func TestLen_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.Len(5)
	expected := "len is 5"

	assert.Equal(t, expected, m.Expected())
}

func TestEmptyMatcher_Format(t *testing.T) {
	t.Parallel()

	actual := fmt.Sprintf("%#v", matcher.IsEmpty())
	expected := "<is empty>"

	assert.Equal(t, expected, actual)
}

func TestEmpty_Match(t *testing.T) {
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

			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestEmpty_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.IsEmpty()
	expected := "is empty"

	assert.Equal(t, expected, m.Expected())
}

func TestNotEmptyMatcher_Format(t *testing.T) {
	t.Parallel()

	actual := fmt.Sprintf("%#v", matcher.IsNotEmpty())
	expected := "<is not empty>"

	assert.Equal(t, expected, actual)
}

func TestNotEmpty_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		value    string
		expected bool
	}{
		{
			scenario: "empty",
		},
		{
			scenario: "not empty",
			value:    "foobar",
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			m := matcher.IsNotEmpty()
			actual, err := m.Match(tc.value)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestNotEmpty_Expected(t *testing.T) {
	t.Parallel()

	m := matcher.IsNotEmpty()
	expected := "is not empty"

	assert.Equal(t, expected, m.Expected())
}

func TestRegex_Match(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		matcher  matcher.Matcher
		actual   any
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
			matcher:  matcher.Regex(".*"),
			actual:   `hello`,
			expected: true,
		},
		{
			scenario: "no match",
			matcher:  matcher.Regex("^[0-9]+$"),
			actual:   "mismatch",
		},
		{
			scenario: "not a string",
			matcher:  matcher.Regex(""),
			actual:   nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := tc.matcher.Match(tc.actual)

			assert.Equal(t, tc.expected, result)
			require.NoError(t, err)
		})
	}
}

func TestCallback(t *testing.T) {
	t.Parallel()

	m := matcher.Callback(func() matcher.Matcher {
		return matcher.Equal("expected")
	})

	assert.Equal(t, matcher.Equal("expected"), m.Matcher())
}

func TestMatch(t *testing.T) {
	t.Parallel()

	reg := regexp.MustCompile(".*")

	testCases := []struct {
		scenario string
		value    any
		expected matcher.Matcher
	}{
		{
			scenario: "matcher",
			value:    matcher.Equal("expected"),
			expected: matcher.Equal("expected"),
		},
		{
			scenario: "[]byte",
			value:    []byte("expected"),
			expected: matcher.Equal([]byte("expected")),
		},
		{
			scenario: "string",
			value:    "expected",
			expected: matcher.Equal("expected"),
		},
		{
			scenario: "int",
			value:    42,
			expected: matcher.Equal(42),
		},
		{
			scenario: "*regexp",
			value:    reg,
			expected: matcher.Regex(".*"),
		},
		{
			scenario: "regexp",
			value:    *reg,
			expected: matcher.Regex(".*"),
		},
		{
			scenario: "fmt.Stringer",
			value:    time.UTC,
			expected: matcher.Equal("UTC"),
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
		return matcher.Equal("expected")
	})

	assert.Equal(t, "expected", m.Expected())

	result, err := m.Match("expected")

	assert.True(t, result)
	require.NoError(t, err)

	result, err = m.Match("mismatched")

	assert.False(t, result)
	require.NoError(t, err)
}
