package plan

import shared "plandex-shared"

// numberBuildPromptFiles keeps file numbering consistent across build prompts.
func numberBuildPromptFiles(original, proposed string) (shared.LineNumberedTextType, shared.LineNumberedTextType) {
	return shared.AddLineNums(original), shared.AddLineNumsWithPrefix(proposed, "pdx-new-")
}
