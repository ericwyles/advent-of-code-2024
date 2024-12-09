package main

import (
	_ "embed"
	"fmt"
	"strconv"
)

//go:embed input.txt
var embeddedFile string

type FileMetadata struct {
	fileId   int
	location int
	size     int
}

func main() {
	fileId := 0
	diskSize := 0
	var diskMap []int
	var originalFileLocations []FileMetadata
	var emptyLocations []FileMetadata

	for _, r := range embeddedFile {
		j, _ := strconv.Atoi(string(r))
		diskMap = append(diskMap, j)
		diskSize += j
	}

	wholeDisk := make([]int, diskSize)
	wholeDiskIndex := 0

	for i := range diskMap {
		size := diskMap[i]
		if i%2 == 0 {
			originalFileLocations = append(originalFileLocations, FileMetadata{fileId: fileId, location: wholeDiskIndex, size: size})
			setRange(wholeDisk, fileId, wholeDiskIndex, size)
			fileId++
		} else {
			emptyLocations = append(emptyLocations, FileMetadata{fileId: -1, location: wholeDiskIndex, size: size})
			setRange(wholeDisk, -1, wholeDiskIndex, size)
		}

		wholeDiskIndex += size
	}

	for i := len(originalFileLocations) - 1; i >= 0; i-- {
		fileToMove := originalFileLocations[i]
		emptySpace := findNextSufficientEmptySpace(&emptyLocations, fileToMove.size)
		if emptySpace == -1 {
			continue // can't move this because no space
		} else if emptySpace >= fileToMove.location {
			continue // can't move this because new spot is after original spot
		}

		moveFile(wholeDisk, fileToMove, emptySpace)
	}

	fmt.Printf("\nChecksum: %d\n", calcCheckSum(wholeDisk))
}

func moveFile(wholeDisk []int, fileToMove FileMetadata, newLocation int) {
	setRange(wholeDisk, fileToMove.fileId, newLocation, fileToMove.size)
	setRange(wholeDisk, -1, fileToMove.location, fileToMove.size)
}

func findNextSufficientEmptySpace(emptyLocations *[]FileMetadata, spaceNeeded int) int {
	for i, file := range *emptyLocations {
		if file.size >= spaceNeeded {
			location := file.location

			if spaceNeeded == file.size {
				// Remove the entry from availableSpace slice, this space is used up
				*emptyLocations = append((*emptyLocations)[:i], (*emptyLocations)[i+1:]...)
			} else {
				// Update the entry: increment location and decrement size, encapsulation kind of trash here but it's fast
				(*emptyLocations)[i].location += spaceNeeded
				(*emptyLocations)[i].size -= spaceNeeded
			}

			return location
		}
	}

	// No sufficient space found
	return -1
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
