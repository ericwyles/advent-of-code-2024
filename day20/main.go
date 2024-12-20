package main

import (
	"bufio"
	"container/heap"
	_ "embed"
	"fmt"
	"os"
)

const (
	START     = 'S'
	END       = 'E'
	WALL      = '#'
	EMPTY     = '.'
	STEP_COST = 1
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

type Cheat struct {
	start Coordinate
	end   Coordinate
}

type Coordinate struct {
	row, col int
}

type State struct {
	pos    Coordinate
	cheats []Cheat
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

func dijkstra(maze [][]rune, start Coordinate) (int, []Cheat) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Initial state: start facing EAST
	startState := State{pos: start}
	heap.Push(&pq, &Item{state: startState, cost: 0})

	// Visited position map (track lowest cost to each position)
	visited := make(map[Coordinate]int)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		state := current.state
		cost := current.cost
		cheats := state.cheats

		// Stop if we've reached the END
		if maze[state.pos.row][state.pos.col] == END {
			return cost, cheats
		}

		// Skip if already visited with lower cost
		if prevCost, exists := visited[state.pos]; exists && cost >= prevCost {
			continue
		}
		visited[state.pos] = cost

		for dir := range 4 {
			cheatStart := getNextPosition(state.pos, dir)
			cheatEnd := getNextPosition(cheatStart, dir)
			if _, exists := visited[cheatEnd]; !exists {
				if isWall(maze, cheatStart) && isValidMove(maze, cheatEnd) {
					cheats = append(cheats, Cheat{start: cheatStart, end: cheatEnd})
				}
			}
		}

		// Explore moves: up, down, left, right
		for dir := range 4 {
			nextPos := getNextPosition(state.pos, dir)
			if isValidMove(maze, nextPos) {
				nextState := State{pos: nextPos, cheats: cheats}
				heap.Push(&pq, &Item{state: nextState, cost: cost + STEP_COST})
			}
		}
	}

	return -1, nil
}

func replace(track [][]rune, pos Coordinate, new rune) {
	track[pos.row][pos.col] = new
}

func getNextPosition(pos Coordinate, dir int) Coordinate {
	return Coordinate{pos.row + directions[dir].row, pos.col + directions[dir].col}
}

func isValidMove(maze [][]rune, pos Coordinate) bool {
	return pos.row >= 0 && pos.row < len(maze) && pos.col >= 0 && pos.col < len(maze[0]) && maze[pos.row][pos.col] != WALL
}

func isWall(maze [][]rune, pos Coordinate) bool {
	return pos.row >= 0 && pos.row < len(maze) && pos.col >= 0 && pos.col < len(maze[0]) && maze[pos.row][pos.col] == WALL
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
	track := readInput()

	printGrid(track, "Initial Track")

	start := findStart(track)
	time, cheats := dijkstra(track, start)

	fmt.Printf("Picoseconds to reach the end: %d\n", time)
	fmt.Printf("Potential Cheats: %d\n", len(cheats))

	numFastCheats := 0
	for _, cheat := range cheats {
		replace(track, cheat.start, EMPTY)
		cheatTime, _ := dijkstra(track, start)
		replace(track, cheat.start, WALL)

		timeSaving := time - cheatTime
		//fmt.Printf("Cheat [%03d] at [%v] saves %d picoseconds\n", i, cheat.start, timeSaving)
		if timeSaving >= 100 {
			numFastCheats++
		}
	}

	fmt.Printf("Number of fast cheats found: %d\n", numFastCheats)
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

func readInput() [][]rune {
	scanner := bufio.NewScanner(os.Stdin)

	var track [][]rune

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			track = append(track, []rune(line))
		}
	}

	return track
}
