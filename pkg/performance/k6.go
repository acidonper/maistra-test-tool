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

	var p95 string

	if res.Metrics.Checks.Fails > 0 {
		err = fmt.Errorf("There were %v fails in the tests", res.Metrics.Checks.Fails)
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	} else {
		p95 = fmt.Sprintf("%f", res.Metrics.HTTPReqReceiving.P95)
	}

	if p95 > reqAvg95pAcceptanceTime {
		util.Log.Error("Acceptance time exceeded")
	}
}
