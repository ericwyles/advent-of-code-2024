package main

import (
	"container/heap"
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var embeddedFile string

const (
	START     = 'S'
	END       = 'E'
	WALL      = '#'
	EMPTY     = '.'
	STEP_COST = 1
	TURN_COST = 1000
)

const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

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
		if maze[state.pos.row][state.pos.col] == END {
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
	return pos.row >= 0 && pos.row < len(maze) && pos.col >= 0 && pos.col < len(maze[0]) && maze[pos.row][pos.col] != WALL
}

func findStart(maze [][]rune) Coordinate {
	for r, row := range maze {
		for c, cell := range row {
			if cell == START {
				return Coordinate{row: r, col: c}
			}
		}
	}
	return Coordinate{-1, -1}
}

func main() {
	lines := strings.Split(embeddedFile, "\n")
	var maze [][]rune
	for _, line := range lines {
		if len(line) > 0 {
			maze = append(maze, []rune(line))
		}
	}

	start := findStart(maze)
	minCost := dijkstra(maze, start)

	fmt.Printf("Minimum cost to reach the end: %d\n", minCost)
}
