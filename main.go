// Package main controls the user interaction logic for the application
package main

import (
	"flag"
	"ptt/pkg/mask"
	"ptt/pkg/models"
	"ptt/pkg/transform"
	"strings"
)

var version = "1.1.0"
var transformationModeArray models.MultiString
var transformationFiles models.MultiString
var wordRange models.IntRange
var intRange models.IntRange
var lenRange models.IntRange

func main() {
	flag.Var(&transformationModeArray, "t", "Transformation mode to be used. Can be specified multiple times.")
	flag.Var(&intRange, "i", "Starting index for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&wordRange, "w", "Number of words for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&lenRange, "l", "Only output items of a certain length. Accepts ranges separated by '-'.")
	replacementMask := flag.String("m", "uldsbt", "Mask for transformations if applicable.")
	debugMode := flag.Bool("d", false, "Enable debug mode.")
	verbose := flag.Bool("v", false, "Show verbose report output. Warning: loads information into memory.")
	flag.Parse()

	if *debugMode {
		models.DebugMode = true
	}

	if intRange.Start >= 0 && intRange.End >= 0 {
		models.OperationStart = intRange.Start
		models.OperationEnd = intRange.End
	}

	if wordRange.Start >= 0 && wordRange.End >= 0 {
		models.WordStart = wordRange.Start
		models.WordEnd = wordRange.End
	}

	if lenRange.Start >= 0 && lenRange.End >= 0 {
		models.LengthStart = lenRange.Start
		models.LengthEnd = lenRange.End
	}

	if lenRange.Start < 1 {
		models.LengthStart = 1
	}

	if lenRange.End < 1 {
		models.LengthEnd = 99
	}

	if replacementMask != nil {
		models.GlobalMask = *replacementMask
		replacements := mask.ConstructReplacements(*replacementMask)
		models.MaskReplacer = strings.NewReplacer(replacements...)
	}

	if *verbose {
		models.VerboseOutput = make(map[string]int)
		models.Verbose = true
	}

	transform.ReadReturnStandardInput(transformationModeArray)
}
