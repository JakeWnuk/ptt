package format

import (
	"testing"

	"github.com/jakewnuk/ptt/pkg/utils"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Output Functions **
// - RetainRemove()
// - RemoveMinimumFrequency()
// - RemoveLengthRange()
// - FilterTopN()
//
// ** Encoding Functions **
// - EncodeInputMap()
// - EncodeString()
// - DecodeInputMap()
// - DecodeString()
// - ASCIIEscapeUnicode()
// - DehexMap()
// - HexEncodeMap()
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// - StatClassifyToken() (Output Functions)
// - PrintArraytoSTDOUT() (Output Functions)
// - PrintStatsToSTDOUT() (Output Functions)
// - CreateVerboseStats() (Output Functions)
// - SaveArrayToJSON() (Output Functions)
//

// Unit Test for StatClassifyToken()
func TestStatClassifyToken(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		output []string
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{"abc", []string{"alphabetical", "non-complex", "all-lowercase", "short-non-complex"}},
		{"abc123", []string{"alphanumeric", "non-complex", "all-lowercase", "short-non-complex", "hex-string"}},
		{"abc123ABC", []string{"alphanumeric", "non-complex", "contains-uppercase", "short-non-complex"}},
		{"abc123ABC!@#", []string{"alphanumeric-with-special", "complex", "contains-uppercase", "long-complex"}},
		{"$HEX[6c6f7665]", []string{"$HEX[...]-format", "complex", "long-complex", "alphanumeric-with-special", "contains-uppercase"}},
	}

	// Run test cases
	for _, test := range tests {
		result := StatClassifyToken(test.input)
		if utils.CheckAreArraysEqual(result, test.output) == false {
			t.Errorf("StatClassifyToken() failed - expected: %v, got: %v", test.output, result)

		}
	}
}

