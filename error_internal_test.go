package matcher

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovered(t *testing.T) {
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
			scenario: "error",
			input:    errors.New("foobar"),
			expected: "foobar",
		},
		{
			scenario: "integer",
			input:    42,
			expected: "42",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := recovered(tc.input)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
