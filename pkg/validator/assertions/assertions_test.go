package assertions_test

import (
	"testing"

	"github.com/meganwall/dbvalidator/pkg/validator/assertions"
)

func makeRS(cols []string, rows []map[string]interface{}) assertions.ResultSet {
	return assertions.ResultSet{Columns: cols, Rows: rows}
}

func TestRowCount_Pass(t *testing.T) {
	rs := makeRS([]string{"id"}, []map[string]interface{}{{"id": int64(1)}, {"id": int64(2)}})
	r := assertions.CheckRowCount(rs, 2)
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestRowCount_Fail(t *testing.T) {
	rs := makeRS([]string{"id"}, []map[string]interface{}{{"id": int64(1)}})
	r := assertions.CheckRowCount(rs, 3)
	if r.Passed {
		t.Error("expected fail")
	}
	if r.Message == "" {
		t.Error("expected failure message")
	}
}

func TestRowCount_Empty(t *testing.T) {
	rs := makeRS([]string{"id"}, nil)
	r := assertions.CheckRowCount(rs, 0)
	if !r.Passed {
		t.Errorf("expected pass for empty result with expected=0: %s", r.Message)
	}
}

func TestColumnExists_Pass(t *testing.T) {
	rs := makeRS([]string{"id", "name", "email"}, nil)
	r := assertions.CheckColumnExists(rs, "name")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestColumnExists_Fail(t *testing.T) {
	rs := makeRS([]string{"id", "name"}, nil)
	r := assertions.CheckColumnExists(rs, "missing_col")
	if r.Passed {
		t.Error("expected fail for missing column")
	}
}

func TestNotNull_Pass(t *testing.T) {
	rs := makeRS([]string{"name"}, []map[string]interface{}{
		{"name": "Alice"},
		{"name": "Bob"},
	})
	r := assertions.CheckNotNull(rs, "name")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestNotNull_Fail_NullValue(t *testing.T) {
	rs := makeRS([]string{"name"}, []map[string]interface{}{
		{"name": "Alice"},
		{"name": nil},
	})
	r := assertions.CheckNotNull(rs, "name")
	if r.Passed {
		t.Error("expected fail: row 1 is NULL")
	}
}

func TestNotNull_MissingColumn(t *testing.T) {
	rs := makeRS([]string{"id"}, []map[string]interface{}{{"id": int64(1)}})
	r := assertions.CheckNotNull(rs, "nonexistent")
	if r.Passed {
		t.Error("expected fail for missing column")
	}
}

func TestValue_StringMatch(t *testing.T) {
	rs := makeRS([]string{"name"}, []map[string]interface{}{{"name": "Alice"}})
	r := assertions.CheckValue(rs, "name", 0, "Alice")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestValue_IntMatch(t *testing.T) {
	rs := makeRS([]string{"age"}, []map[string]interface{}{{"age": int64(30)}})
	r := assertions.CheckValue(rs, "age", 0, 30)
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestValue_IntMismatch(t *testing.T) {
	rs := makeRS([]string{"age"}, []map[string]interface{}{{"age": int64(30)}})
	r := assertions.CheckValue(rs, "age", 0, 99)
	if r.Passed {
		t.Error("expected fail for wrong int value")
	}
}

func TestValue_RowOutOfRange(t *testing.T) {
	rs := makeRS([]string{"id"}, []map[string]interface{}{{"id": int64(1)}})
	r := assertions.CheckValue(rs, "id", 5, int64(1))
	if r.Passed {
		t.Error("expected fail for out-of-range row index")
	}
}

func TestValue_FloatMatch(t *testing.T) {
	rs := makeRS([]string{"score"}, []map[string]interface{}{{"score": float64(9.5)}})
	r := assertions.CheckValue(rs, "score", 0, 9.5)
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestFormat_Email_Pass(t *testing.T) {
	rs := makeRS([]string{"email"}, []map[string]interface{}{
		{"email": "alice@example.com"},
		{"email": "bob.smith+tag@sub.domain.org"},
	})
	r := assertions.CheckFormat(rs, "email", "email")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestFormat_Email_Fail(t *testing.T) {
	rs := makeRS([]string{"email"}, []map[string]interface{}{
		{"email": "not-an-email"},
	})
	r := assertions.CheckFormat(rs, "email", "email")
	if r.Passed {
		t.Error("expected fail for invalid email")
	}
}

func TestFormat_UUID_Pass(t *testing.T) {
	rs := makeRS([]string{"id"}, []map[string]interface{}{
		{"id": "550e8400-e29b-41d4-a716-446655440000"},
	})
	r := assertions.CheckFormat(rs, "id", "uuid")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestFormat_UUID_Fail(t *testing.T) {
	rs := makeRS([]string{"id"}, []map[string]interface{}{
		{"id": "not-a-uuid"},
	})
	r := assertions.CheckFormat(rs, "id", "uuid")
	if r.Passed {
		t.Error("expected fail for invalid UUID")
	}
}

func TestFormat_Date_Pass(t *testing.T) {
	rs := makeRS([]string{"created_at"}, []map[string]interface{}{
		{"created_at": "2024-01-15"},
	})
	r := assertions.CheckFormat(rs, "created_at", "date")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestFormat_URL_Pass(t *testing.T) {
	rs := makeRS([]string{"link"}, []map[string]interface{}{
		{"link": "https://example.com/path?query=1"},
	})
	r := assertions.CheckFormat(rs, "link", "url")
	if !r.Passed {
		t.Errorf("expected pass: %s", r.Message)
	}
}

func TestFormat_SkipsNulls(t *testing.T) {
	rs := makeRS([]string{"email"}, []map[string]interface{}{
		{"email": "alice@example.com"},
		{"email": nil},
	})
	r := assertions.CheckFormat(rs, "email", "email")
	if !r.Passed {
		t.Errorf("expected pass: NULLs should be skipped: %s", r.Message)
	}
}
