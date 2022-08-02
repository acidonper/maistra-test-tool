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
	xdsPushTimeValue, err := parseResponse([]byte(xdsPushTime))

	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	util.Log.Info(" If xdsPushCount and xdsPushTime are equal - OK")
	if xdsPushCountValue[0] != xdsPushTimeValue[0] {
		t.Errorf("xdsPushCount (%v) and xdsPushTime (%v) are not equal", xdsPushCountValue, xdsPushTimeValue)
	}
}

func TestIstiodMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodMem")
	istiodMem, err := getMetricPrometheusOCP("istiod_mem", nil)
	istiodMemValue, err := parseResponse([]byte(istiodMem))
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	// Transform values to integers and compare them in bytes
	istiodMemValueInt, err := strconv.Atoi(istiodMemValue[0])
	istiodAcceptanceMemIntBytes, err := strconv.Atoi(istiodAcceptanceMem)

	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	istiodAcceptanceMemIntBytes = istiodAcceptanceMemIntBytes * bytesToMegaBytes

	util.Log.Info(" If istiodMem is lower than ", istiodAcceptanceMem, "MB")
	if istiodMemValueInt > istiodAcceptanceMemIntBytes {
		t.Errorf("Istiod Memory Value is %v. Want something lower than %v", istiodMemValueInt, istiodAcceptanceMemIntBytes)
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

	istiodAcceptanceCpuFloat, err := strconv.ParseFloat(istiodAcceptanceCpu, 32)
	istiodCpuValueFloat, err := strconv.ParseFloat(istiodCpuValue[0], 32)

	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	util.Log.Info(" If istiodCpu is lower than ", istiodAcceptanceCpu)
	if istiodCpuValueFloat > istiodAcceptanceCpuFloat {
		t.Errorf("Istiod CPU Value is %v. Want something lower than %v", istiodCpuValueFloat, istiodAcceptanceCpuFloat)
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

		istiodMemValue, err := parseResponse([]byte(istiodMem))
		istiodMemValueInt, err := strconv.Atoi(istiodMemValue[0])
		istioProxiesAcceptanceMemInt, err := strconv.Atoi(istioProxiesAcceptanceMem)
		istioProxiesAcceptanceMemBytes := istioProxiesAcceptanceMemInt * bytesToMegaBytes

		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		util.Log.Info(" If IstioProxiesMem is lower than ", istioProxiesAcceptanceMem, "MB")
		if istiodMemValueInt > istioProxiesAcceptanceMemBytes {
			t.Errorf("Proxy Memory value (%v) is higher than the acceptance (in bytes): %v", istiodMemValueInt, istioProxiesAcceptanceMemBytes)
		}

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

		istiodCpuValue, err := parseResponse([]byte(istiodCpu))
		istiodCpuValueFloat, err := strconv.ParseFloat(istiodCpuValue[0], 32)
		istioProxiesAcceptanceCpuFloat, err := strconv.ParseFloat(istioProxiesAcceptanceCpu, 32)

		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		util.Log.Info(" If istioProxiesCpu is lower than ", istioProxiesAcceptanceCpu)
		if istiodCpuValueFloat > istioProxiesAcceptanceCpuFloat {
			t.Errorf("Istiod CPU Proxy Value is %v. Want something lower than %v", istiodCpuValueFloat, istioProxiesAcceptanceCpuFloat)
		}

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
		istioIngressProxyMemValue, err := parseResponse([]byte(istioIngressProxyMem))
		istioIngressProxyMemValueInt, err := strconv.Atoi(istioIngressProxyMemValue[0])
		istioIngressProxiesAcceptanceMemInt, err := strconv.Atoi(istioIngressProxiesAcceptanceMem)
		istioIngressProxiesAcceptanceMemBytes := istioIngressProxiesAcceptanceMemInt * bytesToMegaBytes

		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		util.Log.Info(" If istioIngressProxiesAcceptanceMem is lower than ", istioIngressProxiesAcceptanceMem, "MB")
		if istioIngressProxyMemValueInt > istioIngressProxiesAcceptanceMemBytes {
			t.Errorf("Proxy Memory value (%v) is higher than the acceptance (in bytes): %v", istioIngressProxyMemValueInt, istioIngressProxiesAcceptanceMemBytes)
		}

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
		istioIngressProxyCpuValue, err := parseResponse([]byte(istioIngressProxyCpu))
		istioIngressProxyCpuValueFloat, err := strconv.ParseFloat(istioIngressProxyCpuValue[0], 32)
		istioIngressProxiesAcceptanceCpuFloat, err := strconv.ParseFloat(istioIngressProxiesAcceptanceCpu, 32)

		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		util.Log.Info(" If istioIngressProxiesAcceptanceCpu is lower than ", istioIngressProxiesAcceptanceCpu)
		if istioIngressProxyCpuValueFloat > istioIngressProxiesAcceptanceCpuFloat {
			t.Errorf("Istiod CPU Proxy Value is %v. Want something lower than %v", istioIngressProxyCpuValueFloat, istioIngressProxiesAcceptanceCpuFloat)
		}

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
		istioEgressProxyMemValue, err := parseResponse([]byte(istioEgressProxyMem))
		istioEgressProxyMemValueInt, err := strconv.Atoi(istioEgressProxyMemValue[0])
		istioEgressProxiesAcceptanceMemInt, err := strconv.Atoi(istioEgressProxiesAcceptanceMem)
		istioEgressProxiesAcceptanceMemBytes := istioEgressProxiesAcceptanceMemInt * bytesToMegaBytes

		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		util.Log.Info(" If istioEgressProxiesAcceptanceMem is lower than ", istioEgressProxiesAcceptanceMem, "MB")
		if istioEgressProxyMemValueInt > istioEgressProxiesAcceptanceMemBytes {
			t.Errorf("Proxy Egress Memory value (%v) is higher than the acceptance (in bytes): %v", istioEgressProxyMemValueInt, istioEgressProxiesAcceptanceMemBytes)
		}

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
