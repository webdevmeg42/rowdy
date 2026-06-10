package db_test

import (
	"testing"

	"github.com/webdevmeg42/dbvalidator/pkg/validator/db"
)

func TestOpen_Memory(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	conn.Close()
}

func TestSeedAndQuery_TypesPreserved(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	rows := []map[string]interface{}{
		{"id": 1, "name": "Alice", "score": 9.5},
		{"id": 2, "name": "Bob", "score": 7.0},
	}
	if err := db.Seed(conn, "players", rows); err != nil {
		t.Fatalf("seed: %v", err)
	}

	cols, results, err := db.Query(conn, "SELECT * FROM players ORDER BY id")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 rows, got %d", len(results))
	}
	if len(cols) == 0 {
		t.Error("expected columns to be populated")
	}

	// id was seeded as int — SQLite returns int64
	if id, ok := results[0]["id"].(int64); !ok || id != 1 {
		t.Errorf("expected id=1 as int64, got %v (%T)", results[0]["id"], results[0]["id"])
	}
	// name was seeded as string — SQLite returns string
	if name, ok := results[0]["name"].(string); !ok || name != "Alice" {
		t.Errorf("expected name=Alice, got %v (%T)", results[0]["name"], results[0]["name"])
	}
}

func TestSeed_EmptyRows(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if err := db.Seed(conn, "empty", []map[string]interface{}{}); err != nil {
		t.Errorf("unexpected error for empty rows: %v", err)
	}
}

func TestCleanup_DropsTable(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	db.Seed(conn, "tmp", []map[string]interface{}{{"id": 1}})
	if err := db.Cleanup(conn, "tmp"); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	_, _, err = db.Query(conn, "SELECT * FROM tmp")
	if err == nil {
		t.Error("expected error querying dropped table")
	}
}

func TestQuery_SQLError(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	_, _, err = db.Query(conn, "SELECT * FROM nonexistent_table")
	if err == nil {
		t.Error("expected error for missing table")
	}
}
