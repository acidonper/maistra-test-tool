package performance

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestXDSSyncTime(t *testing.T) {
	util.Log.Info("** TEST: TestXDSSyncTime")
	xdsPushCount, err := getMetricPrometheusMesh("xds_ppctc")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	xdsPushTime, err := getMetricPrometheusMesh("xds_ppctb")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	util.Log.Info(xdsPushCount)
	util.Log.Info(xdsPushTime)
	util.Log.Info(" If xdsPushCount and xdsPushTime are equal - OK")
}

func TestIstiodMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodMem")
	istiodMem, err := getMetricPrometheusOCP("istiod_mem")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	util.Log.Info(istiodMem)
	util.Log.Info(" If istiodMem is lower than ", istiodAcceptanceMem)
}

func TestIstiodCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodCpu")
	istiodCpu, err := getMetricPrometheusOCP("istiod_cpu")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	util.Log.Info(istiodCpu)
	util.Log.Info(" If istiodCpu is lower than ", istiodAcceptanceCpu)
}
