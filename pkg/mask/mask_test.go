package mask

import (
	"reflect"
	"testing"

	"github.com/jakewnuk/ptt/pkg/utils"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Mask Generation Functions **
// - ConstructReplacements()
// - MakeMaskedMap()
// - MakeRetainMaskedMap()
// - MakeMaskedString()
//
// ** Mask Conversion Functions **
// - ConvertMultiByteMask()
// - TestMaskComplexity()
// - RemoveMaskedCharacters()
//
// ** Mask Utility Functions **
// - MakeMatchedMaskedMap()
// - BoundarySplitPopMap()
// - ShuffleMap()
// - CalculateKeySpace()
// - IsMaskAFullMask()
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// - CalculateMaskKeyspace()

// Unit Test for ConstructReplacements()
func TestConstructReplacements(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		output []string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"l", []string{"a", "?l", "b", "?l", "c", "?l", "d", "?l", "e", "?l", "f", "?l", "g", "?l", "h", "?l", "i", "?l", "j", "?l", "k", "?l", "l", "?l", "m", "?l", "n", "?l", "o", "?l", "p", "?l", "q", "?l", "r", "?l", "s", "?l", "t", "?l", "u", "?l", "v", "?l", "w", "?l", "x", "?l", "y", "?l", "z", "?l"}},
		{"u", []string{"A", "?u", "B", "?u", "C", "?u", "D", "?u", "E", "?u", "F", "?u", "G", "?u", "H", "?u", "I", "?u", "J", "?u", "K", "?u", "L", "?u", "M", "?u", "N", "?u", "O", "?u", "P", "?u", "Q", "?u", "R", "?u", "S", "?u", "T", "?u", "U", "?u", "V", "?u", "W", "?u", "X", "?u", "Y", "?u", "Z", "?u"}},
		{"d", []string{"0", "?d", "1", "?d", "2", "?d", "3", "?d", "4", "?d", "5", "?d", "6", "?d", "7", "?d", "8", "?d", "9", "?d"}},
		// !\"#$%&\\()*+,-./:;<=>?@[\\]^_`{|}~'
		{"s", []string{" ", "?s", "!", "?s", "\"", "?s", "#", "?s", "$", "?s", "%", "?s", "&", "?s", "\\", "?s", "(", "?s", ")", "?s", "*", "?s", "+", "?s", ",", "?s", "-", "?s", ".", "?s", "/", "?s", ":", "?s", ";", "?s", "<", "?s", "=", "?s", ">", "?s", "?", "?s", "@", "?s", "[", "?s", "\\", "?s", "]", "?s", "^", "?s", "_", "?s", "`", "?s", "{", "?s", "|", "?s", "}", "?s", "~", "?s", "'", "?s"}},
	}

	// Run test cases
	for _, test := range tests {
		output := ConstructReplacements(test.input)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for MakeMaskedMap()
func TestMakeMaskedMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input        map[string]int
		replacements string
		output       map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc": 1, "ABC": 2, "ABCabc123!!!": 3}, "l", map[string]int{"?l?l?l": 1, "ABC": 2, "ABC?l?l?l123!!!": 3}},
		{map[string]int{"abc": 1, "ABC": 2, "ABCabc123!!!": 3}, "u", map[string]int{"abc": 1, "?u?u?u": 2, "?u?u?uabc123!!!": 3}},
		{map[string]int{"abc": 1, "ABC": 2, "ABCabc123!!!": 3}, "d", map[string]int{"abc": 1, "ABC": 2, "ABCabc?d?d?d!!!": 3}},
		{map[string]int{"abc": 1, "ABC": 2, "ABCabc123!!!": 3}, "s", map[string]int{"abc": 1, "ABC": 2, "ABCabc123?s?s?s": 3}},
		{map[string]int{"abc": 1, "ABC": 2, "ABCabc123!!!": 3}, "luds", map[string]int{"?l?l?l": 1, "?u?u?u": 2, "?u?u?u?l?l?l?d?d?d?s?s?s": 3}},
	}

	// Run test cases
	for _, test := range tests {
		output := MakeMaskedMap(test.input, test.replacements, false, false, false)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for MakeRetainMaskedMap()
func TestMakeRetainMaskedMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input        map[string]int
		replacements string
		retain       map[string]int
		output       map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc123": 1, "ABC": 2, "ABCabc123!!!": 3}, "luds", map[string]int{"abc": 1, "ABC": 1}, map[string]int{"abc?d?d?d": 1, "ABC": 2, "ABC?l?l?l?d?d?d?s?s?s": 3, "?u?u?uabc?d?d?d?s?s?s": 3}},
		{map[string]int{"abc123": 1, "ABC": 2, "ABCabc123!!!": 3}, "lud", map[string]int{"abc": 1, "ABC": 1}, map[string]int{"abc?d?d?d": 1, "ABC": 2, "ABC?l?l?l?d?d?d!!!": 3, "?u?u?uabc?d?d?d!!!": 3}},
		{map[string]int{"abc123": 1, "ABC": 2, "ABCabc123!!!": 3}, "ld", map[string]int{"123": 1}, map[string]int{"?l?l?l123": 1, "ABC?l?l?l123!!!": 3}},
		{map[string]int{"abc123": 1, "ABC": 2, "ABCabc123!!!": 3}, "luds", map[string]int{"abc": 1, "ABC": 1}, map[string]int{"abc?d?d?d": 1, "ABC": 2, "ABC?l?l?l?d?d?d?s?s?s": 3, "?u?u?uabc?d?d?d?s?s?s": 3}},
	}

	// Run test cases
	for _, test := range tests {
		output := MakeRetainMaskedMap(test.input, test.replacements, test.retain, false, false, false)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for MakeMaskedString()
func TestMakeMaskedString(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		mask   string
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"abc123", "l", "?l?l?l123"},
		{"abc123", "u", "abc123"},
		{"abc123", "d", "abc?d?d?d"},
		{"abc123", "s", "abc123"},
		{"abc123", "luds", "?l?l?l?d?d?d"},
	}

	// Run test cases
	for _, test := range tests {
		output := MakeMaskedString(test.input, test.mask)
		if output != test.output {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for ConvertMultiByteMask()
func TestConvertMultiByteMask(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"爱test", "?b?b?btest"},
		{"test爱test", "test?b?b?btest"},
		{"test爱", "test?b?b?b"},
		{"爱", "?b?b?b"},
		{"test", "test"},
	}

	// Run test cases
	for _, test := range tests {
		output := ConvertMultiByteMask(test.input)
		if output != test.output {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for TestMaskComplexity()
func TestTestMaskComplexity(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		output int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"?l?l?l123", 2},
		{"?l?l?l?d?d?d", 2},
		{"?l?l?l?d?d?d?s?s?s", 3},
		{"?u?u?l?d?d?d?s?s?s", 4},
	}

	// Run test cases
	for _, test := range tests {
		output := TestMaskComplexity(test.input)
		if output != test.output {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for RemoveMaskedCharacters()
func TestRemoveMaskedCharacters(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		output map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"?l?l?l123": 1, "?l?l?l?d?d?d": 2, "?l?l?l?d?d?d!!!": 3, "ABC?l?l?l?d?d?d!!!": 4}, map[string]int{"123": 1, "": 2, "!!!": 3, "ABC!!!": 4}},
	}

	// Run test cases
	for _, test := range tests {
		output := RemoveMaskedCharacters(test.input, "ulsbd", false, false)
		if !reflect.DeepEqual(output, test.output) {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for MakeMatchedMaskedMap()
func TestMakeMatchedMaskedMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input        map[string]int
		replacements string
		masks        map[string]int
		output       map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc123": 1, "ABC": 2, "ABCabc123!!!": 3}, "luds", map[string]int{"?l?l?l": 1, "?u?u?u": 1}, map[string]int{"ABC": 2}},
		{map[string]int{"123": 1, "123456": 2, "123456789": 3}, "d", map[string]int{"?d?d?d": 1, "?d?d?d?d?d?d": 1}, map[string]int{"123": 1, "123456": 2}},
		{map[string]int{"🙂": 1, "😀": 2, "😁": 3}, "b", map[string]int{"?b?b?b?b": 1}, map[string]int{"🙂": 1, "😀": 2, "😁": 3}},
	}

	// Run test cases
	for _, test := range tests {
		output := MakeMatchedMaskedMap(test.input, test.replacements, test.masks, false, false)
		if utils.CheckAreMapsEqual(output, test.output) == false {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for BoundarySplitPopMap()
func TestBoundarySplitPopMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input        map[string]int
		replacements string
		output       map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc123": 1, "ABC": 2, "ABCabc123!!!": 3}, "luds", map[string]int{"!!!": 1, "ABC": 2, "abc": 2, "123": 2}},
		{map[string]int{"123ABC": 1, "123456ABC": 2, "1Z2X39": 3}, "d", map[string]int{"1": 1, "123": 1, "123456": 1, "2": 1, "39": 1}},
		{map[string]int{"12🙂test": 1, "😀test": 2, "test😁": 3}, "b", map[string]int{"🙂": 1, "😀": 1, "😁": 1}},
	}

	// Run test cases
	for _, test := range tests {
		output := BoundarySplitPopMap(test.input, test.replacements, false, false)
		if !utils.CheckAreMapsEqual(output, test.output) {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for ShuffleMap()
func TestShuffleMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input        map[string]int
		replacements string
		swaps        map[string]int
		output       map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc?d?d?d": 1, "?u?u?u456": 2, "ABC?l?l?l123!!!": 3}, "luds", map[string]int{"DEF": 1, "321": 1, "zxc": 1}, map[string]int{"ABCzxc123!!!": 3, "DEF456": 2, "abc321": 1}},
		{map[string]int{"?d?d?dabc?d?d?d": 1, "123?l?l?l": 2, "ABC?l?l?l123!!!": 3}, "luds", map[string]int{"DEF": 1, "321": 1, "zxc": 1}, map[string]int{"123zxc": 2, "ABCzxc123!!!": 3}},
	}

	// Run test cases
	for _, test := range tests {
		output := ShuffleMap(test.input, test.replacements, test.swaps, false, false)
		if !utils.CheckAreMapsEqual(output, test.output) {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for CalculateKeyspace()
func TestCalculateKeyspace(t *testing.T) {
	// Define a test case struct
	type testCase struct {
		input  string
		output int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"?l?l?l123", 78},
		{"?l?l?l?d?d?d", 108},
		{"?l?l?l?d?d?d?s?s?s", 204},
		{"?u?u?l?d?d?d?s?s?s", 204},
	}

	// Run test cases
	for _, test := range tests {
		output := CalculateKeyspace(test.input)
		if output != test.output {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}

// Unit Test for IsMaskAFullMask()
func TestIsMaskAFullMask(t *testing.T) {
	// Define a test case struct
	type testCase struct {
		input  string
		output bool
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"?l?l?l123", false},
		{"?l?l?l?d?d?d", true},
		{"?l?l?l?d?d?d?s?s?s", true},
		{"?u?u?l?d?d?d?s?s?s", true},
	}

	// Run test cases
	for _, test := range tests {
		output := IsMaskAFullMask(test.input)
		if output != test.output {
			t.Errorf("Test failed: %v inputted, %v expected, %v returned", test.input, test.output, output)
		}
	}
}
