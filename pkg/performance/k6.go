package performance

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func GenerateTrafficLoadK6(t *testing.T) {
	util.Log.Info("** TEST: GenerateTrafficLoadK6")
	pid, err := generateSimpleTrafficLoadK6("http", "bookinfo")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	if pid > reqAvg95pAcceptanceTime {
		util.Log.Error("Acceptance time exceeded")
	}
}
