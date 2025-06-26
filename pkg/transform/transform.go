// Package transform contains logic for transforming input maps
package transform

import (
	"fmt"
	"os"
	"strings"

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
//	debug (int): Different debug levels to use for debugging [0-2]
//	wordRangeStart (int): The starting range for word operations
//	wordRangeEnd (int): The ending range for word operations
//
// Returns:
//
//	(map[string]int): A map of transformed values
func TransformationController(input map[string]int, mode string, startingIndex int, endingIndex int, verbose bool, replacementMask string, transformationFilesMap map[string]int, bypass bool, debug int, wordRangeStart int, wordRangeEnd int) (output map[string]int) {

	functionDebug := false
	if debug > 1 {
		functionDebug = true
	}

	if debug > 0 {
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Starting debug mode:\n")
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Running in mode %s.\n", mode)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Starting index is %d.\n", startingIndex)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Ending index is %d.\n", endingIndex)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Replacement mask is %s.\n", replacementMask)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Bypass is %t.\n", bypass)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Verbose is %t.\n", verbose)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Transformation files map is %v.\n", transformationFilesMap)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Input map is %v.\n", input)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Starting transformation...\n")
	}

	switch mode {
	case "rule-append", "rule-append-remove", "append", "append-remove":
		output = rule.AppendRules(input, mode, bypass, functionDebug)
	case "rule-prepend", "rule-prepend-remove", "rule-prepend-toggle", "prepend", "prepend-remove", "prepend-toggle":
		output = rule.PrependRules(input, mode, bypass, functionDebug)
	case "rule-insert", "insert":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.InsertRules(input, strIndex, endIndex, bypass, functionDebug)
	case "rule-overwrite", "overwrite":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.OverwriteRules(input, strIndex, endIndex, bypass, functionDebug)
	case "rule-toggle", "toggle":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.ToggleRules(input, strIndex, endIndex, bypass, functionDebug)
	case "encode":
		output = format.EncodeInputMap(input, bypass, functionDebug)
	case "decode":
		output = format.DecodeInputMap(input, bypass, functionDebug)
	case "mask":
		output = mask.MakeMaskedMap(input, replacementMask, verbose, bypass, functionDebug)
	case "dehex":
		output = format.DehexMap(input, bypass, functionDebug)
	case "hex":
		output = format.HexEncodeMap(input, bypass, functionDebug)
	case "mask-remove", "remove":
		input = mask.MakeMaskedMap(input, replacementMask, false, false, false)
		output = mask.RemoveMaskedCharacters(input, replacementMask, bypass, functionDebug)
	case "mask-retain", "retain":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Retain masks require use of one or more -tf flags to specify one or more files.\n")
			os.Exit(1)
		}
		output = mask.MakeRetainMaskedMap(input, replacementMask, transformationFilesMap, bypass, functionDebug, verbose)
	case "mask-match", "match":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Match masks require use of one or more -tf flags to specify one or more files.\n")
			os.Exit(1)
		}
		output = mask.MakeMatchedMaskedMap(input, replacementMask, transformationFilesMap, bypass, functionDebug)
	case "swap", "swap-single":
		fmt.Fprintf(os.Stderr, "[*] This transformation mode requires a ':' separated list of keys to swap.\n")
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Swap operations require use of one or more -tf flags to specify one or more files.\n")
			os.Exit(1)
		}
		output = ReplaceKeysInMap(input, transformationFilesMap, bypass, functionDebug)
	case "mask-pop", "pop":
		output = mask.BoundarySplitPopMap(input, replacementMask, bypass, functionDebug)
	case "mask-swap":
		fmt.Fprintf(os.Stderr, "[*] This transformation mode requires a retain mask file to use for swapping.\n")
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Mask-swap operations require use of one or more -tf flags to specify one or more files.\n")
			os.Exit(1)
		}
		output = mask.ShuffleMap(input, replacementMask, transformationFilesMap, bypass, functionDebug)
	case "passphrase":
		fmt.Fprintf(os.Stderr, "[*] This transformation mode expects space separated content.\n")
		if wordRangeStart == 0 {
			fmt.Fprintf(os.Stderr, "[!] Passphrase operations require use of the -w flag to specify the number of words to use.\n")
			os.Exit(1)
		}
		output = MakePassphraseMap(input, bypass, functionDebug, wordRangeStart, wordRangeEnd)
	case "substring":
		output = utils.SubstringMap(input, startingIndex, endingIndex, bypass, functionDebug)
	case "replace-all", "replace":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Replace operations require use of one or more -tf flags to specify one or more files.\n")
			os.Exit(1)
		}
		output = ReplaceAllKeysInMap(input, transformationFilesMap, bypass, functionDebug)
	case "regram":
		fmt.Fprintf(os.Stderr, "[*] This transformation mode expects space separated content.\n")
		if wordRangeStart == 0 {
			fmt.Fprintf(os.Stderr, "[!] Regram operations require use of the -w flag to specify the number of words to use.\n")
			os.Exit(1)
		}
		output = GenerateNGramMap(input, wordRangeStart, wordRangeEnd, bypass, functionDebug)
	case "rule-apply", "apply":
		fmt.Fprintf(os.Stderr, "[*] This transformation mode expects a rule file to apply.\n")
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Apply operations require use of one or more -tf flags to specify one or more files.\n")
			os.Exit(1)
		}
		output = rule.ApplyRulesHCRE(input, transformationFilesMap, bypass, functionDebug)
	case "rule-simplify", "simplify":
		fmt.Fprintf(os.Stderr, "[*] This transformation mode expects rule input to simplify.\n")
		output = rule.SimplifyRules(input, bypass, functionDebug)
	default:
		output = input
	}

	if debug > 0 {
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Output map is %v.\n", output)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Transformation complete. Resuming output.\n")
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	(map[string]int): A new map with the keys replaced
func ReplaceKeysInMap(originalMap map[string]int, replacements map[string]int, bypass bool, debug bool) map[string]int {
	newMap := make(map[string]int)
	for key, value := range originalMap {
		newKeyArray := utils.ReplaceSubstring(key, replacements)
		for _, newKey := range newKeyArray {

			if debug {
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
			}

			if !bypass {
				newMap[newKey] = value
			} else {
				fmt.Println(newKey)
			}
		}
	}
	return newMap
}

