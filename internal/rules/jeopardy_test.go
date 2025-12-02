package rules

import "testing"

func TestJeopardyRegex(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"What is 2 + 2?", true},
		{"what is 3 times 4?", true},
		{"What is 10 divided by 2?", true},
		{"What is 5 minus 3?", true},
		{"What is 1 plus 2 minus 3?", true},
		{"What is 4 times 5 divided by 2?", true},
		{"What is 2 + 2", true},
		{"What is -2 + 3?", true},
		{"What is -1 + 8?", true},
		{"what is 3+-4", true},
		{"what is 3--4", true},
	}

	for _, c := range cases {
		if !jeopardyRegex.MatchString(c.input) {
			t.Errorf("jeopardyRegex did not match input: %q", c.input)
		}
	}
}
