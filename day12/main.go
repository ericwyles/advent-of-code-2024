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

var squareDirections = []Coordinate{
	{row: -1, column: 0}, // UP
	{row: 1, column: 0},  // DOWN
	{row: 0, column: 1},  // RIGHT
	{row: 0, column: -1}, // LEFT
}

var diagonalDirections = []Coordinate{
	{-1, -1}, // up-left
	{-1, 1},  // up-right
	{1, -1},  // down-left
	{1, 1},   // down-right
}

var coordinatesCounted = make(map[Coordinate]struct{})

func main() {
	price := 0
	discountedPrice := 0

	lines := strings.Split(embeddedFile, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			grid = append(grid, []rune(line))
		}
	}

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			loc := Coordinate{row: i, column: j}
			_, ok := coordinatesCounted[loc]
			if !ok {
				area, perimeter, sides := getRegionSize(loc)
				price += (area * perimeter)
				discountedPrice += (area * sides)
			}
		}
	}

	fmt.Printf("Total Price: %d\n", price)
	fmt.Printf("Discounted Price: %d\n", discountedPrice)
}

func getRegionSize(loc Coordinate) (int, int, int) {
	_, alreadyCounted := coordinatesCounted[loc]
	if alreadyCounted {
		return 0, 0, 0 // already counted the fence for this one
	}

	coordinatesCounted[loc] = struct{}{}
	plantType := getPlantType(loc)

	area := 1
	perimeter := 4
	neighborsArea := 0
	neighborsPerimeter := 0
	neighborsSides := 0

	for _, dir := range squareDirections {
		nextLocation := move(loc, dir)
		if !isOutOfBounds(nextLocation) {
			nextPlantType := getPlantType(nextLocation)
			if nextPlantType == plantType {
				perimeter--
				nextArea, nextPerimeter, nextSides := getRegionSize(nextLocation)
				neighborsArea += nextArea
				neighborsPerimeter += nextPerimeter
				neighborsSides += nextSides
			}
		}
	}

	return area + neighborsArea, perimeter + neighborsPerimeter, findCorners(loc) + neighborsSides
}

func matches(c1, c2 Coordinate) bool {
	oob1 := isOutOfBounds(c1)
	oob2 := isOutOfBounds(c2)
	if !oob1 && !oob2 {
		return getPlantType(c1) == getPlantType(c2)
	}

	return oob1 == oob2
}

func move(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row + c2.row,
		column: c1.column + c2.column,
	}
}

func findCorners(loc Coordinate) int {
	corners := 0

	for _, diag := range diagonalDirections {
		diagLocation := move(loc, diag)

		adjHorizontal := move(loc, Coordinate{diag.row, 0})  // Horizontal neighbor
		adjVertical := move(loc, Coordinate{0, diag.column}) // Vertical neighbor

		diagMatchesSelf := matches(loc, diagLocation)
		horizontalMatchesSelf := matches(loc, adjHorizontal)
		verticalMatchesSelf := matches(loc, adjVertical)

		if !diagMatchesSelf && !horizontalMatchesSelf && !verticalMatchesSelf {
			corners++
		} else if diagMatchesSelf && !horizontalMatchesSelf && !verticalMatchesSelf {
			corners++
		} else if !diagMatchesSelf && horizontalMatchesSelf && verticalMatchesSelf {
			corners++
		}
	}

	return corners
}

func isOutOfBounds(c Coordinate) bool {
	return c.row < 0 || c.column < 0 || c.row > len(grid)-1 || c.column > len(grid[0])-1
}

func getPlantType(c Coordinate) rune {
	return grid[c.row][c.column]
}
