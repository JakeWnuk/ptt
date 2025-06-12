// Package transform provides functions to transform strings based on
// various rules and operations.
package transform

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/jakewnuk/ptt/pkg/filter"
	"github.com/jakewnuk/ptt/pkg/mask"
	"github.com/jakewnuk/ptt/pkg/models"
	"github.com/jakewnuk/ptt/pkg/rule"
	"github.com/jakewnuk/ptt/pkg/utils"

	"launchpad.net/hcre"
)

var verboseOutputMutex sync.Mutex

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
	var wg sync.WaitGroup

	for reader.Scan() {
		readText := reader.Text()
		line := readText

		for _, operation := range transformation {
			if strings.Contains(operation, "insert") || strings.Contains(operation, "overwrite") || strings.Contains(operation, "toggle") || strings.Contains(operation, "hex") {
				start := models.OperationStart
				for models.OperationStart < models.OperationEnd+1 {
					line = Apply(line, operation)

					if filter.Pass(line) {
						if models.Verbose {
							if models.VerboseOutput[line] == 0 {
								models.VerboseOutput[line] = 1
							} else {
								models.VerboseOutput[line]++
							}
						} else {
							fmt.Println(line)
						}
					}

					models.OperationStart++
					line = readText
				}
				models.OperationStart = start
			} else if strings.Contains(operation, "pop") || strings.Contains(operation, "passphrase") || strings.Contains(operation, "regram") || strings.Contains(operation, "swap") {
				wg.Add(1)

				go func(line, operation string) {
					defer wg.Done()

					output := Parse(line, operation)
					for item := range output {
						verboseOutputMutex.Lock()
						if filter.Pass(item) {
							if models.Verbose {
								if models.VerboseOutput[item] == 0 {
									models.VerboseOutput[item] = 1
								} else {
									models.VerboseOutput[item]++
								}
							} else {
								fmt.Println(item)
							}
						}
						verboseOutputMutex.Unlock()
					}
				}(line, operation)
			} else if strings.Contains(operation, "swap") {
				wg.Add(1)

				go func(line, operation string) {
					defer wg.Done()
					for i := 0; i < models.TokenSwapCount; i++ {
						output := Parse(line, operation)
						for item := range output {
							if filter.Pass(item) {
								if models.Verbose {
									if models.VerboseOutput[item] == 0 {
										models.VerboseOutput[item] = 1
									} else {
										models.VerboseOutput[item]++
									}
								} else {
									fmt.Println(item)
								}
							}
							verboseOutputMutex.Unlock()
						}
					}
				}(line, operation)

			} else {
				line = Apply(line, operation)

				if filter.Pass(line) {
					if models.Verbose {
						if models.VerboseOutput[line] == 0 {
							models.VerboseOutput[line] = 1
						} else {
							models.VerboseOutput[line]++
						}
					} else {
						fmt.Println(line)
					}
				}

				line = readText
			}
		}
	}
	if models.VeryVerbose {
		utils.PrintStatsToSTDOUT(models.VerboseOutput)
	} else if models.Verbose {
		utils.PrintArrayToSTDOUT(models.VerboseOutput)
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
		return mask.MakeMask(input)
	case "remove", "mask-remove":
		return mask.RemoveMaskedCharacters(mask.MakeMask(input))
	case "dehex", "unhex":
		return mask.DeHex(input)
	default:
		return ""
	}
}

