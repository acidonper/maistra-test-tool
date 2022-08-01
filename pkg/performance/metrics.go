package performance

import (
	"strconv"
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

	xdsPushCountValue, err := parseResponse([]byte(xdsPushCount))
	xdsPushTimeValue, err := parseResponse([]byte(xdsPushTime))

	if xdsPushCountValue[0] != xdsPushTimeValue[0] {
		t.Errorf("xdsPushCount (%v) and xdsPushTime (%v) are not equal", xdsPushCountValue, xdsPushTimeValue)
	}
	util.Log.Info(" If xdsPushCount and xdsPushTime are equal - OK")
}

func TestIstiodMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodMem")
	istiodMem, err := getMetricPrometheusOCP("istiod_mem")
	istiodMemValue, err := parseResponse([]byte(istiodMem))
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	// Transform values to integers and compare them in bytes
	istiodMemValueInt, err := strconv.Atoi(istiodMemValue[0])
	istiodAcceptanceMemIntBytes, err := strconv.Atoi(istiodAcceptanceMem)
	istiodAcceptanceMemIntBytes = istiodAcceptanceMemIntBytes * bytesToMegaBytes

	util.Log.Info(" If istiodMem is lower than ", istiodAcceptanceMem, "MB")
	if istiodMemValueInt > istiodAcceptanceMemIntBytes {
		t.Errorf("Istiod Memory Value is %v. Want something lower than %v", istiodMemValueInt, istiodAcceptanceMemIntBytes)
	}
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
