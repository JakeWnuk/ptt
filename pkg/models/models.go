package models

import "strings"

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
