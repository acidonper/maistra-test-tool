// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package performance

import (
	"fmt"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

var (
	smcpName      string = util.Getenv("SMCPNAME", "basic")
	meshNamespace string = util.Getenv("MESHNAMESPACE", "istio-system")
	appNSPrefix   string = util.Getenv("APPNSPREFIX", "app-perf-test")
	nsCountBundle string = util.Getenv("NSCOUNTBUNDLE", "10,15,20")
	testCPApps    string = util.Getenv("TESTCPAPPS", "2")
	testDPApps    string = util.Getenv("TESTDPAPPS", "5")
)

var (
	nsAcceptanceTime          string = util.Getenv("NSACCEPTANCETIME", "5")
	xdsPushAcceptanceTime     string = util.Getenv("XDSPUSHACCEPTANCETIME", "1")
	istiodAcceptanceMem       string = util.Getenv("ISTIODACCEPTANCEMEM", "1024")
	istiodAcceptanceCpu       string = util.Getenv("ISTIODACCEPTANCECPU", "1000")
	istioProxiesAcceptanceMem string = util.Getenv("ISTIOPROXIESACCEPTANCEMEM", "1024")
	istioProxiesAcceptanceCpu string = util.Getenv("ISTIOPROXIESDACCEPTANCECPU", "1000")
	istioIngressAcceptanceMem string = util.Getenv("ISTIOINGRESPROXIESACCEPTANCEMEM", "1024")
	istioIngressAcceptanceCpu string = util.Getenv("ISTIOINGRESPROXIESACCEPTANCECPU", "1000")
	istioEgressAcceptanceMem  string = util.Getenv("ISTIOEGRESPROXIESACCEPTANCEMEM", "1024")
	istioEgressAcceptanceCpu  string = util.Getenv("ISTIOEGRESPROXIESACCEPTANCECPU", "1000")
)

var (
	prometheusMeshAPIMap map[string]string = map[string]string{
		"xds_ppctc": "pilot_proxy_convergence_time_count",
		"xds_ppctb": "pilot_proxy_convergence_time_bucket{le=\"" + xdsPushAcceptanceTime + "\"}",
	}
	prometheusAPIMap map[string]string = map[string]string{
		"istiod_mem":        "sum(container_memory_working_set_bytes{pod=~\"istiod.*\",namespace=\"" + meshNamespace + "\",container=\"\",}) BY (pod, namespace)",
		"istiod_cpu":        "pod:container_cpu_usage:sum{pod=~\"istiod.*\",namespace=\"" + meshNamespace + "\"}",
		"istio_proxies_mem": "sum(container_memory_working_set_bytes{container='istio-proxy',}) BY (pod, namespace)",
	}
	prometheusAPIMapCustom map[string]string = map[string]string{
		"istio_proxies_mem_custom": "sum(container_memory_working_set_bytes{pod='istio-proxy-pod-name', namespace='istio-proxy-pod-ns',container='',}) BY (pod, namespace)",
		"istio_proxies_cpu_custom": "pod:container_cpu_usage:sum{pod='istio-proxy-pod-name',namespace='istio-proxy-pod-ns',}",
	}
	prometheusAPIMapParams map[string]string = map[string]string{}
)

var (
	bytesToMegaBytes int     = 1000000
	coresToMilicores float64 = 1000
)

var (
	podFailedGet = "Failed_Get"
	statusField  = 2
)

var (
	basedir = "../../objects"
)

var (
	bookinfoYaml           = fmt.Sprintf("%s/bookinfo/bookinfo.yaml", basedir)
	bookinfoGateway        = fmt.Sprintf("%s/bookinfo/bookinfo-gateway.yaml", basedir)
	bookinfoRuleAllYaml    = fmt.Sprintf("%s/bookinfo/destination-rule-all.yaml", basedir)
	bookinfoRuleAllTLSYaml = fmt.Sprintf("%s/bookinfo/destination-rule-all-mtls.yaml", basedir)

	jumpappYaml       = fmt.Sprintf("%s/jumpapp/jumpapp.yaml", basedir)
	jumpappNetworking = fmt.Sprintf("%s/jumpapp/jumpapp-sm.yaml", basedir)
)

type Bookinfo struct {
	Namespace string `json:"namespace,omitempty"`
}

type JumpApp struct {
	Namespace string `json:"namespace,omitempty"`
}

type PromResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name                 string `json:"__name__"`
				App                  string `json:"app"`
				Instance             string `json:"instance"`
				Istio                string `json:"istio"`
				IstioIoRev           string `json:"istio_io_rev"`
				Job                  string `json:"job"`
				KubernetesNamespace  string `json:"kubernetes_namespace"`
				KubernetesPodName    string `json:"kubernetes_pod_name"`
				Le                   string `json:"le"`
				MaistraControlPlane  string `json:"maistra_control_plane"`
				PodTemplateHash      string `json:"pod_template_hash"`
				SidecarIstioIoInject string `json:"sidecar_istio_io_inject"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
