package db

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	return conn, nil
}

func Seed(conn *sql.DB, table string, rows []map[string]interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	cols := sortedKeys(rows[0])
	defs := make([]string, len(cols))
	for i, col := range cols {
		defs[i] = fmt.Sprintf(`"%s" %s`, col, sqliteType(rows[0][col]))
	}
	_, err := conn.Exec(fmt.Sprintf(`CREATE TABLE "%s" (%s)`, table, strings.Join(defs, ", ")))
	if err != nil {
		return fmt.Errorf("creating table %q: %w", table, err)
	}

	placeholders := make([]string, len(cols))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	quotedCols := make([]string, len(cols))
	for i, c := range cols {
		quotedCols[i] = `"` + c + `"`
	}
	stmt := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`,
		table,
		strings.Join(quotedCols, ", "),
		strings.Join(placeholders, ", "),
	)
	for _, row := range rows {
		vals := make([]interface{}, len(cols))
		for i, col := range cols {
			vals[i] = row[col]
		}
		if _, err := conn.Exec(stmt, vals...); err != nil {
			return fmt.Errorf("inserting row into %q: %w", table, err)
		}
	}
	return nil
}

func Cleanup(conn *sql.DB, table string) error {
	_, err := conn.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, table))
	return err
}

func Query(conn *sql.DB, query string) ([]string, []map[string]interface{}, error) {
	rows, err := conn.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range ptrs {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		row := make(map[string]interface{}, len(cols))
		for i, col := range cols {
			row[col] = vals[i]
		}
		results = append(results, row)
	}
	return cols, results, rows.Err()
}

func sqliteType(v interface{}) string {
	switch v.(type) {
	case int, int64:
		return "INTEGER"
	case float64:
		return "REAL"
	case bool:
		return "INTEGER"
	default:
		return "TEXT"
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
