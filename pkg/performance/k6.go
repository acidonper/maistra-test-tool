package performance

import (
	"fmt"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func GenerateTrafficLoadK6(t *testing.T) {
	util.Log.Info("** TEST: GenerateTrafficLoadK6")
	err := generateSimpleTrafficLoadK6(trafficLoadProtocol, trafficLoadApp)
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
	util.Log.Info("OK: k6 load test executed")
}

func AnalyseLoadK6Output(t *testing.T) {
	util.Log.Info("** TEST: AnalyseLoadK6Output")
	dat, err := readK6File()
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}

	res, err := parseK6Response(dat)
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}

	result, err := checkFailedMetrics(res.Metrics.Checks.Fails)
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	} else {
		util.Log.Info(result)
	}

	var p95 float64
	if trafficLoadProtocol == "http" {
		p95 = res.Metrics.HTTPReqDuration.P95
	} else if trafficLoadProtocol == "grpc" {
		p95 = res.Metrics.GrpcReqDuration.P95
	} else {
		util.Log.Error("Protocol not valid")
		t.FailNow()
	}

	result, err = compareP95(fmt.Sprintf("%f", p95), reqAvg95pAcceptanceTime)
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}

	util.Log.Info(result)
}
