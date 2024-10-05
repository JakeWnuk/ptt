package utils

import (
	"testing"

	"github.com/jakewnuk/ptt/pkg/models"
)

// ----------------------------------------------------------------------------
// Functions with Unit Tests
// ----------------------------------------------------------------------------
// ** Loading and Processing Functions **
// - ReadFilesToMap()
// - LoadStdinToMap()
// - CombineMaps()
// - ReadJSONToArray()
//
// ** Transformation Functions **
// - ReverseString()
// - ConvertMultiByteCharToRule()
// - IncrementIteratingRuleCall()
// - ConvertMultiByteCharToIteratingRule()
// - SplitBySeparatorString()
// - ReplaceSubstring()
// - ReplaceAllSubstring()
// - SubstringMap()
// - GenerateNGrams()
// - GeneratePassphrase()
//
// ** Validation Functions **
// - CheckASCIIString()
// - CheckHexString()
// - CheckAreMapsEqual()
// - CheckAreArraysEqual()
// - IsValidURL()
//
// ----------------------------------------------------------------------------
// Functions without Unit Tests
// ----------------------------------------------------------------------------
// - ReadURLsToMap() (Loading and Processing Functions)
// - ProcessURL() (Loading and Processing Functions)
// - ProcessURLFile() (Loading and Processing Functions)
// - GetFilesInDirectory() (Loading and Processing Functions)
// - IsValidFile() (Validation Functions)
// - IsFileSystemDirectory() (Validation Functions)

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

// Unit Test for ReadJSONToArray()
func TestReadJSONToArray(t *testing.T) {
	// Define a test case struct
	type TestCase struct {
		Input  string
		Output []models.TemplateFileOperation
	}

	type TestCases []TestCase

	// Create a mock file system with example files
	mockFs := &models.MockFileSystem{
		Files: map[string][]byte{
			"file1": []byte(`[{"StartIndex":0,"EndIndex":4,"Verbose":true,"ReplacementMask":"uldbs","Bypass":false,"TransformationMode":"append"}]`),
			"file2": []byte(`[{"StartIndex":0,"EndIndex":4,"Verbose":true,"ReplacementMask":"uldbs","Bypass":false,"TransformationMode":"append"},{"StartIndex":0,"EndIndex":4,"Verbose":true,"ReplacementMask":"uldbs","Bypass":false,"TransformationMode":"append"}]`),
			"file3": []byte(`[{"StartIndex":0,"EndIndex":4,"Verbose":true,"ReplacementMask":"uldbs","Bypass":false,"TransformationMode":"append"},{"StartIndex":0,"EndIndex":4,"Verbose":true,"ReplacementMask":"uldbs","Bypass":false,"TransformationMode":"append"},{"StartIndex":0,"EndIndex":4,"Verbose":true,"ReplacementMask":"uldbs","Bypass":false,"TransformationMode":"append"}]`),
		},
	}

	// Define test cases
	testCases := TestCases{
		{"file1", []models.TemplateFileOperation{{StartIndex: 0, EndIndex: 4, Verbose: true, ReplacementMask: "uldbs", Bypass: false, TransformationMode: "append"}}},
		{"file2", []models.TemplateFileOperation{{StartIndex: 0, EndIndex: 4, Verbose: true, ReplacementMask: "uldbs", Bypass: false, TransformationMode: "append"}, {StartIndex: 0, EndIndex: 4, Verbose: true, ReplacementMask: "uldbs", Bypass: false, TransformationMode: "append"}}},
		{"file3", []models.TemplateFileOperation{{StartIndex: 0, EndIndex: 4, Verbose: true, ReplacementMask: "uldbs", Bypass: false, TransformationMode: "append"}, {StartIndex: 0, EndIndex: 4, Verbose: true, ReplacementMask: "uldbs", Bypass: false, TransformationMode: "append"}, {StartIndex: 0, EndIndex: 4, Verbose: true, ReplacementMask: "uldbs", Bypass: false, TransformationMode: "append"}}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := []string{testCase.Input}
		output := testCase.Output

		given := ReadJSONToArray(mockFs, input)
		if given[0].StartIndex != output[0].StartIndex || given[0].EndIndex != output[0].EndIndex || given[0].Verbose != output[0].Verbose || given[0].ReplacementMask != output[0].ReplacementMask || given[0].Bypass != output[0].Bypass || given[0].TransformationMode != output[0].TransformationMode {
			t.Errorf("ReadJSONToArray(%v) = %v; want %v", input, given, output)
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

// Unit Test for ConvertMultiByteCharToRule()
func TestConvertMultiByteCharToRule(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  string
		Output string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"$l $o $v $e", "$l $o $v $e"},
		{"^e ^v ^o ^l", "^e ^v ^o ^l"},
		{"$a $爱", "$a $\\xE7 $\\x88 $\\xB1"},
		{"^b ^愛", "^b ^\\x9B ^\\x84 ^\\xE6"},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		given := ConvertMultiByteCharToRule(input)
		if given != output {
			t.Errorf("ConvertMultiByteCharToRule(%v) = %v; want %v", input, given, output)
		}
	}
}

// Unit Test for IncrementIteratingRuleCall()
func TestIncrementIteratingRuleCall(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  string
		Output string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"i0", "i1"},
		{"i1", "i2"},
		{"i2", "i3"},
		{"i3", "i4"},
		{"iA", "iB"},
		{"iB", "iC"},
		{"iC", "iD"},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		given := IncrementIteratingRuleCall(input)
		if given != output {
			t.Errorf("IncrementIteratingRuleCall(%v) = %v; want %v", input, given, output)
		}
	}
}

