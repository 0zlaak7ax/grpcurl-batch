// Package reporter provides JUnit XML report generation for grpcurl-batch results.
package reporter

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/nicholasgasior/grpcurl-batch/internal/formatter"
)

// JUnitTestSuites is the root XML element.
type JUnitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a collection of test cases.
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Timestamp string          `xml:"timestamp,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase represents a single gRPC call result.
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
}

// JUnitFailure holds failure details.
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

// WriteJUnit writes a JUnit XML report from collected results to w.
func WriteJUnit(w io.Writer, results []formatter.Result) error {
	suite := JUnitTestSuite{
		Name:      "grpcurl-batch",
		Tests:     len(results),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	for _, r := range results {
		tc := JUnitTestCase{
			Name:      r.Method,
			Classname: r.Address,
		}
		if !r.Success {
			suite.Failures++
			tc.Failure = &JUnitFailure{
				Message: fmt.Sprintf("gRPC call failed after %d attempt(s)", r.Attempts),
				Text:    r.Output,
			}
		}
		suite.TestCases = append(suite.TestCases, tc)
	}

	suites := JUnitTestSuites{Suites: []JUnitTestSuite{suite}}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(suites); err != nil {
		return fmt.Errorf("reporter: encoding XML: %w", err)
	}
	return enc.Flush()
}
