package main

import (
	"sort"
)

// Function to generate partitions with a fixed upper bound
func genpart_fixedupper(N int, partitions []int, index int, numVariables int, coefficients []int, upperBound int, result *[][]int) {
	// this case is simmilar to gp_count explanation about budget tolerate and we give permission to algorithem to generate partition with 5 percent more or less than budget.
	if N < 37 && N > -37 && index == numVariables {
		// Store the partition in the result
		temp := make([]int, numVariables)
		copy(temp, partitions)
		*result = append(*result, temp)
		return
	}

	if N <= -37 || index >= numVariables {
		return
	}

	// Set the maximum value for the current variable considering the common upper bound
	maxValue := min(upperBound, N/coefficients[index])
	for i := 0; i <= maxValue; i++ {
		partitions[index] = i
		genpart_fixedupper(N-(i*coefficients[index]), partitions, index+1, numVariables, coefficients, upperBound, result)
	}
}

func gPartitions(N int, numVariables int, coefficients []int, upperBound int) [][]int {
	var partitions [][]int
	partition := make([]int, numVariables)
	genpart_fixedupper(N, partition, 0, numVariables, coefficients, upperBound, &partitions)
	return partitions
}

// a function to calcualte min between two values.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func gen_parts(N int, numVariables int, coefficients []int, upperBound int) [][]int {

	if len(coefficients) != numVariables {
		return make([][]int, 0)
	}

	result := gPartitions(N, numVariables, coefficients, upperBound)

	// Sort partitions in decreasing order (x1 >= x2 >= x3)
	sort.Slice(result, func(i, j int) bool {
		for k := range result[i] {
			if result[i][k] != result[j][k] {
				return result[i][k] > result[j][k]
			}
		}
		return false
	})

	return result
}
