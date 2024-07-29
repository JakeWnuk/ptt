// Package utils provides utility functions for the application.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/jakewnuk/ptt/pkg/models"

	"golang.org/x/net/html"
)

// ----------------------------------------------------------------------------
// Loading and Processing Functions
// ----------------------------------------------------------------------------

// ReadFilesToMap reads the contents of the multiple files and returns a map of words
//
// Args:
//
//	fs (FileSystem): The filesystem to read the files from (used for testing)
//	filenames ([]string): The names of the files to read
//
// Returns:
//
//	(map[string]int): A map of words from the files
func ReadFilesToMap(fs models.FileSystem, filenames []string) map[string]int {
	wordMap := make(map[string]int)

	// Read the contents of the files and add the words to the map
	for _, filename := range filenames {
		data, err := fs.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error reading file %s\n", filename)
			os.Exit(1)
		}

		err = json.Unmarshal(data, &wordMap)
		if err == nil {
			fmt.Fprintf(os.Stderr, "[*] Detected ptt JSON output. Importing...\n")
			continue
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
//	scanner (models.Scanner): The scanner to read from stdin
//
// Returns:
//
//	map[string]int: A map of lines from stdin
//	error: An error if one occurred
func LoadStdinToMap(scanner models.Scanner) (map[string]int, error) {
	m := make(map[string]int)
	pttInput := false
	line0 := false
	reDetect := regexp.MustCompile(`^\d+\s(\d+|\w+)`)
	reParse := regexp.MustCompile(`^\d+`)

	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}

		// Detect ptt -v output
		if matched := reDetect.MatchString(scanner.Text()); matched && pttInput == false && line0 == false {
			fmt.Fprintf(os.Stderr, "[*] Detected ptt -v output. Importing...\n")
			pttInput = true
		}

		if pttInput {
			line := scanner.Text()
			match := reParse.FindString(line)
			value, err := strconv.Atoi(match)
			if err != nil {
				return nil, err
			}
			newLine := strings.TrimSpace(strings.Replace(line, match, "", 1))
			m[newLine] += value

		} else {
			line := scanner.Text()
			m[line]++
		}
		line0 = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

// ReadURLsToMap reads the contents of the multiple URLs and returns a map of words
// from the URLs
//
// Args:
//
//	urls ([]string): The URLs to read
//
// Returns:
//
//	map[string]int: A map of words from the URLs
//	error: An error if one occurred
func ReadURLsToMap(urls []string) (map[string]int, error) {
	wordMap := make(map[string]int)
	var wg sync.WaitGroup

	ch := make(chan string)

	go func() {
		for word := range ch {
			wordMap[word]++
		}
	}()

	for _, url := range urls {
		wg.Add(1)
		go ProcessURL(url, ch, &wg)
	}

	wg.Wait()
	close(ch)

	delete(wordMap, "")

	return wordMap, nil
}

// CombineMaps combines any number of maps into a single map combining values for common keys
// and returning a new map
//
// Args:
// maps ([]map[string]int): The maps to combine
//
// Returns:
// map[string]int: A new map combining the values of the input maps
func CombineMaps(maps ...map[string]int) map[string]int {
	result := make(map[string]int)

	for _, m := range maps {
		for k, v := range m {
			result[k] += v
		}
	}

	return result
}

// ProcessURL reads the contents of a URL and sends each sentence to the channel
//
// Args:
//
//	url (string): The URL to read
//	ch (chan<- string): The channel to send the sentences to
//	wg (*sync.WaitGroup): The WaitGroup to signal when done
//
// Returns:
//
//	None
func ProcessURL(url string, ch chan<- string, wg *sync.WaitGroup) {
	const maxRetries = 4
	defer wg.Done()

	var err error
	var resp *http.Response
	for attempts := 0; attempts <= maxRetries; attempts++ {

		resp, err = http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error fetching URL %s\n", url)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			time.Sleep(time.Second)
			continue
		}

		break
	}

	// Read Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error reading response body from URL %s\n", url)
		os.Exit(1)
	}
	text := string(body)
	text = html.UnescapeString(text)
	var lines []string

	// Check the Content-Type of the response
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		// Parse the HTML
		doc, err := html.Parse(strings.NewReader(text))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error parsing HTML from URL %s\n", url)
			os.Exit(1)
		}

		// Traverse the HTML tree and extract the text
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.TextNode {
				lines = append(lines, n.Data)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
	} else {
		sentences := strings.Split(text, "\n")
		for _, line := range sentences {
			lines = append(lines, line)
		}
	}

	// Iterate over the lines and split them
	for _, line := range lines {
		textMatch, _ := regexp.MatchString(`[^a-zA-Z0-9.,;:!?'"\- ]`, line)
		if strings.Contains(contentType, "text/html") {
			if textMatch {
				continue
			}
		} else {
			if !textMatch {
				continue
			}
		}

		sentences := strings.Split(line, ".")
		for _, sentence := range sentences {
			sentence = strings.TrimSpace(sentence)

			phrases := strings.Split(sentence, ",")
			for _, phrase := range phrases {
				if phrase != "" {
					ch <- phrase
				}
			}

			if sentence != "" {
				ch <- sentence
			}
		}
	}
}

