package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_strVal(t *testing.T) {
	t.Parallel()

	expected := "foobar"

	testCases := []struct {
		scenario string
		input    interface{}
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
		tc := tc
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
		input          interface{}
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
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			result, err := jsonVal(tc.input)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
