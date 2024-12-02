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
	"launchpad.net/hcre"
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

	if output != "" && len(output) < 93 {
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
// returnMap (map[string]int): Map of items to return
func AppendRules(items map[string]int, operation string, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	switch operation {
	// remove will remove characters then append
	case "rule-append-remove", "append-remove":
		for key, value := range items {
			if len(key) > 15 {
				if debug {
					fmt.Fprintf(os.Stderr, "[!] Error: Key is too long for append-remove operation\n")
				}
				continue
			}
			rule := CharToRule(key, "$")
			remove := LenToRule(key, "]")
			appendRemoveRule := FormatCharToRuleOutput(remove, rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] AppendRules (remove):\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "Remove: %s\n", remove)
				fmt.Fprintf(os.Stderr, "AppendRemoveRule: %s\n", appendRemoveRule)
			}

			if appendRemoveRule != "" && !bypass {
				returnMap[appendRemoveRule] = value
			} else if appendRemoveRule != "" && bypass {
				fmt.Println(appendRemoveRule)
			}
		}
		return returnMap
	default:
		for key, value := range items {
			rule := CharToRule(key, "$")
			appendRule := FormatCharToRuleOutput(rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] AppendRules:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "AppendRule: %s\n", appendRule)
			}

			if appendRule != "" && !bypass {
				returnMap[appendRule] = value
			} else if appendRule != "" && bypass {
				fmt.Println(appendRule)
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	returnMap (map[string]int): Map of items to return
func PrependRules(items map[string]int, operation string, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	switch operation {
	// remove will remove characters then prepend
	case "rule-prepend-remove", "prepend-remove":
		for key, value := range items {
			if len(key) > 15 {
				if debug {
					fmt.Fprintf(os.Stderr, "[!] Error: Key is too long for prepend-remove operation\n")
				}
				continue
			}
			rule := CharToRule(utils.ReverseString(key), "^")
			remove := LenToRule(key, "[")
			prependRemoveRule := FormatCharToRuleOutput(remove, rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] PrependRules (remove):\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "Remove: %s\n", remove)
				fmt.Fprintf(os.Stderr, "PrependRemoveRule: %s\n", prependRemoveRule)
			}

			if prependRemoveRule != "" && !bypass {
				returnMap[prependRemoveRule] = value
			} else if prependRemoveRule != "" && bypass {
				fmt.Println(prependRemoveRule)
			}
		}
		return returnMap
	case "rule-prepend-toggle", "prepend-toggle":
		for key, value := range items {
			rule := CharToRule(utils.ReverseString(key), "^")
			toggle := StringToToggleRule("A", "T", len(key))
			prependToggleRule := FormatCharToRuleOutput(rule, toggle)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] PrependRules (toggle):\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "Toggle: %s\n", toggle)
				fmt.Fprintf(os.Stderr, "PrependToggleRule: %s\n", prependToggleRule)
			}

			if prependToggleRule != "" && !bypass {
				returnMap[prependToggleRule] = value
			} else if prependToggleRule != "" && bypass {
				fmt.Println(prependToggleRule)
			}
		}
		return returnMap
	default:
		for key, value := range items {
			rule := CharToRule(utils.ReverseString(key), "^")
			prependRule := FormatCharToRuleOutput(rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] PrependRules:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "PrependRule: %s\n", prependRule)
			}

			if prependRule != "" && !bypass {
				returnMap[prependRule] = value
			} else if prependRule != "" && bypass {
				fmt.Println(prependRule)
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
// returnMap (map[string]int): Map of items to return
func InsertRules(items map[string]int, index string, end string, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		os.Exit(1)
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		os.Exit(1)
	}

	for i < e+1 {
		for key, value := range items {
			rule := CharToIteratingRule(key, "i", i)
			insertRule := FormatCharToIteratingRuleOutput(i, rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] InsertRules:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "InsertRule: %s\n", insertRule)
			}

			if insertRule != "" && !bypass {
				returnMap[insertRule] = value
			} else if insertRule != "" && bypass {
				fmt.Println(insertRule)
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
// bypass (bool): If true, the map is not used for output or filtering
// debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	returnMap (map[string]int): Map of items to return
func OverwriteRules(items map[string]int, index string, end string, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		os.Exit(1)
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		os.Exit(1)
	}

	for i < e+1 {
		for key, value := range items {
			rule := CharToIteratingRule(key, "o", i)
			overwriteRule := FormatCharToIteratingRuleOutput(i, rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] OverwriteRules:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "OverwriteRule: %s\n", overwriteRule)
			}

			if overwriteRule != "" && !bypass {
				returnMap[overwriteRule] = value
			} else if overwriteRule != "" && bypass {
				fmt.Println(overwriteRule)
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
//	bypass (bool): If true, the map is not used for output or filtering
//	debug (bool): If true, print additional debug information to stderr
//
// Returns:
//
//	returnMap (map[string]int): Map of items to return
func ToggleRules(items map[string]int, index string, end string, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	i, err := strconv.Atoi(index)
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		os.Exit(1)
	}

	e, err := strconv.Atoi(end)
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		os.Exit(1)
	}

	for i < e+1 {
		for key, value := range items {
			// if the key is all uppercase just set it to "u"
			rule := ""
			if strings.ToUpper(key) == key {
				rule = "u"
			} else {
				rule = StringToToggleRule(key, "T", i)
			}

			toggleRule := FormatCharToIteratingRuleOutput(i, rule)

			if debug {
				fmt.Fprintf(os.Stderr, "[?] ToggleRules:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "ToggleRule: %s\n", toggleRule)
			}

			if toggleRule != "" && !bypass {
				returnMap[toggleRule] = value
			} else if toggleRule != "" && bypass {
				fmt.Println(toggleRule)
			}
		}

		i++
	}
	return returnMap
}

// ApplyRulesHCRE uses the HCRE library to apply rules to a map of items
// and returns the results
//
// Args:
// items (map[string]int): Items to use in the operation
// rules (map[string]int): Rules to use in the operation
// bypass (bool): If true, the map is not used for output or filtering
// debug (bool): If true, print additional debug information to stderr
//
// Returns:
// returnMap (map[string]int): Map of items to return
func ApplyRulesHCRE(items map[string]int, rules map[string]int, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	for key, value := range items {
		for rule, _ := range rules {

			rr, err := hcre.Compile(rule)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error: %s\n", err)
				os.Exit(1)
			}
			applyRule := rr.Simplify().Apply([]byte(key))

			if debug {
				fmt.Fprintf(os.Stderr, "[?] ApplyRulesHCRE:\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Rule: %s\n", rule)
				fmt.Fprintf(os.Stderr, "ApplyRule: %s\n", applyRule)
			}

			if applyRule != nil && !bypass {
				returnMap[string(applyRule)] = value
			} else if applyRule != nil && bypass {
				fmt.Println(string(applyRule))
			}
		}
	}
	return returnMap
}

// SimplifyRules simplifies rules by simplifying rules to optimized equivalents
// using the HCRE library
//
// Args:
// items (map[string]int): Items to use in the operation
// bypass (bool): If true, the map is not used for output or filtering
// debug (bool): If true, print additional debug information to stderr
//
// Returns:
// returnMap (map[string]int): Map of items to return
func SimplifyRules(items map[string]int, bypass bool, debug bool) (returnMap map[string]int) {
	returnMap = make(map[string]int)
	for key, value := range items {

		rr, err := hcre.Compile(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error: %s\n", err)
			os.Exit(1)
		}
		simplifyRule := rr.Simplify().String()

		if debug {
			fmt.Fprintf(os.Stderr, "[?] SimplifyRules:\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", key)
			fmt.Fprintf(os.Stderr, "SimplifyRule: %s\n", simplifyRule)
		}

		if simplifyRule != "" && !bypass {
			returnMap[simplifyRule] = value
		} else if simplifyRule != "" && bypass {
			fmt.Println(simplifyRule)
		}
	}
	return returnMap
}
