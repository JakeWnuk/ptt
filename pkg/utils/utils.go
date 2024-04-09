package utils

import (
	"os"
	"strings"
)

// ReadFilesToMap reads the contents of the given files and returns a map of words
//
// Args:
//
//	filenames ([]string): The names of the files to read
//
// Returns:
//
//	(map[string]bool): A map of words from the files
func ReadFilesToMap(filenames []string) map[string]bool {
	wordMap := make(map[string]bool)
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		fileWords := strings.Split(string(data), "\n")
		for _, word := range fileWords {
			wordMap[word] = true
		}
	}
	return wordMap
}
