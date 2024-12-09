package main

import (
	_ "embed"
	"fmt"
	"strconv"
)

//go:embed input.txt
var embeddedFile string

func main() {
	fileId := 0
	diskSize := 0
	var diskMap []int

	for _, r := range embeddedFile {
		j, _ := strconv.Atoi(string(r))
		diskMap = append(diskMap, j)
		diskSize += j
	}

	wholeDisk := make([]int, diskSize)
	wholeDiskIndex := 0

	for i := range diskMap {
		if i%2 == 0 {
			setRange(wholeDisk, fileId, wholeDiskIndex, diskMap[i])
			fileId++
		} else {
			setRange(wholeDisk, -1, wholeDiskIndex, diskMap[i])
		}

		wholeDiskIndex += diskMap[i]
	}

	//printIntsAsStrings(wholeDisk)
	for swapLastValue(wholeDisk) {
		// it's not super efficient but we'll see if i get by with it
	}
	//printIntsAsStrings(wholeDisk)

	fmt.Printf("\nChecksum: %d\n", calcCheckSum(wholeDisk))
}

func setRange(disk []int, value int, startIndex int, times int) {
	endIndex := startIndex + times
	for i := startIndex; i < endIndex; i++ {
		disk[i] = value
	}
}

func printIntsAsStrings(ints []int) {
	for _, value := range ints {
		if value == -1 {
			fmt.Print(".")
		} else {
			fmt.Printf("%d", value)
		}
	}
	fmt.Printf("\n")
}

func swapLastValue(wholeDisk []int) bool {
	firstEmpty := getFirstMatch(wholeDisk, -1)
	lastNonEmpty := getLastNonMatch(wholeDisk, -1)

	foundRequiredMove := firstEmpty != -1 && lastNonEmpty != -1 && firstEmpty < lastNonEmpty
	if foundRequiredMove {
		wholeDisk[firstEmpty] = wholeDisk[lastNonEmpty]
		wholeDisk[lastNonEmpty] = -1
	}

	return foundRequiredMove
}

func getFirstMatch(slice []int, target int) int {
	for i, r := range slice {
		if r == target {
			return i
		}
	}
	return -1
}

func getLastNonMatch(slice []int, target int) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] != target {
			return i
		}
	}
	return -1
}

func calcCheckSum(disk []int) int {
	var checkSum int
	for position, value := range disk {
		if value != -1 {
			checkSum += position * value
		}
	}

	return checkSum
}
