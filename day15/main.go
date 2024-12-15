package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

//go:embed input.txt
var embeddedFile string
var originalGrid [][]rune
var scaledGrid [][]rune
var robotPosition Coordinate

type Coordinate struct {
	column int
	row    int
}

type WideBox struct {
	left  Coordinate
	right Coordinate
}

type CoordinatePair struct {
	a Coordinate
	b Coordinate
}

const ROBOT = '@'
const BOX = 'O'
const WALL = '#'
const EMPTY = '.'
const BOX_LEFT = '['
const BOX_RIGHT = ']'

const LEFT = "LEFT"
const RIGHT = "RIGHT"
const UP = "UP"
const DOWN = "DOWN"

var directionMap = map[rune]Coordinate{
	'^': {row: -1, column: 0}, // UP
	'v': {row: 1, column: 0},  // DOWN
	'>': {row: 0, column: 1},  // RIGHT
	'<': {row: 0, column: -1}, // LEFT
}

var directionWordMap = map[string]Coordinate{
	UP:    {row: -1, column: 0}, // UP
	DOWN:  {row: 1, column: 0},  // DOWN
	RIGHT: {row: 0, column: 1},  // RIGHT
	LEFT:  {row: 0, column: -1}, // LEFT
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
			originalGrid = append(originalGrid, []rune(line))
			scaledGrid = append(scaledGrid, []rune(scaleUp(line))) // for part 2
		}
	}

	movements := strings.TrimSpace(chunks[1])
	reg, _ := regexp.Compile("\\s+") //compile
	movements = reg.ReplaceAllString(movements, "")

	// part 1
	robotPosition = findRobot(originalGrid)
	printGrid(originalGrid, "Initial state:")

	for _, m := range movements {
		robotPosition, _ = move(originalGrid, robotPosition, directionMap[m])
	}

	printGrid(originalGrid, "After all moves:")
	gpsSum := 0
	for row, line := range originalGrid {
		for column, cell := range line {
			if cell == BOX {
				gpsSum += (row * 100) + column
			}
		}
	}
	fmt.Printf("Sum of GPS Coordinates: %d\n", gpsSum)

	if gpsSum != 10092 && gpsSum != 2028 && gpsSum != 1476771 && gpsSum != 908 && gpsSum != 1816 {
		panic("!!!!! Broke part 1 bruh !!!!!")
	}

	// part 2
	robotPosition = findRobot(scaledGrid)
	printGrid(scaledGrid, "Initial state:")
	fmt.Printf("Robot position: %v\n", robotPosition)
	for _, m := range movements {
		//fmt.Printf("Performing move %c:\n", m)
		robotPosition, _ = move(scaledGrid, robotPosition, directionMap[m])
	}

	printGrid(scaledGrid, "Final State")
	scaledGpsSum := 0
	for row, line := range scaledGrid {
		for column, cell := range line {
			if cell == BOX_LEFT {
				scaledGpsSum += (row * 100) + column
			}
		}
	}
	fmt.Printf("Sum of Scaled GPS Coordinates: %d\n", scaledGpsSum)

	if scaledGpsSum != 1468005 {
		panic("!!!!! Broke part 2 bruh !!!!!")
	}

}

func move(grid [][]rune, pos, direction Coordinate) (Coordinate, bool) {
	nextPosition := add(pos, direction)
	//fmt.Println(nextPosition)
	if isWall(grid, nextPosition) {
		return pos, false
	}

	if isEmpty(grid, nextPosition) {
		swap(grid, pos, nextPosition)
		return nextPosition, true
	}

	if isBox(grid, nextPosition) {
		if _, moved := move(grid, nextPosition, direction); moved {
			swap(grid, pos, nextPosition)
			return nextPosition, true
		}
	}

	// going left or right against wide boxes works like part 1
	directionWord := directionWords[direction]
	if directionWord == LEFT || directionWord == RIGHT {
		if isBoxLeft(grid, nextPosition) || isBoxRight(grid, nextPosition) {
			if _, moved := move(grid, nextPosition, direction); moved {
				swap(grid, pos, nextPosition)
				return nextPosition, true
			}
		}
	}

	// going up or down against wide boxes needs special movement
	if directionWord == UP || directionWord == DOWN {
		if isBoxLeft(grid, nextPosition) {
			swapTree := make(map[int][]CoordinatePair)
			box := WideBox{left: nextPosition, right: add(nextPosition, directionWordMap[RIGHT])}
			if wideMove(grid, box, direction, 1, swapTree) {

				executeSwapTree(grid, swapTree)
				swap(grid, pos, nextPosition)
				return nextPosition, true
			}
		}

		if isBoxRight(grid, nextPosition) {
			swapTree := make(map[int][]CoordinatePair)
			box := WideBox{left: add(nextPosition, directionWordMap[LEFT]), right: nextPosition}
			if wideMove(grid, box, direction, 1, swapTree) {
				fmt.Printf("Back and can swap!\n")

				executeSwapTree(grid, swapTree)
				swap(grid, pos, nextPosition)
				return nextPosition, true
			}
		}
	}

	return pos, false
}

