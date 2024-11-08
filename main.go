// Package main controls the user interaction logic for the application
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/jakewnuk/ptt/pkg/format"
	"github.com/jakewnuk/ptt/pkg/models"
	"github.com/jakewnuk/ptt/pkg/transform"
	"github.com/jakewnuk/ptt/pkg/utils"
)

var version = "0.4.0"
var wg sync.WaitGroup
var mutex = &sync.Mutex{}
var retain models.FileArgumentFlag
var remove models.FileArgumentFlag
var readFiles models.FileArgumentFlag
var readURLs models.FileArgumentFlag
var transformationFiles models.FileArgumentFlag
var templateFiles models.FileArgumentFlag
var intRange models.IntRange
var lenRange models.IntRange
var wordRange models.IntRange
var primaryMap map[string]int
var err error

func main() {
	// Parse command line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Password Transformation Tool (ptt) version (%s):\n\n", version)
		fmt.Fprintf(os.Stderr, "ptt [options] [...]\nAccepts standard input and/or additonal arguments.\n\n")
		fmt.Fprintf(os.Stderr, "The -f, -k, -r, -tf, -tp, and -u flags can be used multiple times, together, and with files or directories.\n")
		fmt.Fprintf(os.Stderr, "-------------------------------------------------------------------------------------------------------------\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "These modify or filter the transformation mode.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "-------------------------------------------------------------------------------------------------------------\n")
		fmt.Fprintf(os.Stderr, "Transformation Modes:\n")
		fmt.Fprintf(os.Stderr, "These create or alter based on the selected mode.\n\n")
		modes := map[string]string{
			"rule-append":                           "Transforms input by creating append rules.",
			"rule-append-remove":                    "Transforms input by creating append-remove rules.",
			"rule-prepend":                          "Transforms input by creating prepend rules.",
			"rule-prepend-remove":                   "Transforms input by creating prepend-remove rules.",
			"rule-prepend-toggle":                   "Transforms input by creating prepend-toggle rules.",
			"rule-insert -i [index]":                "Transforms input by creating insert rules starting at index.",
			"rule-overwrite -i [index]":             "Transforms input by creating overwrite rules starting at index.",
			"rule-toggle -i [index]":                "Transforms input by creating toggle rules starting at index.",
			"encode":                                "Transforms input by HTML and Unicode escape encoding.",
			"decode":                                "Transforms input by HTML and Unicode escape decoding.",
			"hex":                                   "Transforms input by encoding strings into $HEX[...] format.",
			"dehex":                                 "Transforms input by decoding $HEX[...] formatted strings.",
			"mask -rm [uldsb] -v":                   "Transforms input by masking characters with provided mask.",
			"mask-remove -rm [uldsb]":               "Transforms input by removing characters with provided mask.",
			"mask-retain -rm [uldsb] -tf [file] -v": "Transforms input by creating masks that still retain strings from file.",
			"mask-pop -rm [uldsbt]":                 "Transforms input by 'popping' tokens from character boundaries using the provided mask.",
			"mask-match -tf [file]":                 "Transforms input by keeping only strings with matching masks from a mask file.",
			"swap-single -tf [file]":                "Transforms input by swapping tokens once per string per replacement with exact matches from a ':' separated file.",
			"mask-swap -tf [file]":                  "Transforms input by swapping tokens from a mask/partial mask input and a transformation file of tokens.",
			"passphrase -w [words]":                 "Transforms input by generating passphrases from sentences with a given number of words.",
			"substring -i [index]":                  "Transforms input by extracting substrings starting at index and ending at index.",
			"replace-all -tf [file]":                "Transforms input by replacing all strings with all matches from a ':' separated file.",
			"regram -w [words]":                     "Transforms input by 'regramming' sentences into new n-grams with a given number of words.",
		}

		// Sort and print transformation modes
		keys := make([]string, 0, len(modes))
		for k := range modes {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(os.Stderr, "  -t %s\n\t%s\n", k, modes[k])
		}
		fmt.Fprintf(os.Stderr, "-------------------------------------------------------------------------------------------------------------\n")

	}

	// Define command line flags
	verbose := flag.Bool("v", false, "Show verbose output when possible. (Can show additional metadata in some modes.)")
	verbose2 := flag.Bool("vv", false, "Show statistics output when possible.")
	verbose3 := flag.Bool("vvv", false, "Show verbose statistics output when possible.")
	minimum := flag.Int("m", 0, "Minimum numerical frequency to include in output.")
	markDownOutput := flag.Bool("md", false, "If Markdown format should be used for output instead.")
	outputVerboseMax := flag.Int("n", 0, "Maximum number of items to return in output.")
	transformation := flag.String("t", "", "Transformation to apply to input.")
	replacementMask := flag.String("rm", "uldsbt", "Replacement mask for transformations if applicable.")
	jsonOutput := flag.String("o", "", "Output to JSON file in addition to stdout. Accepts file names and paths.")
	bypassMap := flag.Bool("b", false, "Bypass map creation and use stdout as primary output. Disables some options.")
	debugMode := flag.Int("d", 0, "Enable debug mode with verbosity levels [0-2].")
	URLParsingMode := flag.Int("p", 0, "Change parsing mode for URL input. [0 = Strict, 1 = Permissive, 2 = Maximum].")
	ignoreCase := flag.Bool("ic", false, "Ignore case when processing output and converts to lowercase.")
	flag.Var(&retain, "k", "Only keep items in a file.")
	flag.Var(&remove, "r", "Only keep items not in a file.")
	flag.Var(&readFiles, "f", "Read additional files for input.")
	flag.Var(&transformationFiles, "tf", "Read additional files for transformations if applicable.")
	flag.Var(&templateFiles, "tp", "Read a template file for multiple transformations and operations. Cannot be used with -t flag.")
	flag.Var(&intRange, "i", "Starting index for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&lenRange, "l", "Only output items of a certain length (does not adjust for rules). Accepts ranges separated by '-'.")
	flag.Var(&wordRange, "w", "Number of words for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&readURLs, "u", "Read additional URLs for input.")
	flag.Parse()

	// Bypass map creation if requested
	if *bypassMap {
		fmt.Fprintf(os.Stderr, "[*] Bypassing map creation and using standard output as primary output. Options are disabled.\n")
	}

	// Print debug information if requested
	if *debugMode > 0 {
		fmt.Fprintf(os.Stderr, "[*] Debug mode enabled with verbosity level %d.\n", *debugMode)
	}

	// Parse any retain, remove, or transformation file arguments
	fs := &models.RealFileSystem{}
	var retainMap map[string]int
	var removeMap map[string]int
	var readFilesMap map[string]int
	var transformationFilesMap map[string]int

	// Read files if provided
	if retain != nil || remove != nil || readFiles != nil || transformationFiles != nil {
		fmt.Fprintf(os.Stderr, "[*] Reading files for input.\n")
	}

	if retain != nil {
		retainMap = utils.ReadFilesToMap(fs, retain)
	}
	if remove != nil {
		removeMap = utils.ReadFilesToMap(fs, remove)
	}
	if readFiles != nil {
		readFilesMap = utils.ReadFilesToMap(fs, readFiles)
	}
	if transformationFiles != nil {
		transformationFilesMap = utils.ReadFilesToMap(fs, transformationFiles)
	}

	transformationTemplateArray := utils.ReadJSONToArray(fs, templateFiles)
	readURLsMap, err := utils.ReadURLsToMap(readURLs, *URLParsingMode, *debugMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Error reading URLs: %s.\n", err)
		return
	}

	// Read from stdin if provided
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		primaryMap, err = utils.LoadStdinToMap(bufio.NewScanner(os.Stdin))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error reading from standard input: %s.\n", err)
			return
		}
	}

	// Combine stdin with any additional files
	if len(primaryMap) == 0 && len(readFilesMap) == 0 && len(readURLsMap) == 0 {
		fmt.Fprintf(os.Stderr, "[!] No input provided. Exiting.\n")
		return
	} else if len(primaryMap) == 0 {
		primaryMap = utils.CombineMaps(readFilesMap, readURLsMap)
	} else {
		primaryMap = utils.CombineMaps(primaryMap, readFilesMap, readURLsMap)
	}

	fmt.Fprintf(os.Stderr, "[*] All input loaded.\n")

	// Apply transformation if provided
	if *transformation != "" && templateFiles == nil {
		primaryMap = transform.TransformationController(primaryMap, *transformation, intRange.Start, intRange.End, *verbose, *replacementMask, transformationFilesMap, *bypassMap, *debugMode, wordRange.Start, wordRange.End)
	} else if templateFiles != nil && *transformation == "" {
		fmt.Fprintf(os.Stderr, "[*] Using template files for multiple transformations.\n")

		// Make a copy of the primary map to avoid modifying the original
		temporaryMap := make(map[string]int)
		for k, v := range primaryMap {
			temporaryMap[k] = v
		}

		// Apply transformations from template files
		for i, template := range transformationTemplateArray {
			if i == 0 {
				temporaryMap = transform.TransformationController(primaryMap, template.TransformationMode, template.StartIndex, template.EndIndex, template.Verbose, template.ReplacementMask, transformationFilesMap, template.Bypass, *debugMode, template.WordRangeStart, template.WordRangeEnd)
			} else {
				temporaryMap = utils.CombineMaps(temporaryMap, transform.TransformationController(primaryMap, template.TransformationMode, template.StartIndex, template.EndIndex, template.Verbose, template.ReplacementMask, transformationFilesMap, template.Bypass, *debugMode, template.WordRangeStart, template.WordRangeEnd))
			}
		}
		primaryMap = temporaryMap

	} else if *transformation != "" && templateFiles != nil {
		fmt.Fprintf(os.Stderr, "[!] Transformation and template flags cannot be used together.\n")
		return
	}

	// Print ignore case if provided
	if *ignoreCase {
		fmt.Fprintf(os.Stderr, "[*] Ignoring case when processing output.\n")
	}

	// Ignore case if provided
	if *ignoreCase {
		primaryMap = format.CreateIgnoreCaseMap(primaryMap)
	}

	// Print remove frequency if provided
	if *minimum > 0 {
		fmt.Fprintf(os.Stderr, "[*] Removing items with frequency less than %d.\n", *minimum)
	}

	// Remove items under minimum frequency if provided
	if *minimum > 0 {
		primaryMap = format.RemoveMinimumFrequency(primaryMap, *minimum)
	}

	// Print length range if provided
	if lenRange.Start > 0 || lenRange.End > 0 {
		fmt.Fprintf(os.Stderr, "[*] Only outputting items between %d and %d characters.\n", lenRange.Start, lenRange.End)
	}

	// Remove items outside of length range if provided
	if lenRange.Start > 0 || lenRange.End > 0 {
		primaryMap = format.RemoveLengthRange(primaryMap, lenRange.Start, lenRange.End)
	}

	// Print retained and removed items if provided
	if len(retainMap) > 0 || len(removeMap) > 0 {
		fmt.Fprintf(os.Stderr, "[*] Retain/remove flags provided. Retaining %d and removing %d items.\n", len(retainMap), len(removeMap))
	}

	// Process retain and remove maps if provided
	if len(retainMap) > 0 || len(removeMap) > 0 {
		primaryMap, err = format.RetainRemove(primaryMap, retainMap, removeMap, *debugMode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error processing retain and remove flags: %s.\n", err)
			return
		}
	}

	// if -n is provided, filter ALL results to only that top amount
	if *outputVerboseMax > 0 {
		primaryMap = format.FilterTopN(primaryMap, *outputVerboseMax)
	}

	fmt.Fprintf(os.Stderr, "[*] Task complete with %d unique results.\n", len(primaryMap))

	// Print in markdown if provided
	if *markDownOutput {
		command := "ptt "
		for _, arg := range os.Args[1:] {
			command += arg + " "
		}

		format.PrintArrayToMarkdown(primaryMap, command)
	}

	// Print output to stdout
	if *verbose3 && !*markDownOutput {
		format.PrintStatsToSTDOUT(primaryMap, *verbose3, *outputVerboseMax)
	} else if *verbose2 && !*markDownOutput {
		format.PrintStatsToSTDOUT(primaryMap, *verbose3, *outputVerboseMax)
	} else if !*markDownOutput {
		format.PrintArrayToSTDOUT(primaryMap, *verbose)
	}

	// Print output location if provided
	if *jsonOutput != "" {
		fmt.Fprintf(os.Stderr, "[*] Saving output to JSON file: %s.\n", *jsonOutput)
	}

	// Save output to JSON if provided
	if *jsonOutput != "" {
		err = format.SaveArrayToJSON(*jsonOutput, primaryMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error saving output to JSON: %s.\n", err)
			return
		}
	}
}
