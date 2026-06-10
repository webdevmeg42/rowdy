package assertions

import (
	"fmt"
	"regexp"
)

type ResultSet struct {
	Columns []string
	Rows    []map[string]interface{}
}

type Result struct {
	Passed  bool
	Message string
}

func CheckRowCount(rs ResultSet, expected int) Result {
	got := len(rs.Rows)
	if got == expected {
		return Result{Passed: true}
	}
	return Result{Message: fmt.Sprintf("row_count: expected %d, got %d", expected, got)}
}

func CheckColumnExists(rs ResultSet, column string) Result {
	for _, col := range rs.Columns {
		if col == column {
			return Result{Passed: true}
		}
	}
	return Result{Message: fmt.Sprintf("column_exists: column %q not found in result set", column)}
}

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

var formatPatterns = map[string]*regexp.Regexp{
	"email": regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
	"uuid":  regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
	"date":  regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
	"url":   regexp.MustCompile(`^https?://`),
}

func CheckFormat(rs ResultSet, column, format string) Result {
	pattern, ok := formatPatterns[format]
	if !ok {
		return Result{Message: fmt.Sprintf("format: unknown format %q", format)}
	}
	for i, row := range rs.Rows {
		val, ok := row[column]
		if !ok {
			return Result{Message: fmt.Sprintf("format: column %q not found in result set", column)}
		}
		if val == nil {
			continue
		}
		s, ok := val.(string)
		if !ok {
			s = fmt.Sprintf("%v", val)
		}
		if !pattern.MatchString(s) {
			return Result{Message: fmt.Sprintf("format: column %q row %d: %q does not match format %q", column, i, s, format)}
		}
	}
	return Result{Passed: true}
}
