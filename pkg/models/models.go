// Package models contains models used by the application
package models

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ----------------------------------------------------------------------------
// Argument Flag Models
// ----------------------------------------------------------------------------
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

// TemplateFileOperation is used to store the transformation operations loaded
// from JSON template files. The intention is to provide a way to define the
// operations in a structured way.
type TemplateFileOperation struct {
	StartIndex         int
	EndIndex           int
	Verbose            bool
	ReplacementMask    string
	Bypass             bool
	TransformationMode string
	WordRangeStart     int
	WordRangeEnd       int
}

// ----------------------------------------------------------------------------
// Output Sorting Models
// ----------------------------------------------------------------------------
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

// ----------------------------------------------------------------------------
// Unit Test Models
// ----------------------------------------------------------------------------
// These models are used to define the test cases that are used to test the
// functions in the application. The intention is to provide a way to define
// the test cases in a structured way.

// FileSystem is an interface that is used to read files from the file system
// This is used to allow the application to read files from the real file system
// or from a mock file system for testing
type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
	Open(filename string) (File, error)
}

// File is an interface that represents a file
type File interface {
	Read(p []byte) (n int, err error)
	Close() error
}

// MockFileSystem is used to read files from the mock file system
type MockFileSystem struct {
	Files map[string][]byte
}

// MockFile represents a mock file
type MockFile struct {
	Data   []byte
	Offset int64
}

// Read reads data from the mock file
func (m *MockFile) Read(p []byte) (n int, err error) {
	if m.Offset >= int64(len(m.Data)) {
		return 0, io.EOF
	}
	n = copy(p, m.Data[m.Offset:])
	m.Offset += int64(n)
	return n, nil
}

// Close closes the mock file (no-op for mock)
func (m *MockFile) Close() error {
	return nil
}

// Open opens a mock file and returns a File interface
func (m *MockFileSystem) Open(filename string) (File, error) {
	if data, ok := m.Files[filename]; ok {
		return &MockFile{Data: data}, nil
	}
	return nil, fmt.Errorf("file not found: %s", filename)
}

// ReadFile Implements the ReadFile method of the FileSystem interface for the MockFileSystem
func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	if data, ok := m.Files[filename]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("file not found: %s", filename)
}

// RealFileSystem is used to read files from the real file system
type RealFileSystem struct{}

// ReadFile is used to read a file from the real file system
func (r *RealFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// Open opens a file and returns a File interface
func (r RealFileSystem) Open(filename string) (File, error) {
	return os.Open(filename)
}

// Scanner is an interface that is used to read lines from a file
type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}

// MockScanner is used to read lines from a mock file for testing
type MockScanner struct {
	Lines []string
	Index int
}

// Scan implements the Scan, Text, and Err methods of the Scanner interface for the MockScanner
func (m *MockScanner) Scan() bool {
	if m.Index < len(m.Lines) {
		m.Index++
		return true
	}
	return false
}

// Text implements the Text method of the Scanner interface for the MockScanner
func (m *MockScanner) Text() string {
	return m.Lines[m.Index-1]
}

// Err implements the Err method of the Scanner interface for the MockScanner
func (m *MockScanner) Err() error {
	return nil
}
