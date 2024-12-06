package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var embeddedFile string
var grid [][]rune

type Coordinate struct {
	row    int
	column int
}

var directions = []rune{'^', '>', 'V', '<'}

var directionMap = map[rune]int{
	'^': 0,
	'>': 1,
	'V': 2,
	'<': 3,
}

var coordinateMap = map[rune]Coordinate{
	'^': {-1, 0},
	'>': {0, 1},
	'V': {1, 0},
	'<': {0, -1},
}

func main() {
	distinctLocationsVisited := make(map[Coordinate]struct{})

	lines := strings.Split(embeddedFile, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			grid = append(grid, []rune(line))
		}
	}

	// find the guard and which direction they are facing to get started
	foundGuard := false
	var guardPosition Coordinate
	var guardDirection rune
	for row := range grid {
		fmt.Printf("Line %s\n", string(grid[row]))
		if !foundGuard {
			for column := range grid[row] {
				if _, exists := directionMap[grid[row][column]]; exists {
					guardDirection = grid[row][column]
					guardPosition = Coordinate{row: row, column: column}
				}
			}
		}
	}

	fmt.Printf("Found guard '%c' at position %v\n", guardDirection, guardPosition)

	walkItOut(guardDirection, guardPosition, distinctLocationsVisited)

	fmt.Printf("Distinct Locations Visited: %v\n\n", distinctLocationsVisited)
	fmt.Printf("Number of distinct locations: %d\n", len(distinctLocationsVisited))
}

func walkItOut(guardDirection rune, guardPosition Coordinate, distinctLocationsVisited map[Coordinate]struct{}) {
	distinctLocationsVisited[guardPosition] = struct{}{}

	nextPosition := addCoordinates(guardPosition, coordinateMap[guardDirection])

	if isOutOfBounds(nextPosition) {
		return
	}

	if '#' == grid[nextPosition.row][nextPosition.column] {
		// turn right but stay here, recursion takes care of it
		guardDirection = turnRight(guardDirection)
		nextPosition = guardPosition
	}

	walkItOut(guardDirection, nextPosition, distinctLocationsVisited)
}

func addCoordinates(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row + c2.row,
		column: c1.column + c2.column,
	}
}

func isOutOfBounds(c Coordinate) bool {
	return c.row < 0 || c.column < 0 || c.row > len(grid)-1 || c.column > len(grid[0])-1
}

func turnRight(guardDirection rune) rune {
	return directions[(directionMap[guardDirection]+1)%4]
}
