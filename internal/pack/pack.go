package pack

import (
	"maps"
)

var (
	sizes = []int{5000, 2000, 1000, 500, 250}
)

// we can set sizes from the ui using this function
func SetSizes(newSizes []int) []int {
	sizes = newSizes
	return sizes
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
