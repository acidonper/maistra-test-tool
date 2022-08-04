package performance

import (
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
			t.Errorf("xdsPushCount (%v) and xdsPushTime (%v) are not equal and some requests are higher than %s seconds", xdsPushCountValue, xdsPushTimeValue, xdsPushAcceptanceTime)
			t.FailNow()
		} else {
			util.Log.Info("OK: ", xdsPushCountValue[i], "/", xdsPushTimeValue[i], " correct XDS pushes")
		}
	}
}

func TestIstiodMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodMem")

	meshPods, err := getMeshProxies("istiod")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}

	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istiodProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		istiodMemValue, err := parseResponse([]byte(istiodProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istiodMemValue[0], istiodAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}

func TestIstiodCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodCpu")

	meshPods, err := getMeshProxies("istiod")
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
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istiodCpuValue[0], istiodAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}

func TestIstioProxiesMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioProxiesMem")

	meshPods, err := getMeshProxies("proxy")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		istioProxyMemValue, err := parseResponse([]byte(istioProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istioProxyMemValue[0], istioProxiesAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
	util.Log.Info("OK: Istio proxies memory are lower than ", istioProxiesAcceptanceMem, " in Megabytes")
}

func TestIstioProxiesCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioProxiesCpu")

	meshPods, err := getMeshProxies("proxy")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioProxyCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		istioProxyCpuValue, err := parseResponse([]byte(istioProxyCpu))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istioProxyCpuValue[0], istioProxiesAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
	util.Log.Info("OK: Istio proxies CPU are lower than ", istioProxiesAcceptanceCpu, " in CPUs")
}

func TestIstioIngressMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioIngressMem")

	meshPods, err := getMeshProxies("ingress")
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

		istioIngressMemValue, err := parseResponse([]byte(istioIngressProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istioIngressMemValue[0], istioIngressAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}

func TestIstioIngressCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioIngressCpu")

	meshPods, err := getMeshProxies("ingress")
	if err != nil {
		util.Log.Error(err)
		t.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioIngressCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		istioIngressCpuValue, err := parseResponse([]byte(istioIngressCpu))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istioIngressCpuValue[0], istioIngressAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}

func TestIstioEgressMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioEgressMem")

	meshPods, err := getMeshProxies("egress")
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

		istioEgressMemValue, err := parseResponse([]byte(istioEgressProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istioEgressMemValue[0], istioEgressAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}

func TestIstioEgressCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioEgressCpu")

	meshPods, err := getMeshProxies("egress")
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

		istioEgressCpuValue, err := parseResponse([]byte(istioEgressProxyCpu))
		if err != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istioEgressCpuValue[0], istioEgressAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
			t.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}
