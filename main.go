// Package main controls the user interaction logic for the application
package main

import (
	"flag"
	"fmt"
	"os"
	"ppt/pkg/models"
	"ppt/pkg/utils"
	"sync"
)

var version = "0.0.0"
var wg sync.WaitGroup
var mutex = &sync.Mutex{}
var retain models.FileArgumentFlag
var remove models.FileArgumentFlag
var readFiles models.FileArgumentFlag

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Password Transformation Tool (ptt) version (%s):\n\n", version)
		fmt.Fprintf(os.Stderr, "ptt [options] [URLS/FILES] [...]\nAccepts standard input and/or additonal arguments.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	verbose := flag.Bool("v", false, "Show verbose output when possible.")
	minimum := flag.Int("m", 0, "Minimum numerical frequency to include in output.")
	flag.Var(&retain, "k", "Only keep items in a file.")
	flag.Var(&remove, "r", "Only keep items not in a file.")
	flag.Var(&readFiles, "f", "Read additonal files for input.")
	flag.Parse()

	// Parse any retain or remove file arguments
	retainMap := utils.ReadFilesToMap(retain)
	removeMap := utils.ReadFilesToMap(remove)
	readFilesMap := utils.ReadFilesToMap(readFiles)

	fmt.Println(verbose, minimum, retainMap, removeMap, readFilesMap)
}
