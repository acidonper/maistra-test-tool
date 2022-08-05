package performance

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func GenerateTrafficLoadK6(t *testing.T) {
	util.Log.Info("** TEST: GenerateTrafficLoadK6")
	_, err := generateSimpleTrafficLoadK6("http", "bookinfo")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
}
