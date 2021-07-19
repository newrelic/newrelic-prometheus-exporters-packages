package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_removeTrailingCommas(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    `{"a":1,}`,
			expected: `{"a":1}`,
		},
		{
			input:    `{"a":1, }`,
			expected: `{"a":1}`,
		},
		{
			input:    `{"a":1,, }`,
			expected: `{"a":1,}`,
		},
		{
			input:    `{"a":"1,," }`,
			expected: `{"a":"1,," }`,
		},
		{
			input:    `{"a":20, }`,
			expected: `{"a":20}`,
		},
		{
			input:    `["1,," ]`,
			expected: `["1,," ]`,
		},
		{
			input:    `["1,,", ]`,
			expected: `["1,,"]`,
		},
		{
			input:    `[{\"a\":1,},]`,
			expected: `[{\"a\":1}]`,
		},
	}
	for i := range cases {
		c := cases[i]
		output := removeTrailingCommas(c.input)
		assert.Equal(t, c.expected, output)
	}
}
