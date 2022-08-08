package tests

import (
	"testing"
)

// this function is used for matching command line argument <test case name>,
// e.g. `go test -run <test case name>` with the names in the test_cases.go file.
func matchString(a, b string) (bool, error) {
	return a == b, nil
}

func TestMain(m *testing.M) {
	testing.Main(matchString, performanceCases, nil, nil)
}
