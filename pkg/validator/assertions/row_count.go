package assertions

import "fmt"

func CheckRowCount(rs ResultSet, expected int) Result {
	got := len(rs.Rows)
	if got == expected {
		return Result{Passed: true}
	}
	return Result{Message: fmt.Sprintf("row_count: expected %d, got %d", expected, got)}
}
