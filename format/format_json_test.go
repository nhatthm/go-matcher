//go:build go1.27

package format_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/matcher/v3/format"
)

func TestSprintf_JSONRawMessage(t *testing.T) {
	t.Parallel()

	const payload = `{"foo":"bar"}`

	testCases := provideFormatValueTestCases(t, jsontext.Value(payload), formatValueTestCaseExpects{
		expectS:      `{"foo":"bar"}`,
		expectPlusS:  `{"foo":"bar"}`,
		expectSharpS: `{"foo":"bar"}`,
		expectV:      `jsontext.Value({"foo":"bar"})`,
		expectPlusV:  `jsontext.Value({"foo":"bar"})`,
		expectSharpV: `jsontext.Value({"foo":"bar"})`,
		expectQ:      `{"foo":"bar"}`,
		expectSharpQ: `{"foo":"bar"}`,
	})

	for _, tc := range testCases {
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
