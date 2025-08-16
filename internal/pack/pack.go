package pack

import (
	"context"
	"maps"
	"sort"
)

var (
	sizes = []int{5000, 2000, 1000, 500, 250}
)

type CalculateParams struct {
	Amount int64
}

type SinglePackResult struct {
	PackSize int64 `json:"pack_size"`
	Count    int64 `json:"count"`
}

type CalculationResult struct {
	Amount     int                `json:"amount"`
	TotalItems int64              `json:"total_items"`
	TotalPacks int64              `json:"total_packs"`
	Packs      []SinglePackResult `json:"packs"`
}

// we can set sizes from the ui using this function
func SetSizes(newSizes []int) []int {
	sizes = newSizes
	return sizes
}

func CalculateNeededPackSizes(ctx context.Context, amount int) CalculationResult {
	// Get all active pack sizes
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	// Calculate optimal packs
	packResult, totalItems := calculateOptimalPacks(amount, sizes)

	// Count total packs
	totalPacks := int64(0)
	for _, pr := range packResult {
		totalPacks += pr.Count
	}

	// Create calculation result
	result := CalculationResult{
		Amount:     amount,
		TotalItems: int64(totalItems),
		TotalPacks: totalPacks,
		Packs:      packResult,
	}
	return result
}

// calculateOptimalPacks calculates the optimal packs for an order
func calculateOptimalPacks(itemsOrdered int, packSizes []int) ([]SinglePackResult, int) {
	if len(packSizes) == 0 {
		return []SinglePackResult{}, 0
	}

	// Sort pack sizes in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	// if itemsOrdered is less than or equal to the smallest pack size
	smallestPackSize := packSizes[len(packSizes)-1]
	if itemsOrdered <= smallestPackSize {
		return []SinglePackResult{
			{PackSize: int64(smallestPackSize), Count: 1},
		}, smallestPackSize
	}

	// Determine a reasonable upper bound for our DP array
	maxDPSize := itemsOrdered * 2
	if maxDPSize > 1000000 {
		maxDPSize = 1000000
	}

	// Initialize dp array to track minimum items needed for each quantity
	dp := make([]int, maxDPSize+1)
	for i := range dp {
		dp[i] = -1 // -1 means not possible
	}
	dp[0] = 0 // Base case: 0 items needed for 0 items ordered

	// Track which pack size was used for each quantity
	packUsed := make([]int, maxDPSize+1)

	// Fill the dp array
	for i := 1; i <= maxDPSize; i++ {
		for _, size := range packSizes {
			if i >= size && dp[i-size] != -1 {
				// If we can use this pack size and it results in fewer items
				newTotal := dp[i-size] + size
				if dp[i] == -1 || newTotal < dp[i] {
					dp[i] = newTotal
					packUsed[i] = size
				}
			}
		}

		if i >= itemsOrdered && dp[i] != -1 {
			if dp[itemsOrdered] != -1 {
				break
			}
		}
	}

	// Find the minimum valid quantity that satisfies the order
	minValidQuantity := itemsOrdered
	while := true
	for while && minValidQuantity <= maxDPSize {
		if dp[minValidQuantity] != -1 {
			while = false
		} else {
			minValidQuantity++
		}
	}

	// If no valid solution found in our DP array
	if minValidQuantity > maxDPSize || dp[minValidQuantity] == -1 {
		return packAllocation(itemsOrdered, packSizes)
	}

	packCounts := make(map[int]int)
	remaining := minValidQuantity

	for remaining > 0 {
		packSize := packUsed[remaining]
		packCounts[packSize]++
		remaining -= packSize
	}

	// Convert to result format
	var result []SinglePackResult
	for size, count := range packCounts {
		result = append(result, SinglePackResult{
			PackSize: int64(size),
			Count:    int64(count),
		})
	}

	// Sort result by pack size in descending order for readability
	sort.Slice(result, func(i, j int) bool {
		return result[i].PackSize > result[j].PackSize
	})

	return result, dp[minValidQuantity]
}

// packAllocation is a fallback method that uses a greedy approach
// for very large orders that can't be handled by the DP algorithm
func packAllocation(amount int, packSizes []int) ([]SinglePackResult, int) {
	// Ensure pack sizes are sorted in descending order
	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	remaining := amount
	packCounts := make(map[int]int)
	totalItems := 0

	// Use as many large packs as possible
	for _, size := range packSizes {
		if remaining <= 0 {
			break
		}

		count := remaining / size
		if count > 0 {
			packCounts[size] = count
			remaining -= count * size
			totalItems += count * size
		}
	}

	// If there's still a remainder, add one more of the smallest pack
	if remaining > 0 {
		smallestSize := packSizes[len(packSizes)-1]
		packCounts[smallestSize]++
		totalItems += smallestSize
	}

	// Convert to result format
	var result []SinglePackResult
	for size, count := range packCounts {
		result = append(result, SinglePackResult{
			PackSize: int64(size),
			Count:    int64(count),
		})
	}

	// Sort result by pack size in descending order
	sort.Slice(result, func(i, j int) bool {
		return result[i].PackSize > result[j].PackSize
	})

	return result, totalItems
}

// Correct returns a map[size]count covering x using available sizes.
// It greedily fills from largest to smallest, adds one smallest pack if a remainder exists,
// then calls optimizePacks to combine smaller packs into larger ones.
// Precondition: sizes must be sorted descending.
// Example:
//
//	Correct(1) // -> map[int]int{250:1}
func Correct(x int) map[int]int {
	packs := make(map[int]int)
	for _, size := range sizes {
		if x <= 0 {
			break
		}
		if cnt := x / size; cnt > 0 {
			packs[size] = cnt
			x -= cnt * size
		}
	}
	if x > 0 {
		smallest := sizes[len(sizes)-1]
		packs[smallest]++

	}
	optimize(packs)
	return packs
}

func optimize(packs map[int]int) {
	for i := len(sizes) - 1; i > 0; i-- {
		small := sizes[i]
		large := sizes[i-1]
		requiredSmall := (large + small - 1) / small
		if requiredSmall <= 1 {
			continue
		}
		if have := packs[small]; have >= requiredSmall {
			convert := have / requiredSmall
			packs[small] -= convert * requiredSmall
			packs[large] += convert
			if packs[small] == 0 {
				delete(packs, small)
			}
		}
	}
}

// InCorrect returns a list of all incorrect pack combinations for a given ordered amount.
// It generates all possible combinations of packs, calculates the correct combination,
// and then filters out the correct one from the list of all combinations.
// Example:
//
//	InCorrect(1) // -> []map[int]int{{500:1}, {250:2}, {1000:1}, ...}
//
// Note: This function assumes that the pack sizes are sorted in descending order.
// It generates combinations based on the available pack sizes and the ordered amount.
// It returns a slice of maps, where each map represents a combination of pack sizes and their counts.
// The function uses a greedy approach to find the correct combination and then filters out that combination from
// the list of all combinations to return only the incorrect ones.package pack
func InCorrect(x int) []map[int]int {
	incorrect := []map[int]int{}
	all := []map[int]int{}
	for _, size := range sizes {
		if size >= x {
			all = append(all, map[int]int{size: 1})
		} else {
			count := (x / size) + 1
			all = append(all, map[int]int{size: count})
		}
	}
	packs := Correct(x)
	for _, combination := range all {
		if !maps.Equal(combination, packs) {
			incorrect = append(incorrect, combination)
		}
	}
	return incorrect
}
