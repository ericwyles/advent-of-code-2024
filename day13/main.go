package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var embeddedFile string

type Coordinate struct {
	x int64
	y int64
}

type ClawMachine struct {
	a     Coordinate
	b     Coordinate
	prize Coordinate
}

const COSTA = 3
const COSTB = 1

func main() {
	var clawMachines []ClawMachine

	// Split the input into chunks representing each claw machine
	chunks := strings.Split(embeddedFile, "\n\n")

	for _, chunk := range chunks {
		if strings.TrimSpace(chunk) == "" {
			continue
		}

		// Split each chunk into lines
		lines := strings.Split(chunk, "\n")
		var clawMachine ClawMachine

		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "Button A:") {
				clawMachine.a = parseCoordinate(line, "Button A:")
			} else if strings.HasPrefix(line, "Button B:") {
				clawMachine.b = parseCoordinate(line, "Button B:")
			} else if strings.HasPrefix(line, "Prize:") {
				clawMachine.prize = parseCoordinate(line, "Prize:")
			}
		}

		// Append the populated claw machine to the slice
		clawMachines = append(clawMachines, clawMachine)
	}

	var totalCost int64 = 0
	var totalCostPart2 int64 = 0
	for _, machine := range clawMachines {
		totalCost += calculateCost(machine.a, machine.b, machine.prize, false)
		totalCostPart2 += calculateCost(machine.a, machine.b, machine.prize, true)
	}

	fmt.Printf("Total Cost: %d\n", totalCost)
	fmt.Printf("Total Cost: %d\n", totalCostPart2)
}

func calculateCost(buttonA, buttonB, prize Coordinate, part2 bool) int64 {
	// got here with help of various resources.
	// couldn't quite put my finger on the algebra :lolcry:
	if part2 {
		prize.x += 10000000000000
		prize.y += 10000000000000
	}

	// we are working with a system of two equations
	//
	// aPresses(a.x) + bPresses(b.x) = prize.x
	// aPresses(a.y) + bPresses(b.y) = prize.y
	//
	// solving for aPresses and bPresses

	buttonAX, buttonAY := buttonA.x, buttonA.y
	buttonBX, buttonBY := buttonB.x, buttonB.y
	prizeX, prizeY := prize.x, prize.y

	// multiply the second equation by b.x
	// aPresses(a.y)(b.x) + bPresses(b.y)(b.x) = prize.y(b.x)
	buttonAY *= buttonB.x
	buttonBY *= buttonB.x
	prizeY *= buttonB.x

	// multiply the first equation by b.y
	// aPresses(a.x)(b.y) + bPresses(b.x)(b.y) = prize.x(b.y)
	buttonAX *= buttonB.y
	buttonBX *= buttonB.y
	prizeX *= buttonB.y

	//Resulting two equations
	// aPresses(a.x)(b.y) + bPresses(b.x)(b.y) = prize.x(b.y)
	// aPresses(a.y)(b.x) + bPresses(b.y)(b.x) = prize.y(b.x)
	//
	// bPresses is the same in both, so we can (implicitly subtract it from both sides)
	//
	// aPresses(a.x)(b.y) = prize.x(b.y)
	// aPresses(a.y)(b.x) = prize.y(b.x)

	// Let:
	// - buttonADelta = abs((a.x)(b.y) - (a.y)(b.x))  (coefficient difference for aPresses)
	// - prizeDelta = abs(prize.x(b.y) - prize.y(b.x))        (difference in adjusted prizes)
	buttonADelta := abs(buttonAX - buttonAY)
	prizeDelta := abs(prizeX - prizeY)

	// Check if prizeDelta is divisible by buttonADelta.
	// If not, no integer solution exists for aPresses and bPresses.
	if !(prizeDelta%buttonADelta == 0) {
		return 0
	}

	// Solve for aPresses:
	// aPresses = prizeDelta / buttonADelta
	aPresses := prizeDelta / buttonADelta

	// Verify if bPresses can also be an integer by substituting aPresses back into one of the original equations.
	// Rearrange for bPresses:
	// bPresses = (prize.x - aPresses * a.x) / b.x
	if !((prizeX-aPresses*buttonAX)%buttonBX == 0) {
		return 0
	}
	bPresses := (prizeX - aPresses*buttonAX) / buttonBX

	if part2 || (aPresses <= 100 && bPresses <= 100) {
		return aPresses*COSTA + bPresses*COSTB
	}

	return 0
}

func abs(x int64) int64 {
	if x < 0 {
		return -1 * x
	}
	return x
}

func parseCoordinate(line, prefix string) Coordinate {
	line = strings.TrimSpace(strings.TrimPrefix(line, prefix))

	parts := strings.Split(line, ", ")
	var coord Coordinate

	for _, part := range parts {
		if strings.Contains(part, "X") {
			value := extractNumber(part)
			coord.x = value
		} else if strings.Contains(part, "Y") {
			value := extractNumber(part)
			coord.y = value
		}
	}
	return coord
}

func extractNumber(part string) int64 {
	// Remove non-numeric characters before parsing the number
	part = strings.Map(func(r rune) rune {
		if r == '-' || r == '+' || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, part)

	value, _ := strconv.Atoi(part)
	return int64(value)
}
