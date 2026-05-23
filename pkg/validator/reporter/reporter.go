package reporter

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/meganwall/dbvalidator/pkg/validator/runner"
)

func Terminal(w io.Writer, results []runner.TestResult) {
	useColor := isTerminal(w)
	pass := color.New(color.FgGreen, color.Bold)
	fail := color.New(color.FgRed, color.Bold)
	if !useColor {
		pass.DisableColor()
		fail.DisableColor()
	}

	passed, failed := 0, 0
	for _, r := range results {
		if r.Passed {
			fmt.Fprintf(w, "%s  %s\n", pass.Sprint("PASS"), r.Name)
			passed++
		} else {
			fmt.Fprintf(w, "%s  %s\n", fail.Sprint("FAIL"), r.Name)
			for _, f := range r.Failures {
				fmt.Fprintf(w, "      %s\n", f)
			}
			failed++
		}
	}
	fmt.Fprintf(w, "\n%d test cases — %d passed, %d failed\n", len(results), passed, failed)
}

type junitTestSuites struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []junitSuite `xml:"testsuite"`
}

type junitSuite struct {
	Name      string      `xml:"name,attr"`
	Tests     int         `xml:"tests,attr"`
	Failures  int         `xml:"failures,attr"`
	TestCases []junitCase `xml:"testcase"`
}

type junitCase struct {
	Name    string        `xml:"name,attr"`
	Failure *junitFailure `xml:"failure,omitempty"`
}

type junitFailure struct {
	Message string `xml:"message,attr"`
	Body    string `xml:",chardata"`
}

func JUnit(w io.Writer, results []runner.TestResult) error {
	var cases []junitCase
	failures := 0
	for _, r := range results {
		tc := junitCase{Name: r.Name}
		if !r.Passed && len(r.Failures) > 0 {
			failures++
			tc.Failure = &junitFailure{
				Message: r.Failures[0],
				Body:    strings.Join(r.Failures, "\n"),
			}
		} else if !r.Passed {
			failures++
			tc.Failure = &junitFailure{Message: "test failed"}
		}
		cases = append(cases, tc)
	}
	suites := junitTestSuites{
		Suites: []junitSuite{{
			Name:      "dbvalidator",
			Tests:     len(results),
			Failures:  failures,
			TestCases: cases,
		}},
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(suites); err != nil {
		return fmt.Errorf("encoding JUnit XML: %w", err)
	}
	return enc.Flush()
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
