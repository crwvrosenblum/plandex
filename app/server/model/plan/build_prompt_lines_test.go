package plan

import (
	"fmt"
	"strings"
	"testing"

	"plandex-server/model/prompts"
)

func TestBuildPromptFileNumbering(t *testing.T) {
	for _, tc := range []struct {
		name, original, proposed string
	}{
		{"edit", "status=pending\nkeep=unchanged", "status=complete\nkeep=unchanged"},
		{"reference comment", "func main() {\n    original()\n}", "func main() {\n    // ... existing code ...\n    added()\n}"},
		{"literal prefix", "pdx-1: literal", "pdx-new-1: literal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original, proposed := numberBuildPromptFiles(tc.original, tc.proposed)
			whole, _ := prompts.GetWholeFilePrompt("fixture.txt", original, proposed, "Apply the proposed edit", "")
			validation, _ := prompts.GetValidationReplacementsXmlPrompt(prompts.ValidationPromptParams{
				Path:                 "fixture.txt",
				OriginalWithLineNums: original,
				ProposedWithLineNums: proposed,
				Desc:                 "Apply the proposed edit",
				Diff:                 "fixture diff",
			})

			for name, prompt := range map[string]string{"whole file": whole, "validation": validation} {
				t.Run(name, func(t *testing.T) {
					assertNumberedPromptFile(t, prompt, "Original file", "pdx-", tc.original)
					assertNumberedPromptFile(t, prompt, "Proposed changes", "pdx-new-", tc.proposed)
				})
			}
		})
	}
}

// Check the model-facing section label together with the actual numbered lines.
// Files have no final newline here: newline preservation is a separate concern.
func assertNumberedPromptFile(t *testing.T, prompt, section, prefix, original string) {
	t.Helper()
	heading := fmt.Sprintf("%s (with line nums prefixed with '%s'):\n>>>\n", section, prefix)
	_, rest, ok := strings.Cut(prompt, heading)
	if !ok {
		t.Fatalf("missing %s section heading", section)
	}
	body, _, ok := strings.Cut(rest, "\n<<<")
	if !ok {
		t.Fatalf("missing %s section end", section)
	}
	lines := strings.Split(strings.TrimSuffix(body, "\n"), "\n")
	wantLines := strings.Split(original, "\n")
	if len(lines) != len(wantLines) {
		t.Fatalf("%s line count = %d, want %d", section, len(lines), len(wantLines))
	}
	for i, line := range lines {
		want := fmt.Sprintf("%s%d: %s", prefix, i+1, wantLines[i])
		if line != want {
			t.Errorf("%s line %d = %q, want %q", section, i+1, line, want)
		}
	}
}
