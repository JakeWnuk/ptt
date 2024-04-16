// Package main controls the user interaction logic for the application
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"ptt/pkg/format"
	"ptt/pkg/models"
	"ptt/pkg/transform"
	"ptt/pkg/utils"
	"sync"
)

var version = "0.0.0"
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
		fmt.Fprintf(os.Stderr, "\nThe -f, -k, -r, -tf, and -u flags can be used multiple times.\n")
		fmt.Fprintf(os.Stderr, "\nTransformation Modes:\n")
		fmt.Fprintf(os.Stderr, "  -t append\n\tTransforms input into append rules.\n")
		fmt.Fprintf(os.Stderr, "  -t append-remove\n\tTransforms input into append-remove rules.\n")
		fmt.Fprintf(os.Stderr, "  -t append-shift\n\tTransforms input into append-shift rules.\n")
		fmt.Fprintf(os.Stderr, "  -t prepend\n\tTransforms input into prepend rules.\n")
		fmt.Fprintf(os.Stderr, "  -t prepend-remove\n\tTransforms input into prepend-remove rules.\n")
		fmt.Fprintf(os.Stderr, "  -t prepend-shift\n\tTransforms input into prepend-shift rules.\n")
		fmt.Fprintf(os.Stderr, "  -t insert -i [index]\n\tTransforms input into insert rules starting at index.\n")
		fmt.Fprintf(os.Stderr, "  -t overwrite -i [index]\n\tTransforms input into overwrite rules starting at index.\n")
		fmt.Fprintf(os.Stderr, "  -t toggle -i [index]\n\tTransforms input into toggle rules starting at index.\n")
		fmt.Fprintf(os.Stderr, "  -t encode\n\tTransforms input by URL, HTML, and Unicode escape encoding.\n")
		fmt.Fprintf(os.Stderr, "  -t decode\n\tTransforms input by URL, HTML, and Unicode escape decoding.\n")
		fmt.Fprintf(os.Stderr, "  -t dehex\n\tTransforms input by decoding $HEX[...] formatted strings.\n")
		fmt.Fprintf(os.Stderr, "  -t hex\n\tTransforms input by encoding strings into $HEX[...] format.\n")
		fmt.Fprintf(os.Stderr, "  -t mask -rm [uldsb] -v\n\tTransforms input by masking characters with provided mask.\n")
		fmt.Fprintf(os.Stderr, "  -t remove -rm [uldsb] -v\n\tTransforms input by removing characters with provided mask characters.\n")
		fmt.Fprintf(os.Stderr, "  -t retain -rm [uldsb] -tf [file]\n\tTransforms input by creating masks that still retain strings from file.\n")
		fmt.Fprintf(os.Stderr, "  -t match -tf [file]\n\tTransforms input by keeping only strings with matching masks from a mask file.\n")
		fmt.Fprintf(os.Stderr, "  -t fuzzy-swap -tf [file]\n\tTransforms input by swapping tokens with fuzzy matches from another file.\n")
		fmt.Fprintf(os.Stderr, "  -t swap -tf [file]\n\tTransforms input by swapping tokens with exact matches from a ':' separated file.\n")
		fmt.Fprintf(os.Stderr, "  -t pop -rm [uldsb]\n\tTransforms input by generating tokens excluding characters not part of the mask.\n")
	}

	// Define command line flags
	verbose := flag.Bool("v", false, "Show verbose output when possible.")
	verbose2 := flag.Bool("vv", false, "Show statistics output when possible.")
	verbose3 := flag.Bool("vvv", false, "Show verbose statistics output when possible.")
	minimum := flag.Int("m", 0, "Minimum numerical frequency to include in output.")
	verboseStatsMax := flag.Int("n", 25, "Maximum number of items to display in verbose statistics output.")
	transformation := flag.String("t", "", "Transformation to apply to input.")
	replacementMask := flag.String("rm", "uldsb", "Replacement mask for transformations if applicable.")
	flag.Var(&retain, "k", "Only keep items in a file.")
	flag.Var(&remove, "r", "Only keep items not in a file.")
	flag.Var(&readFiles, "f", "Read additional files for input.")
	flag.Var(&transformationFiles, "tf", "Read additional files for transformations if applicable.")
	flag.Var(&intRange, "i", "Starting index for transformations if applicable. Accepts ranges separated by '-'. (default 0)")
	flag.Var(&lenRange, "l", "Keeps output equal to or within a range of lengths. Accepts ranges separated by '-'. (default 0)")
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
}