// ReadJSONToArray reads the contents of a transformation template file and
// returns a slice of template structs.
//
// Args:
//
//	fs (FileSystem): The filesystem to read the file from (used for testing)
//	fileArray ([]string): The name of the files to read
//
// Returns:
//
//	templates ([]models.TemplateFileOperation): The slice of template structs
func ReadJSONToArray(fs models.FileSystem, filenames []string) []models.TemplateFileOperation {
	var templates []models.TemplateFileOperation

	for _, filename := range filenames {
		// Check to see if a directory was passed
		// If so, read all files in the directory and append them to the filenames
		// slice

		data, err := fs.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error reading file %s\n", filename)
			os.Exit(1)
		}

		err = json.Unmarshal(data, &templates)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error unmarshalling JSON file %s\n", filename)
			os.Exit(1)
		}

		alphaRe := regexp.MustCompile(`[a-zA-Z]`)
		numRe := regexp.MustCompile(`[0-9]`)

		for _, template := range templates {
			if !numRe.MatchString(fmt.Sprintf("%v", template.StartIndex)) || !numRe.MatchString(fmt.Sprintf("%v", template.EndIndex)) {
				fmt.Fprintf(os.Stderr, "[!] Error: StartIndex and EndIndex must be integers\n")
				os.Exit(1)
			}

			if !alphaRe.MatchString(fmt.Sprintf("%v", template.Verbose)) {
				fmt.Fprintf(os.Stderr, "[!] Error: Verbose must be a boolean\n")
				os.Exit(1)
			}

			if !alphaRe.MatchString(fmt.Sprintf("%v", template.ReplacementMask)) {
				fmt.Fprintf(os.Stderr, "[!] Error: ReplacementMask must be a string\n")
				os.Exit(1)
			}

			if !alphaRe.MatchString(fmt.Sprintf("%v", template.Bypass)) {
				fmt.Fprintf(os.Stderr, "[!] Error: Bypass must be a boolean\n")
				os.Exit(1)
			}

			if !alphaRe.MatchString(fmt.Sprintf("%v", template.TransformationMode)) {
				fmt.Fprintf(os.Stderr, "[!] Error: TransformationMode must be a string\n")
				os.Exit(1)
			}

			if !numRe.MatchString(fmt.Sprintf("%v", template.PassphraseWords)) {
				fmt.Fprintf(os.Stderr, "[!] Error: PassphraseWords must be an integer\n")
				os.Exit(1)
			}
		}
	}

	return templates
}

// ----------------------------------------------------------------------------
// Transformation Functions
// ----------------------------------------------------------------------------

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

// IncrementIteratingRuleCall increments the last character of a string for
// rules.CharToIteratingRules functions
//
// For example, "i4" will be incremented to "i5", "iA" will be incremented to
// "IB"
//
// Args:
//
//	s (string): Input string to increment
//
// Returns:
//
//	output (string): Incremented string
func IncrementIteratingRuleCall(s string) string {
	if len(s) == 0 {
		return s
	}

	lastChar := s[len(s)-1]
	incChar := lastChar + 1

	// Replace the last character with the incremented character
	output := s[:len(s)-1] + string(incChar)

	return output
}

// ConvertMultiByteCharToIteratingRule converts non-ascii characters to a hashcat valid format
// for rule.CharToIteratingRule functions
//
// Args:
//
//	index (int): Index to start the iteration
//	str (string): Input string to transform
//
// Returns:
//
//	returnStr (string): Converted string
func ConvertMultiByteCharToIteratingRule(index int, str string) string {
	output := ""
	lastIterationSeen := fmt.Sprintf("%s%d", string([]rune(str)[0]), index)

	re := regexp.MustCompile(`[io][\dA-Z]`)

	for _, word := range strings.Split(str, " ") {
		for _, c := range word {
			if c > 127 {
				// Convert to UTF-8 bytes
				bytes := []byte(string(c))
				firstByteOut := true
				// Convert each byte to its hexadecimal representation
				for i, b := range bytes {
					if firstByteOut {
						output += fmt.Sprintf("\\x%X ", b)
						firstByteOut = false
						continue
					}
					lastIterationSeen = IncrementIteratingRuleCall(lastIterationSeen)
					if i == len(bytes)-1 {
						output += fmt.Sprintf("%s\\x%X", lastIterationSeen, b)
					} else {
						output += fmt.Sprintf("%s\\x%X ", lastIterationSeen, b)
					}
				}
			} else {
				output += string(c)
				if len(output) > 2 && re.MatchString(output[len(output)-2:]) {
					lastIterationSeen = output[len(output)-2:]
				}
			}
		}
		output += " "
	}

	return output
}

