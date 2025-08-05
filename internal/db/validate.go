package db

// columnsMatch checks if two string slices are equal (order and content).
// Used for validating that INSERT columns match the table schema exactly.
func ColumnsMatch(expected, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i := range expected {
		if expected[i] != actual[i] {
			return false
		}
	}
	return true
}

// columnsExist checks if all columns in 'requested' exist in 'schema'.
// Used for validating that SELECT columns exist in the table schema.
func ColumnsExist(schema, requested []string) bool {
	schemaSet := make(map[string]struct{}, len(schema))
	for _, col := range schema {
		schemaSet[col] = struct{}{}
	}
	for _, col := range requested {
		if _, ok := schemaSet[col]; !ok {
			return false
		}
	}
	return true
}

// IsSubset checks if all elements of 'subset' are in 'set'.
// Can be used for general column validation.
func IsSubset(set, subset []string) bool {
	setMap := make(map[string]struct{}, len(set))
	for _, v := range set {
		setMap[v] = struct{}{}
	}
	for _, v := range subset {
		if _, ok := setMap[v]; !ok {
			return false
		}
	}
	return true
}
