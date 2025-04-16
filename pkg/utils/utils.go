// Package utils provides utility functions for various tasks
package utils

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/jakewnuk/ptt/pkg/mask"
	"github.com/jakewnuk/ptt/pkg/models"
	"github.com/jakewnuk/ptt/pkg/validation"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ReverseString will return a string in reverse
//
// Args:
// str (string): Input string to transform
//
// Returns:
// (string): Transformed string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// GeneratePassphrase generates a passphrase from a string of text
// and returns a slice of passphrases
//
// Args:
// text (string): The text to generate passphrases from
// n (int): The number of words in the passphrase
//
// Returns:
// []string: A slice of passphrases
func GeneratePassphrase(text string, n int) []string {
	text = strings.ReplaceAll(text, ".", "")
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, ";", "")
	words := strings.Fields(text)
	var passphrases []string

	if len(words) != n {
		return passphrases
	}

	var titleCaseWords string
	var turkTitleCaseWords string
	var CAPSlowerWords []string
	var lowerCAPSWords []string
	tick := false
	titleCaseWords = cases.Title(language.Und, cases.NoLower).String(text)
	turkTitleCaseWords = cases.Upper(language.Turkish, cases.NoLower).String(text)

	for _, word := range words {

		if tick {
			CAPSlowerWords = append(CAPSlowerWords, strings.ToUpper(word))
			lowerCAPSWords = append(lowerCAPSWords, strings.ToLower(word))
		} else {
			CAPSlowerWords = append(CAPSlowerWords, strings.ToLower(word))
			lowerCAPSWords = append(lowerCAPSWords, strings.ToUpper(word))
		}

		tick = !tick

	}

	passphrases = append(passphrases, strings.ReplaceAll(titleCaseWords, " ", ""))
	passphrases = append(passphrases, strings.ReplaceAll(turkTitleCaseWords, " ", ""))
	passphrases = append(passphrases, strings.ReplaceAll(titleCaseWords, " ", "-"))
	passphrases = append(passphrases, strings.ReplaceAll(turkTitleCaseWords, " ", "-"))
	passphrases = append(passphrases, strings.ReplaceAll(titleCaseWords, " ", "_"))
	passphrases = append(passphrases, strings.ReplaceAll(turkTitleCaseWords, " ", "_"))
	passphrases = append(passphrases, strings.Join(CAPSlowerWords, " "))
	passphrases = append(passphrases, strings.Join(CAPSlowerWords, ""))
	passphrases = append(passphrases, strings.Join(lowerCAPSWords, " "))
	passphrases = append(passphrases, strings.Join(lowerCAPSWords, ""))
	passphrases = append(passphrases, strings.ReplaceAll(text, " ", ""))
	passphrases = append(passphrases, titleCaseWords)
	passphrases = append(passphrases, turkTitleCaseWords)

	return passphrases
}

// GenerateNGrams generates n-grams from a string of text
// and returns a slice of n-grams
//
// Args:
// text (string): The text to generate n-grams from
// n (int): The number of words in each n-gram
//
// Returns:
// []string: A slice of n-grams
func GenerateNGrams(text string, n int) []string {
	words := strings.Fields(text)
	var nGrams []string

	for i := 0; i <= len(words)-n; i++ {
		nGram := strings.Join(words[i:i+n], " ")
		nGrams = append(nGrams, nGram)
	}

	return nGrams
}