// SplitBySeparatorString splits a string by a separator string and returns a slice
// with the separator string included
//
// Args:
//
//	s (string): The string to split
//	sep (string): The separator string
//
// Returns:
//
//	[]string: A slice of strings with the separator string included
func SplitBySeparatorString(s string, sep string) []string {
	if !strings.Contains(s, sep) {
		return []string{s}
	}

	// Limit to 2 to ensure we only split on the first instance of the separator
	parts := strings.SplitN(s, sep, 2)
	parts = append(parts[:1], append([]string{sep}, parts[1:]...)...)
	return parts
}

// ReplaceSubstring replaces all instances of a substring in a string with a new
// substring if the substring is found in the original string. The new substring
// is determined by the key in the replacements map separated by a colon
// character.
//
// Args:
//
//	original (string): The original string
//	replacements (map[string]int): A map of substrings to replace
//
// Returns:
//
//	[]string: The original string with all instances of the substring replaced
func ReplaceSubstring(original string, replacements map[string]int) []string {
	var newStrings []string
	for newSubstr := range replacements {
		// Split the new substring into the old and new strings by the colon character
		if !strings.Contains(newSubstr, ":") {
			continue
		}
		oldStr, newStr := strings.Split(newSubstr, ":")[0], strings.Split(newSubstr, ":")[1]
		if strings.Contains(original, oldStr) {
			newStrings = append(newStrings, strings.Replace(original, oldStr, newStr, -1))
		}
	}
	return newStrings
}

// SubstringMap returns a map of substrings from a map of strings starting at
// the start index and ending at the end index. If the bypass flag is set to
// true, the function will print to stdout and return an empty map. If the
// end index is greater than the length of the string, the function will use
// the length of the string as the end index for that string.
//
// Args:
//
//	sMap (map[string]int): The map of substrings
//	sIndex (int): The start index of the substring
//	eIndex (int): The end index of the substring
//	bypass (bool): Skip returning the map and print to stdout
//	debug (bool): A flag to print debug information
//
// Returns:
//
//	map[string]int: A map of substrings
func SubstringMap(sMap map[string]int, sIndex int, eIndex int, bypass bool, debug bool) map[string]int {
	newMap := make(map[string]int)
	for s := range sMap {
		maxLen := eIndex
		if sIndex > len(s) {
			if debug {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("[!] Error: Start index is out of bounds: %s\n", s))
			}
			continue
		} else if eIndex > len(s) {
			maxLen = len(s)
		}

		if bypass {
			fmt.Fprintf(os.Stdout, "%s\n", s[sIndex:maxLen])
			continue
		}
		newMap[s[sIndex:maxLen]]++
	}
	return newMap
}

// ----------------------------------------------------------------------------
// Validation Functions
// ----------------------------------------------------------------------------

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

// CheckHexString is used to identify plaintext in the $HEX[...] format
//
// Args:
//
//	s (str): The string to be evaluated
//
// Returns:
//
//	(bool): Returns true if it matches and false if it did not
func CheckHexString(s string) bool {
	var validateInput = regexp.MustCompile(`^\$HEX\[[a-zA-Z0-9]*\]$`).MatchString
	if validateInput(s) == false {
		return false
	}
	return true
}

// CheckAreMapsEqual checks if two maps are equal by comparing the length of the maps
// and the values of the keys in the maps. If the maps are equal, the function returns
// true, otherwise it returns false.
//
// Args:
//
//	a (map[string]int): The first map to compare
//	b (map[string]int): The second map to compare
//
// Returns:
//
//	bool: True if the maps are equal, false otherwise
func CheckAreMapsEqual(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}

// CheckAreArraysEqual checks if two arrays are equal by comparing the length of the arrays
// and the values of the elements in the arrays. If the arrays are equal, the function returns
// true, otherwise it returns false.
//
// Args:
// a ([]string): The first array to compare
// b ([]string): The second array to compare
//
// Returns:
// bool: True if the arrays are equal, false otherwise
func CheckAreArraysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// IsFileSystemDirectory checks to see if a string is a valid file system
// directory by checking if the path exists and if it is a directory
//
// Args:
//  path (string): The path to check
//
// Returns:
//  bool: True if the path is a directory, false otherwise
func IsFileSystemDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return fileInfo.IsDir()
}
