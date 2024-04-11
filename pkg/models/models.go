package models

import "fmt"

// ----------------------------------------
// Argument Flag Models
// ----------------------------------------
// These models are used to define the flags that are used in the command line
// interface. The intention is to provide a way to define the flags can be used
// for multiple file arguments.

type FileArgumentFlag []string

// This method is used to implement the flag.Value interface
func (w *FileArgumentFlag) String() string {
	return fmt.Sprint(*w)
}

// This method is used to implement the flag.Value interface
func (w *FileArgumentFlag) Set(value string) error {
	*w = append(*w, value)
	return nil
}
