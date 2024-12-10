package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var embeddedFile string
var grid [][]int

type Coordinate struct {
	row    int
	column int
}

const SUMMIT = 9

var directions = []Coordinate{
	{row: -1, column: 0}, // UP
	{row: 1, column: 0},  // DOWN
	{row: 0, column: 1},  // RIGHT
	{row: 0, column: -1}, // LEFT
}

func main() {
	var trailheads []Coordinate
	totalScore := 0
	totalRating := 0

	lines := strings.Split(embeddedFile, "\n")
	for i, line := range lines {
		if len(line) > 0 {
			var gridline []int
			for j, char := range line {
				// Convert rune (character) to its integer representation
				num, _ := strconv.Atoi(string(char))
				gridline = append(gridline, num)

				if num == 0 {
					coordinate := Coordinate{row: i, column: j}
					trailheads = append(trailheads, coordinate)
				}
			}
			grid = append(grid, gridline)
		}
	}

	for _, trailhead := range trailheads {
		uniqueSummitLocations := make(map[Coordinate]struct{})
		rating := 0

		exploreTrail(trailhead, uniqueSummitLocations, &rating)

		totalScore += len(uniqueSummitLocations)
		totalRating += rating
	}

	fmt.Printf("Total score: %d, Total rating: %d\n", totalScore, totalRating)
}

func exploreTrail(location Coordinate, uniqueSummitLocations map[Coordinate]struct{}, rating *int) {
	currentHeight := height(location)
	if currentHeight == SUMMIT {
		uniqueSummitLocations[location] = struct{}{}
		*rating++
		return
	}

	for _, dir := range directions {
		nextLocation := move(location, dir)
		if !isOutOfBounds(nextLocation) {
			nextHeight := height(nextLocation)
			if nextHeight == currentHeight+1 {
				exploreTrail(nextLocation, uniqueSummitLocations, rating)
			}
		}
	}

	return
}

func move(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row + c2.row,
		column: c1.column + c2.column,
	}
}

func height(c Coordinate) int {
	return grid[c.row][c.column]
}

func isOutOfBounds(c Coordinate) bool {
	return c.row < 0 || c.column < 0 || c.row > len(grid)-1 || c.column > len(grid[0])-1
}
