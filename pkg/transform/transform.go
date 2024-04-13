package transform

import (
	"fmt"
	"os"
	"ppt/pkg/format"
	"ppt/pkg/mask"
	"ppt/pkg/rule"
	"ppt/pkg/utils"
)

// ----------------------------------------------------------------------------
// TransformationController
// ----------------------------------------------------------------------------

// TransformationController is the main entry point for the CLI
// application. Operates a switch statement to determine the
// mode to use.
//
// Args:
//
//      input (map[string]int): A map of input values
//      mode (string): The mode to run the CLI in
//      startingIndex (int): The starting index for the transformation if applicable
//      verbose (bool): If true, the verbose information is printed when available
//      replacementMask (string): The mask characters to use for masking operations
//      transformationFilesMap (map[string]int): A map of transformation files to
//      use for modes like retain-mask
//      swapTolerance (int): The swap tolerance for substitution operations. This
//      will allow for flexibility in the number of characters that can be swapped.
//      For example, a swap tolerance of 1 will allow for one character off to
//      be considered a match.
//
// Returns:
//
//      (map[string]int): A map of transformed values

func TransformationController(input map[string]int, mode string, startingIndex int, verbose bool, replacementMask string, transformationFilesMap map[string]int, swapTolerance int) (output map[string]int) {
	strIndex := fmt.Sprintf("%d", startingIndex)
	switch mode {
	case "append", "append-remove", "append-shift":
		output = rule.AppendRules(input, mode)
	case "prepend", "prepend-remove", "prepend-shift":
		output = rule.PrependRules(input, mode)
	case "insert":
		output = rule.InsertRules(input, strIndex)
	case "overwrite":
		output = rule.OverwriteRules(input, strIndex)
	case "toggle":
		output = rule.ToggleRules(input, strIndex)
	case "encode":
		output = format.EncodeInputMap(input)
	case "mask", "partial-mask", "partial":
		output = mask.MakeMaskedMap(input, replacementMask, verbose)
	case "dehex", "unhex":
		output = format.DehexMap(input)
	case "hex", "rehex":
		output = format.HexEncodeMap(input)
	case "remove", "remove-all", "delete", "delete-all":
		input = mask.MakeMaskedMap(input, replacementMask, verbose)
		output = mask.RemoveMaskedCharacters(input)
	case "retain-mask", "retain":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Retain masks require use of one or more -tf flags to specify one or more files")
			os.Exit(1)
		}
		output = mask.MakeRetainMaskedMap(input, replacementMask, transformationFilesMap)
	case "match-mask", "match":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Match masks require use of one or more -tf flags to specify one or more files")
			os.Exit(1)
		}
		output = mask.MakeMatchedMaskedMap(input, replacementMask, transformationFilesMap)
	case "fuzzy-swap", "fuzzy-replace", "fuzzy", "fuzz", "mutate":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Swap operations require use of one or more -tf flags to specify one or more files")
			os.Exit(1)
		}
		output = FuzzyReplaceKeysInMap(input, transformationFilesMap)
	}

	return output
}

// ----------------------------------------------------------------------------
// Generation Functions
// ----------------------------------------------------------------------------

// FuzzyReplaceKeysInMap takes a map of keys and values and replaces the keys
// with fuzzy replacements based on the replacement map. This is useful for
// performing operations like mutation or swapping.
//
// Args:
//
//	originalMap (map[string]int): The original map to replace keys in
//	replacements (map[string]int): The map of replacements to use
//
// Returns:
//
//	(map[string]int): A new map with the keys replaced
func FuzzyReplaceKeysInMap(originalMap map[string]int, replacements map[string]int) map[string]int {
	newMap := make(map[string]int)
	for key, value := range originalMap {
		newKeyArray := utils.FuzzyReplaceSubstring(key, replacements)
		for _, newKey := range newKeyArray {
			newMap[newKey] = value
		}
	}
	return newMap
}
