package utils

// ReverseString will return a string in reverse
//
// Args:
// str (string): Input string to transform
//
// Returns:
// (string): Transformed string
func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
