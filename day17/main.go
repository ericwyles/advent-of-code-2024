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

var bitSegments []int // Array to store bit segments

func main() {
	parseInput()

	i := 0

	fmt.Printf("Program: %v\n", program)
	fmt.Printf("Program string: %s\n", programString)
	printState(i)

	output = ""

	runProgram(program, false, "")

	fmt.Println("PART 1:")
	fmt.Println(output)

	fmt.Println("PART 2:")
	bitSegments = make([]int, len(program))

	// Attempt to brute-force every output from last to first
	reconstructOutputBits(0)

	// Connect all segments to create the initial value for register A
	var initialRegisterA int = 0
	for i := 0; i < len(bitSegments); i++ {
		initialRegisterA = initialRegisterA << 3
		initialRegisterA += bitSegments[i]
	}

	fmt.Println("Initial Register A:", initialRegisterA)
	runProgram2(initialRegisterA)
	fmt.Println(output)
	fmt.Println(programString)
}

func reconstructOutputBits(depth int) bool {
	if depth == len(bitSegments) {
		return true
	}

	// Compute previous values
	var previousValues = 0
	for i := 0; i < depth; i++ {
		previousValues += bitSegments[i]
		previousValues = previousValues << 3
	}

	// Attempt to determine bits for this output
	for i := 0; i < 8; i++ {
		a := previousValues + i
		resultingOutput := runcalc(a)
		if resultingOutput == program[len(program)-1-depth] {
			bitSegments[depth] = i
			if reconstructOutputBits(depth + 1) {
				return true
			}
		}
	}

	return false
}

func runProgram2(a int) { // this is a super simplified version of my input program logic hard-coded so it can be used to check values
	output = ""
	for a > 0 {
		calc := runcalc(a)
		if len(output) > 0 {
			output += ","
		}
		output += fmt.Sprintf("%d", calc)
		a = a / 8
	}
}

func runcalc(a int) int {
	mod8 := a % 8
	return ((mod8 ^ 1) ^ (a / (1 << (mod8 ^ 2)))) % 8
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
