package main

import (
	"sort"
)

func calculateRank(arr []float64) []int {
	// Create a map to store the original indices of the array elements
	indexMap := make(map[float64]int)
	// Copy the array and sort it
	sortedArr := make([]float64, len(arr))
	copy(sortedArr, arr)
	sort.Float64s(sortedArr)
	sort.Slice(sortedArr, func(i, j int) bool {
		return sortedArr[i] > sortedArr[j]
	})

	// Assign ranks to the elements in the sorted array
	rank := 1
	for i := 0; i < len(sortedArr); i++ {
		if i > 0 && sortedArr[i] != sortedArr[i-1] {
			rank++
		}
		indexMap[sortedArr[i]] = rank
	}

	// Calculate the rank list for the original array
	rankList := make([]int, len(arr))
	for i, num := range arr {
		rankList[i] = indexMap[num]
	}

	return rankList
}
