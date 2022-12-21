package pwapi

import (
	"fmt"
	"net/url"
	"testing"
)

func TestPascalCaseToSnakeCase(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "single word", input: "Hello", expected: "hello"},
		{name: "two words", input: "HelloWorld", expected: "hello_world"},
		{name: "two words with numbers", input: "HelloWorld3", expected: "hello_world3"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := PascalCaseToSnakeCase(test.input)
			if output != test.expected {
				t.Errorf("Expected output to be '%s', but got '%s'", test.expected, output)
			}
		})
	}
}

func TestEncodeSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    url.Values
		expected string
	}{
		{"One param", url.Values{"FooBar": {"baz"}}, "foo_bar=baz"},
		{"One param with two value", url.Values{"FooBar": {"baz", "qux"}}, "foo_bar=baz&foo_bar=qux"},
		{"One param value empty", url.Values{"FooBar": {}, "BarBaz": {"qux"}}, "bar_baz=qux"},
		{"Two params unsorted", url.Values{"FooBar": {"baz"}, "BarBaz": {"qux"}}, "bar_baz=qux&foo_bar=baz"},
		{"Empty", url.Values{}, ""},
		{"Nil", nil, ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			output := EncodeSnakeCase(test.input)
			if output != test.expected {
				t.Errorf("Expected output to be '%s', but got '%s'", test.expected, output)
			}
		})
	}
}
