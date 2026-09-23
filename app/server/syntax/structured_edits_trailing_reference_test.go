package syntax

import (
	"context"
	"testing"
)

func TestApplyChangesSyntheticEndReferencePreservesNewlines(t *testing.T) {
	tests := []struct{ name, original, proposed, want string }{
		{"full file one newline", "first\nstatus=pending\nlast\n", "first\nstatus=verified\nlast\n", "first\nstatus=verified\nlast\n"},
		{"full file two newlines", "first\nstatus=pending\nlast\n\n", "first\nstatus=verified\nlast\n\n", "first\nstatus=verified\nlast\n\n"},
		{"full file three newlines", "first\nstatus=pending\nlast\n\n\n", "first\nstatus=verified\nlast\n\n\n", "first\nstatus=verified\nlast\n\n\n"},
		{"partial proposal preserves remaining source", "first\nstatus=pending\nanchor\nlast\n", "first\nstatus=verified\nanchor\n", "first\nstatus=verified\nanchor\nlast\n"},
		{"partial proposal preserves real blank before remainder", "first\nstatus=pending\nanchor\n\nlast\n", "first\nstatus=verified\nanchor\n", "first\nstatus=verified\nanchor\n\nlast\n"},
		{"full file CRLF", "first\r\nstatus=pending\r\nlast\r\n", "first\r\nstatus=verified\r\nlast\r\n", "first\r\nstatus=verified\r\nlast\r\n"},
		{"partial proposal preserves two real blanks", "first\nstatus=pending\nanchor\n\n\nlast\n", "first\nstatus=verified\nanchor\n", "first\nstatus=verified\nanchor\n\n\nlast\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ApplyChanges(context.Background(), ApplyChangesParams{Original: tt.original, Proposed: tt.proposed, Desc: "Type: replace\nReplace: line 2", AddMissingStartEndRefs: true})
			if res.NewFile != tt.want {
				t.Fatalf("got %q\nwant %q", res.NewFile, tt.want)
			}
		})
	}
}
