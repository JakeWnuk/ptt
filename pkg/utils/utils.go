package utils

import (
	"bufio"
	"fmt"
	"os"
	"ptt/pkg/models"
	"ptt/pkg/transform"
)

// ReadReturnStandardInput reads from standard input and returns the input back
// to standard output.
//
// Args:
// transformationArray (models.MultiString): an array of strings to be
// used for the transformation modes
//
// Returns:
// none
func ReadReturnStandardInput(transformation models.MultiString) {
	reader := bufio.NewScanner(os.Stdin)

	for reader.Scan() {
		line := reader.Text()

		for _, operation := range transformation {
			line = transform.Apply(line, operation)
		}
		fmt.Println(line)
	}

	if err := reader.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
