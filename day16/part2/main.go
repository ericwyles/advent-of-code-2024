package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type direction struct {
	dy, dx int
}

var dirs = []direction{
	{-1, 0}, // NORTH
	{0, 1},  // EAST
	{1, 0},  // SOUTH
	{0, -1}, // WEST
}

type coordinate struct {
	y, x int
}

type position struct {
	y, x, cost, d int
	path          []coordinate
}

var (
	mapGrid             [][]byte
	height, width, best int
	ypos, xpos          int
	visited             [][]int
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		mapGrid = append(mapGrid, []byte(line))
		if pos := strings.IndexByte(line, 'S'); pos != -1 {
			ypos = len(mapGrid) - 1
			xpos = pos
		}
	}

	height = len(mapGrid)
	width = len(mapGrid[0])

	// Initialize visited array with a large number (equivalent to INT_MAX)
	visited = make([][]int, height)
	for i := range visited {
		visited[i] = make([]int, width)
		for j := range visited[i] {
			visited[i][j] = math.MaxInt
		}
	}
	visited[ypos][xpos] = 0

	best = math.MaxInt
	var newpos, pos []position
	var reachedEnd []position

	start := position{
		y:    ypos,
		x:    xpos,
		cost: 0,
		d:    1, // EAST
		path: []coordinate{{y: ypos, x: xpos}},
	}
	pos = append(pos, start)

	for len(pos) > 0 {
		for _, p := range pos {
			if visited[p.y][p.x] < p.cost-1000 {
				continue
			}
			visited[p.y][p.x] = p.cost

			if mapGrid[p.y][p.x] == 'E' {
				if p.cost < best {
					best = p.cost
				}
				reachedEnd = append(reachedEnd, p)
				continue
			}

			for _, delta := range []int{0, -1, 1} {
				nd := (p.d + delta + 4) % 4
				ny := p.y + dirs[nd].dy
				nx := p.x + dirs[nd].dx
				if ny >= 0 && ny < height && nx >= 0 && nx < width && mapGrid[ny][nx] != '#' {
					newPath := append([]coordinate(nil), p.path...)
					newPath = append(newPath, coordinate{y: ny, x: nx})
					newCost := p.cost
					if delta != 0 {
						newCost += 1001
					} else {
						newCost++
					}
					newpos = append(newpos, position{
						y:    ny,
						x:    nx,
						cost: newCost,
						d:    nd,
						path: newPath,
					})
				}
			}
		}
		pos = newpos
		newpos = nil
	}

	fmt.Printf("part 1: %d\n", best)

	seats := 0
	for _, p := range reachedEnd {
		if p.cost != best {
			continue
		}
		for _, c := range p.path {
			if mapGrid[c.y][c.x] != 'O' {
				seats++
			}
			mapGrid[c.y][c.x] = 'O'
		}
	}

	fmt.Printf("part 2: %d\n", seats)
}
