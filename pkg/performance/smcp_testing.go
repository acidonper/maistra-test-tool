package performance

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestSMCPInstalled(t *testing.T) {
	util.Log.Info("** TEST: TestSMCPInstalled")
	TestSMCP(t)
}

func TestSMMRInstalled(t *testing.T) {
	util.Log.Info("** TEST: TestSMMRInstalled")
	TestSMMR(t)
}
