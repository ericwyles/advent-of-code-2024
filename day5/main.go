package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var embeddedFile string

func main() {
	// Split the embedded content into lines
	lines := strings.Split(embeddedFile, "\n")

	var orderingRules []string
	var pagesToProduce []string
	foundBlankLine := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check if the line is blank
		if line == "" {
			foundBlankLine = true
			continue
		}

		// Append lines to the appropriate array
		if !foundBlankLine {
			orderingRules = append(orderingRules, line)
		} else {
			pagesToProduce = append(pagesToProduce, line)
		}
	}

	// Print the results
	fmt.Println("Ordering Rules:")
	fmt.Println(orderingRules)

	fmt.Println("Pages to Produce:")
	fmt.Println(pagesToProduce)

	requiredBeforeMap := make(map[int][]int)
	for _, orderingRule := range orderingRules {
		page, requiredBefore := readOrderingRule(orderingRule)
		requiredBeforeMap[page] = append(requiredBeforeMap[page], requiredBefore)
	}

	fmt.Println("Required before map")
	fmt.Println(requiredBeforeMap)

	middleSum := 0
	badMiddleSum := 0
	for _, pageNumbers := range pagesToProduce {
		pages := strings.Split(pageNumbers, ",")
		pageNums := readPages(pages)

		fmt.Printf("Checking update for correct order %v:\n", pageNums)

		goodInput := true
		var pagesAlreadyProcessed []int
		for _, pageNum := range pageNums {
			//fmt.Printf("   Checking if required pages (%v) already processed for page %d\n", requiredBeforeMap[pageNum], pageNum)

			if goodInput {
				for _, requiredPage := range requiredBeforeMap[pageNum] {
					if contains(pageNums, requiredPage) {
						// this only matters if the requiredPage is in this input
						if !contains(pagesAlreadyProcessed, requiredPage) {
							fmt.Printf("    Rule %d|%d broken\n", requiredPage, pageNum)
							goodInput = false
						}
					}
				}
				pagesAlreadyProcessed = append(pagesAlreadyProcessed, pageNum)
			}

		}

		if goodInput {
			middleSum += getMiddleValue(pageNums)
			fmt.Printf("    Good input found: %v.\n", pageNums)
		} else {
			badMiddleSum += getMiddleValue(reorderSlice(pageNums, requiredBeforeMap))
			fmt.Printf("    ERROR bad input found: %v\n", pages)
		}

	}

	fmt.Printf("Sum of good middle values is: %d\n", middleSum)
	fmt.Printf("Sum of bad middle values is: %d\n", badMiddleSum)
}

func readOrderingRule(orderingRule string) (int, int) {
	parts := strings.Split(orderingRule, "|")
	page, _ := strconv.Atoi(parts[1])
	requiredBefore, _ := strconv.Atoi(parts[0])
	return page, requiredBefore
}

func readPages(pages []string) []int {
	var pagesInt []int
	for _, part := range pages {
		pageNum, _ := strconv.Atoi(part)
		pagesInt = append(pagesInt, pageNum)
	}

	return pagesInt
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func reorderSlice(slice []int, requiredBeforeMap map[int][]int) []int {
	// Create a copy of the slice to avoid modifying the original
	result := make([]int, len(slice))
	copy(result, slice)

	// Use the sort.Slice function with the custom comparison logic
	sort.Slice(result, func(i, j int) bool {
		return isRequiredBefore(result[i], result[j], slice, requiredBeforeMap)
	})

	return result
}

func isRequiredBefore(page int, otherPage int, allPages []int, requiredBeforeMap map[int][]int) bool {
	for _, requiredPage := range requiredBeforeMap[otherPage] {
		if contains(allPages, requiredPage) {
			// this only matters if the requiredPage is in this input
			if page == requiredPage {
				return true
			}
		}
	}
	return false
}

func getMiddleValue(pages []int) int {
	return pages[len(pages)/2]
}
