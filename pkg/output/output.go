package output

import (
	"fmt"
	"sort"
)

// FrequencySortMap sorts a map of items by frequency
// in descending order
//
// Args:
// freq (map[string]int): A map of item frequencies
//
// Returns:
// []string: A slice of items sorted by frequency
func FrequencySortMap(freq map[string]int) []string {
	items := make([]string, 0, len(freq))

	for item := range freq {
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return freq[items[i]] > freq[items[j]]
	})

	return items
}

// PrintArrayToSTDOUT prints an array of items to stdout
//
// Args:
//
//	items ([]string): The array of items to print
//	freq (map[string]int): A map of item frequencies
//	verbose (bool): If true, the count is printed
//	min (int): Minimum frequency of an item to print
//
// Returns:
//
//	None
func PrintArrayToSTDOUT(items []string, freq map[string]int, verbose bool, min int) {
	count := 0
	for _, item := range items {
		if freq[item] >= min {
			if verbose {
				fmt.Printf("%d %s\n", freq[item], item)
			} else {
				fmt.Println(item)
			}
			count++
		}
	}

}
