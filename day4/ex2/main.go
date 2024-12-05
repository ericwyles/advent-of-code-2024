package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var grid [][]rune

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		grid = append(grid, line)
	}

	totalFound := 0
	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[y])-1; x++ {
			if grid[y][x] == 'A' { // keying off the middle of the pattern and will look at neighbors from here
				// Check for matching patterns
				if isMas(grid[y+1][x-1], grid[y][x], grid[y-1][x+1]) &&
					isMas(grid[y-1][x-1], grid[y][x], grid[y+1][x+1]) {
					totalFound++
				}
			}
		}
	}

	fmt.Printf("Total found: %d\n", totalFound)
}

func isMas(a, b, c rune) bool {
	return (a == 'M' && b == 'A' && c == 'S') || (a == 'S' && b == 'A' && c == 'M')
}
