package transform

import (
	"bufio"
	"fmt"
	"os"
	"ptt/pkg/filter"
	"ptt/pkg/mask"
	"ptt/pkg/models"
	"ptt/pkg/rule"
	"ptt/pkg/utils"
	"ptt/pkg/validation"
	"strings"

	"launchpad.net/hcre"
)

// ReadReturnStandardInput reads from standard input and applies the given
// transformation to each line.
//
// Args:
// transformation (models.MultiString): A list of transformation rules to be
// applied to each line.
//
// Returns:
// None
func ReadReturnStandardInput(transformation models.MultiString) {
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		readText := reader.Text()
		line := readText

		for _, operation := range transformation {
			if strings.Contains(operation, "insert") || strings.Contains(operation, "overwrite") || strings.Contains(operation, "toggle") {
				start := models.OperationStart
				for models.OperationStart < models.OperationEnd+1 {
					line = Apply(line, operation)

					if filter.Pass(line) {
						fmt.Println(line)
					}

					models.OperationStart++
					line = readText
				}
				models.OperationStart = start
			} else {
				line = Apply(line, operation)

				if filter.Pass(line) {
					fmt.Println(line)
				}

				line = readText
			}
		}
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

// Apply applies a transformation function to the string and returns the
// result.
//
// Args:
// input (string): The input string to be transformed.
// transform (string): The transformation function to be applied.
//
// Returns:
// (string): The transformed string.
func Apply(input string, transform string) string {
	switch transform {
	case "append", "rule-append":
		return appendRules(input)
	case "append-remove", "rule-append-remove":
		return appendRemoveRules(input)
	case "prepend", "rule-prepend":
		return prependRules(input)
	case "prepend-remove", "rule-prepend-remove":
		return prependRemoveRules(input)
	case "prepend-toggle", "rule-prepend-toggle":
		return prependToggleRules(input)
	case "insert", "rule-insert":
		return insertRules(input)
	case "overwrite", "rule-overwrite":
		return overwriteRules(input)
	case "toggle", "rule-toggle":
		return toggleRules(input)
	case "simplify", "rule-simplify":
		return ruleSimplify(input)
	case "mask", "rule-mask":
		return makeMask(input)
	case "remove", "mask-remove":
		return mask.RemoveMaskedCharacters(makeMask(input))
	default:
		return ""
	}
}

// AppendRules transforms the input string into an append rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of an append rule.
func appendRules(key string) string {
	keyRule := rule.CharToRule(key, "$")
	appendRule := rule.FormatCharToRuleOutput(keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.appendRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Return: %s\n", appendRule)
	}
	return appendRule
}

// AppendRemoveRules transforms the input string into an append-remove rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of an append-remove rule.
func appendRemoveRules(key string) string {
	if len(key) > 15 {
		if models.DebugMode {
			fmt.Fprintf(os.Stderr, "[!] transform.appendRemoveRules(key):\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", key)
			fmt.Fprintf(os.Stderr, "Error: key is too long\n")
		}
		return ""
	}
	keyRule := rule.CharToRule(key, "$")
	remove := rule.LenToRule(key, "]")
	appendRemoveRule := rule.FormatCharToRuleOutput(remove, keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.appendRemoveRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Remove: %s\n", remove)
		fmt.Fprintf(os.Stderr, "Return: %s\n", appendRemoveRule)
	}
	return appendRemoveRule
}

// prependRules transforms the input string into a prepend rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a prepend rule.
func prependRules(key string) string {
	keyRule := rule.CharToRule(utils.ReverseString(key), "^")
	prependRule := rule.FormatCharToRuleOutput(keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.prependRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Return: %s\n", prependRule)
	}
	return prependRule
}

// prependRemoveRules transforms the input string into a prepend-remove rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a prepend-remove rule.
func prependRemoveRules(key string) string {
	if len(key) > 15 {
		if models.DebugMode {
			fmt.Fprintf(os.Stderr, "[!] transform.prependRemoveRules(key):\n")
			fmt.Fprintf(os.Stderr, "Key: %s\n", key)
			fmt.Fprintf(os.Stderr, "Error: key is too long\n")
			return ""
		}
	}
	keyRule := rule.CharToRule(utils.ReverseString(key), "^")
	remove := rule.LenToRule(key, "[")
	prependRemoveRule := rule.FormatCharToRuleOutput(remove, keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.prependRemoveRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Remove: %s\n", remove)
		fmt.Fprintf(os.Stderr, "Return: %s\n", prependRemoveRule)
	}
	return prependRemoveRule
}

// prependToggleRules transforms the input string into a prepend-toggle rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a prepend-toggle rule.
func prependToggleRules(key string) string {
	keyRule := rule.CharToRule(utils.ReverseString(key), "^")
	toggle := rule.StringToToggleRule("A", "T", len(key))
	prependToggleRule := rule.FormatCharToRuleOutput(keyRule, toggle)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.prependToggleRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Toggle: %s\n", toggle)
		fmt.Fprintf(os.Stderr, "Return: %s\n", prependToggleRule)
	}
	return prependToggleRule
}

// insertRules transforms the input string into an insert rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of an insert rule.
func insertRules(key string) string {
	keyRule := rule.CharToIteratingRule(key, "i", models.OperationStart)
	insertRule := rule.FormatCharToIteratingRuleOutput(models.OperationStart, keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.insertRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Return: %s\n", insertRule)
	}
	return insertRule
}

// overwriteRules transforms the input string into an overwrite rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of an overwrite rule.
func overwriteRules(key string) string {
	keyRule := rule.CharToIteratingRule(key, "o", models.OperationStart)
	overwriteRule := rule.FormatCharToIteratingRuleOutput(models.OperationStart, keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.overwriteRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Return: %s\n", overwriteRule)
	}
	return overwriteRule
}

// toggleRules transforms the input string into a toggle rule.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a toggle rule.
func toggleRules(key string) string {
	// if the key is all uppercase just set it to "u"
	keyRule := ""
	if strings.ToUpper(key) == key {
		keyRule = "u"
	} else {
		keyRule = rule.StringToToggleRule(key, "T", models.OperationStart)
	}

	toggleRule := rule.FormatCharToIteratingRuleOutput(models.OperationStart, keyRule)

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.toggleRules(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Rule: %s\n", keyRule)
		fmt.Fprintf(os.Stderr, "Return: %s\n", toggleRule)
	}

	return toggleRule
}

// ruleSimplify simplifies the input string into a simplified rule using the
// HCRE library.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a simplified rule.
func ruleSimplify(key string) string {
	rr, err := hcre.Compile(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] transform.ruleSimplify(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	simplifyRule := rr.Simplify().String()

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.ruleSimplify(key):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Return: %s\n", simplifyRule)
	}

	return simplifyRule
}

// makeMask transforms the input string into a mask using the global mask
// defined in the models package.
//
// Args:
// key (string): The input string to be transformed.
//
// Returns:
// (string): The transformed string in the form of a mask.
func makeMask(key string) string {
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
