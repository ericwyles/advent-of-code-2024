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

type State struct {
	position  Coordinate
	direction rune
}

var phantomDistinctLocationsVisited = make(map[State]bool)
var distinctLocationsVisited = make(map[Coordinate]struct{})
var testedObstacleLocations = make(map[Coordinate]bool)
var numObstacles = 0

const OBSTACLE = '#'
const CLEAR = '.'

func main() {

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
		//fmt.Printf("Line %s\n", string(grid[row]))
		if !foundGuard {
			for column := range grid[row] {
				if _, exists := directionMap[grid[row][column]]; exists {
					guardDirection = grid[row][column]
					guardPosition = Coordinate{row: row, column: column}
				}
			}
		}
	}

	//fmt.Printf("Found guard '%c' at position %v\n", guardDirection, guardPosition)

	walkItOut(guardDirection, guardPosition, false)

	fmt.Printf("Number of distinct locations: %d\n", len(distinctLocationsVisited))
	fmt.Printf("Number of possible obstacle locations: %d\n", numObstacles)
}

func walkItOut(guardDirection rune, guardPosition Coordinate, isPhantomRealm bool) bool {
	currentState := State{position: guardPosition, direction: guardDirection}
	if isPhantomRealm {
		if phantomDistinctLocationsVisited[currentState] {
			return true // Found a loop
		}
		phantomDistinctLocationsVisited[currentState] = true
	} else {
		phantomDistinctLocationsVisited = make(map[State]bool)
		distinctLocationsVisited[guardPosition] = struct{}{}
	}

	nextPosition := addCoordinates(guardPosition, coordinateMap[guardDirection])

	if isOutOfBounds(nextPosition) {
		return false // found an exit
	}

	if OBSTACLE == grid[nextPosition.row][nextPosition.column] {
		// turn right but stay here, recursion takes care of it
		guardDirection = turnRight(guardDirection)
		return walkItOut(guardDirection, guardPosition, isPhantomRealm)

	} else if !isPhantomRealm && CLEAR == grid[nextPosition.row][nextPosition.column] && !testedObstacleLocations[nextPosition] {
		// if we haven't already working in the phantom realm,
		//    we'll put an OBSTACLE right in front of us and see if there is a loop
		testedObstacleLocations[nextPosition] = true

		grid[nextPosition.row][nextPosition.column] = OBSTACLE
		if walkItOut(guardDirection, guardPosition, true) {
			numObstacles++
		}
		grid[nextPosition.row][nextPosition.column] = CLEAR
	}

	return walkItOut(guardDirection, nextPosition, isPhantomRealm)
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

func containsRune(r rune, slice []rune) bool {
	for _, item := range slice {
		if item == r {
			return true
		}
	}
	return false
}
