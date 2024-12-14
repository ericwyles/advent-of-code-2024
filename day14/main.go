package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var embeddedFile string

type Coordinate struct {
	x int
	y int
}

type Robot struct {
	p Coordinate
	v Coordinate
}

const WIDTH = 101
const HEIGHT = 103

func main() {
	var robots []Robot

	lines := strings.Split(strings.TrimSpace(embeddedFile), "\n") // Split input into lines
	for _, line := range lines {
		var px, py, vx, vy int

		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &px, &py, &vx, &vy)

		robot := Robot{
			p: Coordinate{x: px, y: py},
			v: Coordinate{x: vx, y: vy},
		}
		robots = append(robots, robot)
	}

	// part 1
	quadrantMap := make(map[int]int)
	for _, r := range robots {
		newpos := move(r, 100)
		q := positionToQuadrant(newpos)
		if _, exists := quadrantMap[q]; exists {
			quadrantMap[q]++
		} else {
			quadrantMap[q] = 1
		}
	}

	safetyFactor := quadrantMap[1] * quadrantMap[2] * quadrantMap[3] * quadrantMap[4]
	fmt.Printf("safetyFactor %d\n", safetyFactor)

	// part 2
	for seconds := range 10000 {
		coordinateMap := make(map[Coordinate]struct{})
		for _, r := range robots {
			newpos := move(r, seconds)
			coordinateMap[newpos] = struct{}{}
		}
		if printRobots(coordinateMap, seconds) {
			break
		}
	}
}

func positionToQuadrant(p Coordinate) int {
	hBound := WIDTH / 2
	vBound := HEIGHT / 2

	if p.x < hBound && p.y < vBound {
		return 1
	}

	if p.x > hBound && p.y > vBound {
		return 4
	}

	if p.x > hBound && p.y < vBound {
		return 2
	}

	if p.x < hBound && p.y > vBound {
		return 3
	}

	return 0
}

func move(r Robot, seconds int) Coordinate {
	delta := multiply(r.v, seconds)
	newpos := add(r.p, delta)
	t := teleport(newpos)
	return t
}

func add(p, a Coordinate) Coordinate {
	return Coordinate{x: p.x + a.x, y: p.y + a.y}
}

func multiply(p Coordinate, times int) Coordinate {
	return Coordinate{x: p.x * times, y: p.y * times}
}

func teleport(p Coordinate) Coordinate {
	t := Coordinate{x: p.x % WIDTH, y: p.y % HEIGHT}
	if t.y < 0 {
		t.y = HEIGHT + t.y
	}
	if t.x < 0 {
		t.x = WIDTH + t.x
	}

	return t
}

func printRobots(coordinateMap map[Coordinate]struct{}, seconds int) bool {
	fullMap := ""

	for y := range WIDTH {
		line := ""
		for x := range HEIGHT {
			r := ' '
			if _, exists := coordinateMap[Coordinate{y: y, x: x}]; exists {
				r = '*'
			}

			line += fmt.Sprintf("%c", r)
		}

		fullMap += line + "\n"
	}

	if strings.Contains(fullMap, "*********") { // took a guess here that I could just look for a group of consecutive * in the output
		fmt.Printf(fullMap)
		fmt.Printf("\n^^^AFTER %04d SECONDS^^^\n", seconds)
		return true
	}

	return false
}
