// Package utils provides utility functions for the application.
package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/jakewnuk/ptt/pkg/models"

	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ----------------------------------------------------------------------------
// Loading and Processing Functions
// ----------------------------------------------------------------------------

// TrackLoadTime tracks the time it takes to load the input and prints the time
//
// Args:
// done (chan bool): channel to use to track tasks
// work (string): string used in status printing
//
// Returns:
// None
func TrackLoadTime(done <-chan bool, work string) {
	start := time.Now()
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-done:
			ticker.Stop()
			fmt.Fprintf(os.Stderr, "[-] Total %s Time: %02d:%02d:%02d.\n", work, int(time.Since(start).Hours()), int(time.Since(start).Minutes())%60, int(time.Since(start).Seconds())%60)
			return
		case t := <-ticker.C:
			elapsed := t.Sub(start)
			memUsage := GetMemoryUsage()
			fmt.Fprintf(os.Stderr, "[-] Please wait loading. Elapsed: %02d:%02d:%02d.%03d. Memory Usage: %.2f MB.\n", int(t.Sub(start).Hours()), int(t.Sub(start).Minutes())%60, int(t.Sub(start).Seconds())%60, elapsed.Milliseconds()%1000, memUsage)
		}
	}
}

// GetMemoryUsage returns the current memory usage in megabytes
func GetMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

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
	// 1 GB read buffer
	chunkSize := int64(1 * 1024 * 1024 * 1024)

	i := 0
	for i < len(filenames) {
		filename := filenames[i]
		if IsFileSystemDirectory(filename) {
			files, err := GetFilesInDirectory(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error reading the directory %v: %v.\n", filename, err)
				os.Exit(1)
			}
			filenames = append(filenames, files...)
		} else {
			file, err := fs.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error opening file %s.\n", filename)
				os.Exit(1)
			}
			defer file.Close()

			buffer := make([]byte, chunkSize)
			for {
				bytesRead, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					fmt.Fprintf(os.Stderr, "[!] Error reading file %s.\n", filename)
					os.Exit(1)
				}
				if bytesRead == 0 {
					break
				}

				data := buffer[:bytesRead]

				err = json.Unmarshal(data, &wordMap)
				if err == nil {
					fmt.Fprintf(os.Stderr, "[*] Detected ptt JSON output. Importing...\n")
					continue
				}

				fileWords := strings.Split(string(data), "\n")
				for _, word := range fileWords {
					wordMap[word]++
				}

				if err == io.EOF {
					break
				}
			}
		}
		i++
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
	reDetect := regexp.MustCompile(`^\d+\s(\w+|\W+)$`)
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
// from the URLs. Supports files or directories containing URLs.
//
// Args:
//
//	urls ([]string): The URLs to read
//	parsingMode (int): Change parsing mode for URL input. [0 = Strict, 1 = Permissive, 2 = Maximum] [0-2].
//	debugMode (int): A flag to print debug information
//
// Returns:
//
//	map[string]int: A map of words from the URLs
//	error: An error if one occurred
func ReadURLsToMap(urls []string, parsingMode int, debugMode int) (map[string]int, error) {
	wordMap := make(map[string]int)
	var wg sync.WaitGroup

	ch := make(chan string)

	go func() {
		for word := range ch {
			wordMap[word]++
		}
	}()

	sleepOnStart := true
	for _, iURL := range urls {
		if IsValidURL(iURL) {
			wg.Add(1)
			go ProcessURL(iURL, ch, &wg, parsingMode, debugMode, sleepOnStart)

		} else if IsFileSystemDirectory(iURL) {
			files, err := GetFilesInDirectory(iURL)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				wg.Add(1)
				go ProcessURLFile(file, ch, &wg, parsingMode, debugMode)
			}
		} else if IsValidFile(iURL) {
			wg.Add(1)
			go ProcessURLFile(iURL, ch, &wg, parsingMode, debugMode)
		} else {
			fmt.Fprintf(os.Stderr, "[!] Rejected URL or file: %s.\n", iURL)
			return nil, fmt.Errorf("invalid input: %s", iURL)
		}
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
	var result sync.Map

	for _, m := range maps {
		for k, v := range m {
			if val, ok := result.Load(k); ok {
				result.Store(k, val.(int)+v)
			} else {
				result.Store(k, v)
			}
		}
	}

	finalResult := make(map[string]int)
	result.Range(func(k, v interface{}) bool {
		finalResult[k.(string)] = v.(int)
		return true
	})

	return finalResult
}

