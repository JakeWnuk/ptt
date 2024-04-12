// Package output contains functions for printing and formatting output to the console
package output

import (
	"fmt"
	"ppt/pkg/models"
	"ppt/pkg/rule"
	"sort"
)

// ----------------------------------------------------------------------------
// Controllers
// ----------------------------------------------------------------------------
// TransformationController is the main entry point for the CLI
// application. Operates a switch statement to determine the
// mode to use.
//
// Modes:
// -
//
// Args:
//
//	input (map[string]int): A map of input values.
//	mode (string): The mode to run the CLI in.
//
// Returns:
//
//	None
func TransformationController(input map[string]int, mode string) (output map[string]int) {
	switch mode {
	case "append", "append-remove", "append-shift":
		output = rule.AppendRules(input, mode)
	case "prepend", "prepend-remove", "prepend-shift":
		output = rule.PrependRules(input, mode)
	case "insert":
		// TODO make length configurable
		output = rule.InsertRules(input, "0")
	case "overwrite":
		output = rule.OverwriteRules(input, "0")
	}
	return output
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// PrintArrayToSTDOUT prints an array of items to stdout
//
// Args:
//
//	items (map[string]int): A map of item frequencies
//	verbose (bool): If true, the count is printed
//
// Returns:
//
//	None
func PrintArrayToSTDOUT(freq map[string]int, verbose bool) {
	p := make(models.PairList, len(freq))
	i := 0
	for k, v := range freq {
		p[i] = models.Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	for _, pair := range p {
		if verbose {
			fmt.Printf("%d %s\n", pair.Value, pair.Key)
		} else {
			fmt.Printf("%s\n", pair.Key)
		}
	}
}

// RetainRemove compares a string against a list of words to retain and remove
//
// Args:
//
//	text (map[string]int): The text to process
//	retainMap (map[string]int): A map of words to retain
//	removeMap (map[string]int): A map of words to remove
//
// Returns:
//
//	(string, error): The processed and filtered text, or an error if the text is not valid
func RetainRemove(textMap, retainMap, removeMap map[string]int) (map[string]int, error) {
	result := make(map[string]int)

	// If the retain and remove maps are empty, return the textMap as is
	if len(retainMap) == 0 && len(removeMap) == 0 {
		return textMap, nil
	}

	for k, v := range textMap {
		// If the key is in the retain map and not in the remove map, add it to the result
		if _, ok := retainMap[k]; ok || len(retainMap) == 0 {
			if _, ok := removeMap[k]; !ok || len(removeMap) == 0 {
				result[k] = v
			}
		}
	}

	// If the result map is empty, return an error
	if len(result) == 0 {
		return nil, fmt.Errorf("no keys in the text map are in the retain list")
	}

	return result, nil
}

// RemoveMinimumFrequency removes items from a map that are below a minimum frequency
// threshold and returns a new map
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//	min (int): The minimum frequency threshold
//
// Returns:
//
//	(map[string]int): A new map of item frequencies above the minimum threshold
func RemoveMinimumFrequency(freq map[string]int, min int) map[string]int {
	newFreq := make(map[string]int)
	for key, value := range freq {
		if value >= min {
			newFreq[key] = value
		}
	}
	return newFreq
}
