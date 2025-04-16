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
var DebugMode bool = false

// OperationStart is a string flag to specify the start of an operation
var OperationStart int

// OperationEnd is a string flag to specify the end of an operation
var OperationEnd int

// MaskReplacer is a strings.Replacer to specify the mask replacement string
// for operations
var MaskReplacer *strings.Replacer

// Mask is a string flag to specify the mask replacement string
// for operations
var GlobalMask string

// WordStart is a string flag to specify the start of the number of words for
// operations
var WordStart int

// WordEnd is a string flag to specify the end of the number of words for
// operations
var WordEnd int

// ----------------------------------------------------------------------------
// Argument Flag Models
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
