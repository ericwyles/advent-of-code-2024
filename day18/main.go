package main

import (
	"bufio"
	"container/heap"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//go:embed input.txt
var embeddedFile string

const (
	BYTE      = '#'
	EMPTY     = '.'
	STEP_COST = 1
	TURN_COST = 0
)

// const (
// 	MEMORY_SIZE = 7
// 	BYTES       = 12
// )

const (
	MEMORY_SIZE = 71
	BYTES       = 1024
)

const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

var start = Coordinate{row: 0, col: 0}
var end = Coordinate{row: MEMORY_SIZE - 1, col: MEMORY_SIZE - 1}

var directions = []Coordinate{
	{-1, 0}, // NORTH
	{0, 1},  // EAST
	{1, 0},  // SOUTH
	{0, -1}, // WEST
}

type Coordinate struct {
	row, col int
}

type State struct {
	pos Coordinate
	dir int
}

type Item struct {
	state State
	cost  int
	index int // for priority queue
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func dijkstra(maze [][]rune, start Coordinate) int {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Initial state: start facing EAST
	startState := State{pos: start, dir: EAST}
	heap.Push(&pq, &Item{state: startState, cost: 0})

	// Visited states map (track lowest cost to each state)
	visited := make(map[State]int)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		state := current.state
		cost := current.cost

		// Stop if we've reached the END
		if state.pos == end {
			return cost
		}

		// Skip if already visited with lower cost
		if prevCost, exists := visited[state]; exists && cost >= prevCost {
			continue
		}
		visited[state] = cost

		// Explore moves: forward, turn clockwise, turn counterclockwise
		// 1. Move Forward
		nextPos := Coordinate{state.pos.row + directions[state.dir].row, state.pos.col + directions[state.dir].col}
		if isValidMove(maze, nextPos) {
			forwardState := State{pos: nextPos, dir: state.dir}
			heap.Push(&pq, &Item{state: forwardState, cost: cost + STEP_COST})
		}

		// 2. Turn Clockwise
		clockwiseState := State{pos: state.pos, dir: (state.dir + 1) % 4}
		heap.Push(&pq, &Item{state: clockwiseState, cost: cost + TURN_COST})

		// 3. Turn Counterclockwise
		counterClockwiseState := State{pos: state.pos, dir: (state.dir + 3) % 4} // (dir - 1 + 4) % 4
		heap.Push(&pq, &Item{state: counterClockwiseState, cost: cost + TURN_COST})
	}

	return -1 // No path found
}

func isValidMove(maze [][]rune, pos Coordinate) bool {
	return pos.row >= 0 && pos.row < len(maze) && pos.col >= 0 && pos.col < len(maze[0]) && maze[pos.row][pos.col] != BYTE
}

func main() {
	var maze [][]rune

	maze = make([][]rune, MEMORY_SIZE)
	resetMaze(maze)

	bytesToPlace := readInput()
	start := Coordinate{0, 0}

	//part1 - place all bytes up to limit and find minimum cost
	for i := 0; i < BYTES; i++ {
		placeByte(maze, bytesToPlace[i])
	}
	printGrid(maze, fmt.Sprintf("Initial maze after %d bytes have fallen", BYTES))
	minSteps := dijkstra(maze, start)
	fmt.Printf("After %d bytes, steps to reach the end: %d\n", BYTES, minSteps)

	// part2, find the block that makes it so there is no solution
	resetMaze(maze)
	for i, byte := range bytesToPlace {
		placeByte(maze, byte)
		if dijkstra(maze, start) == -1 {
			fmt.Printf("No path found. Byte [%d] - %d,%d\n", i+1, byte.col, byte.row)
			break
		}
	}
}

func placeByte(maze [][]rune, pos Coordinate) {
	maze[pos.row][pos.col] = BYTE
}

func resetMaze(maze [][]rune) {
	for i := range maze {
		maze[i] = make([]rune, MEMORY_SIZE)
		for j := range maze[i] {
			maze[i][j] = EMPTY
		}
	}
}

func readInput() []Coordinate {
	scanner := bufio.NewScanner(os.Stdin)

	var bytesToPlace []Coordinate

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > 0 {
			parts := strings.Split(line, ",")
			column, _ := strconv.Atoi(parts[0])
			row, _ := strconv.Atoi(parts[1])
			bytesToPlace = append(bytesToPlace, Coordinate{row: row, col: column})
		}
	}

	return bytesToPlace
}

func printGrid(grid [][]rune, header string) {
	fmt.Printf("%s\n", header)
	fmt.Printf("    ")
	for i := range grid[0] {
		if i >= 100 {
			fmt.Printf("%d", i/100%10)
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
	fmt.Printf("    ")
	for i := range grid[0] {
		if i >= 10 && i%10 == 0 {
			fmt.Printf("%d", i/10%10)
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
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