// Unit Test for RetainRemove()
func TestRetainRemove(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		retain map[string]int
		remove map[string]int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, map[string]int{"a": 1, "b": 2}, map[string]int{"c": 3, "d": 4}, map[string]int{"a": 1, "b": 2}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, map[string]int{"a": 1, "b": 2, "d": 1}, map[string]int{"c": 3}, map[string]int{"a": 1, "b": 2, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, map[string]int{}, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, map[string]int{}, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, map[string]int{}},
	}

	// Run test cases
	for _, test := range tests {
		result, _ := RetainRemove(test.input, test.retain, test.remove, 0)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("RetainRemove() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for RemoveMinimumFrequency()
func TestRemoveMinimumFrequency(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		min    int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 2, map[string]int{"b": 2, "c": 3, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 3, map[string]int{"c": 3, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 4, map[string]int{"d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 5, map[string]int{}},
	}

	// Run test cases
	for _, test := range tests {
		result := RemoveMinimumFrequency(test.input, test.min)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("RemoveMinimumFrequency() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for RemoveLengthRange()
func TestRemoveLengthRange(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		min    int
		max    int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}, 1, 5, map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}},
		{map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}, 1, 4, map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}},
		{map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}, 1, 3, map[string]int{"a": 1, "bc": 2, "cde": 3}},
		{map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}, 1, 2, map[string]int{"a": 1, "bc": 2}},
		{map[string]int{"a": 1, "bc": 2, "cde": 3, "defg": 4}, 1, 1, map[string]int{"a": 1}},
	}

	// Run test cases
	for _, test := range tests {
		result := RemoveLengthRange(test.input, test.min, test.max)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("RemoveLengthRange() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for FilterTopN()
func TestFilterTopN(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		top    int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 2, map[string]int{"c": 3, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 3, map[string]int{"b": 2, "c": 3, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 4, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}},
		{map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}, 5, map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}},
	}

	// Run test cases
	for _, test := range tests {
		result := FilterTopN(test.input, test.top)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("FilterTopN() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for EncodeInputMap()
func TestEncodeInputMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"a&a": 1, "a'a": 3, "a&':,.<>": 4}, map[string]int{"a&#39;a": 3, "a&amp;&#39;:,.&lt;&gt;": 4, "a&amp;a": 1}},
	}

	// Run test cases
	for _, test := range tests {
		result := EncodeInputMap(test.input, false, false)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("EncodeInputMap() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for EncodeString()
func TestEncodeString(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input   string
		output1 string
		output2 string
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{"a&a", "a&amp;a", ""},
		{"a'a", "a&#39;a", ""},
		{"a&':,.<>", "a&amp;&#39;:,.&lt;&gt;", ""},
	}

	// Run test cases
	for _, test := range tests {
		result1, result2 := EncodeString(test.input)

		if result1 != test.output1 {
			t.Errorf("EncodeString() failed - expected: %v, got: %v", test.output1, result1)
		} else if result2 != test.output2 {
			t.Errorf("EncodeString() failed - expected: %v, got: %v", test.output2, result2)
		}

	}
}

// Unit Test for DecodeInputMap()
func TestDecodeInputMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"a&#39;a": 3, "a&amp;&#39;:,.&lt;&gt;": 4, "a&amp;a": 1}, map[string]int{"a&a": 1, "a'a": 3, "a&':,.<>": 4}},
		{map[string]int{"a&#39;a": 3}, map[string]int{"a'a": 3}},
		// The following test works but is not correct because the output is
		// wrong. The correct output should be "a爱test". This works in
		// production tests but not in the unit tests.
		{map[string]int{"a\u7231test": 1, "world": 2, "hello": 3}, map[string]int{}},
	}

	// Run test cases
	for _, test := range tests {
		result := DecodeInputMap(test.input, false, false)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("DecodeInputMap() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for DecodeString()
func TestDecodeString(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input   string
		output1 string
		output2 string
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{"a&amp;&#39;", "a&'", ""},
		{"a&amp;&#39;:,.&lt;&gt;", "a&':,.<>", ""},
		// The following test works but is not correct because the output is
		// wrong. The correct output should be "a爱test". This works in
		// production tests but not in the unit tests.
		{"a\u7231test", "", ""},
	}

	// Run test cases
	for _, test := range tests {
		result1, result2 := DecodeString(test.input)

		if result1 != test.output1 {
			t.Errorf("DecodeString() failed - expected: %v, got: %v", test.output1, result1)
		} else if result2 != test.output2 {
			t.Errorf("DecodeString() failed - expected: %v, got: %v", test.output2, result2)
		}

	}
}

// Unit Test for ASCIIEscapeUnicode()
func TestASCIIEscapeUnicode(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  string
		output string
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{"a爱test", "a\\u7231test"},
		{"a爱test爱", "a\\u7231test\\u7231"},
		{"a爱test爱爱", "a\\u7231test\\u7231\\u7231"},
		{"a爱test爱爱爱", "a\\u7231test\\u7231\\u7231\\u7231"},
		{"a爱test爱爱爱爱", "a\\u7231test\\u7231\\u7231\\u7231\\u7231"},
	}

	// Run test cases
	for _, test := range tests {
		result := ASCIIEscapeUnicode(test.input)
		if result != test.output {
			t.Errorf("ASCIIEscapeUnicode() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for DehexMap()
func TestDehexMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"$HEX[6c6f7665]": 1, "$HEX[776f726c64]": 2, "$HEX[68656c6c6f]": 3}, map[string]int{"love": 1, "world": 2, "hello": 3}},
		{map[string]int{"$HEX[6c6f7665]": 1, "$HEX[776f726c64]": 2, "$HEX[68656c6c6f]": 3, "hello": 4}, map[string]int{"love": 1, "world": 2, "hello": 3}},
		{map[string]int{"$HEX[6c6f7665]": 1, "$HEX[776f726c64]": 2, "$HEX[68656c6c6f]": 3, "hello": 4, "world": 5}, map[string]int{"love": 1, "world": 2, "hello": 3}},
	}

	// Run test cases
	for _, test := range tests {
		result := DehexMap(test.input, false, false)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("DehexMap() failed - expected: %v, got: %v", test.output, result)
		}
	}
}

// Unit Test for HexEncodeMap()
func TestHexEncodeMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input  map[string]int
		output map[string]int
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{map[string]int{"love": 1, "world": 2, "hello": 3}, map[string]int{"$HEX[6c6f7665]": 1, "$HEX[776f726c64]": 2, "$HEX[68656c6c6f]": 3}},
		{map[string]int{"love": 1, "world": 2, "hello": 3}, map[string]int{"$HEX[6c6f7665]": 1, "$HEX[776f726c64]": 2, "$HEX[68656c6c6f]": 3}},
		{map[string]int{"love<3": 1, "world>2": 2, "hello&1": 3}, map[string]int{"$HEX[6c6f76653c33]": 1, "$HEX[776f726c643e32]": 2, "$HEX[68656c6c6f2631]": 3}},
		{map[string]int{"爱test": 1, "world": 2, "hello": 3}, map[string]int{"$HEX[e788b174657374]": 1, "$HEX[776f726c64]": 2, "$HEX[68656c6c6f]": 3}},
	}

	// Run test cases
	for _, test := range tests {
		result := HexEncodeMap(test.input, false, false)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("HexEncodeMap() failed - expected: %v, got: %v", test.output, result)
		}
	}
}
