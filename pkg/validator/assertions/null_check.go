package assertions

import "fmt"

func CheckNotNull(rs ResultSet, column string) Result {
	for i, row := range rs.Rows {
		val, ok := row[column]
		if !ok {
			return Result{Message: fmt.Sprintf("not_null: column %q not found in result set", column)}
		}
		if val == nil {
			return Result{Message: fmt.Sprintf("not_null: column %q is NULL in row %d", column, i)}
		}
	}
	return Result{Passed: true}
}
