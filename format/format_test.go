package format_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.nhat.io/matcher/v3/format"
)

func TestFprintf(t *testing.T) {
	t.Parallel()

	buf := new(bytes.Buffer)
	_, err := format.Fprintf(buf, "%#v", "foobar")

	require.NoError(t, err)

	expected := `string("foobar")`

	assert.Equal(t, expected, buf.String())
}

func TestSprintf_String(t *testing.T) {
	t.Parallel()

	testCases := provideFormatValueTestCases(t, "foobar", formatValueTestCaseExpects{
		expectS:      "foobar",
		expectPlusS:  "foobar",
		expectSharpS: `"foobar"`,
		expectV:      "string(foobar)",
		expectPlusV:  "string(foobar)",
		expectSharpV: `string("foobar")`,
		expectQ:      `"foobar"`,
		expectSharpQ: `string("foobar")`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_String_Alias(t *testing.T) {
	t.Parallel()

	type String string

	testCases := provideFormatValueTestCases(t, String("foobar"), formatValueTestCaseExpects{
		expectS:      "foobar",
		expectPlusS:  "foobar",
		expectSharpS: `"foobar"`,
		expectV:      "format_test.String(foobar)",
		expectPlusV:  "format_test.String(foobar)",
		expectSharpV: `format_test.String("foobar")`,
		expectQ:      `format_test.String(foobar)`,
		expectSharpQ: `format_test.String("foobar")`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_String_Slice(t *testing.T) {
	t.Parallel()

	testCases := provideFormatValueTestCases(t, []string{"foobar"}, formatValueTestCaseExpects{
		expectS:      "[foobar]",
		expectPlusS:  "[foobar]",
		expectSharpS: `[]string{"foobar"}`,
		expectV:      "[]string([foobar])",
		expectPlusV:  "[]string([foobar])",
		expectSharpV: `[]string{"foobar"}`,
		expectQ:      `[]string([foobar])`,
		expectSharpQ: `[]string{"foobar"}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_String_SliceAlias(t *testing.T) {
	t.Parallel()

	type Strings []string

	testCases := provideFormatValueTestCases(t, Strings{"foobar"}, formatValueTestCaseExpects{
		expectS:      "[foobar]",
		expectPlusS:  "[foobar]",
		expectSharpS: `format_test.Strings{"foobar"}`,
		expectV:      "format_test.Strings([foobar])",
		expectPlusV:  "format_test.Strings([foobar])",
		expectSharpV: `format_test.Strings{"foobar"}`,
		expectQ:      `format_test.Strings([foobar])`,
		expectSharpQ: `format_test.Strings{"foobar"}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Byte_Slice(t *testing.T) {
	t.Parallel()

	testCases := provideFormatValueTestCases(t, []byte("foobar"), formatValueTestCaseExpects{
		expectS:      "[102 111 111 98 97 114]",
		expectPlusS:  "[102 111 111 98 97 114]",
		expectSharpS: `[]byte{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72}`,
		expectV:      "[]uint8([102 111 111 98 97 114])",
		expectPlusV:  "[]uint8([102 111 111 98 97 114])",
		expectSharpV: `[]byte{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72}`,
		expectQ:      `[]uint8([102 111 111 98 97 114])`,
		expectSharpQ: `[]byte{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Byte_SliceAlias(t *testing.T) {
	t.Parallel()

	type Raw []byte

	testCases := provideFormatValueTestCases(t, Raw("foobar"), formatValueTestCaseExpects{
		expectS:      "[102 111 111 98 97 114]",
		expectPlusS:  "[102 111 111 98 97 114]",
		expectSharpS: `format_test.Raw{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72}`,
		expectV:      "format_test.Raw([102 111 111 98 97 114])",
		expectPlusV:  "format_test.Raw([102 111 111 98 97 114])",
		expectSharpV: `format_test.Raw{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72}`,
		expectQ:      `format_test.Raw([102 111 111 98 97 114])`,
		expectSharpQ: `format_test.Raw{0x66, 0x6f, 0x6f, 0x62, 0x61, 0x72}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Int(t *testing.T) {
	t.Parallel()

	testCases := provideFormatValueTestCases(t, 42, formatValueTestCaseExpects{
		expectS:      "42",
		expectPlusS:  "42",
		expectSharpS: "42",
		expectV:      "int(42)",
		expectPlusV:  "int(42)",
		expectSharpV: "int(42)",
		expectQ:      "int(42)",
		expectSharpQ: "int(42)",
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Int_Alias(t *testing.T) {
	t.Parallel()

	type Int int

	testCases := provideFormatValueTestCases(t, Int(42), formatValueTestCaseExpects{
		expectS:      "42",
		expectPlusS:  "42",
		expectSharpS: "42",
		expectV:      "format_test.Int(42)",
		expectPlusV:  "format_test.Int(42)",
		expectSharpV: `format_test.Int(42)`,
		expectQ:      `format_test.Int(42)`,
		expectSharpQ: `format_test.Int(42)`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Int_Slice(t *testing.T) {
	t.Parallel()

	testCases := provideFormatValueTestCases(t, []int{42}, formatValueTestCaseExpects{
		expectS:      "[42]",
		expectPlusS:  "[42]",
		expectSharpS: `[]int{42}`,
		expectV:      "[]int([42])",
		expectPlusV:  "[]int([42])",
		expectSharpV: `[]int{42}`,
		expectQ:      `[]int([42])`,
		expectSharpQ: `[]int{42}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Int_SliceAlias(t *testing.T) {
	t.Parallel()

	type Ints []int

	testCases := provideFormatValueTestCases(t, Ints{42}, formatValueTestCaseExpects{
		expectS:      "[42]",
		expectPlusS:  "[42]",
		expectSharpS: `format_test.Ints{42}`,
		expectV:      "format_test.Ints([42])",
		expectPlusV:  "format_test.Ints([42])",
		expectSharpV: `format_test.Ints{42}`,
		expectQ:      `format_test.Ints([42])`,
		expectSharpQ: `format_test.Ints{42}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_JSONRawMessage(t *testing.T) {
	t.Parallel()

	const payload = `{"foo":"bar"}`

	testCases := provideFormatValueTestCases(t, json.RawMessage(payload), formatValueTestCaseExpects{
		expectS:      `{"foo":"bar"}`,
		expectPlusS:  `{"foo":"bar"}`,
		expectSharpS: `{"foo":"bar"}`,
		expectV:      `json.RawMessage({"foo":"bar"})`,
		expectPlusV:  `json.RawMessage({"foo":"bar"})`,
		expectSharpV: `json.RawMessage({"foo":"bar"})`,
		expectQ:      `{"foo":"bar"}`,
		expectSharpQ: `{"foo":"bar"}`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSprintf_Regexp(t *testing.T) {
	t.Parallel()

	testCases := provideFormatValueTestCases(t, regexp.MustCompile(`.*`), formatValueTestCaseExpects{
		expectS:      ".*",
		expectPlusS:  ".*",
		expectSharpS: `".*"`,
		expectV:      "*regexp.Regexp(.*)",
		expectPlusV:  "*regexp.Regexp(.*)",
		expectSharpV: `*regexp.Regexp(".*")`,
		expectQ:      `".*"`,
		expectSharpQ: `*regexp.Regexp(".*")`,
	})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			actual := format.Sprintf(tc.format, tc.value)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

type formatValueTestCase struct {
	scenario string
	format   string
	value    any
	expected string
}

type formatValueTestCaseExpects struct {
	expectS      string
	expectPlusS  string
	expectSharpS string
	expectV      string
	expectPlusV  string
	expectSharpV string
	expectQ      string
	expectSharpQ string
}

func provideFormatValueTestCases(t *testing.T, value any, expects formatValueTestCaseExpects) []formatValueTestCase {
	t.Helper()

	return []formatValueTestCase{
		{
			scenario: formatValueTestCaseName(t, "%s"),
			format:   "%s",
			value:    value,
			expected: expects.expectS,
		},
		{
			scenario: formatValueTestCaseName(t, "%+s"),
			format:   "%+s",
			value:    value,
			expected: expects.expectPlusS,
		},
		{
			scenario: formatValueTestCaseName(t, "%#s"),
			format:   "%#s",
			value:    value,
			expected: expects.expectSharpS,
		},
		{
			scenario: formatValueTestCaseName(t, "%v"),
			format:   "%v",
			value:    value,
			expected: expects.expectV,
		},
		{
			scenario: formatValueTestCaseName(t, "%+v"),
			format:   "%+v",
			value:    value,
			expected: expects.expectPlusV,
		},
		{
			scenario: formatValueTestCaseName(t, "%#v"),
			format:   "%#v",
			value:    value,
			expected: expects.expectSharpV,
		},
		{
			scenario: formatValueTestCaseName(t, "%q"),
			format:   "%q",
			value:    value,
			expected: expects.expectQ,
		},
		{
			scenario: formatValueTestCaseName(t, "%#q"),
			format:   "%#q",
			value:    value,
			expected: expects.expectSharpQ,
		},
	}
}

func formatValueTestCaseName(t *testing.T, format string) string {
	t.Helper()

	return fmt.Sprintf("%s - %s", t.Name(), format)
}
