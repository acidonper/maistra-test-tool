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
		t.FailNow()
	}
	xdsPushTime, err := getMetricPrometheusMesh("xds_ppctb")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}

	xdsPushCountValue, err := parsePromResponse([]byte(xdsPushCount))
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
	xdsPushTimeValue, err := parsePromResponse([]byte(xdsPushTime))
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}

	for i := 0; i < len(xdsPushCountValue); i++ {
		if xdsPushCountValue[i] != xdsPushTimeValue[i] {
			util.Log.Errorf("xdsPushCount (%v) and xdsPushTime (%v) are not equal and some requests are higher than %s seconds", xdsPushCountValue, xdsPushTimeValue, xdsPushAcceptanceTime)
			t.FailNow()
		} else {
			util.Log.Info("OK: ", xdsPushCountValue[i], "/", xdsPushTimeValue[i], " correct XDS pushes")
		}
	}
}

func TestXDSErrors(t *testing.T) {
	util.Log.Info("** TEST: TestXDSErrors")
	// xds_reject_metrics := map[string]string{"xds_cdsrej": "Cluster Discovery Service", "xds_edsrej": "Endpoint Discovery Service", "xds_rdsrej": "Route Discovery Service", "xds_ldsrej": "Listener Discovery Service", "xds_write_timeouts": "Pilot XDS response write timeouts", "pilot_total_xds_internal_errors": "Internal XDS errors in pilot", "pilot_total_xds_rejects": ""}
	xds_reject_metrics := map[string]string{"xds_cdsrej": "Cluster Discovery Service", "xds_edsrej": "Endpoint Discovery Service", "xds_rdsrej": "Route Discovery Service", "xds_ldsrej": "Listener Discovery Service"}
	for metric, usage := range xds_reject_metrics {
		xdsError, err := getMetricPrometheusMesh(metric)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		xdsErrorValue, err := parsePromResponse([]byte(xdsError))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		for i := 0; i < len(xdsErrorValue); i++ {
			if xdsErrorValue[i] != "0" {
				util.Log.Error("XDS has errors. Failure in metric: ", metric)
				t.FailNow()
			} else {
				util.Log.Info("OK: ", xdsErrorValue[i], " errors -- ", usage)
			}
		}
	}
}

func TestIstiodMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstiodMem")

	meshPods, err := getMeshProxies("istiod")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istiodProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istiodMemValue, err := parsePromResponse([]byte(istiodProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istiodMemValue[0], istiodAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
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
		t.FailNow()
	}

	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istiodCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istiodCpuValue, err := parsePromResponse([]byte(istiodCpu))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istiodCpuValue[0], istiodAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
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
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istioProxyMemValue, err := parsePromResponse([]byte(istioProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		if len(istioProxyMemValue) > 0 {
			result, errComp := comparePodsMem(istioProxyMemValue[0], istioProxiesAcceptanceMem)
			if errComp != nil {
				util.Log.Error(err)
				t.FailNow()
			} else {
				util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
			}
		} else {
			util.Log.Info("WARNING: Memory metric for (pod/namespace: ", name, "/", namespace, ") not found")
		}
	}
	util.Log.Info("OK: Istio proxies memory are lower than ", istioProxiesAcceptanceMem, " in Megabytes")
}

func TestIstioProxiesCpu(t *testing.T) {
	util.Log.Info("** TEST: TestIstioProxiesCpu")

	meshPods, err := getMeshProxies("proxy")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioProxyCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istioProxyCpuValue, err := parsePromResponse([]byte(istioProxyCpu))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		if len(istioProxyCpuValue) > 0 {
			result, errComp := comparePodsCpu(istioProxyCpuValue[0], istioProxiesAcceptanceCpu)
			if errComp != nil {
				util.Log.Error(err)
				t.FailNow()
			} else {
				util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
			}
		} else {
			util.Log.Info("WARNING: CPU metric for (pod/namespace: ", name, "/", namespace, ") not found")
		}
	}
	util.Log.Info("OK: Istio proxies CPU are lower than ", istioProxiesAcceptanceCpu, " in CPUs")
}

func TestIstioIngressMem(t *testing.T) {
	util.Log.Info("** TEST: TestIstioIngressMem")

	meshPods, err := getMeshProxies("ingress")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}

	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioIngressProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istioIngressMemValue, err := parsePromResponse([]byte(istioIngressProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istioIngressMemValue[0], istioIngressAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
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
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioIngressCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istioIngressCpuValue, err := parsePromResponse([]byte(istioIngressCpu))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istioIngressCpuValue[0], istioIngressAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
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
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioEgressProxyMem, err := getMetricPrometheusOCP("istio_proxies_mem_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istioEgressMemValue, err := parsePromResponse([]byte(istioEgressProxyMem))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsMem(istioEgressMemValue[0], istioEgressAcceptanceMem)
		if errComp != nil {
			util.Log.Error(err)
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
		t.FailNow()
	}
	for name, namespace := range meshPods {
		prometheusAPIMapParams["istio-proxy-pod-name"] = name
		prometheusAPIMapParams["istio-proxy-pod-ns"] = namespace
		istioEgressProxyCpu, err := getMetricPrometheusOCP("istio_proxies_cpu_custom", prometheusAPIMapParams)
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		istioEgressCpuValue, err := parsePromResponse([]byte(istioEgressProxyCpu))
		if err != nil {
			util.Log.Error(err)
			t.FailNow()
		}

		result, errComp := comparePodsCpu(istioEgressCpuValue[0], istioEgressAcceptanceCpu)
		if errComp != nil {
			util.Log.Error(err)
			t.FailNow()
		} else {
			util.Log.Info(result, " (pod/namespace: ", name, "/", namespace, ")")
		}
	}
}
