package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var embeddedFile string

func main() {
	calibrationResult := 0

	lines := strings.Split(embeddedFile, "\n")
	lines = lines[:len(lines)-1]
	for _, line := range lines {
		testValue, operands := parseLine(line)
		if canProduceTestValue(testValue, operands) {
			calibrationResult += testValue
		}
	}

	fmt.Printf("Calibration Result: %d\n", calibrationResult)
}

func canProduceTestValue(testValue int, operands []int) bool {
	if operands[0] > testValue {
		return false
	}

	if len(operands) == 1 {
		return operands[0] == testValue
	}

	return canProduceTestValue(testValue, append([]int{operands[0] + operands[1]}, operands[2:]...)) ||
		canProduceTestValue(testValue, append([]int{operands[0] * operands[1]}, operands[2:]...)) ||
		canProduceTestValue(testValue, append([]int{concatInts(operands[0], operands[1])}, operands[2:]...))
}

func concatInts(i1 int, i2 int) int {
	result, _ := strconv.Atoi(strconv.Itoa(i1) + strconv.Itoa(i2))
	return result
}

func parseLine(line string) (int, []int) {
	parts := strings.SplitN(line, ":", 2)
	testValue, _ := strconv.Atoi(strings.TrimSpace(parts[0]))

	// Split the numbers after ':' into a slice of integers
	operandStrings := strings.Fields(parts[1])
	var operands []int

	for _, str := range operandStrings {
		num, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Failed to parse operand: %s\n", str)
			continue
		}
		operands = append(operands, num)
	}

	return testValue, operands
}
