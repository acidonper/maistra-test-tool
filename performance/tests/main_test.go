package tests

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
	"github.com/sirupsen/logrus"
)

// this function is used for matching command line argument <test case name>,
// e.g. `go test -run <test case name>` with the names in the test_cases.go file.
func matchString(a, b string) (bool, error) {
	return a == b, nil
}

func TestMain(m *testing.M) {
	util.Log.SetLevel(logrus.InfoLevel) // To be able to change logging level. Leave InfoLevel by default
	testing.Main(matchString, performanceCases, nil, nil)
}
