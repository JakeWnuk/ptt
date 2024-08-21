package rule

import (
	"testing"

	"github.com/jakewnuk/ptt/pkg/utils"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Transformation Functions **
// - LenToRule()
// - CharToRule()
// - CharToIteratingRule()
// - StringToToggleRule()
//
// ** Output Functions **
// - FormatCharToRuleOutput()
// - FormatCharToIteratingRuleOutput()
// - AppendRules()
// - PrependRules()
// - InsertRules()
// - OverwriteRules()
// - ToggleRules()
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// -

// Unit Test for LenToRule()
func TestLenToRule(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		insert string
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"test", "]", "] ] ] ]"},
		{"test", "[", "[ [ [ ["},
		{"test", "}", "} } } }"},
		{"test", "{", "{ { { {"},
	}

	// Run test cases
	for _, test := range tests {
		result := LenToRule(test.input, test.insert)
		if result != test.output {
			t.Errorf("Expected %v, but got %v", test.output, result)
		}
	}
}

// Unit Test for CharToRule()
func TestCharToRule(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		insert string
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"test", "$", "$t $e $s $t"},
		{"test", "^", "^t ^e ^s ^t"},
		{"test", "@", "@t @e @s @t"},
	}

	// Run test cases
	for _, test := range tests {
		result := CharToRule(test.input, test.insert)
		if result != test.output {
			t.Errorf("Expected %v, but got %v", test.output, result)
		}
	}
}

// Unit Test for CharToIteratingRule()
func TestCharToIteratingRule(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		insert string
		index  int
		output string
	}

	type testCases []testCase

	// Define test cases
	// Note: Multibyte characters are not handled in this function
	tests := testCases{
		{"test", "i", 0, "i0t i1e i2s i3t"},
		{"爱test", "i", 1, "i1爱 i4t i5e i6s i7t"},
		{"爱test", "i", 2, "i2爱 i5t i6e i7s i8t"},
		{"test123", "i", 3, "i3t i4e i5s i6t i71 i82 i93"},
	}

	// Run test cases
	for _, test := range tests {
		result := CharToIteratingRule(test.input, test.insert, test.index)
		if result != test.output {
			t.Errorf("Expected %v, but got %v", test.output, result)
		}
	}
}

// Unit Test for StringToToggleRule()
func TestStringToToggleRule(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		insert string
		index  int
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"tEst", "t", 0, "t1"},
		{"爱Test", "t", 1, "t4"},
		{"爱tesT", "t", 2, "t8"},
		{"TEST123", "t", 0, "t0 t1 t2 t3"},
	}

	// Run test cases
	for _, test := range tests {
		result := StringToToggleRule(test.input, test.insert, test.index)
		if result != test.output {
			t.Errorf("Expected %v, but got %v", test.output, result)
		}
	}
}

// Unit Test for FormatCharToRuleOutput()
func TestFormatCharToRuleOutput(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{"$t $e $s $t", "$t $e $s $t"},
		{"^t ^e ^s ^t", "^t ^e ^s ^t"},
		{"$爱 $t $e $s $t", "$\\xE7 $\\x88 $\\xB1 $t $e $s $t"},
		{"^爱 ^t ^e ^s ^t", "^\\xB1 ^\\x88 ^\\xE7 ^t ^e ^s ^t"},
	}

	// Run test cases
	for _, test := range tests {
		result := FormatCharToRuleOutput(test.input)
		if result != test.output {
			t.Errorf("Expected %v, but got %v", test.output, result)
		}
	}
}

// Unit Test for FormatCharToIteratingRuleOutput()
func TestFormatCharToIteratingRuleOutput(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		index  int
		input  string
		output string
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{0, "i0t i1e i2s i3t", "i0t i1e i2s i3t"},
		{1, "i1爱 i4t i5e i6s i7t", "i1\\xE7 i2\\x88 i3\\xB1 i4t i5e i6s i7t"},
		{2, "i2爱 i5t i6e i7s i8t", "i2\\xE7 i3\\x88 i4\\xB1 i5t i6e i7s i8t"},
		{3, "i3t i4e i5s i6t i71 i82 i93", "i3t i4e i5s i6t i71 i82 i93"},
	}

	// Run test cases
	for _, test := range tests {
		result := FormatCharToIteratingRuleOutput(test.index, test.input)
		if result != test.output {
			t.Errorf("Expected %v, but got %v", test.output, result)
		}
	}
}

// Unit Test for AppendRules()
func TestAppendRules(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		items     map[string]int
		operation string
		output    map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc": 1, "efg": 2}, "append", map[string]int{"$a $b $c": 1, "$e $f $g": 2}},
		{map[string]int{"abc": 1, "efg": 2}, "append-remove", map[string]int{"] ] ] $a $b $c": 1, "] ] ] $e $f $g": 2}},
		{map[string]int{"爱test": 1, "a爱test": 2}, "append", map[string]int{"$\\xE7 $\\x88 $\\xB1 $t $e $s $t": 1, "$a $\\xE7 $\\x88 $\\xB1 $t $e $s $t": 2}},
		{map[string]int{"test ": 1}, "append", map[string]int{"$t $e $s $t $ :": 1}},
	}

	// Run test cases
	for _, test := range tests {
		given := AppendRules(test.items, test.operation, false, false)
		if !utils.CheckAreMapsEqual(given, test.output) {
			t.Errorf("Expected %v, but got %v", test.output, given)
		}
	}
}

