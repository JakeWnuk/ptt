// Package utils provides utility functions for the CLI application.
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// ----------------------------------------
// Loading and processing files function
// ----------------------------------------

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

// ----------------------------------------
// Transformation Functions
// ----------------------------------------

// ReverseString will return a string in reverse
//
// Args:
//
//	str (string): Input string to transform
//
// Returns:
//
//	(string): Transformed string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// ConvertMultiByteCharToRule converts non-ascii characters to a hashcat valid format
// for rule.CharToRule functions
//
// Args:
//
//	str (string): Input string to transform
//
// Returns:
//
//	returnStr (string): Converted string
func ConvertMultiByteCharToRule(str string) string {
	returnStr := ""
	deletedChar := ``
	for i, r := range str {
		if r > 127 {
			if i > 0 {
				deletedChar = string(returnStr[len(returnStr)-1])
				returnStr = returnStr[:len(returnStr)-1]
			}
			byteArr := []byte(string(r))
			if deletedChar == "^" {
				for j := len(byteArr) - 1; j >= 0; j-- {
					b := byteArr[j]
					if j == 0 {
						returnStr += fmt.Sprintf("%s\\x%X", deletedChar, b)
					} else {
						returnStr += fmt.Sprintf("%s\\x%X ", deletedChar, b)
					}
				}
			} else {
				for j, b := range byteArr {
					if j == len(byteArr)-1 {
						returnStr += fmt.Sprintf("%s\\x%X", deletedChar, b)
					} else {
						returnStr += fmt.Sprintf("%s\\x%X ", deletedChar, b)
					}
				}
			}
		} else {
			returnStr += fmt.Sprintf("%c", r)
		}
	}
	return returnStr
}

// ConvertMultiByteCharToIteratingRule converts non-ascii characters to a hashcat valid format
// for rule.CharToIteratingRule functions
//
// Args:
//
//	str (string): Input string to transform
//
// Returns:
//
//	returnStr (string): Converted string
func ConvertMultiByteCharToIteratingRule(str string) string {
	returnStr := ""
	deletedChar := ``
	for i, r := range str {
		if r > 127 {
			if i > 0 {
				deletedChar = string(returnStr[len(returnStr)-1])
				returnStr = returnStr[:len(returnStr)-1]
			}
			byteArr := []byte(string(r))
			if deletedChar == "^" {
				for j := len(byteArr) - 1; j >= 0; j-- {
					b := byteArr[j]
					if j == 0 {
						returnStr += fmt.Sprintf("%s\\x%X", deletedChar, b)
					} else {
						returnStr += fmt.Sprintf("%s\\x%X ", deletedChar, b)
					}
				}
			} else {
				for j, b := range byteArr {
					if j == len(byteArr)-1 {
						returnStr += fmt.Sprintf("%s\\x%X", deletedChar, b)
					} else {
						returnStr += fmt.Sprintf("%s\\x%X ", deletedChar, b)
					}
				}
			}
		} else {
			returnStr += fmt.Sprintf("%c", r)
		}
	}
	return returnStr
}

// ----------------------------------------
// Validation Functions
// ----------------------------------------

// CheckASCIIString checks to see if a string only contains ascii characters
//
// Args:
//
//	str (string): Input string to check
//
// Returns:
//
//	(bool): If the string only contained ASCII characters
func CheckASCIIString(str string) bool {
	if utf8.RuneCountInString(str) != len(str) {
		return false
	}
	return true
}
