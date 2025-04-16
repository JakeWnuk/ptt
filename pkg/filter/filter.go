// Package filter provides a function to check if a string passes a filter
package filter

import "github.com/jakewnuk/ptt/pkg/models"

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

	if len(input) < models.LengthStart || len(input) > models.LengthEnd {
		return false
	}

	return true
}
