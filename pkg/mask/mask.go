// Package mask handles the transformation of a string into a masks
package mask

import (
	"fmt"
	"ptt/pkg/utils"
	"strings"
)

// ----------------------------------------------------------------------------
// Mask Generation Functions
// ----------------------------------------------------------------------------

// ConstructReplacements create an array mapping which characters to replace
//
// This function accepts the characters "ulds" in order to generate a map
// - u for uppercase characters
// - l for lowercase characters
// - d for numerical characters
// - s for special characters
//
// Args:
//
//	str (string): Input string
//
// Returns:
//
//	args ([]string): Map of replacement characters
func ConstructReplacements(str string) []string {
	var lowerArgs, upperArgs, digitArgs, args []string
	for c := 'a'; c <= 'z'; c++ {
		lowerArgs = append(lowerArgs, string(c), "?l")
	}
	for c := 'A'; c <= 'Z'; c++ {
		upperArgs = append(upperArgs, string(c), "?u")
	}
	for c := '0'; c <= '9'; c++ {
		digitArgs = append(digitArgs, string(c), "?d")
	}
	specialChars := " !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~'"
	specialArgs := make([]string, len(specialChars)*2)
	for i, c := range specialChars {
		specialArgs[i*2] = string(c)
		specialArgs[i*2+1] = "?s"
	}

	if strings.Contains(str, "l") {
		args = append(args, lowerArgs...)
	}

	if strings.Contains(str, "u") {
		args = append(args, upperArgs...)
	}

	if strings.Contains(str, "d") {
		args = append(args, digitArgs...)
	}

	if strings.Contains(str, "s") {
		args = append(args, specialArgs...)
	}

	return args
}

// MakeMaskedMap replaces all characters in the input maps key with the values
// in the input map
//
// Args:
//
//	input (map[string]int): Map to mask
//	replacements ([]string): Array of characters to replace
//	verbose (bool): Verbose information if true
//
// Returns:
//
// maskedMap (map[string]int): Masked map
func MakeMaskedMap(input map[string]int, replacementMask string, verbose bool) map[string]int {
	maskedMap := make(map[string]int)
	replacements := ConstructReplacements(replacementMask)
	replacer := strings.NewReplacer(replacements...)

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) && strings.Contains(replacementMask, "b") {
			newKey = ConvertMultiByteMask(newKey)
		}

		if verbose {
			newKey = fmt.Sprintf("%s:%d:%d\n", newKey, len(key), TestMaskComplexity(newKey))
		}

		if oldValue, exists := maskedMap[newKey]; exists {
			maskedMap[newKey] = oldValue + value
		} else {
			maskedMap[newKey] = value
		}
	}
	return maskedMap
}

// MakeRetainMaskedMap replaces all characters in the input maps key with mask
// values in the input map but retains keywords provided in the retain list
//
// Args:
//
//	input (map[string]int): Map to mask
//	replacementMask (string): Mask characters to apply
//	retain (map[string]int): Map of keywords to retain
//
// Returns:
//
//	maskedMap (map[string]int): Masked retain map
func MakeRetainMaskedMap(input map[string]int, replacementMask string, retain map[string]int) map[string]int {
	maskedMap := make(map[string]int)
	replacements := ConstructReplacements(replacementMask)
	replacer := strings.NewReplacer(replacements...)

	for key, value := range input {
		for retainKey := range retain {
			newKey := ""
			if strings.Contains(key, retainKey) {
				parts := utils.SplitBySeparatorString(key, retainKey)

				// if the part is not the key replace it using replacer
				for _, part := range parts {
					if part != retainKey {
						newPart := replacer.Replace(part)
						if !utils.CheckASCIIString(newPart) && strings.Contains(replacementMask, "b") {
							newPart = ConvertMultiByteMask(newPart)
						}
						newKey += newPart
					} else {
						newKey += part
					}
				}

			} else {
				// if the key is not in the string continue
				continue
			}

			if oldValue, exists := maskedMap[newKey]; exists {
				maskedMap[newKey] = oldValue + value
			} else {
				maskedMap[newKey] = value
			}
		}
	}
	return maskedMap
}

// MakeMaskedString replaces all characters in the input string with the mask
// values in the mask string
//
// Args:
//
//	input (string): Input string
//	replacementMask (string): Mask characters to apply
//
// Returns:
//
//	(string): Masked string
func MakeMaskedString(input string, replacementMask string) string {
	replacements := ConstructReplacements(replacementMask)
	replacer := strings.NewReplacer(replacements...)
	newKey := replacer.Replace(input)

	if !utils.CheckASCIIString(newKey) && strings.Contains(replacementMask, "b") {
		newKey = ConvertMultiByteMask(newKey)
	}

	return newKey
}

// ----------------------------------------------------------------------------
// Mask Conversion Functions
// ----------------------------------------------------------------------------

// ConvertMultiByteMask converts non-ascii characters to a valid format
//
// Args:
//
//	str (string): Input string
//
// Returns:
//
//	returnStr (string): Converted string
func ConvertMultiByteMask(str string) string {
	returnStr := ""
	for _, r := range str {
		if r > 127 {
			byteArr := []byte(string(r))
			for j := range byteArr {
				if j == len(byteArr)-1 {
					returnStr += fmt.Sprintf("?b")
				} else {
					returnStr += fmt.Sprintf("?b")
				}
			}
		} else {
			returnStr += fmt.Sprintf("%c", r)
		}
	}
	return returnStr
}

// TestMaskComplexity tests the complexity of an input mask
//
// Args:
//
//	str (string): Input string to test
//
// Returns:
//
//	(int): Complexity score as an integer
func TestMaskComplexity(str string) int {
	complexity := 0
	charTypes := []string{"?u", "?l", "?d", "?s", "?b"}
	for _, charType := range charTypes {
		if strings.Contains(str, charType) {
			complexity++
		}
	}
	return complexity
}

// RemoveMaskedCharacters removes masked characters from the input map
// and returns a new map
//
// Args:
//
//	input (map[string]int): Input map
//
// Returns:
//
//	(map[string]int): Masked map
func RemoveMaskedCharacters(input map[string]int) map[string]int {
	maskedMap := make(map[string]int)
	replacer := strings.NewReplacer("?u", "", "?l", "", "?d", "", "?b", "", "?s", "")

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) {
			newKey = ConvertMultiByteMask(newKey)
		}

		if oldValue, exists := maskedMap[newKey]; exists {
			maskedMap[newKey] = oldValue + value
		} else {
			maskedMap[newKey] = value
		}
	}
	return maskedMap
}

// ----------------------------------------------------------------------------
// Mask Utility Functions
// ----------------------------------------------------------------------------

// MakeMatchedMaskedMap returns a map from the input map where the keys matched
// the keys in the mask map after applying the mask to the input map. The
// original keys and values are retained.
//
// Args:
//
//	input (map[string]int): Input map
//	replacementMask (string): Mask characters to apply
//	maskMap (map[string]int): Mask map
//
// Returns:
// (map[string]int): Matched masked map
func MakeMatchedMaskedMap(input map[string]int, replacementMask string, maskMap map[string]int) map[string]int {
	maskedMap := make(map[string]int)
	replacements := ConstructReplacements(replacementMask)
	replacer := strings.NewReplacer(replacements...)

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) && strings.Contains(replacementMask, "b") {
			newKey = ConvertMultiByteMask(newKey)
		}

		if _, exists := maskMap[newKey]; exists {
			if oldValue, exists := maskedMap[newKey]; exists {
				maskedMap[newKey] = oldValue + value
			} else {
				maskedMap[newKey] = value
			}
		}
	}
	return maskedMap
}
