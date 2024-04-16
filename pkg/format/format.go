// Package format contains functions for printing and formatting output
package format

import (
	"encoding/hex"
	"fmt"
	"html"
	"net/url"
	"os"
	"ptt/pkg/models"
	"sort"
	"strings"
	"unicode"
)

// ----------------------------------------------------------------------------
// Output Functions
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

// PrintStatsToSTDOUT prints statistics about the frequency map to stdout
// including several statistics about the frequency map. If verbose is true,
// additional information is printed and increased number of items are
// printed. Items are printed in graph format with a # for each unit of item.
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//	verbose (bool): If true, additional information is printed
//	max (int): The maximum number of items to print
//
// Returns:
//
//	None
func PrintStatsToSTDOUT(freq map[string]int, verbose bool, max int) {

	// Set the count and max values
	count := 0
	if !verbose {
		max = 10
	}

	// Sort by frequency
	p := make(models.PairList, len(freq))
	normalizedP := make(models.PairList, len(freq))
	i := 0
	for k, v := range freq {
		p[i] = models.Pair{k, v}
		normalizedP[i] = models.Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	sort.Sort(sort.Reverse(normalizedP))

	// Get the largest frequency value by getting the first item
	largest := p[0].Value

	// Print the statistics
	if verbose {
		fmt.Fprintf(os.Stderr, "Starting statistics generation. Please wait...\n")
		fmt.Println(fmt.Sprintf("Verbose Statistics: max=%d", max))
		fmt.Println("--------------------------------------------------")
		fmt.Println(CreateVerboseStats(freq))
		fmt.Println("--------------------------------------------------")
	}

	// Use the largest frequency value to normalize the graph
	for index, value := range normalizedP {
		normalizedValue := value.Value * 50 / largest
		normalizedP[index].Value = normalizedValue
	}

	// Print the top items
	for index, value := range p {
		if value.Value == 1 && index == 0 {
			fmt.Println("No items with a frequency greater than 1!")
			break
		}

		if value.Value == 1 {
			continue
		}

		if count < max {
			fmt.Printf("%s [%d]%s\n", value.Key, value.Value, strings.Repeat("=", normalizedP[index].Value))
			count++
		} else {
			count = 0
			break
		}
	}
}

// CreateVerboseStats creates a string of verbose statistics about the frequency map
// including several statistics about the frequency map.
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//
// Returns:
//
//	string: A string of verbose statistics
func CreateVerboseStats(freq map[string]int) string {
	var stats string
	// Sort by frequency
	p := make(models.PairList, len(freq))
	i := 0
	for k, v := range freq {
		p[i] = models.Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))

	// Pull stats
	totalChars := 0
	totalWords := 0
	categoryCounts := make(map[string]int)
	for k := range freq {
		totalChars += len(k)
		totalWords += len(strings.Fields(k))
		categories := StatClassifyToken(k)
		for _, category := range categories {
			categoryCounts[category]++
		}
	}
	stats += "General Stats:\n"
	stats += fmt.Sprintf("Total Unique items: %d\n", len(p))
	stats += fmt.Sprintf("Total Characters: %d\n", totalChars)
	stats += fmt.Sprintf("Total Words: %d\n", totalWords)
	stats += fmt.Sprintf("Average Characters Per Item: %d\n", totalChars/len(freq))
	stats += fmt.Sprintf("Average Words Per Item: %d\n", totalWords/len(freq))
	stats += fmt.Sprintf("Largest frequency: %d\n", p[0].Value)
	stats += fmt.Sprintf("Smallest frequency: %d\n", p[len(p)-1].Value)
	stats += fmt.Sprintf("Mean frequency: %d\n", p[len(p)/2].Value)
	stats += fmt.Sprintf("Median frequency: %d\n", p[len(p)/2].Value)
	stats += fmt.Sprintf("Mode frequency: %d", p[0].Value)

	stats += "\n\nCategory Counts:\n"
	for category, count := range categoryCounts {
		stats += fmt.Sprintf("%s: %d\n", category, count)
	}

	return stats
}

