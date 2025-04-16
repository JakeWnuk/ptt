// Package rule provides functions to convert strings into rules
package rule

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/jakewnuk/ptt/pkg/validation"
)

// LenToRule converts a string to a rule by its length
//
// Args:
//
//	str (string): Input string to transform
//	rule (string): Rule to insert per length
//
// Returns:
//
//	(string): Transformed string
func LenToRule(str string, rule string) string {
	return strings.TrimSpace(strings.Repeat(rule+" ", len(str)))
}

// CharToRule converts a string to a rule by its characters
//
// Args:
//
//	str (string): Input string to transform
//	rule (string): Rule to insert per character
//
// Returns:
//
//	(string): Transformed string
func CharToRule(str string, rule string) string {
	return rule + strings.Join(strings.Split(str, ""), " "+rule)
}

// CharToIteratingRule converts a string to a rule by its characters but
// increments along with each character
//
// Args:
//
//	str (string): Input string to transform
//	rule (string): Rule to insert per length
//	index (int): Index to start at
//
// Returns:
//
//	(string): Transformed string
func CharToIteratingRule(str string, rule string, index int) string {
	var result strings.Builder
	for i, r := range str {
		if i+index < 10 {
			result.WriteString(fmt.Sprintf("%s%d%c ", rule, i+index, r))
		} else if i+index-10 < 26 {
			result.WriteString(fmt.Sprintf("%s%c%c ", rule, 'A'+i+index-10, r))
		}
	}
	return strings.TrimSpace(result.String())
}

// StringToToggleRule converts a string to toggle rules by looking for upper chars
//
// Args:
//
//	str (string): Input string to transform
//	rule (string): Rule to insert per length
//	index (int): Index to start at
//
// Returns:
//
//	(string): Transformed string
func StringToToggleRule(str string, rule string, index int) string {
	var result strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i+index < 10 {
				result.WriteString(fmt.Sprintf("%s%d ", rule, i+index))
			} else if i+index-10 < 26 {
				result.WriteString(fmt.Sprintf("%s%c ", rule, 'A'+i+index-10))
			}
		}
	}
	return strings.TrimSpace(result.String())
}

// FormatCharToRuleOutput handles formatting of rule output
// for CharToRule functions
//
// Args:
//
//	strs (...string): Input strings to print
//
// Returns:
//
//	output (string): Formatted output
func FormatCharToRuleOutput(strs ...string) (output string) {
	output = ""
	for _, str := range strs {
		if validation.CheckASCIIString(str) {
			output += str + " "
		} else {
			output += validation.ConvertMultiByteCharToRule(str)
		}
	}

	// if the string ends in '$ ' add a ':' to the end
	if strings.HasSuffix(output, "$  ") {
		output = output[:len(output)-1] + ":"
	}

	if output != "" && len(output) < 93 {
		return strings.TrimSpace(output)
	}

	return ""
}

// FormatCharToIteratingRuleOutput handles formatting of rule output
// for CharToIteratingRule functions
//
// Args:
//
//	index (int): Index to start at
//	strs (...string): Input strings to print
//
// Returns:
//
//	output (string): Formatted output
func FormatCharToIteratingRuleOutput(index int, strs ...string) (output string) {
	output = ""
	for _, str := range strs {
		if validation.CheckASCIIString(str) {
			output += str + " "
		} else {
			output += validation.ConvertMultiByteCharToIteratingRule(index, str)
		}
	}

	if len(output)-3 >= 0 {
		// allow for the last character to be a space for overwrite and insert rules
		if output[len(output)-3:len(output)-2] == "o" || output[len(output)-3:len(output)-2] == "i" {
			output = output + ":"
		}
	}

	if output != "" && len(output) < 93 {
		return strings.TrimSpace(output)
	}

	return ""
}
