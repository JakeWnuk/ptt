package transform

import (
	"ptt/pkg/utils"
	"testing"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Generation Functions **
// - ReplaceKeysInMap
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// - TransformationController (TransformationController)
// - FuzzyReplaceKeysInMap (Generation Functions)

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
		result := ReplaceKeysInMap(test.input, test.replace)
		if utils.CheckAreMapsEqual(result, test.output) == false {
			t.Errorf("Test case failed. Expected %v, got %v", test.output, result)
		}
	}
}
