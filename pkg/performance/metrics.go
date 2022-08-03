package performance

import (
	"strconv"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestXDSPushes(t *testing.T) {
	util.Log.Info("** TEST: TestXDSPushes")
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
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	xdsPushTimeValue, err := parseResponse([]byte(xdsPushTime))
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	for i := 0; i < len(xdsPushCountValue); i++ {
		if xdsPushCountValue[i] != xdsPushTimeValue[i] {
			util.Log.Error(err)
			t.Errorf("xdsPushCount (%v) and xdsPushTime (%v) are not equal", xdsPushCountValue, xdsPushTimeValue)
			t.FailNow()
		} else {
			util.Log.Info("OK: ", xdsPushCountValue[i], "/", xdsPushTimeValue[i], " correct XDS pushes")
		}
	}
}

func TestIstiodMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodMem")
	istiodMem, err := getMetricPrometheusOCP("istiod_mem", nil)
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	istiodMemValue, err := parseResponse([]byte(istiodMem))
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	for i := 0; i < len(istiodMemValue); i++ {

		istiodMemValueInt, errConver1 := strconv.Atoi(istiodMemValue[i])
		if errConver1 != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		istiodMemValueIntMegaBytes := istiodMemValueInt / bytesToMegaBytes

		istiodAcceptanceMemIntMegabytes, errConver2 := strconv.Atoi(istiodAcceptanceMem)
		if errConver2 != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		if istiodMemValueIntMegaBytes > istiodAcceptanceMemIntMegabytes {
			t.Errorf("Istiod Memory Value is %v. Want something lower than %v", istiodMemValueIntMegaBytes, istiodAcceptanceMemIntMegabytes)
			t.FailNow()
		} else {
			util.Log.Info("OK: ", istiodMemValueIntMegaBytes, " is lower than ", istiodAcceptanceMemIntMegabytes, " in MBs")
		}
	}
}

func TestIstiodCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodCpu")
	istiodCpu, err := getMetricPrometheusOCP("istiod_cpu", nil)
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	istiodCpuValue, err := parseResponse([]byte(istiodCpu))
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	for i := 0; i < len(istiodCpuValue); i++ {

		istiodCpuValueFloat, errConver1 := strconv.ParseFloat(istiodCpuValue[0], 32)
		if errConver1 != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		istiodCpuValueFloatMilicores := istiodCpuValueFloat * coresToMilicores

		istiodAcceptanceCpuFloatMilicores, errConver2 := strconv.ParseFloat(istiodAcceptanceCpu, 32)
		if errConver2 != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		if istiodCpuValueFloatMilicores > istiodAcceptanceCpuFloatMilicores {
			t.Errorf("Istiod CPU Value is %v. Want something lower than %v", istiodCpuValueFloatMilicores, istiodAcceptanceCpuFloatMilicores)
			t.FailNow()
		} else {
			util.Log.Info("OK: ", istiodCpuValueFloatMilicores, " is lower than ", istiodAcceptanceCpuFloatMilicores, " in Milicores")
		}
	}
}

func TestIstioProxiesMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioProxiesMem")

	meshPods, err := getMeshProxies()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istiodMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		util.Log.Info(istiodMem)
		util.Log.Info(" If IstioProxiesMem is lower than ", istioProxiesAcceptanceMem)
	}
}

func TestIstioProxiesCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioProxiesCpu")

	meshPods, err := getMeshProxies()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istiodCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		util.Log.Info(istiodCpu)
		util.Log.Info(" If istioProxiesCpu is lower than ", istioProxiesAcceptanceCpu)
	}
}

func TestIstioIngressProxiesMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioIngressProxiesMem")

	meshPods, err := getMeshIngressProxies()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioIngressProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		util.Log.Info(istioIngressProxyMem)
		util.Log.Info(" If istioIngressProxiesAcceptanceMem is lower than ", istioIngressProxiesAcceptanceMem)
	}
}

func TestIstioIngressProxiesCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioIngressProxiesCpu")

	meshPods, err := getMeshIngressProxies()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioIngressProxyCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		util.Log.Info(istioIngressProxyCpu)
		util.Log.Info(" If istioIngressProxiesAcceptanceCpu is lower than ", istioIngressProxiesAcceptanceCpu)
	}
}

func TestIstioEgressProxiesMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioEgressProxiesMem")

	meshPods, err := getMeshEgressProxies()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioEgressProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		util.Log.Info(istioEgressProxyMem)
		util.Log.Info(" If istioEgressProxiesAcceptanceMem is lower than ", istioEgressProxiesAcceptanceMem)
	}
}

func TestIstioEgressProxiesCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioEgressProxiesCpu")

	meshPods, err := getMeshEgressProxies()
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioEgressProxyCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}
		util.Log.Info(istioEgressProxyCpu)
		util.Log.Info(" If istioEgressProxiesAcceptanceCpu is lower than ", istioEgressProxiesAcceptanceCpu)
	}
}
