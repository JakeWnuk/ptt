// Package validation provides functions to validate and transform strings
package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// CheckASCIIString checks to see if a string only contains ascii characters
//
// Args:
//
//	str (string): Input string to check
//
// Returns:
//
//	(bool): If the string only contained ASCII characters
func CheckASCIIString(str string) bool {
	if utf8.RuneCountInString(str) != len(str) {
		return false
	}
	return true
}

// ConvertMultiByteCharToRule converts non-ascii characters to a hashcat valid format
// for rule.CharToRule functions
//
// Args:
//
//	str (string): Input string to transform
//
// Returns:
//
//	returnStr (string): Converted string
func ConvertMultiByteCharToRule(str string) string {
	returnStr := ""
	deletedChar := ``
	for i, r := range str {
		if r > 127 {
			if i > 0 {
				deletedChar = string(returnStr[len(returnStr)-1])
				returnStr = returnStr[:len(returnStr)-1]
			}
			byteArr := []byte(string(r))
			if deletedChar == "^" {
				for j := len(byteArr) - 1; j >= 0; j-- {
					b := byteArr[j]
					if j == 0 {
						returnStr += fmt.Sprintf("%s\\x%X", deletedChar, b)
					} else {
						returnStr += fmt.Sprintf("%s\\x%X ", deletedChar, b)
					}
				}
			} else {
				for j, b := range byteArr {
					if j == len(byteArr)-1 {
						returnStr += fmt.Sprintf("%s\\x%X", deletedChar, b)
					} else {
						returnStr += fmt.Sprintf("%s\\x%X ", deletedChar, b)
					}
				}
			}
		} else {
			returnStr += fmt.Sprintf("%c", r)
		}
	}
	return returnStr
}

// ConvertMultiByteCharToIteratingRule converts non-ascii characters to a hashcat valid format
// for rule.CharToIteratingRule functions
//
// Args:
//
//	index (int): Index to start the iteration
//	str (string): Input string to transform
//
// Returns:
//
//	returnStr (string): Converted string
func ConvertMultiByteCharToIteratingRule(index int, str string) string {
	output := ""
	lastIterationSeen := fmt.Sprintf("%s%d", string([]rune(str)[0]), index)

	re := regexp.MustCompile(`[io][\dA-Z]`)

	for _, word := range strings.Split(str, " ") {
		for _, c := range word {
			if c > 127 {
				// Convert to UTF-8 bytes
				bytes := []byte(string(c))
				firstByteOut := true
				// Convert each byte to its hexadecimal representation
				for i, b := range bytes {
					if firstByteOut {
						output += fmt.Sprintf("\\x%X ", b)
						firstByteOut = false
						continue
					}
					lastIterationSeen = IncrementIteratingRuleCall(lastIterationSeen)
					if i == len(bytes)-1 {
						output += fmt.Sprintf("%s\\x%X", lastIterationSeen, b)
					} else {
						output += fmt.Sprintf("%s\\x%X ", lastIterationSeen, b)
					}
				}
			} else {
				output += string(c)
				if len(output) > 2 && re.MatchString(output[len(output)-2:]) {
					lastIterationSeen = output[len(output)-2:]
				}
			}
		}
		output += " "
	}

	return output
}

// IncrementIteratingRuleCall increments the last character of a string for
// rules.CharToIteratingRules functions
//
// For example, "i4" will be incremented to "i5", "iA" will be incremented to
// "IB"
//
// Args:
//
//	s (string): Input string to increment
//
// Returns:
//
//	output (string): Incremented string
func IncrementIteratingRuleCall(s string) string {
	if len(s) == 0 {
		return s
	}

	lastChar := s[len(s)-1]
	incChar := lastChar + 1

	// Replace the last character with the incremented character
	output := s[:len(s)-1] + string(incChar)

	return output
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
