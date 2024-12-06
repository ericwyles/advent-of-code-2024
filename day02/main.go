package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	numSafe := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var levels []int
		line := scanner.Text()

		fields := strings.Fields(line)

		for i := range fields {
			level, _ := strconv.Atoi(fields[i])
			levels = append(levels, level)
		}

		fmt.Println(levels)
		safe := checkLevelSafety(levels)

		if safe {
			numSafe++
		}
	}

	fmt.Printf("Total number safe [%d].\n", numSafe)

}

func checkLevelSafety(levels []int) bool {
	if checkDampenedLevelSafety(levels) {
		return true // if it works without removing anything it's good
	}

	for i := 0; i < len(levels); i++ {
		newLevels := append([]int{}, levels[:i]...)
		newLevels = append(newLevels, levels[i+1:]...)
		if checkDampenedLevelSafety(newLevels) {
			fmt.Printf("It's safe!\n")
			return true // if it worked like this it's fine
		}
	}
	return false // if we made it here, it is bad
}

func checkDampenedLevelSafety(levels []int) bool {
	direction := 0
	for i := 0; i < len(levels)-1; i++ {
		diff := levels[i] - levels[i+1]
		absDiff, sign := absInt(diff)
		if direction != 0 {
			if sign != direction {
				//fmt.Printf("Unsafe because sign of [%d] between pair [%d,%d] changed from previous.\n", sign, levels[i], levels[i+1])
				return false
			}
		}
		if absDiff < 1 || absDiff > 3 {
			//fmt.Printf("Unsafe because diff of [%d] between pair [%d,%d] is out of tolerance.\n", absDiff, levels[i], levels[i+1])
			return false // unsafe because amount of increase
		}
		direction = sign
	}
	return true
}

func absInt(n int) (int, int) {
	if n < 0 {
		return -n, -1
	}
	return n, 1
}
