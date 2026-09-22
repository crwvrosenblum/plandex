package shared

import "testing"

func TestLineNumbersPreserveFileBytes(t *testing.T) {
	cases := []struct {
		name, input, expected string
	}{
		{"empty", "", ""},
		{"no final newline", "alpha", "pdx-1: alpha"},
		{"one final newline", "alpha\n", "pdx-1: alpha\n"},
		{"intentional blank line", "alpha\n\n", "pdx-1: alpha\npdx-2: \n"},
		{"empty line", "\n", "pdx-1: \n"},
		{"leading blank line", "\nalpha\n", "pdx-1: \npdx-2: alpha\n"},
		{"CRLF", "alpha\r\nbeta\r\n", "pdx-1: alpha\r\npdx-2: beta\r\n"},
		{"unicode without final newline", "café\n世界", "pdx-1: café\npdx-2: 世界"},
		{"configuration file", "project=example\nstatus=pending\nkeep=unchanged\n", "pdx-1: project=example\npdx-2: status=pending\npdx-3: keep=unchanged\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			numbered := AddLineNums(tc.input)
			if string(numbered) != tc.expected {
				t.Errorf("numbered text = %q, want %q", numbered, tc.expected)
			}
			if restored := RemoveLineNums(numbered); restored != tc.input {
				t.Errorf("round trip = %q, want exact original %q", restored, tc.input)
			}
			proposed := AddLineNumsWithPrefix(tc.input, "pdx-new-")
			if restored := RemoveLineNumsWithPrefix(proposed, "pdx-new-"); restored != tc.input {
				t.Errorf("proposed round trip = %q, want %q", restored, tc.input)
			}
		})
	}
}
