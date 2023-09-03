package lib

import (
	"bufio"
	"os"
	"strings"
)

func ReadAndParseFile(filePath string, delimiter string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result [][]string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		// Read the line
		line := scanner.Text()

		// Split the line into an array using the specified delimiter
		parts := strings.Split(line, delimiter)

		// Append the parsed array to the result
		result = append(result, parts)
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
