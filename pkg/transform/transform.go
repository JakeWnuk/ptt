// Package transform contains logic for transforming input maps
package transform

import (
	"fmt"
	"os"

	"github.com/jakewnuk/ptt/pkg/format"
	"github.com/jakewnuk/ptt/pkg/mask"
	"github.com/jakewnuk/ptt/pkg/rule"
	"github.com/jakewnuk/ptt/pkg/utils"
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
//	input (map[string]int): A map of input values
//	mode (string): The mode to run the CLI in
//	startingIndex (int): The starting index for the transformation if applicable
//	endIndex (int): The ending index for the transformation if applicable
//	verbose (bool): If true, the verbose information is printed when available
//	replacementMask (string): The mask characters to use for masking operations
//	transformationFilesMap (map[string]int): A map of transformation files to
//	use for modes like retain-mask
//	bypass (bool): If true, the map is not used for output or filtering
//
// Returns:
//
//	(map[string]int): A map of transformed values
func TransformationController(input map[string]int, mode string, startingIndex int, endingIndex int, verbose bool, replacementMask string, transformationFilesMap map[string]int, bypass bool) (output map[string]int) {
	switch mode {
	case "append", "append-remove", "append-shift", "a":
		output = rule.AppendRules(input, mode, bypass)
	case "prepend", "prepend-remove", "prepend-shift", "ar":
		output = rule.PrependRules(input, mode)
	case "insert", "i":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.InsertRules(input, strIndex, endIndex)
	case "overwrite", "o":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.OverwriteRules(input, strIndex, endIndex)
	case "toggle", "t":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.ToggleRules(input, strIndex, endIndex)
	case "encode", "e":
		output = format.EncodeInputMap(input)
	case "decode", "de":
		output = format.DecodeInputMap(input)
	case "mask", "partial-mask", "partial", "m":
		output = mask.MakeMaskedMap(input, replacementMask, verbose)
	case "dehex", "unhex", "dh":
		output = format.DehexMap(input)
	case "hex", "rehex":
		output = format.HexEncodeMap(input)
	case "remove", "remove-all", "delete", "delete-all", "rm":
		input = mask.MakeMaskedMap(input, replacementMask, false)
		output = mask.RemoveMaskedCharacters(input)
	case "retain-mask", "retain", "r":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Retain masks require use of one or more -tf flags to specify one or more files")
			os.Exit(1)
		}
		output = mask.MakeRetainMaskedMap(input, replacementMask, transformationFilesMap)
	case "match-mask", "match", "mt":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Match masks require use of one or more -tf flags to specify one or more files")
			os.Exit(1)
		}
		output = mask.MakeMatchedMaskedMap(input, replacementMask, transformationFilesMap)
	case "swap", "replace", "s":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Swap operations require use of one or more -tf flags to specify one or more files")
			fmt.Println("This transformation mode requres a ':' separated list of keys to swap")
			os.Exit(1)
		}
		output = ReplaceKeysInMap(input, transformationFilesMap)
	case "pop", "split", "boundary-split", "boundary-pop", "pop-split", "split-pop", "po":
		output = mask.BoundarySplitPopMap(input, replacementMask)
	case "mask-swap", "shuffle", "shuf", "token-swap", "ms":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Mask-swap operations require use of one or more -tf flags to specify one or more files")
			fmt.Println("This transformation mode requres a retain mask file to use for swapping")
			os.Exit(1)
		}
		output = mask.ShuffleMap(input, replacementMask, transformationFilesMap)
	default:
		output = input
	}

	return output
}

// ----------------------------------------------------------------------------
// Generation Functions
// ----------------------------------------------------------------------------

// ReplaceKeysInMap takes a map of keys and values and replaces the keys
// with replacements based on the replacement map. This is useful for
// exact key swaps.
//
// Args:
//
//	originalMap (map[string]int): The original map to replace keys in
//	replacements (map[string]int): The map of replacements to use
//
// Returns:
//
//	(map[string]int): A new map with the keys replaced
func ReplaceKeysInMap(originalMap map[string]int, replacements map[string]int) map[string]int {
	newMap := make(map[string]int)
	for key, value := range originalMap {
		newKeyArray := utils.ReplaceSubstring(key, replacements)
		for _, newKey := range newKeyArray {
			newMap[newKey] = value
		}
	}
	return newMap
}