// PrintStatsToSTDOUT prints statistics about the frequency map to stdout
// including several statistics about the frequency map. If verbose is true,
// additional information is printed and increased number of items are
// printed. Items are printed in graph format with a # for each unit of item.
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//
// Returns:
//
//	None
func PrintStatsToSTDOUT(freq map[string]int) {

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

	maxItems := 75
	if maxItems > len(p) {
		maxItems = len(p)
	}

	if len(p) == 0 {
		fmt.Println("[!] No items to print!")
		return
	}

	// Print the statistics
	fmt.Fprintf(os.Stderr, "[*] Starting statistics generation. Please wait...\n")
	fmt.Println(fmt.Sprintf("Verbose Statistics: max=%d", maxItems))
	fmt.Println("--------------------------------------------------")
	fmt.Println(createVerboseStats(freq))
	fmt.Println("--------------------------------------------------")

	// Use the largest frequency value to normalize the graph
	largest := p[0].Value
	for index, value := range normalizedP {
		normalizedValue := value.Value * 50 / largest
		normalizedP[index].Value = normalizedValue
	}

	// Use the longest key to normalize padding for the graph
	longest := 0
	for _, value := range p[0:maxItems] {
		if len(value.Key) > longest {
			longest = len(value.Key)
		}
	}

	// Print the top items
	for index, value := range p[0:maxItems] {
		if value.Value == 1 && index == 0 {
			fmt.Println("[!] No items with a frequency greater than 1!")
			break
		}
		padding := longest - len(value.Key)
		fmt.Printf("%s%s [%d]%s\n", value.Key, strings.Repeat(" ", padding), value.Value, strings.Repeat("=", normalizedP[index].Value))
	}
}

// createVerboseStats creates a string of verbose statistics about the frequency map
// including several statistics about the frequency map.
//
// Args:
//
//	freq (map[string]int): A map of item frequencies
//
// Returns:
//
//	string: A string of verbose statistics
func createVerboseStats(freq map[string]int) string {
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
	totalWords := 0
	totalItems := 0
	lengths := make([]int, 0)
	frequencies := make([]int, 0)
	complexities := make([]int, 0)
	categoryCounts := make(map[string]int)
	for k, v := range freq {
		totalWords += len(strings.Fields(k))
		categories := StatClassifyToken(k)
		frequencies = append(frequencies, v)
		totalItems += v

		for i := 0; i < v; i++ {
			lengths = append(lengths, len(k))
		}

		m := mask.MakeMask(k)
		complexity := mask.TestMaskComplexity(m)
		complexities = append(complexities, complexity)

		for _, category := range categories {
			categoryCounts[category]++
		}
	}
	stats += "General Stats:\n"
	stats += fmt.Sprintf("Total Items: %d\n", totalItems)
	stats += fmt.Sprintf("Total Unique items: %d\n", len(p))
	stats += fmt.Sprintf("Total Words: %d\n", totalWords)
	stats += fmt.Sprintf("Largest frequency: %d\n", p[0].Value)
	stats += fmt.Sprintf("Smallest frequency: %d\n", p[len(p)-1].Value)

	stats += "\nPlots:\n"
	plot, minBW, q1, q2, q3, maxBW := CreateBoxAndWhiskersPlot(lengths)
	stats += fmt.Sprintf("Item Length: %s\n", plot)
	stats += fmt.Sprintf("Min: %d, Q1: %d, Q2: %d, Q3: %d, Max: %d\n", minBW, q1, q2, q3, maxBW)
	plot, minBW, q1, q2, q3, maxBW = CreateBoxAndWhiskersPlot(frequencies)
	stats += fmt.Sprintf("Item Frequency: %s\n", plot)
	stats += fmt.Sprintf("Min: %d, Q1: %d, Q2: %d, Q3: %d, Max: %d\n", minBW, q1, q2, q3, maxBW)
	plot, minBW, q1, q2, q3, maxBW = CreateBoxAndWhiskersPlot(complexities)
	stats += fmt.Sprintf("Item Complexity: %s\n", plot)
	stats += fmt.Sprintf("Min: %d, Q1: %d, Q2: %d, Q3: %d, Max: %d\n", minBW, q1, q2, q3, maxBW)

	stats += "\nCategory Counts:\n"
	for category, count := range categoryCounts {
		stats += fmt.Sprintf("%s: %d\n", category, count)
	}

	return stats
}

