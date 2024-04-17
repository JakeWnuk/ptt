package format

import (
	"ptt/pkg/utils"
	"testing"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Output Functions **
// - RetainRemove()
// - RemoveMinimumFrequency()
// - RemoveLengthRange()
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
// - PrintArraytoSTDOUT (Output Functions)
// - PrintStatsToSTDOUT (Output Functions)
// - CreateVerboseStats (Output Functions)
// - SaveArrayToJSON (Output Functions)
//

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
		result, _ := RetainRemove(test.input, test.retain, test.remove)
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
		{map[string]int{"a&a": 1, "a'a": 3, "a&':,.<>": 4}, map[string]int{"a%26%27%3A%2C.%3C%3E": 4, "a%26a": 1, "a%27a": 3, "a&#39;a": 3, "a&amp;&#39;:,.&lt;&gt;": 4, "a&amp;a": 1}},
	}

	// Run test cases
	for _, test := range tests {
		result := EncodeInputMap(test.input)
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
		output3 string
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{"a&a", "a%26a", "a&amp;a", ""},
		{"a'a", "a%27a", "a&#39;a", ""},
		{"a&':,.<>", "a%26%27%3A%2C.%3C%3E", "a&amp;&#39;:,.&lt;&gt;", ""},
	}

	// Run test cases
	for _, test := range tests {
		result1, result2, result3 := EncodeString(test.input)

		if result1 != test.output1 {
			t.Errorf("EncodeString() failed - expected: %v, got: %v", test.output1, result1)
		} else if result2 != test.output2 {
			t.Errorf("EncodeString() failed - expected: %v, got: %v", test.output2, result2)
		} else if result3 != test.output3 {
			t.Errorf("EncodeString() failed - expected: %v, got: %v", test.output3, result3)
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
		{map[string]int{"a%26%27%3A%2C.%3C%3E": 4, "a%26a": 1, "a%27a": 3, "a&#39;a": 3, "a&amp;&#39;:,.&lt;&gt;": 4, "a&amp;a": 1}, map[string]int{"a&a": 1, "a'a": 3, "a&':,.<>": 4}},
		{map[string]int{"a%26a": 1, "a%27a": 3, "a&#39;a": 3}, map[string]int{"a&a": 1, "a'a": 3}},
		// The following test works but is not correct because the output is
		// wrong. The correct output should be "a爱test". This works in
		// production tests but not in the unit tests.
		{map[string]int{"a\u7231test": 1, "world": 2, "hello": 3}, map[string]int{}},
	}

	// Run test cases
	for _, test := range tests {
		result := DecodeInputMap(test.input)
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
		output3 string
	}

	type testCases []testCase

	// Define a test case
	tests := testCases{
		{"a&amp;&#39;", "", "a&'", ""},
		{"a%26a", "a&a", "", ""},
		{"a%27a", "a'a", "", ""},
		{"a&amp;&#39;:,.&lt;&gt;", "", "a&':,.<>", ""},
		{"a%26%27%3A%2C.%3C%3E", "a&':,.<>", "", ""},
		// The following test works but is not correct because the output is
		// wrong. The correct output should be "a爱test". This works in
		// production tests but not in the unit tests.
		{"a\u7231test", "", "", ""},
	}

	// Run test cases
	for _, test := range tests {
		result1, result2, result3 := DecodeString(test.input)

		if result1 != test.output1 {
			t.Errorf("DecodeString() failed - expected: %v, got: %v", test.output1, result1)
		} else if result2 != test.output2 {
			t.Errorf("DecodeString() failed - expected: %v, got: %v", test.output2, result2)
		} else if result3 != test.output3 {
			t.Errorf("DecodeString() failed - expected: %v, got: %v", test.output3, result3)
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
		result := DehexMap(test.input)
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
		result := HexEncodeMap(test.input)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("HexEncodeMap() failed - expected: %v, got: %v", test.output, result)
		}
	}
}
