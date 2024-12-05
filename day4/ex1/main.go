package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var horizontalAndVertical []string
	var diagonal []string
	for lineNum := range len(grid) {
		fmt.Printf("line is: %s\n", string(grid[lineNum]))

		// horizontal string from this pos
		horizontalAndVertical = append(horizontalAndVertical, string(grid[lineNum]))

	}

	// add the vertical strings
	for columnPos := range len(grid[0]) {
		verticalString := findStringFromPos(grid, 0, columnPos, +1, 0)

		// for rowPos := range len(grid) {
		// 	verticalString += string(grid[rowPos][columnPos])
		// }
		horizontalAndVertical = append(horizontalAndVertical, verticalString)
	}

	totalFound := 0
	for i := range horizontalAndVertical {
		search := horizontalAndVertical[i]
		fmt.Printf("searchable string is: %s\n", search)

		totalFound += findNumStrings(search, "XMAS")
		totalFound += findNumStrings(search, "SAMX")
	}

	for x := len(grid[0]) - 1; x >= 0; x-- {
		// get diagonals slopping down and to right (y+1, x+1)
		diagonal = append(diagonal, findStringFromPos(grid, 0, x, +1, +1))
	}
	for y := 1; y < len(grid); y++ {
		// get diagonals slopping down and to right (y+1, x+1)
		diagonal = append(diagonal, findStringFromPos(grid, y, 0, +1, +1))
	}
	for x := len(grid[0]) - 1; x >= 0; x-- {
		// get diagonals slopping up and to right (y-1, x+1)
		diagonal = append(diagonal, findStringFromPos(grid, len(grid)-1, x, -1, +1))
	}
	for y := len(grid[0]) - 2; y >= 0; y-- {
		// get diagonals slopping up and to right (y-1, x+1)
		diagonal = append(diagonal, findStringFromPos(grid, y, 0, -1, +1))
	}

	for i := range diagonal {
		search := diagonal[i]
		fmt.Printf("diagonal string is: %s\n", search)

		totalFound += findNumStrings(search, "XMAS")
		totalFound += findNumStrings(search, "SAMX")
	}

	fmt.Printf("Total found: %d\n", totalFound)
}

func findNumStrings(text, substring string) int {
	count := 0
	for i := 0; ; {
		index := strings.Index(text[i:], substring)
		if index == -1 {
			break
		}
		count++
		i += index + 1
	}
	return count
}

func findStringFromPos(grid [][]rune, y int, x int, yDirection int, xDirection int) string {
	fmt.Printf("y=%d, len(grid)-1=%d, x=%d, len(grid[0])-1=%d\n", y, len(grid)-1, x, len(grid[0])-1)
	if y > len(grid)-1 || x > len(grid[0])-1 ||
		y < 0 || x < 0 {
		// at this point we are out of bounds
		return ""
	}

	return string(grid[y][x]) + findStringFromPos(grid, y+yDirection, x+xDirection, yDirection, xDirection)
}
