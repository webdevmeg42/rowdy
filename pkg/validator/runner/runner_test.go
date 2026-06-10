package runner_test

import (
	"testing"

	"github.com/webdevmeg42/rowdy/pkg/validator/config"
	"github.com/webdevmeg42/rowdy/pkg/validator/db"
	"github.com/webdevmeg42/rowdy/pkg/validator/runner"
)

func intPtr(i int) *int { return &i }

func TestRun_PassingCase(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tc := config.TestCase{
		Name:  "count test",
		Query: "SELECT 1 as n",
		Assertions: []config.Assertion{
			{Type: "row_count", Expected: intPtr(1)},
		},
	}
	results := runner.Run(conn, []config.TestCase{tc})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Passed {
		t.Errorf("expected pass, failures: %v", results[0].Failures)
	}
}

func TestRun_FailingAssertion(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tc := config.TestCase{
		Name:  "wrong count",
		Query: "SELECT 1",
		Assertions: []config.Assertion{
			{Type: "row_count", Expected: intPtr(5)},
		},
	}
	results := runner.Run(conn, []config.TestCase{tc})
	if results[0].Passed {
		t.Error("expected fail")
	}
	if len(results[0].Failures) == 0 {
		t.Error("expected at least one failure message")
	}
}

func TestRun_SQLError_SurfacedAsFailure(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tc := config.TestCase{
		Name:       "bad sql",
		Query:      "SELECT * FROM table_that_does_not_exist",
		Assertions: []config.Assertion{{Type: "row_count", Expected: intPtr(0)}},
	}
	results := runner.Run(conn, []config.TestCase{tc})
	if results[0].Passed {
		t.Error("expected fail for SQL error")
	}
}

func TestRun_SeedAndAssertions_AllPass(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tc := config.TestCase{
		Name: "users seed test",
		Seed: &config.SeedConfig{
			Table: "users",
			Rows: []map[string]interface{}{
				{"id": 1, "name": "Alice", "age": 30, "email": "alice@example.com"},
				{"id": 2, "name": "Bob", "age": 25, "email": "bob@example.com"},
			},
		},
		Query: "SELECT * FROM users WHERE name = 'Alice'",
		Assertions: []config.Assertion{
			{Type: "row_count", Expected: intPtr(1)},
			{Type: "column_exists", Column: "email"},
			{Type: "not_null", Column: "name"},
			{Type: "value", Column: "name", Row: 0, Equals: "Alice"},
			{Type: "value", Column: "age", Row: 0, Equals: 30},
			{Type: "format", Column: "email", Format: "email"},
		},
	}
	results := runner.Run(conn, []config.TestCase{tc})
	if !results[0].Passed {
		t.Errorf("expected all assertions to pass, failures: %v", results[0].Failures)
	}
}

func TestRun_SeedCleanedUpAfterFailure(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	seed := &config.SeedConfig{
		Table: "cleanup_test",
		Rows:  []map[string]interface{}{{"id": 1}},
	}

	// First run — will fail assertions but seed should still be cleaned up
	tc1 := config.TestCase{
		Name:       "failing case",
		Seed:       seed,
		Query:      "SELECT * FROM cleanup_test",
		Assertions: []config.Assertion{{Type: "row_count", Expected: intPtr(99)}},
	}
	runner.Run(conn, []config.TestCase{tc1})

	// Second run using same table name — should succeed (table was cleaned up)
	tc2 := config.TestCase{
		Name:       "re-seed same table",
		Seed:       seed,
		Query:      "SELECT * FROM cleanup_test",
		Assertions: []config.Assertion{{Type: "row_count", Expected: intPtr(1)}},
	}
	results := runner.Run(conn, []config.TestCase{tc2})
	if !results[0].Passed {
		t.Errorf("expected pass after re-seed: %v", results[0].Failures)
	}
}

func TestRun_AllCasesRunDespiteFailure(t *testing.T) {
	conn, err := db.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	cases := []config.TestCase{
		{Name: "fail", Query: "SELECT 1", Assertions: []config.Assertion{{Type: "row_count", Expected: intPtr(99)}}},
		{Name: "pass", Query: "SELECT 1", Assertions: []config.Assertion{{Type: "row_count", Expected: intPtr(1)}}},
	}
	results := runner.Run(conn, cases)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Passed {
		t.Error("expected first case to fail")
	}
	if !results[1].Passed {
		t.Errorf("expected second case to pass: %v", results[1].Failures)
	}
}
