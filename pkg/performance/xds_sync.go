package performance

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestXDSSyncTime(t *testing.T) {
	util.Log.Info("** TEST: TestXDSSyncTime")
	metric, err := getMetricPrometheus("xds_ppctb")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	util.Log.Info(metric)
}
