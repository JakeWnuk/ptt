// Package mask provides functions to handle masking and unmasking of strings
package mask

import (
	"fmt"
	"os"
	"ptt/pkg/models"
	"ptt/pkg/validation"
	"strings"
	"unicode"
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

// MakeMask transforms the input string into a mask using the global mask
// defined in the models package.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a mask.
func MakeMask(key string) string {
	newKey := models.MaskReplacer.Replace(key)

	if !validation.CheckASCIIString(newKey) && strings.Contains(models.GlobalMask, "b") {
		newKey = validation.ConvertMultiByteMask(newKey)
	}

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.MakeMask(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", models.GlobalMask)
		fmt.Fprintf(os.Stderr, "Return: %s\n", newKey)
	}

	return newKey
}

// RemoveMaskedCharacters removes masked characters from the input map
// and returns a new map
//
// Args:
// key (string): Masked key to be replaced
//
// Returns:
// (string): Replaced key
func RemoveMaskedCharacters(key string) string {
	replacer := strings.NewReplacer("?u", "", "?l", "", "?d", "", "?b", "", "?s", "")
	newKey := replacer.Replace(key)

	if !validation.CheckASCIIString(newKey) && strings.Contains(models.GlobalMask, "b") {
		newKey = validation.ConvertMultiByteMask(newKey)
	}

	return newKey
}

// TestMaskComplexity tests the complexity of an input full mask or a partial
// mask string and returns a score
//
// Args:
//
//	str (string): Input string to test
//
// Returns:
//
//	(int): Complexity score as an integer
func TestMaskComplexity(str string) int {
	score := 0
	lowerBool := false
	upperBool := false
	digitBool := false
	specialBool := false
	byteBool := false
	for i := 0; i < len(str); i++ {
		if str[i] == '?' {
			if i+1 < len(str) {
				switch str[i+1] {
				case 'l':
					lowerBool = true
				case 'u':
					upperBool = true
				case 'd':
					digitBool = true
				case 's':
					specialBool = true
				case 'b':
					byteBool = true
				}
			}
		} else {
			if unicode.IsLower(rune(str[i])) {
				lowerBool = true
			}
			if unicode.IsUpper(rune(str[i])) {
				upperBool = true
			}
			if unicode.IsDigit(rune(str[i])) {
				digitBool = true
			}
			if strings.ContainsRune(" !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~'", rune(str[i])) {
				specialBool = true
			}
			if str[i] > 127 {
				byteBool = true
			}
		}
	}
	if lowerBool {
		score++
	}
	if upperBool {
		score++
	}
	if digitBool {
		score++
	}
	if specialBool {
		score++
	}
	if byteBool {
		score++
	}
	return score
}
