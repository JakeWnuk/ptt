// Package mask handles the transformation of a string into a masks
package mask

import (
	"fmt"
	"ppt/pkg/utils"
	"strings"
)

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
func MakeMaskedMap(input map[string]int, replacements []string, verbose bool) map[string]int {
	maskedMap := make(map[string]int)
	replacer := strings.NewReplacer(replacements...)

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) {
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