// ReplaceAllKeysInMap takes a map of keys and values and replaces the keys
// with replacements based on the replacement map. This is useful for
// replacing all instances of a key with a new key.
//
// Args:
//
//	originalMap (map[string]int): The original map to replace keys in
//	replacements (map[string]int): The map of replacements to use
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
//	Returns:
//
//	(map[string]int): A new map with the keys replaced
func ReplaceAllKeysInMap(originalMap map[string]int, replacements map[string]int, bypass bool, debug bool) map[string]int {
	newMap := make(map[string]int)
	for key, value := range originalMap {
		newKeyArray := utils.ReplaceAllSubstring(key, replacements)
		for _, newKey := range newKeyArray {

			if debug {
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
			}

			if !bypass {
				newMap[newKey] = value
			} else {
				fmt.Println(newKey)
			}
		}
	}
	return newMap
}

// MakePassphraseMap takes a map of keys and creates a new map with new
// passphrases for each key.
//
// Args:
//
//	input (map[string]int): The original map to replace keys in
//	use for constructing the passphrases
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//	wordRangeStart (int): The starting number of words to use for passphrases
//	wordRangeEnd (int): The ending iteration number of words to use for passphrases
//
// Returns:
//
//	(map[string]int): A new map with the keys replaced
func MakePassphraseMap(input map[string]int, bypass bool, debug bool, wordRangeStart int, wordRangeEnd int) map[string]int {
	newMap := make(map[string]int)
	for key, value := range input {

		for i := wordRangeStart; i <= wordRangeEnd; i++ {
			newKeyArray := utils.GeneratePassphrase(key, i)
			for _, newKey := range newKeyArray {

				if debug {
					fmt.Fprintf(os.Stderr, "Key: %s\n", key)
					fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
				}

				if !bypass {
					if newMap[newKey] == 0 {
						newMap[newKey] = value
					} else {
						newMap[newKey] += value
					}
				} else {
					fmt.Println(newKey)
				}
			}
		}
	}

	return newMap
}

// GenerateNGramMap takes a map of keys and values and generates a new map
// using the utils.GenerateNGrams function and combines the results. This
// function is used to generate n-grams from the input map for the regram
// transformation mode.
//
// Args:
//
//	input (map[string]int): The original map to generate n-grams from
//	wordRangeStart (int): The starting number of words to use for n-grams
//	wordRangeEnd (int): The ending iteration number of words to use for n-grams
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	(map[string]int): A new map with the n-grams generated
func GenerateNGramMap(input map[string]int, wordRangeStart int, wordRangeEnd int, bypass bool, debug bool) map[string]int {
	newMap := make(map[string]int)
	for key, value := range input {
		for i := wordRangeStart; i <= wordRangeEnd; i++ {
			newKeyArray := utils.GenerateNGrams(key, i)
			for _, newKey := range newKeyArray {

				if debug {
					fmt.Fprintf(os.Stderr, "Key: %s\n", key)
					fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
				}

				newKey = strings.TrimSpace(newKey)
				newKey = strings.TrimLeft(newKey, ",")
				newKey = strings.TrimRight(newKey, ",")
				newKey = strings.TrimLeft(newKey, " ")

				if !bypass {
					if newMap[newKey] == 0 {
						newMap[newKey] = value
					} else {
						newMap[newKey] += value
					}
				} else {
					fmt.Println(newKey)
				}
			}
		}
	}
	return newMap
}
