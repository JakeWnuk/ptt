package mask

import (
	"ptt/pkg/models"
	"ptt/pkg/validation"
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
