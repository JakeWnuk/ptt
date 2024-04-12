package utils

import (
	"bufio"
	"os"
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

// LoadStdinToMap reads the contents of stdin and returns a map[string]int
// where the key is the line and the value is the frequency of the line
// in the input
//
// Args:
//
//	None
//
// Returns:
//
//	map[string]int: A map of lines from stdin
//	error: An error if one occurred
func LoadStdinToMap() (map[string]int, error) {
	m := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		m[line]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// CombineMaps combines two maps into a single map combining values for common keys
// and returning a new map
//
// Args:
// map1 (map[string]int): The first map
// map2 (map[string]int): The second map
//
// Returns:
// map[string]int: A new map combining the values of the two input maps
func CombineMaps(map1, map2 map[string]int) map[string]int {
	result := make(map[string]int)

	for k, v := range map1 {
		result[k] = v
	}

	for k, v := range map2 {
		result[k] += v
	}

	return result
}
