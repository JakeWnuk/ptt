// Package mask handles the transformation of a string into a masks
package mask

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/jakewnuk/ptt/pkg/utils"
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
// maskedMap (map[string]int): Masked map
func MakeMaskedMap(input map[string]int, replacementMask string, verbose bool, bypass bool, debug bool) map[string]int {
	maskedMap := make(map[string]int)
	replacements := ConstructReplacements(replacementMask)
	replacer := strings.NewReplacer(replacements...)

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) && strings.Contains(replacementMask, "b") {
			newKey = ConvertMultiByteMask(newKey)
		}

		if verbose {
			newKey = fmt.Sprintf("%s:%d:%d:%d", newKey, len(key), TestMaskComplexity(newKey), CalculateMaskKeyspace(newKey))
		}

		if debug {
			fmt.Fprintf(os.Stderr, "[?] MakeMaskedMap:\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", key)
			fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
			fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", replacementMask)
		}

		switch bypass {
		case false:

			if oldValue, exists := maskedMap[newKey]; exists {
				maskedMap[newKey] = oldValue + value
			} else {
				maskedMap[newKey] = value
			}

		case true:
			fmt.Println(newKey)
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//	verbose (bool): Verbose information if true
//
// Returns:
//
//	maskedMap (map[string]int): Masked retain map
func MakeRetainMaskedMap(input map[string]int, replacementMask string, retain map[string]int, bypass bool, debug bool, verbose bool) map[string]int {
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

			if verbose {
				newKey = fmt.Sprintf("%s:%d:%d:%d", newKey, len(key), TestMaskComplexity(newKey), CalculateMaskKeyspace(newKey))
			}

			if debug {
				fmt.Fprintf(os.Stderr, "[?] MakeRetainMaskedMap:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Retain Key: %s\n", retainKey)
				fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
				fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", replacementMask)
			}

			switch bypass {
			case false:
				if oldValue, exists := maskedMap[newKey]; exists {
					maskedMap[newKey] = oldValue + value
				} else {
					maskedMap[newKey] = value
				}
			case true:
				fmt.Println(newKey)
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

// RemoveMaskedCharacters removes masked characters from the input map
// and returns a new map
//
// Args:
//
//	input (map[string]int): Input map
//	replacementMask (string): Mask characters to apply
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	(map[string]int): Masked map
func RemoveMaskedCharacters(input map[string]int, replacementMask string, bypass bool, debug bool) map[string]int {
	maskedMap := make(map[string]int)
	replacer := strings.NewReplacer("?u", "", "?l", "", "?d", "", "?b", "", "?s", "")

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) && strings.Contains(replacementMask, "b") {
			newKey = ConvertMultiByteMask(newKey)
		}

		if debug {
			fmt.Fprintf(os.Stderr, "[?] RemoveMaskedCharacters:\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", key)
			fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
		}

		switch bypass {
		case false:
			if oldValue, exists := maskedMap[newKey]; exists {
				maskedMap[newKey] = oldValue + value
			} else {
				maskedMap[newKey] = value
			}
		case true:
			fmt.Println(newKey)
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
// (map[string]int): Matched masked map
func MakeMatchedMaskedMap(input map[string]int, replacementMask string, maskMap map[string]int, bypass bool, debug bool) map[string]int {
	maskedMap := make(map[string]int)
	replacements := ConstructReplacements(replacementMask)
	replacer := strings.NewReplacer(replacements...)

	for key, value := range input {
		newKey := replacer.Replace(key)

		if !utils.CheckASCIIString(newKey) && strings.Contains(replacementMask, "b") {
			newKey = ConvertMultiByteMask(newKey)
		}

		if debug {
			fmt.Fprintf(os.Stderr, "[?] MakeMatchedMaskedMap:\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", key)
			fmt.Fprintf(os.Stderr, "New Key: %s\n", newKey)
			fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", replacementMask)
		}

		switch bypass {
		case false:
			if _, exists := maskMap[newKey]; exists {
				if oldValue, exists := maskedMap[newKey]; exists {
					maskedMap[key] = oldValue + value
				} else {
					maskedMap[key] = value
				}
			}
		case true:
			if _, exists := maskMap[newKey]; exists {
				fmt.Println(key)
			}
		}
	}
	return maskedMap
}

// BoundarySplitPopMap splits the index of the input map into tokens based on
// the provided mask string provided and returns a new map with the tokens
// as keys and the values as the values
//
// Args:
//
//	input (map[string]int): Input map
//	replacementMask (string): Mask characters to apply
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	(map[string]int): Boundary split map
func BoundarySplitPopMap(input map[string]int, replacementMask string, bypass bool, debug bool) map[string]int {
	result := make(map[string]int)
	for s := range input {
		token := ""
		var lastRuneType rune
		var runeType rune
		for _, r := range s {
			switch {
			case unicode.IsLower(r):
				runeType = 'l'
			case unicode.IsUpper(r):
				runeType = 'u'
			case unicode.IsDigit(r):
				runeType = 'd'
			// !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~'
			case strings.ContainsRune("!\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~'", r):
				runeType = 's'
			default:
				runeType = 'b'
			}

			if (lastRuneType != 0 && lastRuneType != runeType) || !strings.ContainsRune(replacementMask, runeType) {
				if strings.ContainsRune(replacementMask, 't') && lastRuneType == 'u' && runeType == 'l' {
					// do nothing so the token continues
				} else if token != "" {
					result[token]++
					token = ""
				}
			}
			if strings.ContainsRune(replacementMask, runeType) {
				token += string(r)
			}
			lastRuneType = runeType
		}

		if debug {
			fmt.Fprintf(os.Stderr, "[?] BoundarySplitPopMap:\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", s)
			fmt.Fprintf(os.Stderr, "Token: %s\n", token)
			fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", replacementMask)
		}

		if token != "" {
			if !bypass {
				result[token]++
			} else {
				fmt.Println(token)
			}
		}
	}
	return result
}

// ShuffleMap shuffles the input map keys and replaces partially the masked
// parts of the keys with matching mask keys from the input map. This function
// resembles 'token-swapping' where the mask value is used to swap key words
// into another.
//
// Args:
//
//	input (map[string]int): Input map
//	replacementMask (string): Mask characters to apply
//	swapMap (map[string]int): Items to swap with
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
// (map[string]int): Shuffled map with swapped keys
func ShuffleMap(input map[string]int, replacementMask string, swapMap map[string]int, bypass bool, debug bool) map[string]int {
	shuffleMap := make(map[string]int)
	re := regexp.MustCompile(`^(\?u|\?l|\?d|\?s|\?b)*$`)
	reParser := regexp.MustCompile("(\\?[ludsb])")

	for key, value := range input {
		newKey := ""
		// Make a new key with the masked parts
		chars := reParser.FindAllString(key, -1)
		match := strings.Join(chars, "")

		if re.MatchString(match) {
			newKey = match
		}

		// Check if the new key is in the swap map
		for swapKey := range swapMap {
			if debug {
				fmt.Fprintf(os.Stderr, "[?] ShuffleMap:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Match: %s\n", match)
				fmt.Fprintf(os.Stderr, "Swap Token: %s\n", swapKey)
				fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", replacementMask)
			}

			maskedSwapKey := MakeMaskedString(swapKey, replacementMask)
			if maskedSwapKey == newKey {

				var shufKey string
				shufKey = strings.Replace(key, newKey, swapKey, 1)

				if debug {
					fmt.Fprintf(os.Stderr, "[?][?] Swap performed:\n")
					fmt.Fprintf(os.Stderr, "Swap Token Mask: %s\n", maskedSwapKey)
					fmt.Fprintf(os.Stderr, "Swap Result: %s\n", shufKey)
				}

				if shufKey == key {
					if debug {
						fmt.Fprintf(os.Stderr, "[?][?] Swap failed identical keys:\n")
						fmt.Fprintf(os.Stderr, "Key: %s\n", key)
						fmt.Fprintf(os.Stderr, "Swap Result: %s\n", shufKey)
					}

					continue
				}

				// if the line ends or starts with ?[uldbs] then the swap failed
				if strings.HasPrefix(shufKey, "?") || strings.HasSuffix(shufKey[len(shufKey)-2:len(shufKey)-1], "?") {

					if strings.ContainsRune("uldbs", rune(shufKey[1])) && strings.HasPrefix(shufKey, "?") || strings.ContainsRune("uldbs", rune(shufKey[len(shufKey)-1])) && strings.HasSuffix(shufKey[len(shufKey)-2:len(shufKey)-1], "?") {

						if debug {
							fmt.Fprintf(os.Stderr, "[?][?] Swap failed invalid key:\n")
							fmt.Fprintf(os.Stderr, "Key: %s\n", key)
							fmt.Fprintf(os.Stderr, "Swap Result: %s\n", shufKey)
						}

						continue

					}
				}

				switch bypass {
				case false:
					if oldValue, exists := shuffleMap[shufKey]; exists {
						shuffleMap[shufKey] = oldValue + value
					} else {
						shuffleMap[shufKey] = value
					}
				case true:
					fmt.Println(shufKey)
				}

			}
		}
	}
	return shuffleMap
}

// CalculateMaskKeyspace accepts a mask or partial mask string and returns the
// estimated keyspace of the mask
//
// Args:
//
// input (string): Input mask or partial mask
//
// Returns:
// (int): Estimated mask complexity
func CalculateMaskKeyspace(input string) int {
	if IsMaskAFullMask(input) {
		return CalculateKeyspace(input)
	}

	keyspace := 0
	// find the first "?" character in the mask
	for i := 0; i < len(input); i++ {
		if input[i] == '?' {
			if i+1 < len(input) {
				switch input[i+1] {
				case 'l':
					keyspace += CalculateMaskKeyspace(input[i : i+2])
				case 'u':
					keyspace += CalculateMaskKeyspace(input[i : i+2])
				case 'd':
					keyspace += CalculateMaskKeyspace(input[i : i+2])
				case 's':
					keyspace += CalculateMaskKeyspace(input[i : i+2])
				case 'b':
					keyspace += CalculateMaskKeyspace(input[i : i+2])
				case 'a':
					keyspace += CalculateMaskKeyspace(input[i : i+2])
				}
			}
		}
	}
	return keyspace
}

// CalculateKeyspace accepts a mask or partial mask string and returns the
// estimated keyspace of the mask
//
// Args:
// input (string): Input mask or partial mask
//
// Returns:
// (int): Estimated mask complexity
func CalculateKeyspace(input string) int {
	keyspace := 0
	for i := 0; i < len(input); i += 2 {
		if input[i] == '?' && input[i+1] == 'u' {
			keyspace += 26
		} else if input[i] == '?' && input[i+1] == 'l' {
			keyspace += 26
		} else if input[i] == '?' && input[i+1] == 'd' {
			keyspace += 10
		} else if input[i] == '?' && input[i+1] == 's' {
			keyspace += 32
		} else if input[i] == '?' && input[i+1] == 'b' {
			keyspace += 256
		} else if input[i] == '?' && input[i+1] == 'a' {
			keyspace += 95
		}
	}
	return keyspace
}

// IsMaskAFullMask accepts a mask string and returns the type of mask
// either a full mask or a partial mask
// A full mask is a mask that contains all mask characters in the mask and
// a partial mask is a mask that contains some mask characters in the mask
//
// Args:
// input (string): Input mask or partial mask
//
// Returns:
// (bool): True if the mask is a full mask, false if the mask is a partial mask
func IsMaskAFullMask(input string) bool {
	if len(input)%2 != 0 {
		return false
	}

	for i := 0; i < len(input); i += 2 {
		if input[i] != '?' {
			return false
		}

		if i+1 >= len(input) {
			return false
		}

		if input[i+1] != 'u' && input[i+1] != 'l' && input[i+1] != 'd' && input[i+1] != 's' && input[i+1] != 'b' && input[i+1] != 'a' {
			return false
		}
	}

	return true
}
