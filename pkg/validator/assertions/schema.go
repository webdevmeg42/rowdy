package assertions

import "fmt"

func CheckColumnExists(rs ResultSet, column string) Result {
	for _, col := range rs.Columns {
		if col == column {
			return Result{Passed: true}
		}
	}
	return Result{Message: fmt.Sprintf("column_exists: column %q not found in result set", column)}
}
