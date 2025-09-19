package main

import (
	"fmt"
)

// Function to count the number of partitions with a fixed upper bound
func countpart_fixedupper(N int, partitions []int, index int, numVariables int, coefficients []int, upperBound int) int {
	// Base case: if N is between -37 and 37 and we have reached the last variable it's assumptions that I talked about it in 'rahnama file' we give permission to alg to iterate with 5 percent tolerant in budget.
	if N < 37 && N > -37 && index == numVariables {
		return 1
	}
	// Base case: if N is less than or equal to -37 or we have reached the last variable
	if N <= -37 || index >= numVariables {
		return 0
	}

	count := 0

	// Set the maximum value for the current variable considering the common upper bound
	// Iterate over possible values for the current variable
	maxValue := min(upperBound, N/coefficients[index])
	for i := 0; i <= maxValue; i++ {
		partitions[index] = i
		// count the number of partitions for the remaining variables
		count += countpart_fixedupper(N-(i*coefficients[index]), partitions, index+1, numVariables, coefficients, upperBound)
	}

	return count
}

func countPartitions(N int, numVariables int, coefficients []int, upperBound int) int {
	partition := make([]int, numVariables)
	return countpart_fixedupper(N, partition, 0, numVariables, coefficients, upperBound)
}

// this function print number of partitions
func count_gen_parts(N int, numVariables int, coefficients []int, upperBound int) {
	// N := 37
	// numVariables := 3
	// coefficients := []int{2, 3, 4}
	// upperBound := 3 // Set the fixed upper bound for all variables

	result := countPartitions(N, numVariables, coefficients, upperBound)
	fmt.Println("Number of partitions:", result)
}
