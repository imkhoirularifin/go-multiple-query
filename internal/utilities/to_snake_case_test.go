package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	// Test cases
	tests := []struct {
		input    string
		expected string
	}{
		{input: "helloWorld", expected: "hello_world"},
		{input: "fooBarBaz", expected: "foo_bar_baz"},
		{input: "camelCase", expected: "camel_case"},
		{input: "snake_case", expected: "snake_case"},
		{input: "UPPER_CASE", expected: "upper_case"},
	}

	// Run tests
	for _, test := range tests {
		result := ToSnakeCase(test.input)
		assert.Equal(t, test.expected, result)
	}
}
