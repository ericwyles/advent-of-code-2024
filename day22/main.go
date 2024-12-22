package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Sequence struct {
	a, b, c, d int
}

func main() {
	numbers := readInput()

	start1 := time.Now()
	sumSecrets := 0
	buyerOffers := make([][]int, len(numbers))
	for i, num := range numbers {
		secret, offers := rotate(num, 2000)
		sumSecrets += secret
		buyerOffers[i] = offers
	}
	elapsed1 := time.Since(start1)
	fmt.Printf("Part 1 - Sum of new secrets: %d\n", sumSecrets)

	start2 := time.Now()
	buyerSequenceMaps := make([]map[Sequence]int, len(numbers))
	uniqueSequences := make(map[Sequence]bool)
	for i, offers := range buyerOffers {
		buyerSequences := getOfferSequences(offers)
		buyerSequenceMaps[i] = buyerSequences
		fmt.Printf("Buyer %d Sequences: %d\n", i, len(buyerSequences))
	}

	for _, buyerSequences := range buyerSequenceMaps {
		for k := range buyerSequences {
			uniqueSequences[k] = true
		}
	}

	fmt.Printf("Total Unique Sequences: %d\n", len(uniqueSequences))
	bestTotalOffer := 0
	var bestSequence Sequence
	for k, _ := range uniqueSequences {
		offer := getTotalOffers(k, buyerSequenceMaps)
		if offer > bestTotalOffer {
			bestTotalOffer = offer
			bestSequence = k
		}
	}

	// i know from previous submissions that 2236 is too low and 2276 is too high
	///    so i have a subtle bug above somehow
	fmt.Printf("Best sequence: %v Offer: %d\n", bestSequence, bestTotalOffer)
	elapsed2 := time.Since(start2)

	fmt.Printf("Calculation time: Part 1 [%d]ms, Part 2 [%d] ms\n", elapsed1.Milliseconds(), elapsed2.Milliseconds())
}

func getTotalOffers(seq Sequence, buyerSequenceMaps []map[Sequence]int) int {
	total := 0
	for _, bsm := range buyerSequenceMaps {
		if bestOffer, exists := bsm[seq]; exists {
			total += bestOffer
		}
	}
	return total
}

func getOfferSequences(offers []int) map[Sequence]int {
	var sequenceMap = make(map[Sequence]int)
	for i := 4; i < len(offers); i++ {
		seq, offer := getSequenceAtPosition(offers, i)
		// Add the sequence only if it hasn't been seen before
		if _, exists := sequenceMap[seq]; !exists {
			sequenceMap[seq] = offer
		}
	}
	return sequenceMap
}

// start this at 4 and go to end
func getSequenceAtPosition(offers []int, i int) (Sequence, int) {
	seq := Sequence{a: offers[i-3] - offers[i-4],
		b: offers[i-2] - offers[i-3],
		c: offers[i-1] - offers[i-2],
		d: offers[i] - offers[i-1]}
	return seq, offers[i]
}

func rotate(secret, times int) (int, []int) {
	var offers []int
	var s uint64 = uint64(secret) // Cast to uint64 for safety
	for i := 0; i < times; i++ {
		s = (s ^ (s * 64)) % 16777216
		s = (s ^ (s / 32)) % 16777216
		s = (s ^ (s * 2048)) % 16777216
		offers = append(offers, int(s%10))
	}
	return int(s), offers
}

func readInput() []int {
	scanner := bufio.NewScanner(os.Stdin)

	var numbers []int

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			num, _ := strconv.Atoi(line)
			numbers = append(numbers, num)
		}
	}

	return numbers
}