// Unit Test for ConvertMultiByteCharToIteratingRule()
func TestConvertMultiByteCharToIteratingRule(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Index  int
		Input  string
		Output string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{0, "i0l i1o i2v i3e", "i0l i1o i2v i3e "},
		{0, "i0a i1爱", "i0a i1\\xE7 i2\\x88 i3\\xB1 "},
		{1, "i1a i2愛", "i1a i2\\xE6 i3\\x84 i4\\x9B "},
		{0, "i0爱 i3t i4e i5s i6t", "i0\\xE7 i1\\x88 i2\\xB1 i3t i4e i5s i6t "},
	}

	// Run test cases
	for _, testCase := range testCases {
		index := testCase.Index
		input := testCase.Input
		output := testCase.Output

		given := ConvertMultiByteCharToIteratingRule(index, input)
		if given != output {
			t.Errorf("ConvertMultiByteCharToIteratingRule(%v) = %v; want %v", input, given, output)
		}
	}
}

// Unit Test for SplitBySeparatorString()
func TestSplitBySeparatorString(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 string
		Input2 string
		Output []string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I love you", "love", []string{"I ", "love", " you"}},
		{"I <3 you", "<3", []string{"I ", "<3", " you"}},
		{"I 爱 you", "爱", []string{"I ", "爱", " you"}},
		{"13Teststreet31p", "street", []string{"13Test", "street", "31p"}},
		{"123131asdasd", "131", []string{"123", "131", "asdasd"}},
		{"12313zxczxc", "13", []string{"123", "13", "zxczxc"}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := SplitBySeparatorString(input1, input2)
		if given[0] != output[0] || given[1] != output[1] || given[2] != output[2] {
			t.Errorf("SplitBySeparatorString(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for ReplaceSubstring()
func TestReplaceSubstring(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input        string
		Replacements map[string]int
		Output       []string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I love you", map[string]int{"love:miss": 1}, []string{"I miss you"}},
		{"I <3 you", map[string]int{"<3:heart": 1}, []string{"I heart you"}},
		{"I 爱 you", map[string]int{"爱:love": 1}, []string{"I love you"}},
		{"I love you", map[string]int{"love:爱": 1}, []string{"I 爱 you"}},
		{"13Teststreet31p", map[string]int{"street:road": 1}, []string{"13Testroad31p"}},
		{"123131asdasd", map[string]int{"131:313": 1}, []string{"123313asdasd"}},
		{"12313zxczxc", map[string]int{"13:31": 1}, []string{"12331zxczxc"}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		replacements := testCase.Replacements
		output := testCase.Output

		given := ReplaceSubstring(input, replacements)
		if given[0] != output[0] {
			t.Errorf("ReplaceSubstring(%v, %v) = %v; want %v", input, replacements, given, output)
		}
	}
}

// Unit Test for ReplaceAllSubstring()
func TestReplaceAllSubstring(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input        string
		Replacements map[string]int
		Output       []string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I love you", map[string]int{"love:miss": 1}, []string{"I miss you"}},
		{"I <3 you", map[string]int{"<3:heart": 1}, []string{"I heart you"}},
		{"I 爱 you", map[string]int{"爱:love": 1}, []string{"I love you"}},
		{"I love you", map[string]int{"love:爱": 1}, []string{"I 爱 you"}},
		{"13Teststreet31p", map[string]int{"street:road": 1}, []string{"13Testroad31p"}},
		{"123131asdasd", map[string]int{"131:313": 1}, []string{"123313asdasd"}},
		{"12313zxczxc", map[string]int{"13:31": 1}, []string{"12331zxczxc"}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		replacements := testCase.Replacements
		output := testCase.Output

		given := ReplaceAllSubstring(input, replacements)
		if given[0] != output[0] {
			t.Errorf("ReplaceAllSubstring(%v, %v) = %v; want %v", input, replacements, given, output)
		}
	}

}

// Unit Test for SubstringMap()
func TestSubstringMap(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		input  map[string]int
		start  int
		end    int
		output map[string]int
	}

	// Define test cases
	testCases := []TestCase{
		{map[string]int{"love1": 1, "love2": 1, "love3": 1}, 0, 4, map[string]int{"love": 3}},
		{map[string]int{"<31": 1, "<32": 1, "<33": 1}, 1, 2, map[string]int{"3": 3}},
		{map[string]int{"爱1": 1, "爱2": 1, "爱3": 1}, 0, 3, map[string]int{"爱": 3}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.input
		start := testCase.start
		end := testCase.end
		output := testCase.output

		given := SubstringMap(input, start, end, false, false)
		if CheckAreMapsEqual(given, output) == false {
			t.Errorf("SubstringMap(%v, %v, %v) = %v; want %v", input, start, end, given, output)
		}
	}
}

// Unit Test for CheckASCIIString()
func TestCheckASCIIString(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  string
		Output bool
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I love you", true},
		{"I <3 you", true},
		{"I 爱 you", false},
		{"13Teststreet31p", true},
		{"123131asdasd", true},
		{"12313zxczxc", true},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		given := CheckASCIIString(input)
		if given != output {
			t.Errorf("CheckASCIIString(%v) = %v; want %v", input, given, output)
		}
	}
}

// Unit Test for CheckHexString()
func TestCheckHexString(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  string
		Output bool
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I love you", false},
		{"I 爱 you", false},
		{"$HEX[6c6f7665]", true},
		{"$HEX[3c33]", true},
		{"$HEX[e723]", true},
		{"$HEX[616d6f72]", true},
		{"$HEX[616d6f7572]", true},
		{"$HEX[e6b581]", true},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		given := CheckHexString(input)
		if given != output {
			t.Errorf("CheckHexString(%v) = %v; want %v", input, given, output)
		}
	}
}

// Unit Test for CheckAreMapsEqual()
func TestCheckAreMapsEqual(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 map[string]int
		Input2 map[string]int
		Output bool
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{map[string]int{"love1": 1, "love2": 1, "love3": 1}, map[string]int{"love1": 1, "love2": 1, "love3": 1}, true},
		{map[string]int{"<31": 1, "<32": 1, "<33": 1}, map[string]int{"<31": 1, "<32": 1, "<33": 1}, true},
		{map[string]int{"爱1": 1, "爱2": 1, "爱3": 1}, map[string]int{"爱1": 1, "爱2": 1, "爱3": 1}, true},
		{map[string]int{"amor1": 1, "amor2": 1, "amor3": 1}, map[string]int{"amor1": 1, "amor2": 1, "amor3": 1}, true},
		{map[string]int{"amour1": 1, "amour2": 1, "amour3": 1}, map[string]int{"amour1": 1, "amour2": 1, "amour3": 1}, true},
		{map[string]int{"愛1": 1, "愛2": 1, "愛3": 1}, map[string]int{"愛1": 1, "愛2": 1, "愛3": 1}, true},
		{map[string]int{"love1": 1, "love2": 1, "love3": 1}, map[string]int{"love1": 1, "love2": 1, "love3": 2}, false},
		{map[string]int{"<31": 1, "<32": 1, "<33": 1}, map[string]int{"<31": 1, "<32": 1, "<33": 2}, false},
		{map[string]int{"爱1": 1, "爱2": 1, "爱3": 1}, map[string]int{"爱1": 1, "爱2": 1, "爱3": 2}, false},
		{map[string]int{"amor1": 1, "amor2": 1, "amor3": 1}, map[string]int{"amor1": 1, "amor2": 1, "amor3": 2}, false},
		{map[string]int{"amour1": 1, "amour2": 1, "amour3": 1}, map[string]int{"amour1": 1, "amour2": 1, "amour3": 2}, false},
		{map[string]int{"愛1": 1, "愛2": 1, "愛3": 1}, map[string]int{"愛1": 1, "愛2": 1, "愛3": 2}, false},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := CheckAreMapsEqual(input1, input2)
		if given != output {
			t.Errorf("CheckAreMapsEqual(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for CheckAreArraysEqual()
func TestCheckAreArraysEqual(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 []string
		Input2 []string
		Output bool
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{[]string{"love1", "love2", "love3"}, []string{"love1", "love2", "love3"}, true},
		{[]string{"<31", "<32", "<33"}, []string{"<31", "<32", "<33"}, true},
		{[]string{"爱1", "爱2", "爱3"}, []string{"爱1", "爱2", "爱3"}, true},
		{[]string{"amor1", "amor2", "amor3"}, []string{"amor1", "amor2", "amor3"}, true},
		{[]string{"amour1", "amour2", "amour3"}, []string{"amour1", "amour2", "amour3"}, true},
		{[]string{"愛1", "愛2", "愛3"}, []string{"愛1", "愛2", "愛3"}, true},
		{[]string{"love1", "love2", "love3"}, []string{"love1", "love2", "love3", "love4"}, false},
		{[]string{"<31", "<32", "<33"}, []string{"<31", "<32", "<33", "<34"}, false},
		{[]string{"爱1", "爱2", "爱3"}, []string{"爱1", "爱2", "爱3", "爱4"}, false},
		{[]string{"amor1", "amor2", "amor3"}, []string{"amor1", "amor2", "amor3", "amor4"}, false},
		{[]string{"amour1", "amour2", "amour3"}, []string{"amour1", "amour2", "amour3", "amour4"}, false},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := CheckAreArraysEqual(input1, input2)
		if given != output {
			t.Errorf("CheckAreArraysEqual(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for GenerateNGrams()
func TestGenerateNGrams(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 string
		Input2 int
		Output []string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I love you", 2, []string{"I love", "love you"}},
		{"I <3 you", 2, []string{"I <3", "<3 you"}},
		{"I 爱 you", 2, []string{"I 爱", "爱 you"}},
		{"I love you very much", 3, []string{"I love you", "love you very", "you very much"}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := GenerateNGrams(input1, input2)
		if CheckAreArraysEqual(given, output) == false {
			t.Errorf("GenerateNGrams(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for GeneratePassphrase()
func TestGeneratePassphrase(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input1 string
		Input2 int
		Output []string
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"I <3 you", 3, []string{"I<3You", "I<3YOU", "I-<3-You", "I-<3-YOU", "I_<3_You", "I_<3_YOU", "i <3 you", "i<3you", "I <3 YOU", "I<3YOU", "I<3you", "I <3 You", "I <3 YOU"}},
	}

	// Run test cases
	for _, testCase := range testCases {
		input1 := testCase.Input1
		input2 := testCase.Input2
		output := testCase.Output

		given := GeneratePassphrase(input1, input2)
		if CheckAreArraysEqual(given, output) == false {
			t.Errorf("GeneratePassphrase(%v, %v) = %v; want %v", input1, input2, given, output)
		}
	}
}

// Unit Test for IsValidURL()
func TestIsValidURL(t *testing.T) {

	// Define a test case struct
	type TestCase struct {
		Input  string
		Output bool
	}

	type TestCases []TestCase

	// Define test cases
	testCases := TestCases{
		{"https://www.google.com", true},
		{"http://www.google.com", true},
		{"www.google.com", false},
		{"google.com", false},
		{"google", false},
	}

	// Run test cases
	for _, testCase := range testCases {
		input := testCase.Input
		output := testCase.Output

		given := IsValidURL(input)
		if given != output {
			t.Errorf("IsValidURL(%v) = %v; want %v", input, given, output)
		}
	}
}
