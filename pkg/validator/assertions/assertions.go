package assertions

// ResultSet holds the columns and rows returned by a SQL query.
type ResultSet struct {
	Columns []string
	Rows    []map[string]interface{}
}

// Result is the outcome of a single assertion check.
type Result struct {
	Passed  bool
	Message string // non-empty on failure
}