// StatClassifyToken classifies a token into a set of categories
// based on the token's content. "Short" and "long" are relative
// to ten characters currently.
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
		if strings.IndexFunc(s, isSpecial) >= 0 && (strings.IndexFunc(s, isAlpha) >= 0 || strings.IndexFunc(s, isDigit) >= 0) {
			categories = append(categories, "alphanumeric-with-special")
		} else if strings.IndexFunc(s, isAlpha) >= 0 && strings.IndexFunc(s, isDigit) >= 0 && strings.IndexFunc(s, isSpecial) == -1 {
			categories = append(categories, "alphanumeric")
		} else if strings.IndexFunc(s, isAlpha) >= 0 && strings.IndexFunc(s, isSpecial) >= 0 && strings.IndexFunc(s, isDigit) == -1 {
			categories = append(categories, "alphabetical-with-special")
		} else {
			categories = append(categories, "alphabetical")
		}
	}

	if strings.IndexFunc(s, isDigit) >= 0 && strings.IndexFunc(s, isAlpha) == -1 {
		if strings.IndexFunc(s, isDigit) >= 0 && strings.IndexFunc(s, isSpecial) >= 0 {
			categories = append(categories, "numeric-with-special")
		} else {
			categories = append(categories, "numeric")
		}
	}

	if strings.Count(s, " ") >= 2 && strings.IndexFunc(s, isAlpha) >= 0 {
		categories = append(categories, "phrase")
	}

	digitCount := 0
	for _, c := range s {
		if isDigit(c) {
			digitCount++
		}
	}
	if digitCount > len(s)*5/7 {
		categories = append(categories, "high-numeric-ratio")
	}

	if strings.HasPrefix(s, "$HEX[") && strings.HasSuffix(s, "]") {
		categories = append(categories, "$HEX[...]-format")
	}

	if _, err := hex.DecodeString(s); err == nil {
		categories = append(categories, "hex-string")
	}

	if _, err := url.ParseRequestURI(s); err == nil && strings.Contains(s, "://") {
		categories = append(categories, "URL")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x1F600 && r <= 0x1F64F }) >= 0 {
		categories = append(categories, "emoji-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x4E00 && r <= 0x9FFF }) >= 0 {
		categories = append(categories, "CJK-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x0400 && r <= 0x04FF }) >= 0 {
		categories = append(categories, "cyrillic-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x0600 && r <= 0x06FF }) >= 0 {
		categories = append(categories, "arabic-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x0590 && r <= 0x05FF }) >= 0 {
		categories = append(categories, "hebrew-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x0E00 && r <= 0x0E7F }) >= 0 {
		categories = append(categories, "thai-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x0370 && r <= 0x03FF }) >= 0 {
		categories = append(categories, "greek-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0xAC00 && r <= 0xD7AF }) >= 0 {
		categories = append(categories, "korean-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x3040 && r <= 0x309F }) >= 0 {
		categories = append(categories, "hiragana-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x30A0 && r <= 0x30FF }) >= 0 {
		categories = append(categories, "katakana-characters")
	}

	if strings.IndexFunc(s, func(r rune) bool { return r >= 0x31F0 && r <= 0x31FF }) >= 0 {
		categories = append(categories, "katakana-extended-characters")
	}

	if strings.IndexFunc(s, unicode.IsUpper) >= 0 {
		if strings.IndexFunc(s, unicode.IsUpper) == 0 {
			categories = append(categories, "starts-uppercase")
		} else {
			categories = append(categories, "contains-uppercase")
		}
	}

	if strings.IndexFunc(s, unicode.IsLower) == -1 && strings.IndexFunc(s, unicode.IsUpper) >= 0 {
		categories = append(categories, "all-uppercase")
	}

	if strings.IndexFunc(s, unicode.IsUpper) == -1 && strings.IndexFunc(s, unicode.IsLower) >= 0 {
		categories = append(categories, "all-lowercase")
	}

	if strings.IndexFunc(s, unicode.IsUpper) >= 0 && strings.IndexFunc(s, unicode.IsLower) >= 0 && strings.IndexFunc(s, isDigit) >= 0 && strings.IndexFunc(s, isSpecial) >= 0 {
		categories = append(categories, "complex")
		if len(s) > 10 {
			categories = append(categories, "long-complex")
		} else {
			categories = append(categories, "short-complex")
		}
	} else {
		categories = append(categories, "non-complex")
		if len(s) > 10 {
			categories = append(categories, "long-non-complex")
		} else {
			categories = append(categories, "short-non-complex")
		}
	}

	for _, c := range s {
		if c > 127 {
			categories = append(categories, "non-ASCII")
			break
		}
	}

	return categories
}

// CalculateQuartiles calculates the first, second, and third quartiles of a
// list of integers and returns the values.
//
// Args:
// data ([]int): A list of integers
//
// Returns:
// int: The first quartile value
// int: The second quartile value
// int: The third quartile value
func CalculateQuartiles(data []int) (int, int, int) {
	sort.Ints(data)
	n := len(data)

	q1 := data[n/4]
	q2 := data[n/2]
	q3 := data[3*n/4]

	return q1, q2, q3
}

// CreateBoxAndWhiskersPlot creates a box and whiskers plot from a list of
// integers.
//
// Args:
//
//	data ([]int): A list of integers
//
//	Returns:
//	string: A string representation of the box and whiskers plot
//	int: The minimum value
//	int: The first quartile value
//	int: The second quartile value
//	int: The third quartile value
//	int: The maximum value
func CreateBoxAndWhiskersPlot(data []int) (string, int, int, int, int, int) {
	q1, q2, q3 := CalculateQuartiles(data)
	minBW := data[0]
	maxBW := data[len(data)-1]

	// Normalize the plot
	largest := maxBW
	normalizedQ1 := q1 * 50 / largest
	normalizedQ2 := q2 * 50 / largest
	normalizedQ3 := q3 * 50 / largest
	normalizedMin := minBW * 50 / largest
	normalizedMax := maxBW * 50 / largest

	plot := fmt.Sprintf("|%s[%s|%s]%s|", strings.Repeat("-", normalizedQ1-normalizedMin), strings.Repeat("=", normalizedQ2-normalizedQ1), strings.Repeat("=", normalizedQ3-normalizedQ2), strings.Repeat("-", normalizedMax-normalizedQ3))
	return plot, minBW, q1, q2, q3, maxBW
}

// MakeRetainMaskedMap replaces all characters in the input maps key with mask
// values in the input map but retains keywords provided in the retain list
//
// Args:
//
//	input (string): String to mask
//	retain (map[string]int): Map of keywords to retain
//
// Returns:
//
//	maskedMap (map[string]int): Masked retain map
func MakeRetainMaskedMap(input string, retain map[string]int) map[string]int {
	maskedMap := make(map[string]int)

	for retainKey := range retain {
		newKey := ""
		if strings.Contains(input, retainKey) {
			parts := splitBySeparatorString(input, retainKey)

			// if the part is not the key replace it using replacer
			for _, part := range parts {
				if part != retainKey {
					newPart := models.MaskReplacer.Replace(part)
					if !validation.CheckASCIIString(newPart) && strings.Contains(models.GlobalMask, "b") {
						newPart = validation.ConvertMultiByteMask(newPart)
					}
					newKey += newPart
				} else {
					newKey += part
				}
			}

		} else {
			// if the key is not in the string continue
			continue
		}

		if oldValue, exists := maskedMap[newKey]; exists {
			maskedMap[newKey] = oldValue + 1
		} else {
			maskedMap[newKey] = 1
		}
	}
	return maskedMap
}

// splitBySeparatorString splits a string by a separator string and returns a slice
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
func splitBySeparatorString(s string, sep string) []string {
	if !strings.Contains(s, sep) {
		return []string{s}
	}

	// Limit to 2 to ensure we only split on the first instance of the separator
	parts := strings.SplitN(s, sep, 2)
	parts = append(parts[:1], append([]string{sep}, parts[1:]...)...)
	return parts
}

// ShuffleMap shuffles the input map keys and replaces partially the masked
// parts of the keys with matching mask keys from the input map. This function
// resembles 'token-swapping' where the mask value is used to swap key words
// into another.
//
// Args:
//
//	input (map[string]int): Input map
//	swapMap (map[string]int): Items to swap with
//
// Returns:
// (map[string]int): Shuffled map with swapped keys
func ShuffleMap(input map[string]int, swapMap map[string]int) map[string]int {
	shuffleMap := make(map[string]int)
	re := regexp.MustCompile(`^(\?u|\?l|\?d|\?s|\?b)*$`)
	reParser := regexp.MustCompile("(\\?[ludsb])")

	for key, value := range input {
		newKey := ""
		// Make a new key with the masked parts
		chars := reParser.FindAllString(key, -1)
		match := strings.Join(chars, "")

		if re.MatchString(match) {
			newKey = match
		}

		// Check if the new key is in the swap map
		for swapKey := range swapMap {
			if models.DebugMode {
				fmt.Fprintf(os.Stderr, "[?] utils.ShuffleMap(input, swapMap)\n")
				fmt.Fprintf(os.Stderr, "Key: %s\n", key)
				fmt.Fprintf(os.Stderr, "Match: %s\n", match)
				fmt.Fprintf(os.Stderr, "Swap Token: %s\n", swapKey)
				fmt.Fprintf(os.Stderr, "Replacement Mask: %s\n", models.GlobalMask)
			}

			maskedSwapKey := mask.MakeMask(swapKey)
			if maskedSwapKey == newKey {

				var shufKey string
				shufKey = strings.Replace(key, newKey, swapKey, 1)

				if models.DebugMode {
					fmt.Fprintf(os.Stderr, "[?] utils.ShuffleMap(input, swapMap)\n")
					fmt.Fprintf(os.Stderr, "Swap Token Mask: %s\n", maskedSwapKey)
					fmt.Fprintf(os.Stderr, "Swap Result: %s\n", shufKey)
				}

				if shufKey == key {
					if models.DebugMode {
						fmt.Fprintf(os.Stderr, "[!] Swap failed invalid key:\n")
						fmt.Fprintf(os.Stderr, "Key: %s\n", key)
						fmt.Fprintf(os.Stderr, "Swap Result: %s\n", shufKey)
					}

					continue
				}

				// if the line ends or starts with ?[uldbs] then the swap failed
				if strings.HasPrefix(shufKey, "?") || strings.HasSuffix(shufKey[len(shufKey)-2:len(shufKey)-1], "?") {

					if strings.ContainsRune("uldbs", rune(shufKey[1])) && strings.HasPrefix(shufKey, "?") || strings.ContainsRune("uldbs", rune(shufKey[len(shufKey)-1])) && strings.HasSuffix(shufKey[len(shufKey)-2:len(shufKey)-1], "?") {

						if models.DebugMode {
							fmt.Fprintf(os.Stderr, "[!] Swap failed invalid key:\n")
							fmt.Fprintf(os.Stderr, "Key: %s\n", key)
							fmt.Fprintf(os.Stderr, "Swap Result: %s\n", shufKey)
						}

						continue

					}
				}

				if oldValue, exists := shuffleMap[shufKey]; exists {
					shuffleMap[shufKey] = oldValue + value
				} else {
					shuffleMap[shufKey] = value
				}

			}
		}
	}
	return shuffleMap
}

// TrackLoadTime tracks the load time of a process and prints the elapsed time
func TrackLoadTime(done <-chan bool, work string) {
	start := time.Now()
	interval := 10 * time.Second
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-done:
			ticker.Stop()
			return

		case t := <-ticker.C:
			elapsed := t.Sub(start)
			memUsage := GetMemoryUsage()
			fmt.Fprintf(os.Stderr,
				"[-] Please wait loading. Elapsed: %02d:%02d:%02d.%03d. Memory Usage: %.2f MB.\n",
				int(elapsed.Hours()), int(elapsed.Minutes())%60, int(elapsed.Seconds())%60,
				elapsed.Milliseconds()%1000, memUsage,
			)

			ticker.Stop()
			// Increment the interval by 10 seconds, capped at 10 minutes
			if interval < 10*time.Minute {
				interval += 10 * time.Second
			}
			ticker = time.NewTicker(interval)
		}
	}
}

// GetMemoryUsage returns the current memory usage in megabytes
func GetMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}
