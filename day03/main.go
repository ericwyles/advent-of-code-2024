package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	filePath := "input.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close() // Ensure the file is closed when done

	// Read the entire file into a string
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	input := string(data)
	total := 0

	// remove disabled instructions
	re := regexp.MustCompile(`don't\(\)[\s\S]*?do\(\)`)
	input = re.ReplaceAllString(input, "DISABLED")

	// process the instructions
	indexes := findAllStartIndexes(input, "mul(")

	fmt.Println("Indexes of mul: ", indexes)
	for i := range indexes {
		mul := indexes[i]
		fmt.Printf("processing mul location %d\n", mul)
		closingParen := mul + strings.Index(input[mul:], ")")
		fmt.Printf("    closing parent location is %d\n", closingParen)
		if closingParen != -1 {
			// found a closing paren
			candidateInstruction := input[mul : closingParen+1]
			args := candidateInstruction[4 : len(candidateInstruction)-1]
			fmt.Printf("    found candidate instruction %s, args %s\n", candidateInstruction, args)
			first, second, err := validateAndSplit(args)
			if err != nil {
				fmt.Println("Skipping args because error:", err)
			} else {
				fmt.Printf("First integer: %d, Second integer: %d\n", first, second)
				total += first * second
			}
		}
	}

	fmt.Printf("Total [%d].\n", total)

}

func findAllStartIndexes(s, substr string) []int {
	var indexes []int
	start := 0

	for {
		// Search for the substring starting from the current index
		idx := strings.Index(s[start:], substr)
		if idx == -1 {
			// No more occurrences found
			break
		}
		// Calculate the actual index in the original string
		actualIndex := start + idx
		indexes = append(indexes, actualIndex)
		// Move the start index past the current match
		start = actualIndex + len(substr)
	}

	return indexes
}

func validateAndSplit(args string) (int, int, error) {
	// Split the string on commas
	parts := strings.Split(args, ",")

	// Check if we have exactly two fields
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("input must contain exactly two fields separated by a comma")
	}

	// Parse and validate each field
	first, err1 := strconv.Atoi(parts[0])
	second, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("both fields must be valid integers")
	}

	// Ensure no extra whitespace around the original fields
	if parts[0] != strconv.Itoa(first) || parts[1] != strconv.Itoa(second) {
		return 0, 0, fmt.Errorf("fields must be integers with no whitespace or extra characters")
	}

	return first, second, nil
}
