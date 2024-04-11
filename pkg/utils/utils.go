package utils

import (
	"os"
	"sort"
	"strings"
)

// ReadFilesToMap reads the contents of the multiple files and returns a map of words
//
// Args:
//
//	filenames ([]string): The names of the files to read
//
// Returns:
//
//	(map[string]int): A map of words from the files
func ReadFilesToMap(filenames []string) map[string]int {
	wordMap := make(map[string]int)

	// Read the contents of the files and add the words to the map
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		fileWords := strings.Split(string(data), "\n")
		for _, word := range fileWords {
			wordMap[word]++
		}
	}

	// Remove empty strings from the map
	delete(wordMap, "")

	return wordMap
}

// FrequencySortMap sorts a map of items by frequency
// in descending order
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//
// Returns:
//
//	[]string: A slice of items sorted by frequency
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
