package matcher

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_strVal(t *testing.T) {
	t.Parallel()

	expected := "foobar"

	testCases := []struct {
		scenario string
		input    any
		expected *string
	}{
		{
			scenario: "string",
			input:    "foobar",
			expected: &expected,
		},
		{
			scenario: "[]",
			input:    []byte("foobar"),
			expected: &expected,
		},
		{
			scenario: "not a string or []byte",
			input:    42,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, strVal(tc.input))
		})
	}
}

func Test_jsonVal(t *testing.T) {
	t.Parallel()

	const payload = `{"name":"foobar"}`

	expected := []byte(payload)

	testCases := []struct {
		scenario       string
		input          any
		expectedResult []byte
		expectedError  string
	}{
		{
			scenario:      "chan",
			input:         make(chan struct{}),
			expectedError: `json: unsupported type: chan struct {}`,
		},
		{
			scenario:       "string",
			input:          payload,
			expectedResult: expected,
		},
		{
			scenario:       "[]",
			input:          []byte(payload),
			expectedResult: expected,
		},
		{
			scenario:       "map[string]string",
			input:          map[string]string{"name": "foobar"},
			expectedResult: expected,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := jsonVal(tc.input)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func Test_regexpVal(t *testing.T) {
	t.Parallel()

	expected := regexp.MustCompile(".*")

	testCases := []struct {
		scenario string
		input    any
		expected *regexp.Regexp
	}{
		{
			scenario: "chan",
			input:    make(chan struct{}),
			expected: nil,
		},
		{
			scenario: "string",
			input:    ".*",
			expected: expected,
		},
		{
			scenario: "regexp ptr",
			input:    expected,
			expected: expected,
		},
		{
			scenario: "regexp ptr",
			input:    *expected,
			expected: expected,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result := regexpVal(tc.input)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	errCh := make(chan error, 1)
	errCh <- errors.New("error")

	nonEmptyStr := "foo"
	emptyStr := ""

	testCases := []struct {
		scenario string
		value    any
		expected bool
	}{
		{
			scenario: "nil",
			expected: true,
		},
		{
			scenario: "empty array",
			value:    [0]int{},
			expected: true,
		},
		{
			scenario: "not empty array",
			value:    [1]int{},
		},
		{
			scenario: "empty slice",
			value:    []int{},
			expected: true,
		},
		{
			scenario: "not empty slice",
			value:    []int{1},
		},
		{
			scenario: "empty chan",
			value:    make(chan error, 1),
			expected: true,
		},
		{
			scenario: "not empty chan",
			value:    errCh,
		},
		{
			scenario: "empty map",
			value:    map[string]int{},
			expected: true,
		},
		{
			scenario: "not empty map",
			value:    map[string]int{"id": 1},
		},
		{
			scenario: "empty string",
			value:    "",
			expected: true,
		},
		{
			scenario: "nil interface",
			value:    (*error)(nil),
			expected: true,
		},
		{
			scenario: "empty string ptr",
			value:    &emptyStr,
			expected: true,
		},
		{
			scenario: "not empty string ptr",
			value:    &nonEmptyStr,
		},
		{
			scenario: "not empty string",
			value:    "foobar",
		},
		{
			scenario: "empty int",
			value:    0,
			expected: true,
		},
		{
			scenario: "not empty int",
			value:    42,
		},
		{
			scenario: "false",
			value:    false,
			expected: true,
		},
		{
			scenario: "true",
			value:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, isEmpty(tc.value))
		})
	}
}