// ProcessURL reads the contents of a URL and sends each sentence to the channel
//
// Args:
//
//	url (string): The URL to read
//	ch (chan<- string): The channel to send the sentences to
//	wg (*sync.WaitGroup): The WaitGroup to signal when done
//	parsingMode (int): Change parsing mode for URL input. [0 = Strict,
//	1 = Permissive, 2 = Maximum] [0-2].
//	debugMode (int): A flag to print debug information
//	sleepOnStart (bool): A flag to sleep before starting the request
//
// Returns:
//
//	None
func ProcessURL(url string, ch chan<- string, wg *sync.WaitGroup, parsingMode int, debugMode int, sleepOnStart bool) {
	defer wg.Done()
	var resp *http.Response
	throttleInterval := 90
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	const maxRetries = 3
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15A372 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15A5341f Safari/604.1",
		"Mozilla/5.0 (Linux; Android 11; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.181 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
	}

	if sleepOnStart {
		time.Sleep(time.Second * time.Duration(r.Intn(throttleInterval)))
	}

	for attempts := 0; attempts <= maxRetries; attempts++ {

		// Set a random user agent
		randomUserAgent := userAgents[r.Intn(len(userAgents))]

		// Fetch the URL
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", randomUserAgent)
		resp, err = client.Do(req)
		if err != nil {
			if debugMode >= 2 {
				fmt.Fprintf(os.Stderr, "[!] Error fetching URL %s.\n", url)
			}
			return
		}

		if resp == nil {
			if debugMode >= 2 {
				fmt.Fprintf(os.Stderr, "[!] Error no response from URL %s.\n", url)
			}
			return
		}

		defer resp.Body.Close()

		// Check the response code for throttling
		if resp.StatusCode == http.StatusTooManyRequests {
			throttleInterval += 30
			time.Sleep(time.Second * time.Duration(throttleInterval*(r.Intn(3)+1)))
			fmt.Fprintf(os.Stderr, "[!] Requested %s. Attempt [%d/%d]. Response Code: %s. Waiting %d seconds before retrying. \n", url, attempts, maxRetries, resp.Status, throttleInterval)
		} else {
			fmt.Fprintf(os.Stderr, "[+] Requested %s. Attempt [%d/%d]. Response Code: %s. Content-Type: %s. \n", url, attempts, maxRetries, resp.Status, resp.Header.Get("Content-Type"))
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
			// end the loop if the code is 300, 301, 302, 303, 307, 308, 400, 401, 403, 405, 500, 503
			if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect || resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusServiceUnavailable {

				if debugMode >= 2 {
					fmt.Fprintf(os.Stderr, "[!] Error fetching URL service returned %s. Removing target. %s\n", resp.Status, url)
				}
				return
			}

			if debugMode >= 2 {
				fmt.Fprintf(os.Stderr, "[!] Error unexpected response code %s. Retrying... %s\n", resp.Status, url)
			}
			continue
		}

		break
	}

	if resp == nil {
		if debugMode >= 1 {
			fmt.Fprintf(os.Stderr, "[!] Error no response from URL %s.\n", url)
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if debugMode >= 1 {
			fmt.Fprintf(os.Stderr, "[!] Error reading response body from URL %s.\n", url)
		}
		return
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
			if debugMode >= 1 {
				fmt.Fprintf(os.Stderr, "[!] Error parsing HTML from URL %s.\n", url)
			}
			return
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

	if debugMode == 1 {
		fmt.Fprintf(os.Stderr, "[?] URL: %s\n", url)
		fmt.Fprintf(os.Stderr, "[?] Content-Type: %s\n", contentType)
		fmt.Fprintf(os.Stderr, "[?] Parsing Mode: %d\n", parsingMode)
	} else if debugMode == 2 {
		fmt.Fprintf(os.Stderr, "[?] URL: %s\n", url)
		fmt.Fprintf(os.Stderr, "[?] Content-Type: %s\n", contentType)
		fmt.Fprintf(os.Stderr, "[?] Parsing Mode: %d\n", parsingMode)
		fmt.Fprintf(os.Stderr, "[?] Line Count: %d\n", len(lines))
		fmt.Fprintf(os.Stderr, "[?] Sample Lines:\n")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(os.Stderr, "%s\n", lines[i])
		}
	}

	// Iterate over the lines and split them
	for _, line := range lines {
		if parsingMode == 0 {
			textMatch, _ := regexp.MatchString(`[^a-zA-Z0-9.,;:!? ]`, line)
			if strings.Contains(contentType, "text/html") {
				if textMatch {
					continue
				}
			} else {
				if !textMatch {
					continue
				}
			}
		} else if parsingMode >= 1 {
			textMatch, _ := regexp.MatchString(`[^a-zA-Z0-9.,;:!?'"\- \/+_#@"\[\]]`, line)
			if strings.Contains(contentType, "text/html") {
				if textMatch {
					continue
				}
			} else {
				if !textMatch {
					continue
				}
			}
		}

		sentences := strings.Split(line, ".")
		for _, sentence := range sentences {

			if parsingMode >= 1 {
				phrases := strings.Split(sentence, ",")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				phrases = strings.Split(sentence, ";")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				phrases = strings.Split(sentence, ":")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				phrases = strings.Split(sentence, "!")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				phrases = strings.Split(sentence, "?")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}
			}

			if parsingMode >= 2 {
				phrases := strings.Split(sentence, " ")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				phrases = strings.Split(sentence, "-")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				phrases = strings.Split(sentence, "'")
				for _, phrase := range phrases {
					if phrase != "" {
						phrase = strings.TrimSpace(phrase)
						ch <- phrase
					}
				}

				twoGrams := GenerateNGrams(sentence, 2)
				threeGrams := GenerateNGrams(sentence, 3)
				fourGrams := GenerateNGrams(sentence, 4)
				fiveGrams := GenerateNGrams(sentence, 5)
				sixGrams := GenerateNGrams(sentence, 6)
				sevenGrams := GenerateNGrams(sentence, 7)
				allNGrams := append(twoGrams, append(threeGrams, append(fourGrams, append(fiveGrams, append(sixGrams, sevenGrams...)...)...)...)...)
				for _, nGram := range allNGrams {
					if nGram != "" {
						nGram = strings.TrimSpace(nGram)
						ch <- nGram

						noDot := strings.TrimRight(nGram, ".")
						ch <- noDot

						noComma := strings.TrimRight(nGram, ",")
						ch <- noComma

						frontParse := strings.TrimLeft(nGram, ", ")
						ch <- frontParse

						noSpace := strings.ReplaceAll(nGram, " ", "")
						ch <- noSpace

					}
				}
			}

			if sentence != "" {
				sentence = strings.TrimSpace(sentence)
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
	var combinedTemplate []models.TemplateFileOperation
	var template []models.TemplateFileOperation

	i := 0
	for i < len(filenames) {
		filename := filenames[i]
		// Check to see if a directory was passed
		// If so, read all files in the directory and append them to the filenames
		// slice
		if IsFileSystemDirectory(filename) {
			err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
					filenames = append(filenames, path)
				}
				return nil
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error walking the path %v: %v.\n", filename, err)
				os.Exit(1)
			}
		} else {
			data, err := fs.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error reading file %s.\n", filename)
				os.Exit(1)
			}

			err = json.Unmarshal(data, &template)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] Error unmarshalling JSON file %s.\n", filename)
				os.Exit(1)
			}

			combinedTemplate = append(combinedTemplate, template...)
		}
		i++
	}

	alphaRe := regexp.MustCompile(`[a-zA-Z]`)
	numRe := regexp.MustCompile(`[0-9]`)

	for _, template := range combinedTemplate {
		if !numRe.MatchString(fmt.Sprintf("%v", template.StartIndex)) || !numRe.MatchString(fmt.Sprintf("%v", template.EndIndex)) {
			fmt.Fprintf(os.Stderr, "[!] Error: StartIndex and EndIndex must be integers.\n")
			os.Exit(1)
		}

		if !alphaRe.MatchString(fmt.Sprintf("%v", template.Verbose)) {
			fmt.Fprintf(os.Stderr, "[!] Error: Verbose must be a boolean.\n")
			os.Exit(1)
		}

		if !alphaRe.MatchString(fmt.Sprintf("%v", template.ReplacementMask)) {
			fmt.Fprintf(os.Stderr, "[!] Error: ReplacementMask must be a string.\n")
			os.Exit(1)
		}

		if !alphaRe.MatchString(fmt.Sprintf("%v", template.Bypass)) {
			fmt.Fprintf(os.Stderr, "[!] Error: Bypass must be a boolean.\n")
			os.Exit(1)
		}

		if !alphaRe.MatchString(fmt.Sprintf("%v", template.TransformationMode)) {
			fmt.Fprintf(os.Stderr, "[!] Error: TransformationMode must be a string.\n")
			os.Exit(1)
		}

		if !numRe.MatchString(fmt.Sprintf("%v", template.WordRangeStart)) || !numRe.MatchString(fmt.Sprintf("%v", template.WordRangeEnd)) {
			fmt.Fprintf(os.Stderr, "[!] Error: WordRangeStart and WordRangeEnd must be integers.\n")
			os.Exit(1)
		}
	}

	return combinedTemplate
}

