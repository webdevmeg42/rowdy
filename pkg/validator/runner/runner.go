package runner

import (
	"database/sql"

	"github.com/webdevmeg42/dbvalidator/pkg/validator/assertions"
	"github.com/webdevmeg42/dbvalidator/pkg/validator/config"
	"github.com/webdevmeg42/dbvalidator/pkg/validator/db"
)

type TestResult struct {
	Name     string
	Passed   bool
	Failures []string
}

func Run(conn *sql.DB, testCases []config.TestCase) []TestResult {
	results := make([]TestResult, len(testCases))
	for i, tc := range testCases {
		results[i] = runOne(conn, tc)
	}
	return results
}

func runOne(conn *sql.DB, tc config.TestCase) TestResult {
	result := TestResult{Name: tc.Name}

	if tc.Seed != nil {
		defer db.Cleanup(conn, tc.Seed.Table) //nolint:errcheck — always clean up, even if seed fails
		if err := db.Seed(conn, tc.Seed.Table, tc.Seed.Rows); err != nil {
			result.Failures = append(result.Failures, "seed error: "+err.Error())
			return result
		}
	}

	cols, rows, err := db.Query(conn, tc.Query)
	if err != nil {
		result.Failures = append(result.Failures, "query error: "+err.Error())
		return result
	}

	rs := assertions.ResultSet{Columns: cols, Rows: rows}
	for _, a := range tc.Assertions {
		if r := dispatch(rs, a); !r.Passed {
			result.Failures = append(result.Failures, r.Message)
		}
	}
	result.Passed = len(result.Failures) == 0
	return result
}

func dispatch(rs assertions.ResultSet, a config.Assertion) assertions.Result {
	switch a.Type {
	case "row_count":
		return assertions.CheckRowCount(rs, *a.Expected)
	case "column_exists":
		return assertions.CheckColumnExists(rs, a.Column)
	case "not_null":
		return assertions.CheckNotNull(rs, a.Column)
	case "value":
		return assertions.CheckValue(rs, a.Column, a.Row, a.Equals)
	case "format":
		return assertions.CheckFormat(rs, a.Column, a.Format)
	default:
		return assertions.Result{Message: "unknown assertion type: " + a.Type}
	}
}
