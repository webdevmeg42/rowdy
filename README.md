# Rowdy

A Go CLI tool that validates SQL queries against expected results using YAML-defined test cases. Seed fixture data, run queries, assert on row counts, column shapes, null values, specific values, and data formats — all from a single config file.

## Why I built this

At my previous org, query correctness in our data pipeline was validated manually or not at all. I wanted a lightweight, config-driven tool that could catch regressions — something you could drop into any project without writing Go code. Most SDET tooling focuses on APIs; this fills the gap for DB-heavy workflows.

## What it does

```
PASS  active users query returns one row
PASS  all users have valid emails
PASS  uuid format validation

3 test cases — 3 passed, 0 failed
```

Define test cases in YAML, point the tool at a SQLite database, get clear pass/fail output. Integrates with CI via JUnit XML output.

## Setup

**Requirements:** Go 1.22+

```bash
# Build from source
git clone https://github.com/webdevmeg42/rowdy
cd rowdy
go build -o rowdy ./cmd/rowdy

# Or install directly
go install github.com/webdevmeg42/rowdy/cmd/rowdy@latest
```

## Usage

```bash
# Run with default terminal output
./rowdy --config queries.yaml

# Run with JUnit XML output (for CI)
./rowdy --config queries.yaml --format junit > results.xml

# Override database path from config
./rowdy --config queries.yaml --db ./path/to/test.db
```

Exit codes: `0` = all pass, `1` = assertion failures, `2` = config/DB error.

## Config reference

```yaml
database:
  path: ":memory:"   # or a file path

test_cases:
  - name: "descriptive test name"
    seed:                          # optional: insert rows before query
      table: users
      rows:
        - { id: 1, name: "Alice", age: 30, email: "alice@example.com" }
    query: "SELECT * FROM users WHERE id = 1"
    assertions:
      - type: row_count
        expected: 1

      - type: column_exists
        column: email

      - type: not_null
        column: name

      - type: value
        column: name
        row: 0          # zero-indexed
        equals: "Alice"

      - type: value
        column: age
        row: 0
        equals: 30      # integer comparison

      - type: format
        column: email
        format: email   # email | uuid | date | url
```

### Assertion types

| Type | Fields | Description |
|---|---|---|
| `row_count` | `expected` | Result must contain exactly N rows |
| `column_exists` | `column` | Column must appear in result set |
| `not_null` | `column` | Column must have no NULL values |
| `value` | `column`, `row`, `equals` | Cell must equal value (type-aware) |
| `format` | `column`, `format` | Values must match named format pattern |

### Format patterns

| Format | Validates |
|---|---|
| `email` | RFC-style email address |
| `uuid` | Lowercase UUID v4 (`xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`) |
| `date` | ISO 8601 date (`YYYY-MM-DD`) |
| `url` | HTTP or HTTPS URL |

### Column types

Seed column types are inferred from YAML: integers → `INTEGER`, floats → `REAL`, booleans → `INTEGER`, everything else → `TEXT`.

## CI integration

The tool validates itself on every push using GitHub Actions. The `validate` job builds the binary, runs it against `testdata/sample.yaml`, and uploads the JUnit XML as a test report artifact.

[![CI](https://github.com/webdevmeg42/rowdy/actions/workflows/ci.yml/badge.svg)](https://github.com/webdevmeg42/rowdy/actions/workflows/ci.yml)

## Running the tests

```bash
# Run all tests
go test ./...

# Force re-run, bypassing cache
go test ./... -count=1

# Verbose output (shows each test name and pass/fail)
go test ./... -v

# Run a specific package
go test ./pkg/validator/assertions/...

# Run a specific test by name
go test ./... -run TestRowCount_Pass

# Run with race detector
go test ./... -race
```
