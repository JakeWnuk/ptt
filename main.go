// Package main controls the user interaction logic for the application
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jakewnuk/ptt/pkg/mask"
	"github.com/jakewnuk/ptt/pkg/models"
	"github.com/jakewnuk/ptt/pkg/transform"
	"github.com/jakewnuk/ptt/pkg/utils"
)

var version = "1.2.1"
var transformationModeArray models.MultiString
var transformationFiles models.MultiString
var wordRange models.IntRange
var intRange models.IntRange
var lenRange models.IntRange

func main() {
	// Parse command line arguments
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Password Transformation Tool (ptt) version (%s):\n\n", version)
		fmt.Fprintf(os.Stderr, "ptt [options] [...]\nAccepts standard input and/or additonal arguments.\n\n")
		fmt.Fprintf(os.Stderr, "The -t flag can be used multiple times.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "These modify or filter the transformation mode.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Transformation Modes:\n")
		fmt.Fprintf(os.Stderr, "These create or alter based on the selected mode.\n\n")
		modes := map[string]string{
			"rule-append":               "Transforms input by creating append rules.",
			"rule-append-remove":        "Transforms input by creating append-remove rules.",
			"rule-prepend":              "Transforms input by creating prepend rules.",
			"rule-prepend-remove":       "Transforms input by creating prepend-remove rules.",
			"rule-prepend-toggle":       "Transforms input by creating prepend-toggle rules.",
			"rule-insert -i [index]":    "Transforms input by creating insert rules starting at index.",
			"rule-overwrite -i [index]": "Transforms input by creating overwrite rules starting at index.",
			"rule-toggle -i [index]":    "Transforms input by creating toggle rules starting at index.",
			"mask -m [uldsb]":           "Transforms input by masking characters with provided mask.",
			"mask-remove -m [uldsb]":    "Transforms input by removing characters with provided mask.",
			"mask-pop -m [uldsbt]":      "Transforms input by popping tokens from character boundaries using the provided mask.",
			"token-swap":                "Transforms input by performing token swapping using top tokens from input.",
			"passphrase -w [words]":     "Transforms input by generating passphrases from sentences with a given number of words.",
			"regram -w [words]":         "Transforms input by regramming sentences into new n-grams with a given number of words.",
			"rule-simplify":             "Transforms input by simplifying rules to efficient equivalents using the HCRE library.",
			"dehex":                     "Transforms input by decoding $HEX[...] formatted strings. Will not affect other transformations.",
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
		fmt.Fprintf(os.Stderr, "\n")

	}

	flag.Var(&transformationModeArray, "t", "Transformation mode to be used. Can be specified multiple times.")
	flag.Var(&intRange, "i", "Starting index for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&wordRange, "w", "Number of words for transformations if applicable. Accepts ranges separated by '-'.")
	flag.Var(&lenRange, "l", "Only output items of a certain length. Accepts ranges separated by '-'.")
	replacementMask := flag.String("m", "uldsbt", "Mask for transformations if applicable.")
	debugMode := flag.Bool("d", false, "Enable debug mode.")
	verbose := flag.Bool("v", false, "Show verbose report output. Warning: loads information into memory.")
	veryVerbose := flag.Bool("vv", false, "Show very verbose report output. Warning: loads information into memory.")
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

	if *veryVerbose {
		models.VerboseOutput = make(map[string]int)
		models.Verbose = true
		models.VeryVerbose = true
	}

	models.GlobalTokens = make(map[string]int)
	doneLoad := make(chan bool)
	go utils.TrackLoadTime(doneLoad, "Run")
	transform.ReadReturnStandardInput(transformationModeArray)
	doneLoad <- true
	close(doneLoad)
}
