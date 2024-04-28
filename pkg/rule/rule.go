// Package rule contains the transformation logic to convert input into Hashcat
// rules
package rule

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/jakewnuk/ptt/pkg/utils"
)

// ----------------------------------------------------------------------------
// Transformation Functions
// ----------------------------------------------------------------------------

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
		} else {
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

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

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
		if utils.CheckASCIIString(str) {
			output += str + " "
		} else {
			output += utils.ConvertMultiByteCharToRule(str)
		}
	}

	// if the string ends in '$ ' add a ':' to the end
	if strings.HasSuffix(output, "$  ") {
		output = output[:len(output)-1] + ":"
	}

	if output != "" && len(output) <= 93 {
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
		if utils.CheckASCIIString(str) {
			output += str + " "
		} else {
			output += utils.ConvertMultiByteCharToIteratingRule(index, str)
		}
	}

	if len(output)-3 >= 0 {
		// allow for the last character to be a space for overwrite and insert rules
		if output[len(output)-3:len(output)-2] == "o" || output[len(output)-3:len(output)-2] == "i" {
			output = output + ":"
		}
	}

	if output != "" && len(output) <= 93 {
		return strings.TrimSpace(output)
	}

	return ""
}

// AppendRules transforms input into append rules
//
// Args:
//
//	items (map[string]int): Items to use in the operation
//	operation (string): Operation to use in the function
//
// Returns:
//
// returnMap (map[string]int): Map of items to return
func AppendRules(items map[string]int, operation string) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	switch operation {
	// remove will remove characters then append
	case "append-remove":
		for key, value := range items {
			rule := CharToRule(key, "$")
			remove := LenToRule(key, "]")
			appendRemoveRule := FormatCharToRuleOutput(remove, rule)
			if appendRemoveRule != "" {
				returnMap[appendRemoveRule] = value
			}
		}
		return returnMap
	// shift will shift characters back to front then append
	case "append-shift":
		for key, value := range items {
			rule := CharToRule(key, "$")
			shift := LenToRule(key, "}")
			appendShiftRule := FormatCharToRuleOutput(shift, rule)
			if appendShiftRule != "" {
				returnMap[appendShiftRule] = value
			}
		}
		return returnMap
	default:
		for key, value := range items {
			rule := CharToRule(key, "$")
			appendRule := FormatCharToRuleOutput(rule)
			if appendRule != "" {
				returnMap[appendRule] = value
			}
		}
		return returnMap
	}
}

// PrependRules transforms input into prepend rules
//
// Args:
//
//	items (map[string]int): Items to use in the operation
//	operation (string): Operation to use in the function
//
// Returns:
//
//	returnMap (map[string]int): Map of items to return
func PrependRules(items map[string]int, operation string) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	switch operation {
	// remove will remove characters then prepend
	case "prepend-remove":
		for key, value := range items {
			rule := CharToRule(utils.ReverseString(key), "^")
			remove := LenToRule(key, "[")
			prependRemoveRule := FormatCharToRuleOutput(remove, rule)
			if prependRemoveRule != "" {
				returnMap[prependRemoveRule] = value
			}
		}
		return returnMap
	// shift will shift characters front to back then prepend
	case "prepend-shift":
		for key, value := range items {
			rule := CharToRule(utils.ReverseString(key), "^")
			shift := LenToRule(key, "{")
			prependShiftRule := FormatCharToRuleOutput(shift, rule)
			if prependShiftRule != "" {
				returnMap[prependShiftRule] = value
			}
		}
		return returnMap
	default:
		for key, value := range items {
			rule := CharToRule(utils.ReverseString(key), "^")
			prependRule := FormatCharToRuleOutput(rule)
			if prependRule != "" {
				returnMap[prependRule] = value
			}
		}
		return returnMap
	}
}

// InsertRules transforms input into insert rules by index
//
// Args:
//
//	items (map[string]int): Items to use in the operation
//	index (string): Index to insert at
//	end (string): Index to end at
//
// Returns:
// returnMap (map[string]int): Map of items to return
func InsertRules(items map[string]int, index string, end string) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for i < e+1 {
		for key, value := range items {
			rule := CharToIteratingRule(key, "i", i)
			insertRule := FormatCharToIteratingRuleOutput(i, rule)
			if insertRule != "" {
				returnMap[insertRule] = value
			}
		}
		i++
	}
	return returnMap
}

// OverwriteRules transforms input into overwrite rules by index
//
// Args:
//
// items (map[string]int): Items to use in the operation
// index (string): Index to overwrite at
// end (string): Index to end at
//
// Returns:
//
//	returnMap (map[string]int): Map of items to return
func OverwriteRules(items map[string]int, index string, end string) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for i < e+1 {
		for key, value := range items {
			rule := CharToIteratingRule(key, "o", i)
			overwriteRule := FormatCharToIteratingRuleOutput(i, rule)
			if overwriteRule != "" {
				returnMap[overwriteRule] = value
			}
		}
		i++
	}
	return returnMap
}

// ToggleRules transforms input into  toggle rules starting at an index
//
// Args:
//
//	items (map[string]int): Items to use in the operation
//	index (string): Index to start at
//	end (string): Index to end at
//
// Returns:
//
//	returnMap (map[string]int): Map of items to return
func ToggleRules(items map[string]int, index string, end string) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	for i < e+1 {
		for key, value := range items {
			rule := StringToToggleRule(key, "T", i)
			toggleRule := FormatCharToIteratingRuleOutput(i, rule)
			if toggleRule != "" {
				returnMap[toggleRule] = value
			}
		}

		i++
	}
	return returnMap
}
