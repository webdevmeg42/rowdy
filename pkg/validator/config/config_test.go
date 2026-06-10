package config_test

import (
	"os"
	"testing"

	"github.com/webdevmeg42/rowdy/pkg/validator/config"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases:
  - name: "basic test"
    query: "SELECT 1"
    assertions:
      - type: row_count
        expected: 1
`
	cfg, err := config.Load(writeTempFile(t, yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.TestCases) != 1 {
		t.Errorf("expected 1 test case, got %d", len(cfg.TestCases))
	}
	if cfg.TestCases[0].Assertions[0].Type != "row_count" {
		t.Errorf("unexpected assertion type: %s", cfg.TestCases[0].Assertions[0].Type)
	}
}

func TestLoad_MissingDatabasePath(t *testing.T) {
	yaml := `
test_cases:
  - name: "t"
    query: "SELECT 1"
    assertions: []
`
	_, err := config.Load(writeTempFile(t, yaml))
	if err == nil {
		t.Fatal("expected error for missing database.path")
	}
}

func TestLoad_MissingQuery(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases:
  - name: "t"
    assertions: []
`
	_, err := config.Load(writeTempFile(t, yaml))
	if err == nil {
		t.Fatal("expected error for missing query")
	}
}

func TestLoad_UnknownAssertionType(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases:
  - name: "t"
    query: "SELECT 1"
    assertions:
      - type: bogus
`
	_, err := config.Load(writeTempFile(t, yaml))
	if err == nil {
		t.Fatal("expected error for unknown assertion type")
	}
}

func TestLoad_UnknownFormat(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases:
  - name: "t"
    query: "SELECT 1"
    assertions:
      - type: format
        column: x
        format: fax_number
`
	_, err := config.Load(writeTempFile(t, yaml))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestLoad_RowCountMissingExpected(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases:
  - name: "t"
    query: "SELECT 1"
    assertions:
      - type: row_count
`
	_, err := config.Load(writeTempFile(t, yaml))
	if err == nil {
		t.Fatal("expected error for row_count missing expected")
	}
}

func TestLoad_ColumnExistsMissingColumn(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases:
  - name: "t"
    query: "SELECT 1"
    assertions:
      - type: column_exists
`
	_, err := config.Load(writeTempFile(t, yaml))
	if err == nil {
		t.Fatal("expected error for column_exists missing column")
	}
}

func TestLoad_PathParsed(t *testing.T) {
	yaml := `
database:
  path: ":memory:"
test_cases: []
`
	cfg, err := config.Load(writeTempFile(t, yaml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Database.Path != ":memory:" {
		t.Errorf("unexpected path: %s", cfg.Database.Path)
	}
}
