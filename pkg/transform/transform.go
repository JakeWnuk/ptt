// Package transform contains logic for transforming input maps
package transform

import (
	"fmt"
	"math/rand"
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
//	debug (int): Different debug levels to use for debugging [0-2]
//	passphraseWord (int): The number of words to use for passphrase generation
//
// Returns:
//
//	(map[string]int): A map of transformed values
func TransformationController(input map[string]int, mode string, startingIndex int, endingIndex int, verbose bool, replacementMask string, transformationFilesMap map[string]int, bypass bool, debug int, passphraseWord int) (output map[string]int) {

	functionDebug := false
	if debug > 1 {
		functionDebug = true
	}

	if debug > 0 {
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Starting debug mode:\n")
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Running in mode %s\n", mode)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Starting index is %d\n", startingIndex)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Ending index is %d\n", endingIndex)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Replacement mask is %s\n", replacementMask)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Bypass is %t\n", bypass)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Verbose is %t\n", verbose)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Transformation files map is %v\n", transformationFilesMap)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Input map is %v\n", input)
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Starting transformation...\n")
	}

	switch mode {
	case "append", "append-remove", "append-shift", "a":
		output = rule.AppendRules(input, mode, bypass, functionDebug)
	case "prepend", "prepend-remove", "prepend-shift", "ar":
		output = rule.PrependRules(input, mode, bypass, functionDebug)
	case "insert", "i":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.InsertRules(input, strIndex, endIndex, bypass, functionDebug)
	case "overwrite", "o":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.OverwriteRules(input, strIndex, endIndex, bypass, functionDebug)
	case "toggle", "t":
		strIndex := fmt.Sprintf("%d", startingIndex)
		endIndex := fmt.Sprintf("%d", endingIndex)
		output = rule.ToggleRules(input, strIndex, endIndex, bypass, functionDebug)
	case "encode", "e":
		output = format.EncodeInputMap(input, bypass, functionDebug)
	case "decode", "de":
		output = format.DecodeInputMap(input, bypass, functionDebug)
	case "mask", "partial-mask", "partial", "m":
		output = mask.MakeMaskedMap(input, replacementMask, verbose, bypass, functionDebug)
	case "dehex", "unhex", "dh":
		output = format.DehexMap(input, bypass, functionDebug)
	case "hex", "rehex":
		output = format.HexEncodeMap(input, bypass, functionDebug)
	case "remove", "remove-all", "delete", "delete-all", "rm":
		input = mask.MakeMaskedMap(input, replacementMask, false, false, false)
		output = mask.RemoveMaskedCharacters(input, bypass, functionDebug)
	case "retain-mask", "retain", "r", "mask-retain":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Retain masks require use of one or more -tf flags to specify one or more files\n")
			os.Exit(1)
		}
		output = mask.MakeRetainMaskedMap(input, replacementMask, transformationFilesMap, bypass, functionDebug)
	case "match-mask", "match", "mt", "mask-match":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Match masks require use of one or more -tf flags to specify one or more files\n")
			os.Exit(1)
		}
		output = mask.MakeMatchedMaskedMap(input, replacementMask, transformationFilesMap, bypass, functionDebug)
	case "swap", "replace", "s":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Swap operations require use of one or more -tf flags to specify one or more files\n")
			fmt.Fprintf(os.Stderr, "[!] This transformation mode requires a ':' separated list of keys to swap\n")
			os.Exit(1)
		}
		output = ReplaceKeysInMap(input, transformationFilesMap, bypass, functionDebug)
	case "pop", "split", "boundary-split", "boundary-pop", "pop-split", "split-pop", "po":
		output = mask.BoundarySplitPopMap(input, replacementMask, bypass, functionDebug)
	case "mask-swap", "shuffle", "shuf", "token-swap", "ms":
		if len(transformationFilesMap) == 0 {
			fmt.Fprintf(os.Stderr, "[!] Mask-swap operations require use of one or more -tf flags to specify one or more files")
			fmt.Fprintf(os.Stderr, "[!] This transformation mode requres a retain mask file to use for swapping")
			os.Exit(1)
		}
		output = mask.ShuffleMap(input, replacementMask, transformationFilesMap, bypass, functionDebug)
	case "passphrase", "phrase", "pp":
		if passphraseWord == 0 {
			fmt.Fprintf(os.Stderr, "[!] Passphrase operations require use of the -w flag to specify the number of words to use\n")
			os.Exit(1)
		}
		output = MakePassphraseMap(input, transformationFilesMap, bypass, functionDebug, passphraseWord)
	default:
		output = input
	}

	if debug > 0 {
		fmt.Fprintf(os.Stderr, "[?] TransformationController: Output map is %v\n", output)
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

// MakePassphraseMap takes a map of keys and creates a new map with new
// passphrases for each key. The transformation file is used to insert
// separators between the words. If the replacement mask is set to blank, then
// the words are concatenated together without any separators. Passphrases are
// generated by selecting a random word from the transformation file for each key.
//
// Args:
//
//	input (map[string]int): The original map to replace keys in
//	transformationFilesMap (map[string]int): A map of transformation files to
//	use for constructing the passphrases
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//	passphraseWord (int): The number of words to use for passphrase generation
//
// Returns:
//
//	(map[string]int): A new map with the keys replaced
func MakePassphraseMap(input map[string]int, transformationFilesMap map[string]int, bypass bool, debug bool, passphraseWord int) map[string]int {
	newMap := make(map[string]int)

	// should generate based on the value of the total sum of the input map not
	// each individual key. Repeat the process again for the value of each key

	for key, value := range input {

		for i := 0; i < value; i++ {
			newKeyPhrase := GeneratePassphrase(input, transformationFilesMap, passphraseWord)
			if debug {
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "New Phrase: %s\n", newKeyPhrase)
			}

			if !bypass {
				newMap[newKeyPhrase] = value
			} else {
				fmt.Println(newKeyPhrase)
			}
		}
	}
	return newMap
}

// GeneratePassphrase takes a key and a map of transformation files and
// generates a passphrase based on the number of words specified. The words
// are selected from the transformation files and concatenated together with
// a separator. If the replacement mask is set to blank, then the words are
// concatenated together without any separators.
//
// Args:
//
//	passWords (map[string]int): Content of the passphrase for use as words in
//	the passphrase
//	transformationFilesMap (map[string]int): Content of the transformation
//	files for use as separators between words
//	passphraseWord (int): The number of words to use for passphrase generation
//
// Returns:
//
//	(string): The generated passphrase
func GeneratePassphrase(passWords map[string]int, transformationFilesMap map[string]int, passphraseWord int) string {
	words := make([]string, passphraseWord)

	seps := make([]string, 0, len(transformationFilesMap))
	for k := range transformationFilesMap {
		seps = append(seps, k)
	}

	if len(seps) == 0 {
		seps = append(seps, "")
	}

	keys := make([]string, 0, len(passWords))
	for k := range passWords {
		keys = append(keys, k)
	}

	for i := 0; i < passphraseWord; i++ {
		sep := seps[rand.Intn(len(seps))]
		key := keys[rand.Intn(len(keys))]

		if i+1 >= passphraseWord {
			words[i] = fmt.Sprintf("%s%s", key, "")
		} else {
			words[i] = fmt.Sprintf("%s%s", key, sep)
		}

	}

	var newKeyPhrase string
	for _, word := range words {
		newKeyPhrase += word
	}

	return newKeyPhrase
}
