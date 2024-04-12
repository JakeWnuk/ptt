// Package models contains models used by the application
package models

import "fmt"

// ----------------------------------------
// Argument Flag Models
// ----------------------------------------
// These models are used to define the flags that are used in the command line
// interface. The intention is to provide a way to define the flags can be used
// for multiple file arguments.

// FileArgumentFlag is a custom type that is used to store the file arguments
type FileArgumentFlag []string

// String is used to implement the flag.Value interface
func (w *FileArgumentFlag) String() string {
	return fmt.Sprint(*w)
}

// Set is used to implement the flag.Value interface
func (w *FileArgumentFlag) Set(value string) error {
	*w = append(*w, value)
	return nil
}

// ----------------------------------------
// Output Sorting Models
// ----------------------------------------
// These models are used to define the sorting options that are used to sort the
// final output to the console. The intention is to provide a way to efficiently
// sort the map[string]int that is used to store the word count in o(n log n) time.

// Pair is used to store the key value pair of the map[string]int
type Pair struct {
	Key   string
	Value int
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Value
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
