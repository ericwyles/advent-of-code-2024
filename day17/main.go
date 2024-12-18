package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	opcode  int
	operand int
}

var RegisterA int
var RegisterB int
var RegisterC int
var output string

var programString string
var program []int

func main() {
	parseInput()

	i := 0

	fmt.Printf("Program: %v\n", program)
	fmt.Printf("Program string: %s\n", programString)
	printState(i)

	a := RegisterA
	b := RegisterB
	c := RegisterC
	output = ""

	runProgram(program, false, "")

	fmt.Println("PART 1:")
	fmt.Println(output)

	if len(os.Args) > 1 {
		arg := os.Args[1] // First argument after the program name
		parsedValue, _ := strconv.Atoi(arg)
		a = parsedValue
	}

	RegisterA = a
	RegisterB = b
	RegisterC = c
	output = ""

	// some stuff for part 2 maybe
	for i = len(program) - 2; i >= 0; i-- {
		segment := ""
		for k := i; k < len(program)-1; k++ {
			if len(segment) > 0 {
				segment += ","
			}
			segment += fmt.Sprintf("%d", program[k])
		}
		segment += ",0"
		fmt.Printf("Finding match for program[%d]=[%s]...\n", i, segment)

		for j := 0; j <= 7; j++ {
			RegisterA = zeroLowest3Bits(RegisterA)
			RegisterA += j
			RegisterB, RegisterC = 0, 0
			output = ""

			runProgram(program, false, "")
			if output == segment {
				fmt.Printf("    Found a match %d\n", j)
				RegisterA = RegisterA << 3
			}
		}

	}

	//RegisterA = a
	RegisterB = b
	RegisterC = c
	output = ""
	fmt.Println("PART 2:")
	runProgram(program, true, "")
	fmt.Println(output)
	fmt.Println(programString)
}

func zeroLowest3Bits(n int) int {
	mask := ^7 // 7 in binary is 0b111, ~7 flips all bits
	return n & mask
}

func runProgram(program []int, debug bool, expectedOutput string) {
	checkExpected := len(expectedOutput) > 0

	i := 0
	for i < len(program)-1 {
		instruction := Instruction{opcode: program[i], operand: program[i+1]}

		i = executeInstruction(instruction, i)

		if checkExpected && len(output) > 0 {
			if !strings.HasPrefix(expectedOutput, output) {
				return
			}
		}

		if debug {
			printState(i)
		}
	}
}

func printState(instruction int) {
	fmt.Printf("IP = %02d - Registers: A=%08d, B=%08d, C=%08d - output: [%s]\n", instruction, RegisterA, RegisterB, RegisterC, output)
}

func executeInstruction(instruction Instruction, i int) int {
	j := i

	switch instruction.opcode {
	case 0:
		adv(instruction)
	case 1:
		bxl(instruction)
	case 2:
		bst(instruction)
	case 3:
		i = jnz(instruction, i)
	case 4:
		bxc()
	case 5:
		out(instruction)
	case 6:
		bdv(instruction)
	case 7:
		cdv(instruction)
	default:
		fmt.Printf("Invalid Instruction %v\n", instruction)
	}

	if j == i { // if these are still the same no jump happened
		return i + 2
	} else {
		return i
	}
}

func adv(instruction Instruction) {
	RegisterA = div(instruction)
}

func bdv(instruction Instruction) {
	RegisterB = div(instruction)
}

func cdv(instruction Instruction) {
	RegisterC = div(instruction)
}

func div(instruction Instruction) int {
	numerator := RegisterA
	shift := getComboOperand(instruction) // Operand determines the power of 2
	denominator := 1 << shift

	return numerator / denominator
}

func bst(instruction Instruction) {
	RegisterB = getComboOperand(instruction) % 8
}

func bxl(instruction Instruction) {
	RegisterB = RegisterB ^ instruction.operand
}

func bxc() {
	RegisterB = RegisterB ^ RegisterC
}

func out(instruction Instruction) {
	if len(output) > 0 {
		output += ","
	}

	output += fmt.Sprintf("%d", getComboOperand(instruction)%8)
}

func jnz(instruction Instruction, i int) int {
	if RegisterA == 0 {
		return i
	}

	return instruction.operand
}

func getComboOperand(instruction Instruction) int {
	if instruction.operand <= 3 {
		return instruction.operand
	}

	if instruction.operand == 4 {
		return RegisterA
	}

	if instruction.operand == 5 {
		return RegisterB
	}

	if instruction.operand == 6 {
		return RegisterC
	}

	panic(fmt.Sprintf("Invalid Operand %d", instruction.operand))
}
func parseInput() {
	scanner := bufio.NewScanner(os.Stdin)

	// Read input line by line
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Register A:") {
			RegisterA, _ = parseRegister(line)
		} else if strings.HasPrefix(line, "Register B:") {
			RegisterB, _ = parseRegister(line)
		} else if strings.HasPrefix(line, "Register C:") {
			RegisterC, _ = parseRegister(line)
		} else if strings.HasPrefix(line, "Program:") {
			program = parseProgram(line)
		}
	}
}

// Helper to parse register lines
func parseRegister(line string) (int, error) {
	parts := strings.Split(line, ":")
	if len(parts) == 2 {
		return strconv.Atoi(strings.TrimSpace(parts[1]))
	}
	return 0, nil
}

// Helper to parse the Program line
func parseProgram(line string) []int {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return []int{}
	}

	rawProgram := strings.TrimSpace(parts[1])
	programString = rawProgram
	numStrings := strings.Split(rawProgram, ",")
	var result []int

	for _, numStr := range numStrings {
		if num, err := strconv.Atoi(strings.TrimSpace(numStr)); err == nil {
			result = append(result, num)
		}
	}
	return result
}