// Parse parses the input string based on the specified transformation.
func Parse(input string, transform string) map[string]int {
	switch transform {
	case "pop", "mask-pop":
		return maskPop(input)
	case "passphrase", "phrase":
		return makePassphrase(input)
	case "regram", "n-gram":
		return makeNGrams(input)
	case "swap", "token-swap":
		return tokenSwap(input)
	default:
		return map[string]int{}
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

// maskPop transforms the input string into a map of tokens and their
// occurrences.
//
// Args:
// input (string): The input string to be transformed.
//
// Returns:
// (map[string]int): A map of tokens and their occurrences.
func maskPop(input string) map[string]int {
	result := make(map[string]int)
	token := ""
	var lastRuneType rune
	var runeType rune
	for _, r := range input {
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

		if (lastRuneType != 0 && lastRuneType != runeType) || !strings.ContainsRune(models.GlobalMask, runeType) {
			if strings.ContainsRune(models.GlobalMask, 't') && lastRuneType == 'u' && runeType == 'l' {
				// do nothing so the token continues
			} else if token != "" {
				result[token]++
				token = ""
			}
		}
		if strings.ContainsRune(models.GlobalMask, runeType) {
			token += string(r)
		}
		lastRuneType = runeType
	}

	if models.DebugMode {
		fmt.Fprintf(os.Stderr, "[?] transform.maskPop(input):\n")
		fmt.Fprintf(os.Stderr, "Key: %s\n", input)
		fmt.Fprintf(os.Stderr, "Token: %s\n", token)
		fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", models.GlobalMask)
	}

	if token != "" {
		result[token]++
	}
	return result
}

// makePassphrase transforms the input string into a map of passphrases
// and their occurrences.
//
// Args:
// input (string): The input string to be transformed.
//
// Returns:
// (map[string]int): A map of passphrases and their occurrences.
func makePassphrase(input string) map[string]int {
	newMap := make(map[string]int)
	for i := models.WordStart; i <= models.WordEnd; i++ {
		newKeyArray := utils.GeneratePassphrase(input, i)
		for value, newKey := range newKeyArray {
			if newMap[newKey] == 0 {
				newMap[newKey] = value
			} else {
				newMap[newKey] += value
			}
		}
	}

	return newMap
}

// makeNGrams transforms the input string into a map of n-grams
// and their occurrences.
//
// Args:
// input (string): The input string to be transformed.
//
// Returns:
// (map[string]int): A map of n-grams and their occurrences.
func makeNGrams(input string) map[string]int {
	newMap := make(map[string]int)
	for i := models.WordStart; i <= models.WordEnd; i++ {
		newKeyArray := utils.GenerateNGrams(input, i)
		for value, newKey := range newKeyArray {
			newKey = strings.TrimSpace(newKey)
			newKey = strings.TrimLeft(newKey, ",")
			newKey = strings.TrimRight(newKey, ",")
			newKey = strings.TrimLeft(newKey, " ")

			if newMap[newKey] == 0 {
				newMap[newKey] = value
			} else {
				newMap[newKey] += value
			}
		}
	}
	return newMap
}

// tokenSwap transforms the input string by token swapping with the following method:
//
// 1. Pops tokens and adds them to a global list
// 2. Take top 250k tokens and look for retain masks in the string
// 3. All found retain/partial masks are used for swapping
// 4. swap with ALL of the popped tokens so far
//
// Args:
// input (string): The input string to be transformed
//
// Returns:
// (map[string]int): A map of token-swapped input strings
func tokenSwap(input string) map[string]int {
	models.GlobalTokensMutex.Lock()
	poppedTokens := maskPop(input)
	for i := range poppedTokens {
		if i == "" || i == " " {
			continue
		}
		if models.GlobalTokens[i] == 0 {
			models.GlobalTokens[i] = 1
		} else {
			models.GlobalTokens[i]++
		}
	}
	// Sort by frequency
	p := make(models.PairList, len(models.GlobalTokens))
	i := 0
	for k, v := range models.GlobalTokens {
		p[i] = models.Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))

	// Create an array of the top 250,000 tokens
	topTokens := make(map[string]int)
	for i := 0; i < 250000 && i < len(p); i++ {
		// Increase min token size to 3
		if len(p[i].Key) > 2 {
			topTokens[p[i].Key] = p[i].Value
		}
	}

	// Create retained masks
	retainedMasks := utils.MakeRetainMaskedMap(input, topTokens)

	// Token swap
	swappedTokens := utils.ShuffleMap(retainedMasks, models.GlobalTokens)
	models.GlobalTokensMutex.Unlock()

	return swappedTokens
}
