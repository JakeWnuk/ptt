// Package models provides the data structures and models used in the
// application.
package models

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------------------
// Global Variables
// ----------------------------------------------------------------------------
// These variables are used to store the global state of the application.
// The intention is to provide a way to define the global state of the
// application in a structured way.

// DebugMode is a boolean flag to enable or disable debug mode
var DebugMode = false

// OperationStart is a string flag to specify the start of an operation
var OperationStart int

// OperationEnd is a string flag to specify the end of an operation
var OperationEnd int

// MaskReplacer is a strings.Replacer to specify the mask replacement string
// for operations
var MaskReplacer *strings.Replacer

// GlobalMask is a string flag to specify the mask replacement string
// for operations
var GlobalMask string

// WordStart is a string flag to specify the start of the number of words for
// operations
var WordStart int

// WordEnd is a string flag to specify the end of the number of words for
// operations
var WordEnd int

// LengthStart is a string flag to specify the start of the number of
// characters
var LengthStart int

// LengthEnd is a string flag to specify the end of the number of characters
var LengthEnd int

// Verbose is a boolean flag to enable or disable verbose mode
// for operations
var Verbose bool

// VerboseOutput is a map[string]int to hold the verbose output
// for operations
var VerboseOutput map[string]int

// GlobalTokens is a map[string]int to hold the global tokens for token
// swapping if specified
var GlobalTokens map[string]int

// ----------------------------------------------------------------------------
// Flag Models
// ----------------------------------------------------------------------------
// These models are used to define the flags that are used in the command line
// interface. The intention is to provide a way to define the flags can be used
// for multiple file arguments.

// MultiString is a custom type that implements the flag.Value interface
// to handle multiple string values from command line arguments
type MultiString []string

// String returns the string representation of the MultiString slice
func (m *MultiString) String() string {
	return strings.Join(*m, ",")
}

// Set appends a new string value to the MultiString slice
func (m *MultiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// IntRange is used to store the start and end of the range of integers. These
// are used to generate the insert and overwrite rules. The intention is to
// provide a way to define the range of integers in a structured way.
type IntRange struct {
	Start, End int
}

// NewIntRange is used to create a new IntRange object
func (i *IntRange) String() string {
	return fmt.Sprintf("%d-%d", i.Start, i.End)
}

// Set is used to implement the flag.Value interface
func (i *IntRange) Set(value string) error {
	parts := strings.Split(value, "-")
	if len(parts) == 1 {
		fmt.Sscanf(value, "%d", &i.Start)
		i.End = i.Start
	} else if len(parts) == 2 {
		fmt.Sscanf(value, "%d-%d", &i.Start, &i.End)
	}
	return nil
}

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
