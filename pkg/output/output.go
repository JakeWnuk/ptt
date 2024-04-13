// Package output contains functions for printing and formatting output to the console
package output

import (
	"encoding/hex"
	"fmt"
	"html"
	"net/url"
	"os"
	"ppt/pkg/mask"
	"ppt/pkg/models"
	"ppt/pkg/rule"
	"sort"
	"strings"
)

// ----------------------------------------------------------------------------
// Output Controllers
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
//	verbose (bool): If true, the verbose information is printed when available
//	replacementMask (string): The mask characters to use for masking operations
//	transformationFilesMap (map[string]int): A map of transformation files to
//	use for modes like retain-mask
//
// Returns:
//
//	(map[string]int): A map of transformed values
func TransformationController(input map[string]int, mode string, startingIndex int, verbose bool, replacementMask string, transformationFilesMap map[string]int) (output map[string]int) {
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
		output = EncodeInputMap(input)
	case "mask", "partial-mask", "partial":
		output = mask.MakeMaskedMap(input, replacementMask, verbose)
	case "dehex", "unhex":
		output = DehexMap(input)
	case "hex", "rehex":
		output = HexEncodeMap(input)
	case "remove", "remove-all", "delete", "delete-all":
		input = mask.MakeMaskedMap(input, replacementMask, verbose)
		output = mask.RemoveMaskedCharacters(input)
	case "retain-mask", "retain":
		if len(transformationFilesMap) == 0 {
			fmt.Println("Retain masks require use of one or more -tf flags to specify one or more files")
			os.Exit(1)
		}
		output = mask.MakeRetainMaskedMap(input, replacementMask, transformationFilesMap)
	}

	return output
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// PrintArrayToSTDOUT prints an array of items to stdout
//
// Args:
//
//	items (map[string]int): A map of item frequencies
//	verbose (bool): If true, the count is printed
//
// Returns:
//
//	None
func PrintArrayToSTDOUT(freq map[string]int, verbose bool) {
	p := make(models.PairList, len(freq))
	i := 0
	for k, v := range freq {
		p[i] = models.Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	for _, pair := range p {
		if verbose {
			fmt.Printf("%d %s\n", pair.Value, pair.Key)
		} else {
			fmt.Printf("%s\n", pair.Key)
		}
	}
}

// RetainRemove compares a string against a list of words to retain and remove
//
// Args:
//
//	text (map[string]int): The text to process
//	retainMap (map[string]int): A map of words to retain
//	removeMap (map[string]int): A map of words to remove
//
// Returns:
//
//	(string, error): The processed and filtered text, or an error if the text is not valid
func RetainRemove(textMap, retainMap, removeMap map[string]int) (map[string]int, error) {
	result := make(map[string]int)

	// If the retain and remove maps are empty, return the textMap as is
	if len(retainMap) == 0 && len(removeMap) == 0 {
		return textMap, nil
	}

	for k, v := range textMap {
		// If the key is in the retain map and not in the remove map, add it to the result
		if _, ok := retainMap[k]; ok || len(retainMap) == 0 {
			if _, ok := removeMap[k]; !ok || len(removeMap) == 0 {
				result[k] = v
			}
		}
	}

	// If the result map is empty, return an error
	if len(result) == 0 {
		return nil, fmt.Errorf("no keys in the text map are in the retain list")
	}

	return result, nil
}

// RemoveMinimumFrequency removes items from a map that are below a minimum frequency
// threshold and returns a new map
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//	min (int): The minimum frequency threshold
//
// Returns:
//
//	(map[string]int): A new map of item frequencies above the minimum threshold
func RemoveMinimumFrequency(freq map[string]int, min int) map[string]int {
	newFreq := make(map[string]int)
	for key, value := range freq {
		if value >= min {
			newFreq[key] = value
		}
	}
	return newFreq
}

// ----------------------------------------------------------------------------
// Encoding Functions
// ----------------------------------------------------------------------------

// EncodeInputMap will encode a map of strings to URL, HTML, and unicode escaped strings
// where possible and return a new map of encoded strings
//
// Args:
//
//	input (map[string]int): A map of input strings
//
// Returns:
//
//	(map[string]int): A new map of encoded strings
func EncodeInputMap(input map[string]int) map[string]int {
	output := make(map[string]int)
	for k, v := range input {
		urlEncoded, htmlEncoded, escapeEncoded := EncodeString(k)

		if urlEncoded != "" {
			output[urlEncoded] = v
		}

		if htmlEncoded != "" {
			output[htmlEncoded] = v
		}

		if escapeEncoded != "" {
			output[escapeEncoded] = v
		}
	}
	return output
}

// EncodeString is used to URL and HTML encode a string where possible this
// will only return the encoded string if it is different from the input string
//
// Args:
//
//	s (string): Input string
//
// Returns:
//
//	urlEncoded (string): Input string URL encoded
//	htmlEncoded (string): Input string HTML encoded
//	escapedEncoded (string): Input string unicode escaped encoded
func EncodeString(s string) (string, string, string) {
	urlEncoded := url.QueryEscape(s)
	htmlEncoded := html.EscapeString(s)
	escapedEncoded := ASCIIEscapeUnicode(s)

	if urlEncoded == s {
		urlEncoded = ""
	}

	if htmlEncoded == s {
		htmlEncoded = ""
	}

	if escapedEncoded == s {
		escapedEncoded = ""
	}

	return urlEncoded, htmlEncoded, escapedEncoded
}

// ASCIIEscapeUnicode will convert a string into an unicode escaped format
//
// Args:
//
//	str (string): String to escape
//
// Returns:
//
//	escapedRunes (string): Converted runes in string format
func ASCIIEscapeUnicode(str string) string {
	runes := []rune(str)
	escapedRunes := make([]rune, 0, len(runes))

	for _, r := range runes {
		if r > 127 {
			// The rune is non-ASCII
			escapedRune := []rune(fmt.Sprintf("\\u%04x", r))
			escapedRunes = append(escapedRunes, escapedRune...)
		} else {
			escapedRunes = append(escapedRunes, r)
		}
	}
	return string(escapedRunes)
}

// DehexMap will decode a map of hex encoded strings and return a new map of decoded strings
// where possible
//
// Args:
//
//	input (map[string]int): A map of hex encoded strings
//
// Returns:
//
//	(map[string]int): A new map of decoded strings
func DehexMap(input map[string]int) map[string]int {
	decodedMap := make(map[string]int)

	for k, v := range input {
		k = strings.TrimPrefix(k, "$HEX[")
		k = strings.TrimSuffix(k, "]")
		decoded, err := hex.DecodeString(k)
		if err != nil {
			continue
		}
		decodedStr := string(decoded)
		decodedMap[decodedStr] = v
	}

	return decodedMap
}

// HexEncodeMap will encode a map of strings to hex encoded strings
// where possible and return a new map of encoded strings
//
// Args:
//
//	input (map[string]int): A map of input strings
//
// Returns:
//
//	(map[string]int): A new map of encoded strings
func HexEncodeMap(input map[string]int) map[string]int {
	output := make(map[string]int)
	for k, v := range input {
		encoded := hex.EncodeToString([]byte(k))
		output["$HEX["+encoded+"]"] = v
	}
	return output
}
