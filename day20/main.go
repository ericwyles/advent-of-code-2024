package main

import (
	"bufio"
	"container/heap"
	_ "embed"
	"fmt"
	"os"
	"sort"
)

const (
	START     = 'S'
	END       = 'E'
	WALL      = '#'
	EMPTY     = '.'
	STEP_COST = 1
	// // PART 1:
	// MAX_CHEAT_DISTANCE = 2
	// MIN_CHEAT_SAVINGS  = 100
	// PART 2:
	MAX_CHEAT_DISTANCE = 20
	MIN_CHEAT_SAVINGS  = 100
)

var directions = []Coordinate{
	{-1, 0}, // UP
	{0, 1},  // RIGHT
	{1, 0},  // DOWN
	{0, -1}, // LEFT
}

type Cheat struct {
	start Coordinate
	end   Coordinate
}

type Coordinate struct {
	row, col int
}

type State struct {
	pos Coordinate
	//cheats []Cheat
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

func dijkstra(maze [][]rune, start Coordinate) (int, map[int][]Cheat) {
	costMap := make(map[int]State)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Initial state: start facing EAST
	startState := State{pos: start}
	heap.Push(&pq, &Item{state: startState, cost: 0})
	costMap[0] = startState
	cheatMap := make(map[int][]Cheat)

	// Visited position map (track lowest cost to each position)
	visited := make(map[Coordinate]int)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		state := current.state
		cost := current.cost

		// Skip if already visited with lower cost
		if prevCost, exists := visited[state.pos]; exists && cost >= prevCost {
			continue
		}
		visited[state.pos] = cost

		// find shortcuts that could lead to here
		for i := range cost - MIN_CHEAT_SAVINGS {
			if prevState, ok := costMap[i]; ok {
				taxiDistance := getTaxiDistance(state.pos, prevState.pos)
				if taxiDistance <= MAX_CHEAT_DISTANCE {
					distanceSaved := cost - i - taxiDistance
					if distanceSaved >= MIN_CHEAT_SAVINGS {
						newCheat := Cheat{start: prevState.pos, end: state.pos}
						cheatMap[distanceSaved] = append(cheatMap[distanceSaved], newCheat)
					}
				}
			}
		}

		// Stop if we've reached the END
		if maze[state.pos.row][state.pos.col] == END {
			return cost, cheatMap
		}

		// Explore moves: up, down, left, right
		for dir := range 4 {
			nextPos := getNextPosition(state.pos, dir)
			if isValidMove(maze, nextPos) {
				if _, exists := visited[nextPos]; !exists {
					nextState := State{pos: nextPos}
					nextCost := cost + STEP_COST
					heap.Push(&pq, &Item{state: nextState, cost: nextCost})
					costMap[nextCost] = nextState
				}
			}
		}
	}

	return -1, cheatMap
}

func getTaxiDistance(a, b Coordinate) int {
	return absInt(a.row-b.row) + absInt(a.col-b.col)
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
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

	fmt.Printf("Cheats found:\n")
	var keys []int
	for k := range cheats {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	totalCheats := 0
	for _, k := range keys {
		totalCheats += len(cheats[k])
		fmt.Printf("    There are %d cheats that save %d picoseconds.\n", len(cheats[k]), k)
	}
	fmt.Printf("Total Cheats: %d\n", totalCheats)
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
