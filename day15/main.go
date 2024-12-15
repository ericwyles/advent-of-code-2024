package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed simple2.txt
var embeddedFile string
var grid [][]rune
var scaledGrid [][]rune
var guardPosition Coordinate

type Coordinate struct {
	column int
	row    int
}

const GUARD = '@'
const BOX = 'O'
const WALL = '#'
const EMPTY = '.'
const BOX_LEFT = '['
const BOX_RIGHT = ']'

var directionMap = map[rune]Coordinate{
	'^': {row: -1, column: 0}, // UP
	'v': {row: 1, column: 0},  // DOWN
	'>': {row: 0, column: 1},  // RIGHT
	'<': {row: 0, column: -1}, // LEFT
}

var directionWords = map[Coordinate]string{
	{row: -1, column: 0}: "UP",
	{row: 1, column: 0}:  "DOWN",
	{row: 0, column: 1}:  "RIGHT",
	{row: 0, column: -1}: "LEFT",
}

func main() {
	chunks := strings.Split(embeddedFile, "\n\n")

	lines := strings.Split(chunks[0], "\n")
	for _, line := range lines {
		if len(line) > 0 {
			grid = append(grid, []rune(line))
			scaledGrid = append(scaledGrid, []rune(scaleUp(line))) // for part 2
		}
	}

	movements := strings.TrimSpace(chunks[1])

	// part 1
	guardPosition = findGuard(grid)
	printGrid(grid, "Initial state:")

	for _, m := range movements {
		guardPosition, _ = move(guardPosition, directionMap[m])
	}

	printGrid(grid, "After all moves:")
	gpsSum := 0
	for row, line := range grid {
		for column, cell := range line {
			if cell == BOX {
				gpsSum += (row * 100) + column
			}
		}
	}
	fmt.Printf("Sum of GPS Coordinates: %d\n", gpsSum)

	if gpsSum != 10092 && gpsSum != 2028 && gpsSum != 1476771 && gpsSum != 908 {
		panic("!!!!! Broke part 1 bruh !!!!!")
	}

	// part 2
	guardPosition = findGuard(scaledGrid)
	printGrid(scaledGrid, "Initial state:")
	scaledGpsSum := 0
	for row, line := range scaledGrid {
		for column, cell := range line {
			if cell == BOX_LEFT {
				scaledGpsSum += (row * 100) + column
			}
		}
	}
	fmt.Printf("Guard position: %v\n", guardPosition)
	fmt.Printf("Sum of Scaled GPS Coordinates: %d\n", scaledGpsSum)

}

func move(pos, direction Coordinate) (Coordinate, bool) {
	nextPosition := add(pos, direction)
	if isWall(nextPosition) {
		return pos, false
	}

	if isEmpty(nextPosition) {
		swap(pos, nextPosition)
		return nextPosition, true
	}

	if isBox(nextPosition) {
		if _, moved := move(nextPosition, direction); moved {
			swap(pos, nextPosition)
			return nextPosition, true
		}
	}

	return pos, false
}

func swap(a, b Coordinate) {
	grid[a.row][a.column], grid[b.row][b.column] = grid[b.row][b.column], grid[a.row][a.column]
}

func isWall(pos Coordinate) bool {
	return WALL == grid[pos.row][pos.column]
}

func isBox(pos Coordinate) bool {
	return BOX == grid[pos.row][pos.column]
}

func isEmpty(pos Coordinate) bool {
	return EMPTY == grid[pos.row][pos.column]
}

func add(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row + c2.row,
		column: c1.column + c2.column,
	}
}

func printGrid(grid [][]rune, header string) {
	fmt.Println(header)
	for _, line := range grid {
		for _, cell := range line {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	fmt.Println()
}

func findGuard(grid [][]rune) Coordinate {
	for row, line := range grid {
		for column, cell := range line {
			if cell == GUARD {
				return Coordinate{row: row, column: column}
			}
		}
	}

	return Coordinate{row: -1, column: -1}
}

func scaleUp(line string) string {
	newLine := line

	newLine = strings.ReplaceAll(newLine, "#", "##")
	newLine = strings.ReplaceAll(newLine, "O", "[]")
	newLine = strings.ReplaceAll(newLine, ".", "..")
	newLine = strings.ReplaceAll(newLine, "@", "@.")

	return newLine
}
