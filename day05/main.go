package main

import (
	_ "embed"
	"fmt"
	"reflect"
	"slices"
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

	requiredBeforeMap := make(map[int][]int)
	for _, orderingRule := range orderingRules {
		page, requiredBefore := readOrderingRule(orderingRule)
		requiredBeforeMap[page] = append(requiredBeforeMap[page], requiredBefore)
	}

	middleSum := 0
	badMiddleSum := 0
	for _, pageNumbers := range pagesToProduce {
		pages := strings.Split(pageNumbers, ",")
		pageNums := readPages(pages)
		orderedPageNums := reorderSlice(pageNums, requiredBeforeMap)

		middle := getMiddleValue(orderedPageNums)

		if reflect.DeepEqual(pageNums, orderedPageNums) {
			middleSum += middle
		} else {
			badMiddleSum += middle
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
		if slices.Contains(allPages, requiredPage) {
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
