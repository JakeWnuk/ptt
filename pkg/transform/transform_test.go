package transform

import (
	"testing"

	"github.com/jakewnuk/ptt/pkg/utils"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Generation Functions **
// - ReplaceKeysInMap()
// - ReplaceAllKeysInMap()
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// - TransformationController (TransformationController)
// - MakePassphraseMap (Generation Functions)
// - GeneratePassphrase (Generation Functions)

// Unit Test for ReplaceKeysInMap
func TestReplaceKeysInMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input   map[string]int
		replace map[string]int
		output  map[string]int
	}

	type testCases []testCase

	// Define the test cases
	tests := testCases{
		{map[string]int{"a": 1, "b": 2, "c": 3}, map[string]int{"a:foo": 1, "b:bar": 2, "c:love": 3}, map[string]int{"foo": 1, "bar": 2, "love": 3}},
		{map[string]int{"123Testing": 1, "456Testing": 2, "789Testing": 3}, map[string]int{"123:foo": 1, "456:bar": 2, "789:love": 3}, map[string]int{"fooTesting": 1, "barTesting": 2, "loveTesting": 3}},
		{map[string]int{"爱test": 1}, map[string]int{"爱:foo": 1, "爱:bar": 2, "爱:love": 3}, map[string]int{"bartest": 1, "footest": 1, "lovetest": 1}},
	}

	// Run the test cases
	for _, test := range tests {
		result := ReplaceKeysInMap(test.input, test.replace, false, false)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("Test case failed. Expected %v, got %v", test.output, result)
		}
	}
}

// Unit Test for ReplaceAllKeysInMap
func TestReplaceAllKeysInMap(t *testing.T) {

	// Define a test case struct
	type testCase struct {
		input   map[string]int
		replace map[string]int
		output  map[string]int
	}

	type testCases []testCase

	// Define the test cases
	tests := testCases{
		{map[string]int{"abc": 1}, map[string]int{"a:1": 1, "b:2": 2, "c:3": 3}, map[string]int{"123": 1}},
		{map[string]int{"123Testing": 1, "456Testing": 2, "789Testing": 3}, map[string]int{"123:foo": 1, "456:bar": 2, "789:love": 3}, map[string]int{"fooTesting": 1, "barTesting": 2, "loveTesting": 3}},
	}

	// Run the test cases
	for _, test := range tests {
		result := ReplaceAllKeysInMap(test.input, test.replace, false, false)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("Test case failed. Expected %v, got %v", test.output, result)
		}
	}
}
