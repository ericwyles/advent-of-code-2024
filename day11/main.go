package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var embeddedFile string
var inputStones []int

type StoneBlink struct {
	engraving       int
	blinksRemaining int
}

const BLINKS1 = 25
const BLINKS2 = 75

var rememberingStone = make(map[StoneBlink]int)

func main() {

	parts := strings.Split(embeddedFile, " ")
	for _, part := range parts {
		num, _ := strconv.Atoi(part)
		inputStones = append(inputStones, num)
	}

	part1 := 0
	part2 := 0
	for _, engraving := range inputStones {
		part1 += blink(engraving, BLINKS1)
		part2 += blink(engraving, BLINKS2)
	}
	fmt.Printf("Part 1: %d, Part 2: %d\n", part1, part2)
}

func blink(engraving int, remainingBlinks int) int {
	observation := StoneBlink{engraving: engraving, blinksRemaining: remainingBlinks}

	if remainingBlinks == 0 {
		return 1 // we're done, count self
	}

	memory, remembered := rememberingStone[observation]
	if remembered {
		return memory // remembering stone coming in clutch
	}

	remainingBlinks--
	stonesObserved := 0
	if engraving == 0 {
		stonesObserved = blink(1, remainingBlinks)
	} else {
		digits := numDigits(engraving)
		if digits%2 == 0 {
			left, right := splitNumber(engraving, digits)
			stonesObserved = blink(left, remainingBlinks) +
				blink(right, remainingBlinks)
		} else {
			stonesObserved = blink(engraving*2024, remainingBlinks)
		}
	}

	rememberingStone[observation] = stonesObserved // remember this
	return stonesObserved
}

func numDigits(n int) int {
	if n == 0 {
		return 1
	}

	count := 0
	for n != 0 {
		n /= 10
		count++
	}
	return count
}

func splitNumber(num int, digits int) (int, int) {
	half := digits / 2

	power := 1
	for i := 0; i < half; i++ {
		power *= 10
	}

	left := num / power
	right := num % power
	return left, right
}
