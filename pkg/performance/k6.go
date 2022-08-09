package performance

import (
	"fmt"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func GenerateTrafficLoadK6(t *testing.T) {
	util.Log.Info("** TEST: GenerateTrafficLoadK6")
	err := generateSimpleTrafficLoadK6("http", "bookinfo")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	dat, err := readK6File()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	res, err := parseK6Response(dat)
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	result, err := checkFailedMetrics(res.Metrics.Checks.Fails)
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	} else {
		util.Log.Info(result)
	}

	p95 := res.Metrics.HTTPReqReceiving.P95
	result, err = compareP95(fmt.Sprintf("%f", p95), reqAvg95pAcceptanceTime)
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	} else {
		util.Log.Info(result)
	}
}
