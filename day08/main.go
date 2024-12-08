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

func main() {
	uniqueAntinodeLocations := make(map[Coordinate]struct{})
	antennaMap := make(map[rune][]Coordinate)

	lines := strings.Split(embeddedFile, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			grid = append(grid, []rune(line))
		}
	}

	// find all the antennas mapped by frequency
	for row := range grid {
		for column := range grid[row] {
			if '.' != grid[row][column] {
				frequency := grid[row][column]
				coordinate := Coordinate{row: row, column: column}
				antennaMap[frequency] = append(antennaMap[frequency], coordinate)
			}
		}
	}

	for frequency := range antennaMap {
		antennas := antennaMap[frequency]
		for i, antennaLocationA := range antennas {
			for j, antennaLocationB := range antennas {
				// the problem was ambiguous and i don't fully understand why
				//    every antenna location is now an antinode but it seems
				//    like it is from reading the examples? Fictional world physics > me
				uniqueAntinodeLocations[antennaLocationA] = struct{}{}

				if i != j { // self + self is not a pair
					slope := subtract(antennaLocationA, antennaLocationB)
					candidateLocation := add(antennaLocationA, slope)
					recordAntinodes(candidateLocation, slope, uniqueAntinodeLocations)
				}
			}
		}
	}

	fmt.Printf("Number of antinodes: %d\n", len(uniqueAntinodeLocations))
}

func recordAntinodes(candidateLocation, slope Coordinate, uniqueAntinodeLocations map[Coordinate]struct{}) {
	if isOutOfBounds(candidateLocation) {
		return
	}
	uniqueAntinodeLocations[candidateLocation] = struct{}{}
	recordAntinodes(add(candidateLocation, slope), slope, uniqueAntinodeLocations)
}

func add(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row + c2.row,
		column: c1.column + c2.column,
	}
}

func subtract(c1, c2 Coordinate) Coordinate {
	return Coordinate{
		row:    c1.row - c2.row,
		column: c1.column - c2.column,
	}
}

func isOutOfBounds(c Coordinate) bool {
	return c.row < 0 || c.column < 0 || c.row > len(grid)-1 || c.column > len(grid[0])-1
}
