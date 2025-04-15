package filter

// Pass checks if a string passes the filter function and returns the boolean
// result.
//
// Args:
// input (string): The string to be checked.
//
// Returns:
// bool: true if the string passes the filter, false otherwise.
func Pass(input string) bool {
	if input == "" {
		return false
	}
	return true
}
