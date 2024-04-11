package output

import (
	"fmt"
)

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
	for key, value := range freq {
		if verbose {
			fmt.Printf("%s: %d\n", key, value)
		} else {
			fmt.Printf("%s\n", key)
		}
	}
}

// CompareToRetainRemoveFlags compares a string against a list of words to retain and remove
//
// Args:
//
//	text (string): The text to process
//	retainMap (map[string]int): A map of words to retain
//	removeMap (map[string]int): A map of words to remove
//
// Returns:
//
//	(string, error): The processed and filtered text, or an error if the text is not valid
//
//	TODO: See if a larger object can be passed in to avoid multiple calls to this function
func CompareToRetainRemoveFlags(text string, retainMap map[string]int, removeMap map[string]int) (string, error) {
	// If the retain and remove maps are empty, return the text as is
	if len(retainMap) == 0 && len(removeMap) == 0 {
		return text, nil
	}

	// If the retain map is empty, check the remove map
	if _, ok := retainMap[text]; ok || len(retainMap) == 0 {
		// If the remove map is empty, return the text as is
		if _, ok := removeMap[text]; !ok || len(removeMap) == 0 {
			return text, nil
		}
	}

	// If the text is not in the retain list, return an error
	return "", fmt.Errorf("text is not in the retain list")
}

// RemoveUnderMinFrequency removes items from a map that are below a minimum frequency
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
func RemoveUnderMinFrequency(freq map[string]int, min int) map[string]int {
	newFreq := make(map[string]int)
	for key, value := range freq {
		if value >= min {
			newFreq[key] = value
		}
	}
	return newFreq
}