// ProcessURLFile reads the contents of a file containing URLs and sends each
// URL to the channel
//
// Args:
// filePath (string): The path to the file containing URLs
// ch (chan<- string): The channel to send the URLs to
// wg (*sync.WaitGroup): The WaitGroup to signal when done
// parsingMode (int): Change parsing mode for URL input. [0 = Strict,
// 1 = Permissive, 2 = Maximum] [0-2].
// debugMode (int): A flag to print debug information
//
// Returns:
// None
func ProcessURLFile(filePath string, ch chan<- string, wg *sync.WaitGroup, parsingMode int, debugMode int) {
	defer wg.Done()
	sleepOnStart := false

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error opening file %v: %v.\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if IsValidURL(line) {
			wg.Add(1)
			go ProcessURL(line, ch, wg, parsingMode, debugMode, sleepOnStart)
		} else {
			fmt.Fprintf(os.Stderr, "[!] Rejected URL: %s.\n", line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error reading file %v: %v.\n", filePath, err)
	}
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

// ReplaceAllSubstring replaces all instances of a substring in a string with
// a new substring if the substring is found in the original string. All of the
// replacements are applied to the original string. The new substring is
// determined by the key in the replacements map separated by a colon
// character.
//
// Args:
//
//	original (string): The original string
//	replacements (map[string]int): A map of substrings to replace
//
// Returns:
//
// []string: The original string with all instances of the substring replaced
func ReplaceAllSubstring(original string, replacements map[string]int) []string {
	newStrings := []string{original}
	for newSubstr := range replacements {
		// Split the new substring into the old and new strings by the colon character
		if !strings.Contains(newSubstr, ":") {
			continue
		}
		oldStr, newStr := strings.Split(newSubstr, ":")[0], strings.Split(newSubstr, ":")[1]
		var tempStrings []string
		for _, s := range newStrings {
			tempStrings = append(tempStrings, strings.Replace(s, oldStr, newStr, -1))
		}
		newStrings = tempStrings
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
				fmt.Fprintf(os.Stderr, fmt.Sprintf("[!] Error: Start index is out of bounds: %s.\n", s))
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
//
//	path (string): The path to check
//
// Returns:
//
//	bool: True if the path is a directory, false otherwise
func IsFileSystemDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return fileInfo.IsDir()
}

// IsValidURL checks if a string is a valid URL by parsing the string
//
// Args:
// str (string): The URL to check
//
// Returns:
// bool: True if the URL is valid, false otherwise
func IsValidURL(str string) bool {
	_, err := url.Parse(str)
	if err != nil {
		return false
	}

	if !strings.Contains(str, "http://") && !strings.Contains(str, "https://") {
		return false
	}

	return true
}

// GetFilesInDirectory returns a slice of files in a directory
// by reading the directory and appending the files to a slice
// if they are not directories
//
// Args:
// dir (string): The directory to read
//
// Returns:
// []string: A slice of files in the directory
func GetFilesInDirectory(dir string) ([]string, error) {
	var files []string
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if !item.IsDir() {
			files = append(files, filepath.Join(dir, item.Name()))
		}
	}

	return files, nil
}

// IsValidFile checks if a file exists and is not a directory
// by checking if the file exists
//
// Args:
// path (string): The path to the file
//
// Returns:
// bool: True if the file is valid, false otherwise
func IsValidFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
