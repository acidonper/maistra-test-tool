package performance

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestXDSSyncTime(t *testing.T) {
	util.Log.Info("** TEST: TestXDSSyncTime")
	xdsPushCount, err := getMetricPrometheus("xds_ppctc")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	xdsPushTime, err := getMetricPrometheus("xds_ppctb")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	util.Log.Info(xdsPushCount)
	util.Log.Info(xdsPushTime)
	util.Log.Info(" If xdsPushCount and xdsPushTime are equal - OK")
}
