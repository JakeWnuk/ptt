package utils

import (
	"strings"

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