// Unit Test for PrependRules()
func TestPrependRules(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		items     map[string]int
		operation string
		output    map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc": 1, "efg": 2}, "prepend", map[string]int{"^c ^b ^a": 1, "^g ^f ^e": 2}},
		{map[string]int{"abc": 1, "efg": 2}, "prepend-remove", map[string]int{"[ [ [ ^c ^b ^a": 1, "[ [ [ ^g ^f ^e": 2}},
		{map[string]int{"爱test": 1, "a爱test": 2}, "prepend", map[string]int{"^t ^s ^e ^t ^\\xB1 ^\\x88 ^\\xE7": 1, "^t ^s ^e ^t ^\\xB1 ^\\x88 ^\\xE7 ^a": 2}},
	}

	// Run test cases
	for _, test := range tests {
		given := PrependRules(test.items, test.operation, false, false)
		if !utils.CheckAreMapsEqual(given, test.output) {
			t.Errorf("Expected %v, but got %v", test.output, given)
		}
	}
}

// Unit Test for InsertRules()
func TestInsertRules(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		items      map[string]int
		startIndex string
		endIndex   string
		output     map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc": 1, "efg": 2}, "0", "0", map[string]int{"i0a i1b i2c": 1, "i0e i1f i2g": 2}},
		{map[string]int{"爱test": 1, "a爱test": 2}, "1", "1", map[string]int{"i1\\xE7 i2\\x88 i3\\xB1 i4t i5e i6s i7t": 1, "i1a i2\\xE7 i3\\x88 i4\\xB1 i5t i6e i7s i8t": 2}},
		{map[string]int{"abc": 1, "efg": 2}, "1", "2", map[string]int{"i1a i2b i3c": 1, "i1e i2f i3g": 2, "i2a i3b i4c": 1, "i2e i3f i4g": 2}},
		{map[string]int{"abc": 1}, "6", "8", map[string]int{"i6a i7b i8c": 1, "i8a i9b iAc": 1, "i7a i8b i9c": 1}},
	}

	// Run test cases
	for _, test := range tests {
		given := InsertRules(test.items, test.startIndex, test.endIndex, false, false)
		if !utils.CheckAreMapsEqual(given, test.output) {
			t.Errorf("Expected %v, but got %v", test.output, given)
		}
	}
}

// Unit Test for OverwriteRules()
func TestOverwriteRules(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		items      map[string]int
		startIndex string
		endIndex   string
		output     map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"abc": 1, "efg": 2}, "0", "0", map[string]int{"o0a o1b o2c": 1, "o0e o1f o2g": 2}},
		{map[string]int{"爱test": 1, "a爱test": 2}, "1", "1", map[string]int{"o1\\xE7 o2\\x88 o3\\xB1 o4t o5e o6s o7t": 1, "o1a o2\\xE7 o3\\x88 o4\\xB1 o5t o6e o7s o8t": 2}},
		{map[string]int{"abc": 1, "efg": 2}, "1", "2", map[string]int{"o1a o2b o3c": 1, "o1e o2f o3g": 2, "o2a o3b o4c": 1, "o2e o3f o4g": 2}},
		{map[string]int{"abc": 1}, "6", "8", map[string]int{"o6a o7b o8c": 1, "o8a o9b oAc": 1, "o7a o8b o9c": 1}},
	}

	// Run test cases
	for _, test := range tests {
		given := OverwriteRules(test.items, test.startIndex, test.endIndex, false, false)
		if !utils.CheckAreMapsEqual(given, test.output) {
			t.Errorf("Expected %v, but got %v", test.output, given)
		}
	}
}

// Unit Test for ToggleRules()
func TestToggleRules(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		items      map[string]int
		startIndex string
		endIndex   string
		output     map[string]int
	}

	type testCases []testCase

	// Define test cases
	tests := testCases{
		{map[string]int{"aBc": 1, "EfG": 2}, "0", "0", map[string]int{"T1": 1, "T0 T2": 2}},
		{map[string]int{"爱tesT": 1, "a爱Test": 2}, "0", "0", map[string]int{"T4": 2, "T6": 1}},
		{map[string]int{"aBc": 1, "EfG": 2}, "1", "2", map[string]int{"T2": 1, "T3": 1, "T1 T3": 2, "T2 T4": 2}},
	}

	// Run test cases
	for _, test := range tests {
		given := ToggleRules(test.items, test.startIndex, test.endIndex, false, false)
		if !utils.CheckAreMapsEqual(given, test.output) {
			t.Errorf("Expected %v, but got %v", test.output, given)
		}
	}
}
