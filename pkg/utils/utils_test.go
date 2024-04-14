package utils

import (
	"ptt/pkg/models"
	"testing"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Loading and Processing Functions **
// - ReadFilesToMap()
// - LoadStdinToMap()
// - CombineMaps()
//
// ** Transformation Functions **
// - ReverseString()
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// -

// Unit Test for ReadFilesToMap()
func TestReadFilesToMap(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 string
		Input2 string
		Output map[string]int
	}

	type TestCases []TestCase

	// Create a mock file system with example files
	mockFs := &models.MockFileSystem{
		Files: map[string][]byte{
			"file1": []byte("love1\nlove2\nlove3"),
			"file2": []byte("<31\n<32\n<33"),
			"file3": []byte("爱1\n爱2\n爱3"),
			"file4": []byte("amor1\namor2\namor3"),
			"file5": []byte("amour1\namour2\namour3"),
			"file6": []byte("愛1\n愛2\n愛3"),
		},
	}

	// Define test cases
	testCases := TestCases{
		{"file1", "file2", map[string]int{"love1": 1, "love2": 1, "love3": 1, "<31": 1, "<32": 1, "<33": 1}},
		{"file3", "file4", map[string]int{"爱1": 1, "爱2": 1, "爱3": 1, "amor1": 1, "amor2": 1, "amor3": 1}},
		{"file5", "file6", map[string]int{"amour1": 1, "amour2": 1, "amour3": 1, "愛1": 1, "愛2": 1, "愛3": 1}},
		{"file1", "file1", map[string]int{"love1": 2, "love2": 2, "love3": 2}},
		{"file2", "file2", map[string]int{"<31": 2, "<32": 2, "<33": 2}},
		{"file3", "file3", map[string]int{"爱1": 2, "爱2": 2, "爱3": 2}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := ReadFilesToMap(mockFs, []string{input1, input2})
		if CheckAreMapsEqual(given, output) == false {
			t.Errorf("ReadFilesToMap(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for LoadStdinToMap()
func TestLoadStdinToMap(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  []string
		Output map[string]int
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{[]string{"love1", "love2", "love3", "love2", "love3", "love1"}, map[string]int{"love1": 2, "love2": 2, "love3": 2}},
		{[]string{"<31", "<32", "<33", "<32", "<33", "<31"}, map[string]int{"<31": 2, "<32": 2, "<33": 2}},
		{[]string{"爱1", "爱2", "爱3", "爱2", "爱3", "爱1"}, map[string]int{"爱1": 2, "爱2": 2, "爱3": 2}},
		{[]string{"amor1", "amor2", "amor3", "amor2", "amor3", "amor1"}, map[string]int{"amor1": 2, "amor2": 2, "amor3": 2}},
		{[]string{"amour1", "amour2", "amour3", "amour2", "amour3", "amour1"}, map[string]int{"amour1": 2, "amour2": 2, "amour3": 2}},
		{[]string{"愛1", "愛2", "愛3", "愛2", "愛3", "愛1"}, map[string]int{"愛1": 2, "愛2": 2, "愛3": 2}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		// Create a mock scanner with example lines
		mockScanner := &models.MockScanner{
			Lines: input,
		}

		given, err := LoadStdinToMap(mockScanner)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if CheckAreMapsEqual(given, output) == false {
			t.Errorf("LoadStdinToMap(%v) = %v; want %v", input, given, output)
		}
	}
}

// Unit Test for CombineMaps()
func TestCombineMaps(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 map[string]int
		Input2 map[string]int
		Output map[string]int
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{map[string]int{"love1": 1, "love2": 1, "love3": 1}, map[string]int{"love1": 1, "love2": 1, "love3": 1}, map[string]int{"love1": 2, "love2": 2, "love3": 2}},
		{map[string]int{"<31": 1, "<32": 1, "<33": 1}, map[string]int{"<31": 1, "<32": 1, "<33": 1}, map[string]int{"<31": 2, "<32": 2, "<33": 2}},
		{map[string]int{"爱1": 1, "爱2": 1, "爱3": 1}, map[string]int{"爱1": 1, "爱2": 1, "爱3": 1}, map[string]int{"爱1": 2, "爱2": 2, "爱3": 2}},
		{map[string]int{"amor1": 1, "amor2": 1, "amor3": 1}, map[string]int{"amor1": 1, "amor2": 1, "amor3": 1}, map[string]int{"amor1": 2, "amor2": 2, "amor3": 2}},
		{map[string]int{"amour1": 1, "amour2": 1, "amour3": 1}, map[string]int{"amour1": 1, "amour2": 1, "amour3": 1}, map[string]int{"amour1": 2, "amour2": 2, "amour3": 2}},
		{map[string]int{"愛1": 1, "愛2": 1, "愛3": 1}, map[string]int{"愛1": 1, "愛2": 1, "愛3": 1}, map[string]int{"愛1": 2, "愛2": 2, "愛3": 2}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := CombineMaps(input1, input2)
		if CheckAreMapsEqual(given, output) == false {
			t.Errorf("CombineMaps(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for ReverseString()
func TestReverseString(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  string
		Output string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"love", "evol"},
		{"<31", "13<"},
		{"爱", "爱"},
		{"amor", "roma"},
		{"amour", "ruoma"},
		{"愛", "愛"},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		given := ReverseString(input)
		if given != output {
			t.Errorf("ReverseString(%v) = %v; want %v", input, given, output)
		}
	}
}