// StatClassifyToken classifies a token into a set of categories
// based on the token's content
//
// Args:
//
//	s (string): The token to classify
//
// Returns:
//
//	[]string: A list of categories that the token belongs to
func StatClassifyToken(s string) []string {
	var categories []string

	isAlpha := func(c rune) bool { return unicode.IsLetter(c) }
	isDigit := func(c rune) bool { return unicode.IsDigit(c) }
	isSpecial := func(c rune) bool { return !unicode.IsLetter(c) && !unicode.IsDigit(c) && !unicode.IsSpace(c) }

	if strings.IndexFunc(s, isAlpha) >= 0 {
		categories = append(categories, "alphabetical")
	}

	if strings.IndexFunc(s, isDigit) >= 0 {
		categories = append(categories, "numeric")
	}

	if strings.IndexFunc(s, isAlpha) >= 0 && strings.IndexFunc(s, isDigit) >= 0 {
		categories = append(categories, "alphanumeric")
	}

	if strings.IndexFunc(s, isSpecial) >= 0 && (strings.IndexFunc(s, isAlpha) >= 0 || strings.IndexFunc(s, isDigit) >= 0) {
		categories = append(categories, "alphanumeric with special")
	}

	if strings.Contains(s, " ") {
		categories = append(categories, "phrase")
	}

	digitCount := 0
	for _, c := range s {
		if isDigit(c) {
			digitCount++
		}
	}
	if digitCount > len(s)*2/3 {
		categories = append(categories, "high numeric ratio")
	}

	if strings.IndexFunc(s, unicode.IsUpper) >= 0 && strings.IndexFunc(s, unicode.IsLower) >= 0 && strings.IndexFunc(s, isDigit) >= 0 && strings.IndexFunc(s, isSpecial) >= 0 {
		categories = append(categories, "complex")
	} else {
		categories = append(categories, "non-complex")
	}

	for _, c := range s {
		if c > 127 {
			categories = append(categories, "non-ASCII")
			break
		}
	}

	return categories
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
//	(map[string]int): A new map of words that have been retained and removed
//	(error): An error if the text map is empty
func RetainRemove(textMap map[string]int, retainMap map[string]int, removeMap map[string]int) (map[string]int, error) {
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

// RemoveLengthRange removes items from a map that are outside of a length
// range or not equal to the start of the end if no end is provided. The length
// of the range is inclusive.
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//	start (int): The start of the length range
//	end (int): The end of the length range
//
// Returns:
//
//	(map[string]int): A new map of item frequencies within the length range
func RemoveLengthRange(freq map[string]int, start int, end int) map[string]int {
	newFreq := make(map[string]int)
	for key, value := range freq {
		if end == 0 {
			if len(key) == start {
				newFreq[key] = value
			}
		} else {
			if len(key) >= start && len(key) <= end {
				newFreq[key] = value
			}
		}
	}
	return newFreq
}

// ----------------------------------------------------------------------------
// Encoding Functions
// ----------------------------------------------------------------------------

// EncodeInputMap will encode a map of strings to URL, HTML, and unicode escaped strings
// where possible and return a new map of encoded strings
//
// Args:
//
//	input (map[string]int): A map of input strings
//
// Returns:
//
//	(map[string]int): A new map of encoded strings
func EncodeInputMap(input map[string]int) map[string]int {
	output := make(map[string]int)
	for k, v := range input {
		urlEncoded, htmlEncoded, escapeEncoded := EncodeString(k)

		if urlEncoded != "" {
			output[urlEncoded] = v
		}

		if htmlEncoded != "" {
			output[htmlEncoded] = v
		}

		if escapeEncoded != "" {
			output[escapeEncoded] = v
		}
	}
	return output
}

// EncodeString is used to URL and HTML encode a string where possible this
// will only return the encoded string if it is different from the input string
//
// Args:
//
//	s (string): Input string
//
// Returns:
//
//	urlEncoded (string): Input string URL encoded
//	htmlEncoded (string): Input string HTML encoded
//	escapedEncoded (string): Input string unicode escaped encoded
func EncodeString(s string) (string, string, string) {
	urlEncoded := url.QueryEscape(s)
	htmlEncoded := html.EscapeString(s)
	escapedEncoded := ASCIIEscapeUnicode(s)

	if urlEncoded == s {
		urlEncoded = ""
	}

	if htmlEncoded == s {
		htmlEncoded = ""
	}

	if escapedEncoded == s {
		escapedEncoded = ""
	}

	return urlEncoded, htmlEncoded, escapedEncoded
}

// ASCIIEscapeUnicode will convert a string into an unicode escaped format
//
// Args:
//
//	str (string): String to escape
//
// Returns:
//
//	escapedRunes (string): Converted runes in string format
func ASCIIEscapeUnicode(str string) string {
	runes := []rune(str)
	escapedRunes := make([]rune, 0, len(runes))

	for _, r := range runes {
		if r > 127 {
			// The rune is non-ASCII
			escapedRune := []rune(fmt.Sprintf("\\u%04x", r))
			escapedRunes = append(escapedRunes, escapedRune...)
		} else {
			escapedRunes = append(escapedRunes, r)
		}
	}
	return string(escapedRunes)
}

// DehexMap will decode a map of hex encoded strings and return a new map of decoded strings
// where possible
//
// Args:
//
//	input (map[string]int): A map of hex encoded strings
//
// Returns:
//
//	(map[string]int): A new map of decoded strings
func DehexMap(input map[string]int) map[string]int {
	decodedMap := make(map[string]int)

	for k, v := range input {
		k = strings.TrimPrefix(k, "$HEX[")
		k = strings.TrimSuffix(k, "]")
		decoded, err := hex.DecodeString(k)
		if err != nil {
			continue
		}
		decodedStr := string(decoded)
		decodedMap[decodedStr] = v
	}

	return decodedMap
}

// HexEncodeMap will encode a map of strings to hex encoded strings
// where possible and return a new map of encoded strings
//
// Args:
//
//	input (map[string]int): A map of input strings
//
// Returns:
//
//	(map[string]int): A new map of encoded strings
func HexEncodeMap(input map[string]int) map[string]int {
	output := make(map[string]int)
	for k, v := range input {
		encoded := hex.EncodeToString([]byte(k))
		output["$HEX["+encoded+"]"] = v
	}
	return output
}
