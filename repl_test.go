package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello. world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Go is great!",
			expected: []string{"go", "is", "great"},
		},
		{
			input:    " 123, 456; 789? ",
			expected: []string{"123", "456", "789"},
		},
		{
			input:    "Mixed CASE Words.",
			expected: []string{"mixed", "case", "words"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length of cleanInput(%q) = %d; \n length expected: %d", c.input, len(actual), len(c.expected))
			break
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q)[%d] = %q; \n expected: %q", c.input, i, word, expectedWord)
				break
			}
		}
	}
}
