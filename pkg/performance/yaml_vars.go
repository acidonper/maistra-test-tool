package performance

import (
	"fmt"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

var (
	smcpName       string = util.Getenv("SMCPNAME", "basic")
	meshNamespace  string = util.Getenv("MESHNAMESPACE", "istio-system")
	appNSPrefix    string = util.Getenv("APPNSPREFIX", "app-perf-test")
	nsCountBundle  string = util.Getenv("NSCOUNTBUNDLE", "10,15,20")
	testCPApps     string = util.Getenv("TESTCPAPPS", "2")
	testDPApps     string = util.Getenv("TESTDPAPPS", "5")
	testDPAppsFill string = util.Getenv("TESTDPAPPSFILL", "false")
	testVUs        string = util.Getenv("TESTVUS", "1")
	testDuration   string = util.Getenv("TESTDURATION", "30")
)

var (
	nsAcceptanceTime          string = util.Getenv("NSACCEPTANCETIME", "5")
	xdsPushAcceptanceTime     string = util.Getenv("XDSPUSHACCEPTANCETIME", "1")
	istiodAcceptanceMem       string = util.Getenv("ISTIODACCEPTANCEMEM", "1024")
	istiodAcceptanceCpu       string = util.Getenv("ISTIODACCEPTANCECPU", "1000")
	istioProxiesAcceptanceMem string = util.Getenv("ISTIOPROXIESACCEPTANCEMEM", "1024")
	istioProxiesAcceptanceCpu string = util.Getenv("ISTIOPROXIESDACCEPTANCECPU", "1000")
	istioIngressAcceptanceMem string = util.Getenv("ISTIOINGRESSPROXIESACCEPTANCEMEM", "1024")
	istioIngressAcceptanceCpu string = util.Getenv("ISTIOINGRESSPROXIESACCEPTANCECPU", "1000")
	istioEgressAcceptanceMem  string = util.Getenv("ISTIOEGRESSPROXIESACCEPTANCEMEM", "1024")
	istioEgressAcceptanceCpu  string = util.Getenv("ISTIOEGRESSPROXIESACCEPTANCECPU", "1000")
	reqAvg95pAcceptanceTime   string = util.Getenv("REQAVG95PACCEPTANCETIME", "500")
)

var (
	prometheusMeshAPIMap map[string]string = map[string]string{
		"xds_ppctc":  "pilot_proxy_convergence_time_count",
		"xds_ppctb":  "pilot_proxy_convergence_time_bucket{le=\"" + xdsPushAcceptanceTime + "\"}",
		"xds_cdsrej": "sum(pilot_xds_cds_reject{app='istiod'}) or (absent(pilot_xds_cds_reject{app='istiod'}) - 1)",
		"xds_edsrej": "sum(pilot_xds_eds_reject{app='istiod'}) or (absent(pilot_xds_eds_reject{app='istiod'}) - 1)",
		"xds_rdsrej": "sum(pilot_xds_rds_reject{app='istiod'}) or (absent(pilot_xds_rds_reject{app='istiod'}) - 1)",
		"xds_ldsrej": "sum(pilot_xds_lds_reject{app='istiod'}) or (absent(pilot_xds_lds_reject{app='istiod'}) - 1)",
		// "xds_write_timeouts":              "sum(rate(pilot_xds_write_timeout{app='istiod'}[1m]))", ## Uncomment for OSSM > 2.1
		// "pilot_total_xds_internal_errors": "sum(rate(pilot_total_xds_internal_errors{app='istiod'}[1m]))", ## Uncomment for OSSM > 2.1
		// "pilot_total_xds_rejects":         "sum(rate(pilot_total_xds_rejects{app='istiod'}[1m]))", ## Uncomment for OSSM > 2.1
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
	bytesToMegaBytes     int     = 1000000
	coresToMilicores     float64 = 1000
	MegaBytesToKiloBytes int     = 1000
)

var (
	podFailedGet = "Failed_Get"
	statusField  = 2
)

var (
	basedir = "../../objects"
)

var (
	supportedApps              = []string{"bookinfo", "jumpapp"}
	trafficLoadApp      string = util.Getenv("TRAFFICLOADAPP", "bookinfo")
	trafficLoadProtocol string = util.Getenv("TRAFFICLOADPROTOCOL", "http")

	bookinfoYaml           = fmt.Sprintf("%s/bookinfo/bookinfo.yaml", basedir)
	bookinfoRuleAllYaml    = fmt.Sprintf("%s/bookinfo/destination-rule-all.yaml", basedir)
	bookinfoRuleAllTLSYaml = fmt.Sprintf("%s/bookinfo/destination-rule-all-mtls.yaml", basedir)
	bookinfoTotalCPU       = "600"  // (50 + 50 + 50 + 50 + 50 + 50) * 2
	bookinfoTotalMem       = "1134" // 250 + 250 + 250 + 128 + 128 + 128

	jumpappYaml       = fmt.Sprintf("%s/jumpapp/jumpapp.yaml", basedir)
	jumpappNetworking = fmt.Sprintf("%s/jumpapp/jumpapp-sm.yaml", basedir)
	jumpappTotalCPU   = "1000" // 250 x 4
	jumpappTotalMem   = "1000" // 250 x 4

	bookinfoNSPrefix = "bookinfo-"
	jumpappNSPrefix  = "jumpapp-"

	appName              = bookinfoNSPrefix + "1"
	appFillClusterNumber = 0
	reportFile           = "/tmp/" + appName + ".json"
)
