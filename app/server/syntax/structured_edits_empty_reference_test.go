package syntax

import (
	"context"
	"testing"
)

func TestApplyChangesEmptyReferencePreservesFileBytes(t *testing.T) {
	const original = "project=coreweave-harness-check\nstatus=pending\nkeep=unchanged\n"
	const updated = "project=coreweave-harness-check\nstatus=verified\nkeep=unchanged\n"
	// Captured native description: this also preserves the end of the file.
	const desc = "The current task is to update `status.txt`, replacing only the `status=pending` line with `status=verified`. The loaded file contains three lines — `project=coreweave-harness-check`, `status=pending`, and `keep=unchanged` — followed by a final newline. My approach is a single in-place text edit applied directly as a file update; since no commands are to be executed, `_apply.sh` requires no commands. The `project=coreweave-harness-check` and `keep=unchanged` lines are preserved exactly as-is, and the trailing newline at the end of the file is preserved.\n\n**Updating `status.txt`**\nType: replace\nSummary: Replace the `status=pending` line with `status=verified`\nReplace: line 2\nContext: Located between the `project=coreweave-harness-check` line and the `keep=unchanged` line"
	tests := []struct {
		name, original, proposed, desc, want string
	}{
		{
			name:     "captured native full-file replacement",
			original: original, proposed: updated, desc: desc, want: updated,
		},
		{
			name:     "preserve real leading blank from original reference",
			original: "\n" + original, proposed: updated,
			desc: "Type: replace\nReplace: line 3\nPreserve the end of the file", want: "\n" + updated,
		},
		{
			name:     "preserve intentional leading blank in proposal",
			original: "\n" + original, proposed: "\n" + updated,
			desc: "Type: replace\nReplace: line 3\nPreserve the end of the file", want: "\n" + updated,
		},
		{
			name:     "empty interior reference does not insert a blank",
			original: "alpha\nbeta\nstatus=pending\nomega\n",
			proposed: "alpha\n// ... existing code ...\nbeta\nstatus=verified\nomega\n",
			desc:     "Type: replace\nReplace: line 3\nPreserve the start of the file and end of the file",
			want:     "alpha\nbeta\nstatus=verified\nomega\n",
		},
		{
			name:     "preserve real interior blank from original reference",
			original: "alpha\n\nbeta\nstatus=pending\nomega\n",
			proposed: "alpha\n// ... existing code ...\nbeta\nstatus=verified\nomega\n",
			desc:     "Type: replace\nReplace: line 4\nPreserve the start of the file and end of the file",
			want:     "alpha\n\nbeta\nstatus=verified\nomega\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyChanges(context.Background(), ApplyChangesParams{
				Original: tt.original, Proposed: tt.proposed, Desc: tt.desc,
				AddMissingStartEndRefs: true,
			})
			if result.NewFile != tt.want {
				t.Fatalf("file bytes changed beyond the requested edit:\n got %q\nwant %q", result.NewFile, tt.want)
			}
		})
	}
}
