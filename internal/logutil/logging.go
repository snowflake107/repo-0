package logutil

import (
	"fmt"
	"os"
)

// WriteTextToFile writes the given text to the specified file.
func WriteTextToFile(filePath, text string) error {
	// Open the file for writing, create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Write the text to the file
	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
