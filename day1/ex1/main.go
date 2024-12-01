package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Open the input file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Initialize slices for the two columns
	var column1 []int
	var column2 []int

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line into fields
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("Skipping malformed line: %s\n", line)
			continue
		}

		// Convert fields to integers and append to slices
		num1, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("Error converting first column to int: %v\n", err)
			continue
		}

		num2, err := strconv.Atoi(fields[1])
		if err != nil {
			fmt.Printf("Error converting second column to int: %v\n", err)
			continue
		}

		column1 = append(column1, num1)
		column2 = append(column2, num2)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var column1len = len(column1)
	var column2len = len(column2)

	if column1len != column2len {
		fmt.Printf("Column lengths don't match %d vs %d", column1len, column2len)
		return
	}

	sort.Ints(column1)
	sort.Ints(column2)

	var totalerror int = 0
	for i := 0; i < column1len; i++ {
		totalerror += absInt(column1[i] - column2[i])
	}

	fmt.Printf("Total error is %d\n", totalerror)
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
