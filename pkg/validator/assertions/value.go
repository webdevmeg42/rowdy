package assertions

import "fmt"

func CheckValue(rs ResultSet, column string, rowIdx int, equals interface{}) Result {
	if rowIdx >= len(rs.Rows) {
		return Result{Message: fmt.Sprintf("value: row index %d out of range (result has %d rows)", rowIdx, len(rs.Rows))}
	}
	row := rs.Rows[rowIdx]
	actual, ok := row[column]
	if !ok {
		return Result{Message: fmt.Sprintf("value: column %q not found in result set", column)}
	}
	if valuesEqual(actual, equals) {
		return Result{Passed: true}
	}
	return Result{Message: fmt.Sprintf("value: column %q row %d: expected %v (%T), got %v (%T)",
		column, rowIdx, equals, equals, actual, actual)}
}

func valuesEqual(dbVal, yamlVal interface{}) bool {
	switch y := yamlVal.(type) {
	case int:
		switch d := dbVal.(type) {
		case int64:
			return d == int64(y)
		case float64:
			return d == float64(y)
		}
	case float64:
		switch d := dbVal.(type) {
		case float64:
			return d == y
		case int64:
			return float64(d) == y
		}
	case string:
		if d, ok := dbVal.(string); ok {
			return d == y
		}
		return fmt.Sprintf("%v", dbVal) == y
	}
	return fmt.Sprintf("%v", dbVal) == fmt.Sprintf("%v", yamlVal)
}
