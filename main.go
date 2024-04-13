// Package main controls the user interaction logic for the application
package main

import (
	"flag"
	"fmt"
	"os"
	"ppt/pkg/format"
	"ppt/pkg/models"
	"ppt/pkg/transform"
	"ppt/pkg/utils"
	"sync"
)

var version = "0.0.0"
var wg sync.WaitGroup
var mutex = &sync.Mutex{}
var retain models.FileArgumentFlag
var remove models.FileArgumentFlag
var readFiles models.FileArgumentFlag
var transformationFiles models.FileArgumentFlag
var primaryMap map[string]int
var err error

func main() {
	// Parse command line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Password Transformation Tool (ptt) version (%s):\n\n", version)
		fmt.Fprintf(os.Stderr, "ptt [options] [...]\nAccepts standard input and/or additonal arguments.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// Define command line flags
	verbose := flag.Bool("v", false, "Show verbose output when possible.")
	minimum := flag.Int("m", 0, "Minimum numerical frequency to include in output.")
	transformation := flag.String("t", "", "Transformation to apply to input.")
	startingIndex := flag.Int("i", 0, "Starting index for transformations if applicable.")
	replacementMask := flag.String("rm", "uldsb", "Replacement mask for transformations if applicable.")
	swapTolerance := flag.Int("st", 0, "Swap tolerance for transformations if applicable.")
	flag.Var(&retain, "k", "Only keep items in a file.")
	flag.Var(&remove, "r", "Only keep items not in a file.")
	flag.Var(&readFiles, "f", "Read additonal files for input.")
	flag.Var(&transformationFiles, "tf", "Read additonal files for transformations if applicable.")
	flag.Parse()

	// Parse any retain, remove, or transformation file arguments
	retainMap := utils.ReadFilesToMap(retain)
	removeMap := utils.ReadFilesToMap(remove)
	readFilesMap := utils.ReadFilesToMap(readFiles)
	transformationFilesMap := utils.ReadFilesToMap(transformationFiles)

	// Read from stdin if provided
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		primaryMap, err = utils.LoadStdinToMap()
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}
	}

	// Combine stdin with any additional files
	if len(primaryMap) == 0 && len(readFilesMap) == 0 {
		fmt.Println("No input provided. Exiting.")
		return
	} else if len(primaryMap) == 0 {
		primaryMap = readFilesMap
	} else if len(readFilesMap) > 0 {
		primaryMap = utils.CombineMaps(primaryMap, readFilesMap)
	}

	// Apply transformation if provided
	if *transformation != "" {
		primaryMap = transform.TransformationController(primaryMap, *transformation, *startingIndex, *verbose, *replacementMask, transformationFilesMap, *swapTolerance)
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

	format.PrintArrayToSTDOUT(primaryMap, *verbose)
}
