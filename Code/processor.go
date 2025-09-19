package main

import (
	"fmt"
	"gumbel_test/gumbel"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Define a struct to represent a student with various attributes.
type Student struct {
	Income, EduInv1, Aid, Ability, EFC, OUR_SAT_needAid, OUR_SAT, util_A, util_B, util_C, util_D, util_E float64
	Income_decile, Ability_decile, number                                                                int64
	util_quality_A, util_quality_B, util_quality_C, util_quality_D, util_quality_E                       int
}

// Define a struct to represent income and ability deciles.
type Deciles struct {
	Income_decile, Ability_decile int64
	Count                         int
}

// Define constants for the beta and mu values.
var beta, mu = 0.8, 1
var sigma = 0.008

// Function to calculate the OUR_SAT value based on ability, eduInv1, and aid.
func calculateOurSat(ability, eduInv1, aid float64) float64 {
	// b_bar = 10
	// zeta = 0.35

	// our_sat = []

	// for i, j in enumerate(x_noaid['Ability']):
	// 	if j >= b_bar:
	// 		our_sat.append(b_bar + (x_noaid.iloc[i]['Ability'] ** zeta) * (x_noaid.iloc[i]['edu_inv1'] ** (1 - zeta)))
	// 	else:
	// 		our_sat.append(j)

	var alpha, b_bar = 0.41, 10. //zeta= 0.35

	if ability >= b_bar {
		return (b_bar + math.Pow(ability-b_bar, alpha)) * math.Pow(eduInv1, 1-alpha)
	} else {
		return ability
	}
	// return math.Pow(ability, 0.085) * math.Pow(eduInv1+aid, 0.915) OLD SAT
}

// Function to calculate the utility quality based on multiple factors.
func calculateUtilQuality(income, eduInv1, levelTuition, aid, quality, ability, beta float64) float64 {
	// noise = np.random.gumbel(256181329.3349902, 2000, x_need_aid.shape[0])

	gumbel := gumbel.GumbelRight{
		Mu:   1,
		Beta: 0.008,
	}

	noise := gumbel.Rand()
	return (math.Log(income-eduInv1-levelTuition+aid+1) + math.Log(quality+1) + beta*math.Log(ability+1)) + noise
}

// Function to calculate deciles from CSV rows.
func calDeciles(rows [][]string) []Deciles {
	var deciles []Deciles

	// Extract data from CSV rows and create deciles.
	for _, row := range rows {
		Income_decile, _ := strconv.ParseInt(row[3], 10, 32)
		Ability_decile, _ := strconv.ParseInt(row[4], 10, 32)
		count := 0

		deciles = append(deciles, Deciles{
			Income_decile: Income_decile, Ability_decile: Ability_decile, Count: count,
		})
	}

	// Group and process deciles.
	caledDeciles := groupDeciles(deciles)

	return caledDeciles
}

// Function to process data and calculate various attributes for students.
func Proccess(rows [][]string, aids_list []int) []Student {
	rand.Seed(time.Now().UnixNano())

	var students []Student

	levelTuition := map[string]float64{
		"A": 1.6,
		"B": .8,
		"C": .7,
		"D": .5,
		"E": 0,
	}

	quality_utility := map[string]float64{
		"A": 499.4868378815731,
		"B": 443.852001391096,
		"C": 370.93193653432913,
		"D": 339.1289739989186,
		"E": 345.69276059244614,
	}

	// Process CSV rows and calculate student attributes.
	for _, row := range rows {
		ability, _ := strconv.ParseFloat(row[0], 64)
		income, _ := strconv.ParseFloat(row[1], 64)
		eduInv1, _ := strconv.ParseFloat(row[2], 64)
		Income_decile, _ := strconv.ParseInt(row[3], 10, 32)
		Ability_decile, _ := strconv.ParseInt(row[4], 10, 32)

		if int(Income_decile) == 0 {
			continue
		}

		EFC, _ := strconv.ParseFloat(row[5], 64)
		OUR_SAT_needAid, _ := strconv.ParseFloat(row[6], 64)
		number, _ := strconv.ParseInt(row[7], 10, 32)
		aid_ := calAid(Income_decile, Ability_decile, aids_list) * 1

		uq_A := calculateUtilQuality(income, eduInv1, levelTuition["A"], aid_, quality_utility["A"], ability, beta)
		uq_B := calculateUtilQuality(income, eduInv1, levelTuition["B"], aid_, quality_utility["B"], ability, beta)
		uq_C := calculateUtilQuality(income, eduInv1, levelTuition["C"], aid_, quality_utility["C"], ability, beta)
		uq_D := calculateUtilQuality(income, eduInv1, levelTuition["D"], aid_, quality_utility["D"], ability, beta)
		uq_E := calculateUtilQuality(income, eduInv1, levelTuition["E"], aid_, quality_utility["E"], ability, beta)

		num_qualities := []float64{uq_A, uq_B, uq_C, uq_D, uq_E}
		rankList := calculateRank(num_qualities)

		students = append(students, Student{
			Ability:         ability,
			Income:          income,
			EduInv1:         eduInv1,
			Aid:             aid_,
			Income_decile:   Income_decile,
			Ability_decile:  Ability_decile,
			EFC:             EFC,
			OUR_SAT_needAid: OUR_SAT_needAid,
			number:          number,

			OUR_SAT:        calculateOurSat(ability, eduInv1, aid_),
			util_quality_A: rankList[0],
			util_quality_B: rankList[1],
			util_quality_C: rankList[2],
			util_quality_D: rankList[3],
			util_quality_E: rankList[4],

			util_A: uq_A,
			util_B: uq_B,
			util_C: uq_C,
			util_D: uq_D,
			util_E: uq_E,
		})
	}

	// Sort students by OUR_SAT_needAid in descending order to prirotize people with higher rank in assinging process.
	sort.Slice(students, func(i, j int) bool {
		return students[i].OUR_SAT_needAid > students[j].OUR_SAT_needAid
	})

	return students
}

// We define two following functions that we used in follwing sim_util function.
func cal_Average(arr []float64) float64 {
	sumaa := 0.0
	for _, numbera := range arr {
		sumaa += numbera
	}
	return sumaa / float64(len(arr))
}

func cal_Variance(arr []float64) float64 {
	mean := cal_Average(arr)
	sumSquaredDiff := 0.0
	for _, numbera := range arr {
		diff := numbera - mean
		sumSquaredDiff += diff * diff
	}
	return sumSquaredDiff / float64(len(arr))
}

// Function to calculate the utility for students in group A.
func sim_util(students []Student, A []int64, B []int64, C []int64, D []int64, E []int64) float64 {
	S_A := 0.0
	S_B := 0.0
	S_C := 0.0
	S_D := 0.0
	S_E := 0.0
	S_W := 0.0
	// S_A_mean := 0.0
	// S_B_mean := 0.0
	// S_C_mean := 0.0
	// S_D_mean := 0.0
	// S_E_mean := 0.0

	stu_util_array := []float64{}
	// Loop through each student and accumulate utility values for each group.
	for _, student := range students {
		for i := 0; i < len(A); i++ {
			if A[i] == int64(student.number) {
				stu_util_array = append(stu_util_array, student.util_A)
				S_A += student.util_A
				S_W += student.util_A
				break
			}
		}

		for i := 0; i < len(B); i++ {
			if B[i] == int64(student.number) {
				stu_util_array = append(stu_util_array, student.util_B)
				S_B += student.util_B
				S_W += student.util_B
				break
			}
		}

		for i := 0; i < len(C); i++ {
			if C[i] == int64(student.number) {
				stu_util_array = append(stu_util_array, student.util_C)
				S_C += student.util_C
				S_W += student.util_C
				break
			}
		}

		for i := 0; i < len(D); i++ {
			if D[i] == int64(student.number) {
				stu_util_array = append(stu_util_array, student.util_D)
				S_D += student.util_D
				S_W += student.util_D
				break
			}
		}

		for i := 0; i < len(E); i++ {
			if E[i] == int64(student.number) {
				stu_util_array = append(stu_util_array, student.util_E)
				S_E += student.util_E
				S_W += student.util_E
				break
			}
		}

	}

	return cal_Variance(stu_util_array)
}

// Function to calculate the total utility for each group of students that we call it social welfare function.
func sum_util(students []Student, A []int64, B []int64, C []int64, D []int64, E []int64) float64 {
	S_A := 0.0
	S_B := 0.0
	S_C := 0.0
	S_D := 0.0
	S_E := 0.0

	// Loop through each student and accumulate utility values for each group.
	for _, student := range students {
		for i := 0; i < len(A); i++ {
			if A[i] == int64(student.number) {
				S_A += student.util_A
				break
			}
		}

		for i := 0; i < len(B); i++ {
			if B[i] == int64(student.number) {
				S_B += student.util_B
				break
			}
		}

		for i := 0; i < len(C); i++ {
			if C[i] == int64(student.number) {
				S_C += student.util_C
				break
			}
		}

		for i := 0; i < len(D); i++ {
			if D[i] == int64(student.number) {
				S_D += student.util_D
				break
			}
		}

		for i := 0; i < len(E); i++ {
			if E[i] == int64(student.number) {
				S_E += student.util_E
				break
			}
		}
	}

	return S_A + S_B + S_C + S_D + S_E
}

// We don't use following functions RMSE, RMSE1, ... in code but we keep them to flexible our codeing and use when we like to do.

// Function to calculate the Euclidean distance between every two groups of students.
func dist(dfw []subStudent, df2 []subStudent) float64 {
	S := 0.0

	// Calculate the sum of squared differences between the OUR_SAT_needAid values of corresponding students in the two groups.
	for i := 0; i < len(dfw); i++ {
		S += math.Pow(dfw[i].OUR_SAT_needAid/5-df2[i].OUR_SAT_needAid, 2)
	}

	// Return the square root of the sum as the Euclidean distance.
	return math.Sqrt(S)
}

// Function to calculate the Root Mean Square Error (RMSE) between multiple groups.
func RMSE(dfw []subStudent, dfa []subStudent, dfb []subStudent, dfc []subStudent, dfd []subStudent, dfe []subStudent) float64 {
	// Calculate the distance between the reference group (dfw) and each of the other groups.
	d1 := dist(dfw, dfa)
	d2 := dist(dfw, dfb)
	d3 := dist(dfw, dfc)
	d4 := dist(dfw, dfd)
	d5 := dist(dfw, dfe)

	// Return the sum of these distances as the RMSE.
	return d1 + d2 + d3 + d4 + d5
}

// Another version of RMSE calculation, similar to the previous one.
func RMSE1(dfw []subStudent, dfa []subStudent, dfb []subStudent, dfc []subStudent, dfd []subStudent, dfe []subStudent) float64 {
	d1 := dist(dfw, dfa)
	d2 := dist(dfw, dfb)
	d3 := dist(dfw, dfc)
	d4 := dist(dfw, dfd)
	d5 := dist(dfw, dfe)

	// Return the sum of these distances as the RMSE.
	return d1 + d2 + d3 + d4 + d5
}

// Function to calculate aid based on income and ability deciles, using if-else statements and it will assign to every index value or number of students there.
func cal_aid(i_decile int64, a_decile int64, aids_list []int) float64 {
	// fmt.Println(i_decile, a_decile)

	if i_decile == 1 && a_decile == 5 {
		return float64(aids_list[0])
	} else if i_decile == 1 && a_decile == 4 {
		return float64(aids_list[1])
	} else if i_decile == 2 && a_decile == 5 {
		return float64(aids_list[2])
	} else if i_decile == 2 && a_decile == 4 {
		return float64(aids_list[3])
	} else if i_decile == 1 && a_decile == 3 {
		return float64(aids_list[4])
	} else if i_decile == 3 && a_decile == 5 {
		return float64(aids_list[5])
	} else if i_decile == 2 && a_decile == 3 {
		return float64(aids_list[6])
	} else if i_decile == 3 && a_decile == 4 {
		return float64(aids_list[7])
	} else if i_decile == 3 && a_decile == 3 {
		return float64(aids_list[8])
	} else if i_decile == 4 && a_decile == 4 {
		return float64(aids_list[9])
	} else if i_decile == 4 && a_decile == 3 {
		return float64(aids_list[10])
	} else if i_decile == 4 && a_decile == 2 {
		return float64(aids_list[11])
	} else if i_decile == 3 && a_decile == 2 {
		return float64(aids_list[12])
	} else if i_decile == 2 && a_decile == 2 {
		return float64(aids_list[13])
	} else if i_decile == 1 && a_decile == 2 {
		return float64(aids_list[14])
	} else if i_decile == 4 && a_decile == 1 {
		return float64(aids_list[15])
	} else if i_decile == 3 && a_decile == 1 {
		return float64(aids_list[16])
	} else if i_decile == 2 && a_decile == 1 {
		return float64(aids_list[17])
	} else if i_decile == 1 && a_decile == 1 {
		return float64(aids_list[18])
	}
	return 0.0

}

// Function to calculate aid based on income and ability deciles, using a mathematical formula.
func calAid(i_decile int64, a_decile int64, aids_list []int) float64 {
	// fmt.Println(int(5*i_decile+a_decile-6), i_decile, a_decile, float64(aids_list[int(5*i_decile+a_decile-6)]))
	return float64(aids_list[int(5*i_decile+a_decile-6)])
}

// Function to create a list of all possible income and ability deciles.
func AllDeciles(deciles []Deciles) []Deciles {
	result := make([]Deciles, 0)

	for i := 1; i < 6; i++ {
		for j := 1; j < 6; j++ {
			bool_ := true
			for _, d := range deciles {
				if i == int(d.Income_decile) && j == int(d.Ability_decile) {
					bool_ = false
					result = append(result, Deciles{d.Income_decile, d.Ability_decile, d.Count})
					break
				}
			}
			if bool_ {
				// fmt.Println("-------/-----------")
				result = append(result, Deciles{int64(i), int64(j), 0})
			}
		}
	}

	return result
}

// Function to group deciles based on income and ability, filling in missing combinations.
func groupDeciles(deciles []Deciles) []Deciles {
	groupedData := make(map[string]float64)
	key := ""

	for _, s := range deciles {
		key += fmt.Sprintf("%v - %v", s.Income_decile, s.Ability_decile)
		groupedData[key]++
		key = ""
	}
	return cnvrtMapToDecileDF(groupedData)
}

// Function to convert a map of deciles with counts to a slice of Deciles.
func cnvrtMapToDecileDF(a_map map[string]float64) []Deciles {
	converted := make([]Deciles, 0)
	// N := float64(n)
	for key, count := range a_map {
		parts := strings.Split(key, " - ")
		num1, err1 := strconv.ParseInt(parts[0], 10, 64)
		num2, err2 := strconv.ParseInt(parts[1], 10, 64)
		// newCount, _ := strconv.ParseInt(count, 10, 64)

		if err1 == nil && err2 == nil {
			converted = append(converted, Deciles{Income_decile: num1, Ability_decile: num2, Count: int(count)})
		}
	}

	return converted
}

// Function to extract counts from a list of Deciles and return them as a slice.
func calCoefs(deciles []Deciles) []int {
	res := make([]int, 0)

	for _, item := range deciles {
		res = append(res, item.Count)
	}

	return res
}
