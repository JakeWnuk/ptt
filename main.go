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

var version = "0.1.2"
var wg sync.WaitGroup
var mutex = &sync.Mutex{}
var retain models.FileArgumentFlag
var remove models.FileArgumentFlag
var readFiles models.FileArgumentFlag
var readURLs models.FileArgumentFlag
var transformationFiles models.FileArgumentFlag
var intRange models.IntRange
var lenRange models.IntRange
var primaryMap map[string]int
var err error

func main() {
	// Parse command line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Password Transformation Tool (ptt) version (%s):\n\n", version)
		fmt.Fprintf(os.Stderr, "ptt [options] [...]\nAccepts standard input and/or additonal arguments.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nThe -f, -k, -r, -tf, and -u flags can be used multiple times and together.\n")
		fmt.Fprintln(os.Stderr, "\nTransformation Modes:")
		modes := map[string]string{
			"append":                        "Transforms input into append rules.",
			"append-remove":                 "Transforms input into append-remove rules.",
			"append-shift":                  "Transforms input into append-shift rules.",
			"prepend":                       "Transforms input into prepend rules.",
			"prepend-remove":                "Transforms input into prepend-remove rules.",
			"prepend-shift":                 "Transforms input into prepend-shift rules.",
			"insert -i [index]":             "Transforms input into insert rules starting at index.",
			"overwrite -i [index]":          "Transforms input into overwrite rules starting at index.",
			"toggle -i [index]":             "Transforms input into toggle rules starting at index.",
			"encode":                        "Transforms input by URL, HTML, and Unicode escape encoding.",
			"decode":                        "Transforms input by URL, HTML, and Unicode escape decoding.",
			"hex":                           "Transforms input by encoding strings into $HEX[...] format.",
			"dehex":                         "Transforms input by decoding $HEX[...] formatted strings.",
			"mask -rm [uldsb] -v":           "Transforms input by masking characters with provided mask.",
			"remove -rm [uldsb]":            "Transforms input by removing characters with provided mask characters.",
			"retain -rm [uldsb] -tf [file]": "Transforms input by creating masks that still retain strings from file.",
			"pop -rm [uldsb]":               "Transforms input by generating tokens from popping strings at character boundaries.",
			"match -tf [file]":              "Transforms input by keeping only strings with matching masks from a mask file.",
			"swap -tf [file]":               "Transforms input by swapping tokens with exact matches from a ':' separated file.",
			"mask-swap -tf [file]":          "Transforms input by swapping tokens from a partial mask file and a input file.",
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

	}

	// Define command line flags
	verbose := flag.Bool("v", false, "Show verbose output when possible.")
	verbose2 := flag.Bool("vv", false, "Show statistics output when possible.")
	verbose3 := flag.Bool("vvv", false, "Show verbose statistics output when possible.")
	minimum := flag.Int("m", 0, "Minimum numerical frequency to include in output.")
	verboseStatsMax := flag.Int("n", 25, "Maximum number of items to display in verbose statistics output.")
	transformation := flag.String("t", "", "Transformation to apply to input.")
	replacementMask := flag.String("rm", "uldsb", "Replacement mask for transformations if applicable.")
	jsonOutput := flag.String("o", "", "Output to JSON file in addition to stdout.")
	flag.Var(&retain, "k", "Only keep items in a file.")
	flag.Var(&remove, "r", "Only keep items not in a file.")
	flag.Var(&readFiles, "f", "Read additional files for input.")
	flag.Var(&transformationFiles, "tf", "Read additional files for transformations if applicable.")
	flag.Var(&intRange, "i", "Starting index for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&lenRange, "l", "Keeps output equal to or within a range of lengths. Accepts ranges separated by '-'.")
	flag.Var(&readURLs, "u", "Read additional URLs for input.")
	flag.Parse()

	// Parse any retain, remove, or transformation file arguments
	fs := &models.RealFileSystem{}
	retainMap := utils.ReadFilesToMap(fs, retain)
	removeMap := utils.ReadFilesToMap(fs, remove)
	readFilesMap := utils.ReadFilesToMap(fs, readFiles)
	transformationFilesMap := utils.ReadFilesToMap(fs, transformationFiles)
	readURLsMap, err := utils.ReadURLsToMap(readURLs)
	if err != nil {
		fmt.Println("Error reading URLs:", err)
		return
	}

	// Read from stdin if provided
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		primaryMap, err = utils.LoadStdinToMap(bufio.NewScanner(os.Stdin))
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}
	}

	// Combine stdin with any additional files
	if len(primaryMap) == 0 && len(readFilesMap) == 0 && len(readURLsMap) == 0 {
		flag.Usage()
		return
	} else if len(primaryMap) == 0 {
		primaryMap = utils.CombineMaps(readFilesMap, readURLsMap)
	} else {
		primaryMap = utils.CombineMaps(primaryMap, readFilesMap, readURLsMap)
	}

	// Apply transformation if provided
	if *transformation != "" {
		primaryMap = transform.TransformationController(primaryMap, *transformation, intRange.Start, intRange.End, *verbose, *replacementMask, transformationFilesMap)
	}

	// Process retain and remove maps if provided
	if len(retainMap) > 0 || len(removeMap) > 0 {
		primaryMap, err = format.RetainRemove(primaryMap, retainMap, removeMap)
		if err != nil {
			fmt.Println("Error processing retain and remove flags:", err)
			return
		}
	}

	// Remove items under minimum frequency if provided
	if *minimum > 0 {
		primaryMap = format.RemoveMinimumFrequency(primaryMap, *minimum)
	}

	// Remove items outside of length range if provided
	if lenRange.Start > 0 || lenRange.End > 0 {
		primaryMap = format.RemoveLengthRange(primaryMap, lenRange.Start, lenRange.End)
	}

	// Print output to stdout
	if *verbose3 {
		format.PrintStatsToSTDOUT(primaryMap, *verbose3, *verboseStatsMax)
	} else if *verbose2 {
		format.PrintStatsToSTDOUT(primaryMap, *verbose3, *verboseStatsMax)
	} else {
		format.PrintArrayToSTDOUT(primaryMap, *verbose)
	}

	// Save output to JSON if provided
	if *jsonOutput != "" {
		err = format.SaveArrayToJSON(*jsonOutput, primaryMap)
		if err != nil {
			fmt.Println("Error saving output to JSON:", err)
			return
		}
	}
}
