package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type DesignResult struct {
	possible bool
	total    int
}

var cache = map[string]DesignResult{}

func main() {
	towelPatterns, designs := parseInput()
	fmt.Printf("Towel Patterns: %s\n", towelPatterns)

	c := 0
	d := 0
	for _, design := range designs {
		fmt.Printf("Design: %s", design)
		designResult := checkIfPossible(design, towelPatterns)
		fmt.Printf(" - Possible: %d\n", designResult.total)
		if designResult.possible {
			d += 1
			c += designResult.total
		}
	}

	fmt.Printf("Designs Possible: %d\n", d)
	fmt.Printf("Total Combinations: %d\n", c)
}

func checkIfPossible(design string, towelPatterns []string) DesignResult {
	designResult := DesignResult{possible: false, total: 0}

	if value, exists := cache[design]; exists {
		return value
	}

	for _, pattern := range towelPatterns {
		if design == pattern {
			designResult.possible = true
			designResult.total = designResult.total + 1
		}

		if strings.HasPrefix(design, pattern) {
			result := checkIfPossible(design[len(pattern):], towelPatterns)
			if result.possible {
				designResult.possible = true
				designResult.total = designResult.total + result.total
			}
		}

	}

	cache[design] = designResult
	return designResult
}

func parseInput() ([]string, []string) {
	scanner := bufio.NewScanner(os.Stdin)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	towelPatterns := strings.Split(lines[0], ",")
	for i := range towelPatterns {
		towelPatterns[i] = strings.TrimSpace(towelPatterns[i])
	}

	return towelPatterns, lines[2:]
}
