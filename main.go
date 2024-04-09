// Package main controls the user interaction logic for the application
package main

import (
	"flag"
	"fmt"
	"os"
	"ppt/pkg/utils"
)

type argumentFilesFlag []string

func (w *argumentFilesFlag) String() string {
	return fmt.Sprint(*w)
}

func (w *argumentFilesFlag) Set(value string) error {
	*w = append(*w, value)
	return nil
}

var version = "0.0.0"
var retain argumentFilesFlag
var remove argumentFilesFlag

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of ptt version (%s):\n\n", version)
		fmt.Fprintf(os.Stderr, "ptt [options] [URLS/FILES] [...]\nAccepts standard input and/or additonal arguments.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	count := flag.Bool("c", false, "Show the frequency count of each item")
	minimum := flag.Int("m", 0, "Minimum frequency to include in output. Value should be an integer.")
	flag.Var(&retain, "w", "Only include items in a file.")
	flag.Var(&remove, "v", "Only include items not in a file.")
	flag.Parse()

	// Parse any retain/remove files
	retainMap := utils.ReadFilesToMap(retain)
	removeMap := utils.ReadFilesToMap(remove)

	fmt.Println(count, minimum, retainMap, removeMap)
}
