package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Define a struct to represent a student's filter criteria.
type filterStudent struct {
	Income_decile, Ability_decile, number int64
	OUR_SAT_needAid                       float64
}

// fillUnis function allocates students to universities based on their preferences.
func fillUnis(students []Student) ([]int64, []int64, []int64, []int64, []int64) {
	// Get the number of students.
	n := len(students)

	// Calculate the initial capacity for each university.
	A_cap := n / 5
	B_cap := n / 5
	C_cap := n / 5
	D_cap := n / 5
	E_cap := n / 5

	// Initialize lists to track students assigned to each university.
	A_students_people := []int64{}
	B_students_people := []int64{}
	C_students_people := []int64{}
	D_students_people := []int64{}
	E_students_people := []int64{}
	W_students_people := []int64{}

	// Loop through the students and their preferences.
	for _, student := range students {
		num := student.number

		// Create a list of university preferences for the student.
		preference_list := []int{student.util_quality_A, student.util_quality_B, student.util_quality_C, student.util_quality_D, student.util_quality_E}

		// Iterate through the preferences and assign students to universities.
		for _, prefer := range preference_list {
			W_students_people = append(W_students_people, num)
			if prefer == 1 && A_cap >= 1 {
				A_cap -= 1
				A_students_people = append(A_students_people, num)
				break
			} else if prefer == 2 && B_cap >= 1 {
				B_cap -= 1
				B_students_people = append(B_students_people, num)
				break
			} else if prefer == 3 && C_cap >= 1 {
				C_cap -= 1
				C_students_people = append(C_students_people, num)
				break
			} else if prefer == 4 && D_cap >= 1 {
				D_cap -= 1
				D_students_people = append(D_students_people, num)
				break
			} else if prefer == 5 && E_cap >= 1 {
				E_cap -= 1
				E_students_people = append(E_students_people, num)
				break
			}
		}
	}

	return A_students_people, B_students_people, C_students_people, D_students_people, E_students_people
}

// filterStudents function filters students based on a list of student numbers.
func filterStudents(studentlist []int64, students []Student) []subStudent {
	filteredData := make([]subStudent, 0)

	// Loop through students and filter based on the provided student numbers.
	for _, s := range students {
		for _, num := range studentlist {
			if num == s.number {
				filteredData = append(filteredData, subStudent{Income_decile: s.Income_decile, Ability_decile: s.Ability_decile, OUR_SAT_needAid: s.OUR_SAT_needAid})
			}
		}
	}

	return filteredData
}

// group function groups substudents based on income and ability deciles.
func group(substudents []subStudent, n int) []subStudent {
	groupedData := make(map[string]float64)
	key := ""

	// Loop through substudents and group them based on income and ability deciles.
	for _, s := range substudents {
		key += fmt.Sprintf("%v - %v", s.Income_decile, s.Ability_decile)
		groupedData[key]++
		key = ""
	}

	return cnvrtMapToDF(groupedData, n)
}

// cnvrtMapToDF function converts a map to a slice of substudents.
func cnvrtMapToDF(a_map map[string]float64, n int) []subStudent {
	converted := make([]subStudent, 0)
	N := float64(n)

	// Loop through the map and convert it to substudents.
	for key, count := range a_map {
		parts := strings.Split(key, " - ")
		num1, err1 := strconv.ParseInt(parts[0], 10, 64)
		num2, err2 := strconv.ParseInt(parts[1], 10, 64)

		if err1 == nil && err2 == nil {
			converted = append(converted, subStudent{Income_decile: num1, Ability_decile: num2, OUR_SAT_needAid: count * 100 / N})
		}
	}

	return converted
}
