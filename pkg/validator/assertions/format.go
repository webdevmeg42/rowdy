package assertions

import (
	"fmt"
	"regexp"
)

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
