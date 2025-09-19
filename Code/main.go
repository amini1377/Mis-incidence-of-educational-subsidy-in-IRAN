package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func main() {
	// Open the CSV file "x_need_aid.csv" for reading.
	file, err := os.Open("x_need_aid.csv")
	startTime := time.Now()

	// Check for errors when opening the CSV file.
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader and read all rows from the file.
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()

	// Check for errors when reading the CSV file.
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Calculate deciles based on the data in the CSV file.
	deciles := calDeciles(rows)

	// Group the deciles into categories.
	grpd := AllDeciles(deciles)

	// Calculate coefficients based on grouped deciles.
	coefs := calCoefs(grpd)

	fmt.Println("coefs:", coefs)

	// Define parameters for generating partitions.
	N := 730
	bound := 1

	// Generate and count partitions.
	count_gen_parts(N, 25, coefs, bound)

	// Generate partitions.
	aids := gen_parts(N, 25, coefs, bound)
	l := len(aids)

	// Initialize variables for finding min and max distances.
	min_dist := 1.2e20
	var min_aid = make([]int, 0)

	max_dist := 0.001
	var max_aid = make([]int, 0)

	// Iterate through generated partitions (aids) and evaluate their utility.
	for e, aid_ := range aids {
		if e%10000 == 0 {
			fmt.Println(e, l)
		}

		// Process the rows with the current aid_ partition.
		var data = Proccess(rows, aid_)

		// Fill universities with students based on preferences.
		dfa, dfb, dfc, dfd, dfe := fillUnis(data)

		// Calculate the utility or similarity based on filled universities.
		dist_ := sum_util(data, dfa, dfb, dfc, dfd, dfe) // Social Welafare maximizing policy maker.
		// dist_ := sim_util(data, dfa, dfb, dfc, dfd, dfe) // Dispersion of utility minimizer policy maker.

		// Update min and max distances and corresponding aid_ partitions.
		if dist_ < min_dist {
			min_dist = dist_
			min_aid = aid_
			fmt.Println(dist_)
		} else if dist_ > max_dist {
			max_dist = dist_
			max_aid = aid_
		}
	}

	// Display the results.
	fmt.Println("_____________*******************************_______________")
	fmt.Println("Between", len(aids), "number for different Aids:")
	fmt.Println("Min Utility is for :")
	fmt.Println(min_aid)
	fmt.Println("with rmse : ")
	fmt.Println(min_dist)
	fmt.Println("_____________*******************************_______________")
	fmt.Println("Max Utility is for :")
	fmt.Println(max_aid)
	fmt.Println("with rmse : ")
	fmt.Println(max_dist)

	// Calculate and display the elapsed time.
	elapsedTime := time.Since(startTime)
	fmt.Printf("Elapsed time: %s\n", elapsedTime)
}
