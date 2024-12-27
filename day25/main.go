package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"
)

func main() {
	locks, keys := readInput()

	fits := 0
	for _, lock := range locks {
		for _, key := range keys {
			fit := true
			for i := 0; i < 5; i++ {
				if lock[i]+key[i] > 5 {
					fit = false
					break
				}
			}
			if fit {
				fits++
			}
		}
	}
	fmt.Printf("Fits %d\n", fits)
}

func readInput() ([][]int, [][]int) {
	var locks [][]int
	var keys [][]int

	scanner := bufio.NewScanner(os.Stdin)

	var inputBuilder strings.Builder

	for scanner.Scan() {
		inputBuilder.WriteString(scanner.Text())
		inputBuilder.WriteString("\n")
	}

	input := inputBuilder.String()

	chunks := strings.Split(input, "\n\n")
	for _, chunk := range chunks {
		lines := strings.Split(chunk, "\n")
		pins := countColumns(lines[1:6])
		if lines[0] == "#####" {
			locks = append(locks, pins)
		} else {
			keys = append(keys, pins)
		}
	}

	return locks, keys
}

func countColumns(pins []string) []int {
	counts := make([]int, 5)
	for _, line := range pins {
		for i, rune := range line {
			if rune == '#' {
				counts[i] += 1
			}
		}
	}
	return counts
}