func wideMove(grid [][]rune, box WideBox, direction Coordinate, depth int, swapTree map[int][]CoordinatePair) bool {

	nextPositions := WideBox{left: add(box.left, direction), right: add(box.right, direction)}

	// can't move walls
	if isWall(grid, nextPositions.left) || isWall(grid, nextPositions.right) {
		return false
	}

	// check if we can move one box straight ahead
	if isBoxLeft(grid, nextPositions.left) && isBoxRight(grid, nextPositions.right) {
		newBox := nextPositions
		if !wideMove(grid, newBox, direction, depth+1, swapTree) {
			return false
		}
	}

	// check if we can move one box offset to left
	if isBoxLeft(grid, nextPositions.right) && isEmpty(grid, nextPositions.left) {
		newBox := WideBox{left: nextPositions.right, right: add(nextPositions.right, directionWordMap[RIGHT])}

		if !wideMove(grid, newBox, direction, depth+1, swapTree) {
			return false
		}
	}

	// check if we can move one box offset to right
	if isEmpty(grid, nextPositions.right) && isBoxRight(grid, nextPositions.left) {
		newBox := WideBox{left: add(nextPositions.left, directionWordMap[LEFT]), right: nextPositions.left}

		if !wideMove(grid, newBox, direction, depth+1, swapTree) {
			return false
		}
	}

	// check if we can move two boxes (one offset left, one offset right)
	if isBoxLeft(grid, nextPositions.right) && isBoxRight(grid, nextPositions.left) {
		newBoxLeft := WideBox{left: add(nextPositions.left, directionWordMap[LEFT]), right: nextPositions.left}
		newBoxRight := WideBox{left: nextPositions.right, right: add(nextPositions.right, directionWordMap[RIGHT])}

		if !wideMove(grid, newBoxLeft, direction, depth+1, swapTree) || !wideMove(grid, newBoxRight, direction, depth+1, swapTree) {
			return false
		}
	}

	swapTree[depth] = append(swapTree[depth], CoordinatePair{a: box.left, b: nextPositions.left})
	swapTree[depth] = append(swapTree[depth], CoordinatePair{a: box.right, b: nextPositions.right})
	return true
}

func swap(grid [][]rune, a, b Coordinate) {
	grid[a.row][a.column], grid[b.row][b.column] = grid[b.row][b.column], grid[a.row][a.column]
}

func isWall(grid [][]rune, pos Coordinate) bool {
	return WALL == grid[pos.row][pos.column]
}

func isBox(grid [][]rune, pos Coordinate) bool {
	return BOX == grid[pos.row][pos.column]
}

func isBoxLeft(grid [][]rune, pos Coordinate) bool {
	return BOX_LEFT == grid[pos.row][pos.column]
}

func isBoxRight(grid [][]rune, pos Coordinate) bool {
	return BOX_RIGHT == grid[pos.row][pos.column]
}

func isEmpty(grid [][]rune, pos Coordinate) bool {
	return EMPTY == grid[pos.row][pos.column]
}

func add(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row + c2.row,
		column: c1.column + c2.column,
	}
}

func printGrid(grid [][]rune, header string) {
	fmt.Printf("%s\n", header)
	fmt.Printf("    ")
	for i := range grid[0] {
		fmt.Printf("%d", i%10)
	}
	fmt.Println()
	for l, line := range grid {
		fmt.Printf("%03d ", l)
		for _, cell := range line {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	fmt.Println()
}

func findRobot(grid [][]rune) Coordinate {
	for row, line := range grid {
		for column, cell := range line {
			if cell == ROBOT {
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

func executeSwapTree(grid [][]rune, swapTree map[int][]CoordinatePair) {
	depths := make([]int, 0, len(swapTree))
	for k := range swapTree {
		depths = append(depths, k)
	}

	sort.Slice(depths, func(i, j int) bool {
		return depths[i] > depths[j] // descending
	})

	fmt.Println("Depths to swap in descending order:")
	for _, d := range depths {
		swaps := deduplicate(swapTree[d])
		for _, s := range swaps {
			swap(grid, s.a, s.b)
		}
	}
}

func deduplicate(slice []CoordinatePair) []CoordinatePair {
	seen := make(map[CoordinatePair]bool) // Map to track seen pairs
	var result []CoordinatePair           // Deduplicated result slice

	for _, pair := range slice {
		if !seen[pair] {
			seen[pair] = true             // Mark as seen
			result = append(result, pair) // Add to result if not already seen
		}
	}

	return result
}
