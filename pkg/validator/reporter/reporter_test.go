package reporter_test

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/webdevmeg42/rowdy/pkg/validator/reporter"
	"github.com/webdevmeg42/rowdy/pkg/validator/runner"
)

var mixedResults = []runner.TestResult{
	{Name: "passing test", Passed: true},
	{Name: "failing test", Passed: false, Failures: []string{"row_count: expected 2, got 1", "not_null: column \"name\" is NULL in row 0"}},
}

func TestTerminal_ContainsPassFail(t *testing.T) {
	var buf bytes.Buffer
	reporter.Terminal(&buf, mixedResults)
	out := buf.String()
	if !strings.Contains(out, "PASS") {
		t.Error("expected PASS in output")
	}
	if !strings.Contains(out, "FAIL") {
		t.Error("expected FAIL in output")
	}
}

func TestTerminal_ShowsFailureMessages(t *testing.T) {
	var buf bytes.Buffer
	reporter.Terminal(&buf, mixedResults)
	out := buf.String()
	if !strings.Contains(out, "row_count") {
		t.Error("expected row_count failure message in output")
	}
	if !strings.Contains(out, "not_null") {
		t.Error("expected not_null failure message in output")
	}
}

func TestTerminal_Summary(t *testing.T) {
	var buf bytes.Buffer
	reporter.Terminal(&buf, mixedResults)
	out := buf.String()
	if !strings.Contains(out, "1 passed") {
		t.Error("expected '1 passed' in summary")
	}
	if !strings.Contains(out, "1 failed") {
		t.Error("expected '1 failed' in summary")
	}
}

func TestTerminal_AllPass_Summary(t *testing.T) {
	var buf bytes.Buffer
	reporter.Terminal(&buf, []runner.TestResult{{Name: "t", Passed: true}})
	out := buf.String()
	if !strings.Contains(out, "1 passed") {
		t.Error("expected '1 passed' in summary")
	}
	if !strings.Contains(out, "0 failed") {
		t.Error("expected '0 failed' in summary")
	}
}

func TestJUnit_ValidXML(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.JUnit(&buf, mixedResults); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var v interface{}
	if err := xml.Unmarshal(buf.Bytes(), &v); err != nil {
		t.Errorf("output is not valid XML: %v\n%s", err, buf.String())
	}
}

func TestJUnit_ContainsTestsuitesRoot(t *testing.T) {
	var buf bytes.Buffer
	reporter.JUnit(&buf, mixedResults) //nolint:errcheck
	if !strings.Contains(buf.String(), "testsuites") {
		t.Error("expected testsuites root element")
	}
}

func TestJUnit_FailureElementPresent(t *testing.T) {
	var buf bytes.Buffer
	reporter.JUnit(&buf, mixedResults) //nolint:errcheck
	out := buf.String()
	if !strings.Contains(out, "failure") {
		t.Error("expected failure element for failing test")
	}
	if !strings.Contains(out, "row_count") {
		t.Error("expected failure message in JUnit output")
	}
}
