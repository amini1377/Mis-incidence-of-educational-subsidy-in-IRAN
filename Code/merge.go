package main

// Define a struct to represent the decile factors for income and ability.
type DF struct {
	Income_decile, Ability_decile int64
}

// Define a struct to represent subStudent data, including income, ability, and OUR_SAT_needAid.
type subStudent struct {
	Income_decile, Ability_decile int64
	OUR_SAT_needAid               float64
}

// baseDF function generates a base dataset of DF structs representing all combinations of income and ability deciles.
func baseDF() []DF {
	dfi := []DF{
		{Income_decile: 1},
		{Income_decile: 2},
		{Income_decile: 3},
		{Income_decile: 4},
		{Income_decile: 5},
	}

	dfs := []DF{
		{Ability_decile: 1},
		{Ability_decile: 2},
		{Ability_decile: 3},
		{Ability_decile: 4},
		{Ability_decile: 5},
	}

	mergedData := make([]DF, 0)

	// Create all possible combinations of income and ability deciles.
	for _, s1 := range dfi {
		for _, s2 := range dfs {
			mergedData = append(mergedData, DF{Income_decile: s1.Income_decile, Ability_decile: s2.Ability_decile})
		}
	}

	return mergedData
}

// merge function combines the subStudent data with the base data and fills in missing OUR_SAT_needAid values.
func merge(school []subStudent) []subStudent {

	mergedData := make([]subStudent, 0)
	base := baseDF()

	for _, r1 := range base {
		found := false

		for _, r2 := range school {
			// If a matching income and ability decile is found in school data, use the OUR_SAT_needAid value.
			if r1.Income_decile == r2.Income_decile && r1.Ability_decile == r2.Ability_decile {
				mergedData = append(mergedData, subStudent{Income_decile: r1.Income_decile, Ability_decile: r1.Ability_decile, OUR_SAT_needAid: r2.OUR_SAT_needAid})
				found = true
				break
			}
		}

		// If no matching data is found, set OUR_SAT_needAid to 0.
		if !found {
			mergedData = append(mergedData, subStudent{Income_decile: r1.Income_decile, Ability_decile: r1.Ability_decile, OUR_SAT_needAid: 0})
		}
	}

	// Print merged data for debugging purposes.
	// for _, item := range mergedData {
	//     fmt.Println(item.Income_decile, item.Ability_decile)
	//     fmt.Println("")
	// }

	return mergedData
}
