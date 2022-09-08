package performance

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestSMCP(t *testing.T) {
	msg, _ := util.ShellSilent(`oc wait --for condition=Ready -n %s smcp/%s --timeout 30s`, meshNamespace, smcpName)
	if !strings.Contains(msg, "condition met") {
		util.Log.Error("SMCP not Ready")
		t.Error("SMCP not Ready")
		t.FailNow()
	}
	util.Log.Info("OK - SMCP ", smcpName, " in namespace ", meshNamespace)
}

func TestSMMR(t *testing.T) {
	msg, _ := util.ShellSilent(`oc wait --for condition=Ready -n %s smmr/default --timeout 30s`, meshNamespace)
	if !strings.Contains(msg, "condition met") {
		util.Log.Error("SMMR not Ready")
		t.Error("SMMR not Ready")
		t.FailNow()
	}
	util.Log.Info("OK - SMMR in namespace ", meshNamespace)

}

func getRouteHost(route string, namespace string) (string, error) {
	util.Log.Debug("Getting Route Host", route, " in namespace ", namespace)
	msg, err := util.ShellSilent(`oc get route %s -n %s --template='{{ .spec.host }}'`, route, namespace)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func getOCPNodeNames(role string) (string, error) {
	util.Log.Debug("Getting Node Names with role", role)
	msg, err := util.ShellSilent(`oc get nodes -l node-role.kubernetes.io/%s="" -o go-template='{{range .items}}{{.metadata.name}},{{end}}'`, role)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func getOCPNodeAllocatableCPU(node string) (string, error) {
	util.Log.Debug("Getting Node free CPU: ", node)
	cpu, err := util.ShellSilent(`oc get nodes %s --template='{{ .status.allocatable.cpu }}'`, node)
	if err != nil {
		return "", err
	}
	return cpu, nil
}

func getOCPNodeAllocatableMem(node string) (string, error) {
	util.Log.Debug("Getting Node free Memory: ", node)
	cpu, err := util.ShellSilent(`oc get nodes %s --template='{{ .status.allocatable.memory }}'`, node)
	if err != nil {
		return "", err
	}
	return cpu, nil
}

func getOCPNodeFreeCPU(node string) (int, error) {
	util.Log.Debug("Getting Node free CPU: ", node)

	// Obtain Allocatable CPU
	cpuAllocatable, err := getOCPNodeAllocatableCPU(node)
	if err != nil {
		return 0, err
	}

	// Obtain Consumed CPU
	cpuConsumedtmp, err := util.ShellSilent(`oc adm top node %s --no-headers`, node)
	if err != nil {
		return 0, err
	}
	cpuConsumed := strings.Split(cpuConsumedtmp, "   ")

	// Calculate free CPU
	cpuAllocatableData := strings.Split(cpuAllocatable, "m")
	cpuAllocatableInt, err := strconv.Atoi(strings.TrimSpace(cpuAllocatableData[0]))
	if err != nil {
		return 0, err
	}
	cpuConsumedData := strings.Split(cpuConsumed[1], "m")
	cpuConsumedDataInt, err := strconv.Atoi(strings.TrimSpace(cpuConsumedData[0]))
	if err != nil {
		return 0, err
	}
	cpuFree := cpuAllocatableInt - cpuConsumedDataInt

	// Return Free CPU in milicores
	return cpuFree, nil
}

func getOCPNodeFreeMem(node string) (int, error) {
	util.Log.Debug("Getting Node free Mem: ", node)

	// Obtain Allocatable Mem
	memAllocatable, err := getOCPNodeAllocatableMem(node)
	if err != nil {
		return 0, err
	}

	// Obtain Consumed Mem
	memConsumedtmp, err := util.ShellSilent(`oc adm top node %s --no-headers`, node)
	if err != nil {
		return 0, err
	}
	memConsumed := strings.Split(memConsumedtmp, "   ")

	// Calculate free Mem
	memAllocatableData := strings.Split(memAllocatable, "Ki")
	memAllocatableInt, err := strconv.Atoi(strings.TrimSpace(memAllocatableData[0]))
	if err != nil {
		return 0, err
	}
	memConsumedData := strings.Split(memConsumed[3], "Mi")
	memConsumedDataInt, err := strconv.Atoi(strings.TrimSpace(memConsumedData[0]))
	if err != nil {
		return 0, err
	}
	memFree := memAllocatableInt/MegaBytesToKiloBytes - memConsumedDataInt

	// Rturn Free Memory in Megabytes
	return memFree, nil
}

func getOCPAppsDomain() (string, error) {
	util.Log.Debug("Getting Openshift Apps Domain")
	exampleRoute, err := getRouteHost("console", "openshift-console")
	if err != nil {
		return "", err
	}
	domain := strings.ReplaceAll(exampleRoute, "console-openshift-console.", "")
	return domain, nil
}

func getMeshProxyPods() ([]string, error) {
	meshPods, err := util.ShellSilent(`oc get pods -A -l istio.io/rev=%s --field-selector=status.phase==Running -o go-template='{{range .items}}{{.metadata.name}}/{{.metadata.namespace}},{{end}}'`, smcpName)
	if err != nil {
		return nil, err
	}
	podsList := strings.Split(meshPods, ",")
	return podsList, nil
}

func getMeshIstiodPods() ([]string, error) {
	meshPods, err := util.ShellSilent(`oc get pods -l app=istiod -l istio.io/rev=%s -n %s --field-selector=status.phase==Running -o go-template='{{range .items}}{{.metadata.name}}/{{.metadata.namespace}},{{end}}'`, smcpName, meshNamespace)
	if err != nil {
		return nil, err
	}
	podsList := strings.Split(meshPods, ",")
	return podsList, nil
}

func getMeshIngressPods() ([]string, error) {
	meshPods, err := util.ShellSilent(`oc get pods -A -l maistra-control-plane=%s -l istio=ingressgateway --field-selector=status.phase==Running -o go-template='{{range .items}}{{.metadata.name}}/{{.metadata.namespace}},{{end}}'`, meshNamespace)
	if err != nil {
		return nil, err
	}
	podsList := strings.Split(meshPods, ",")
	return podsList, nil
}

func getMeshEgressPods() ([]string, error) {
	meshPods, err := util.ShellSilent(`oc get pods -A -l maistra-control-plane=%s -l istio=egressgateway --field-selector=status.phase==Running -o go-template='{{range .items}}{{.metadata.name}}/{{.metadata.namespace}},{{end}}'`, meshNamespace)
	if err != nil {
		return nil, err
	}
	podsList := strings.Split(meshPods, ",")
	return podsList, nil
}

func deleteNS(namespace string) error {
	util.Log.Debug("Deleting namespace", namespace)
	_, err := util.ShellSilent(`oc delete project %s --wait=true`, namespace)
	if err != nil {
		return err
	}
	check := false
	for !check {
		_, err := util.ShellSilent(`oc get namespace %s`, namespace)
		if err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return err
			} else {
				check = true
			}
		}
	}
	return nil
}

func createNS(namespace string) error {
	util.Log.Debug("Creating namespace", namespace)
	_, err := util.ShellSilent(`oc new-project %s`, namespace)
	if err != nil {
		return err
	}
	return nil
}

func deleteNSMesh(namespace string) error {
	// Find namespace in members array
	tmp, err := util.ShellSilent(`oc get smmr default -n %s --template='{{ .spec.members }}'`, meshNamespace)
	if err != nil {
		return err
	}
	members := strings.Split(tmp, " ")
	position := arrayPositionFind(members, namespace)

	// If namespace exists
	if position > 0 {
		// Path SMMR
		_, err := util.ShellSilent(`oc patch smmr default -n %s --type='json' -p='[{"op": "remove", "path": "/spec/members/%s"}]'`, meshNamespace, strconv.Itoa(position))
		if err != nil {
			return err
		}

		// Verify SMMR deleting a NS
		configured := false
		util.Log.Debug("Waiting for the namespace to be deleted to the mesh: ", namespace)
		for !configured {
			msgConfigured, errConfigured := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.configuredMembers }}'`, meshNamespace)
			if errConfigured != nil {
				return errConfigured
			}
			msgPending, errPending := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.pendingMembers }}'`, meshNamespace)
			if errPending != nil {
				return errConfigured
			}
			if !strings.Contains(msgConfigured, namespace+" ") && !strings.Contains(msgPending, namespace+"]") {
				if !strings.Contains(msgPending, namespace+" ") && !strings.Contains(msgPending, namespace+"]") {
					err := deleteNS(namespace)
					if err != nil {
						return err
					}
					configured = true
				}
			}
		}
	}
	return nil
}

func createNSMesh(namespace string) error {
	// Create NS
	err := createNS(namespace)
	if err != nil {
		return err
	}

	// Path SMMR
	_, errPatch := util.ShellSilent(`oc patch smmr default -n %s --type='json' -p='[{"op": "add", "path": "/spec/members/-", "value":"%s"}]'`, meshNamespace, namespace)
	if errPatch != nil {
		return err
	}

	// Verify SMMR after adding a new NS
	configured := false
	util.Log.Debug("Waiting for the namespace to be added to the mesh: ", namespace)
	for !configured {
		msg, err := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.configuredMembers }}'`, meshNamespace)
		if errPatch != nil {
			return err
		}
		if strings.Contains(msg, namespace) {
			configured = true
			return nil
		}
	}
	return nil
}

func addNSToMesh(namespace string) error {
	// Path SMMR
	_, err := util.ShellSilent(`oc patch smmr default -n %s --type='json' -p='[{"op": "add", "path": "/spec/members/-", "value":"%s"}]'`, meshNamespace, namespace)
	if err != nil {
		return err
	}

	// Verify SMMR after adding a new NS
	configured := false
	util.Log.Debug("Waiting for the namespace to be added to the mesh: ", namespace)
	for !configured {
		msg, err := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.configuredMembers }}'`, meshNamespace)
		if err != nil {
			return err
		}
		if strings.Contains(msg, namespace) {
			configured = true
			return nil
		}
	}
	return nil
}

func arrayPositionFind(a []string, x string) int {
	for i, n := range a {
		if x == n || "["+x == n || x+"]" == n {
			return i
		}
	}
	return -1
}

func deleteNSBundle(min int, max int, prefix string) error {
	util.Log.Info("Deleting namespaces from ", min, " to ", max)
	for i := min; i < max; i++ {
		nsName := prefix + strconv.Itoa(i)
		err := deleteNSMesh(nsName)
		if err != nil {
			return err
		}
	}
	return nil
}

func createNSBundle(min int, max int, prefix string) error {
	util.Log.Info("Creating namespaces from ", min, " to ", max)
	for i := min; i < max; i++ {
		nsName := prefix + strconv.Itoa(i)
		err := createNSMesh(nsName)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteAppBundle(app string, number int, plane string) error {
	util.Log.Info("Deleting ", app, " applications: ", strconv.Itoa(number))

	// Calculate number of application to fill the cluster if it is required
	if testDPAppsFill == "true" && plane == "dataplane" {
		number = appFillClusterNumber
		util.Log.Info("Deleting ", app, " applications that filled the Openshift clusters: ", strconv.Itoa(number))
	}

	// Creating the respective number of apps
	for i := 1; i <= number; i++ {
		prefix := ""

		switch app {
		case "bookinfo":
			prefix = bookinfoNSPrefix
		case "jumpapp":
			prefix = jumpappNSPrefix
		default:
			return fmt.Errorf("application " + app + " not supported (bookinfo or jumpapp)")
		}

		nsName := prefix + strconv.Itoa(i)
		err := deleteNSMesh(nsName)
		if err != nil {
			return err
		}
	}
	return nil
}

func createAppBundle(app string, number int, plane string) error {
	util.Log.Info("Deploying ", app, " applications: ", strconv.Itoa(number))

	// Calculate number of application to fill the cluster if it is required
	if testDPAppsFill == "true" && plane == "dataplane" {
		apps, _ := calculateAppsFillCluster(app)
		number = apps
		util.Log.Info("Deploying ", app, " applications to fill the Openshift clusters: ", strconv.Itoa(number))

		// Save as global variable de number of applications
		appFillClusterNumber = apps
	}

	// Creating the respective number of apps
	for i := 1; i <= number; i++ {

		nsName := ""

		switch app {
		case "bookinfo":
			nsName = bookinfoNSPrefix + strconv.Itoa(i)
			err := createNSMesh(nsName)
			if err != nil {
				return err
			}
			ocpDomain, err := getOCPAppsDomain()
			if err != nil {
				return err
			}
			host := "bookinfo." + nsName + "." + ocpDomain
			bookinfo := Bookinfo{Name: nsName, Namespace: nsName, Host: host}
			err = bookinfo.BookinfoInstall(false)
			if err != nil {
				return err
			}
		case "jumpapp":
			nsName = bookinfoNSPrefix + strconv.Itoa(i)
			err := createNSMesh(nsName)
			if err != nil {
				return err
			}
			jumpapp := JumpApp{Namespace: nsName}
			err = jumpapp.JumpappInstall()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("application " + app + " not supported (bookinfo or jumpapp)")
		}
	}
	return nil
}

func httpPostQueryAuth(url string, user string, pass string) (*http.Response, error) {
	// Declare http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Declare HTTP Method and Url
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, pass)

	return client.Do(req)
}

func httpPostQueryBearer(url string, token string) (*http.Response, error) {
	// Declare http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Declare HTTP Method and Url
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	token = strings.TrimSuffix(token, "\n")
	req.Header.Add("Authorization", "Bearer "+token)

	return client.Do(req)
}

func obtainPrometheusMeshToken() (string, error) {
	token, err := util.ShellSilent(`oc get secret htpasswd -n istio-system --template='{{ .data.rawPassword }}' | base64 -d`)
	if err != nil {
		return "", err
	}
	return token, nil
}

func obtainPrometheusOCPToken() (string, error) {
	token, err := util.ShellSilent(`oc whoami -t`)
	if err != nil {
		return "", err
	}
	return token, nil
}

func getMetricPrometheus(host string, auth string, secret string, query string) (string, error) {
	// Generate final URL
	promUrl := "https://" + host + "/api/v1/query"
	baseUrl, err := url.Parse(promUrl)
	if err != nil {
		return "", err
	}
	values := baseUrl.Query()
	values.Add("query", query)
	baseUrl.RawQuery = values.Encode()

	// HTTP Post call to Prometheus
	var resp *http.Response
	if auth == "user" {
		user := "internal"
		pass := secret
		resp, err = httpPostQueryAuth(baseUrl.String(), user, pass)
		if err != nil {
			return "", err
		}
	} else if auth == "token" {
		token := secret
		resp, err = httpPostQueryBearer(baseUrl.String(), token)
		if err != nil {
			return "", err
		}
	} else {

		err := fmt.Errorf("auth method not defined: %q", auth)
		return "", err
	}

	// Process HTTP response
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	} else {
		return "", errors.New("HTTP Error " + resp.Status)
	}

}

func getMetricPrometheusMesh(query string) (string, error) {
	// Retrieve Mesh Prometheus Host
	routeHost, err := getRouteHost("prometheus", meshNamespace)
	if err != nil {
		return "", err
	}

	// Retrive respective query
	prometheusQuery := prometheusMeshAPIMap[query]

	// Obtain token to connect to Prometheus
	pass, err := obtainPrometheusMeshToken()
	if err != nil {
		return "", err
	}

	// HTTP Post call to Prometheus
	resp, err := getMetricPrometheus(routeHost, "user", pass, prometheusQuery)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func getMetricPrometheusOCP(query string, params map[string]string) (string, error) {
	// Retrieve Mesh Prometheus Host
	routeHost, err := getRouteHost("prometheus-k8s", "openshift-monitoring")
	if err != nil {
		return "", err
	}

	// Obtain token to connect to Prometheus
	token, err := obtainPrometheusOCPToken()
	if err != nil {
		return "", err
	}

	// Retrive respective query
	prometheusQuery := ""
	if params != nil {
		prometheusQuery = prometheusAPIMapCustom[query]
		for key, value := range params {
			prometheusQuery = strings.Replace(prometheusQuery, key, value, -1)
		}
	} else {
		prometheusQuery = prometheusAPIMap[query]
	}

	// HTTP Post call to Prometheus
	resp, err := getMetricPrometheus(routeHost, "token", token, prometheusQuery)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func getMeshProxies(role string) (map[string]string, error) {

	var podsList []string
	var err error

	// Get pods depending on type of role
	switch role {
	case "proxy":
		podsList, err = getMeshProxyPods()
	case "istiod":
		podsList, err = getMeshIstiodPods()
	case "ingress":
		podsList, err = getMeshIngressPods()
	case "egress":
		podsList, err = getMeshEgressPods()
	default:
		err = fmt.Errorf("mesh pod role not defined (Supported: proxy, istiod, ingress or egress)")
	}

	// Find proxies in the mesh
	if err != nil {
		return nil, err
	}

	// generate the respective map
	pods := make(map[string]string)
	for _, s := range podsList {
		if s != "" {
			podMetadata := strings.Split(s, "/")
			pods[podMetadata[0]] = podMetadata[1]
		}
	}

	return pods, nil
}

func GetPodStatus(n, pod string) (string, error) {
	status, err := util.ShellSilent("kubectl -n %s get pods %s --no-headers", n, pod)
	if err != nil {
		status = podFailedGet
	}
	f := strings.Fields(status)
	if len(f) > statusField {
		return f[statusField], nil
	}
	return "", err
}

// GetPodName gets the pod name for the given namespace and label selector
func GetPodName(n, labelSelector string) (pod string, err error) {
	pod, err = util.ShellSilent("kubectl -n %s get pod -l %s -o jsonpath='{.items[0].metadata.name}'", n, labelSelector)
	if err != nil {
		return "", fmt.Errorf("could not get %s pod: %v", labelSelector, err)
	}
	return
}

func CheckPodRunning(n, name string) error {
	retry := util.Retrier{
		BaseDelay: 30 * time.Second,
		MaxDelay:  30 * time.Second,
		Retries:   6,
	}

	retryFn := func(_ context.Context, i int) error {
		pod, err := GetPodName(n, name)
		if err != nil {
			return err
		}
		ready := true
		status, errStatusPod := GetPodStatus(n, pod)
		if errStatusPod != nil {
			return err
		}
		if status != "Running" {
			util.Log.Debug("%s in namespace %s is not running: %s", pod, n, status)
			ready = false
		}
		if !ready {
			return fmt.Errorf("pod %s is not ready", pod)
		}
		return nil
	}
	ctx := context.Background()
	_, err := retry.Retry(ctx, retryFn)
	if err != nil {
		return err
	}
	util.Log.Debug("Got the pod name=%s running!", name)
	return nil
}

func parsePromResponse(response []byte) ([]string, error) {

	var newResponse PromResponse

	err := json.Unmarshal(response, &newResponse)

	if err != nil {
		return nil, err
	}

	var values []string

	// Retrieve all values from the array of results
	for i := 0; i < len(newResponse.Data.Result); i++ {
		var value = fmt.Sprintf("%v", newResponse.Data.Result[i].Value[1])
		values = append(values, value)
	}

	return values, nil
}

func comparePodsMem(value1 string, value2 string) (string, error) {

	value1Int, errConver1 := strconv.Atoi(value1)
	if errConver1 != nil {
		return "", errConver1
	}
	value1IntMegaBytes := value1Int / bytesToMegaBytes

	value2IntMegabytes, errConver2 := strconv.Atoi(value2)
	if errConver2 != nil {
		return "", errConver1
	}

	if value1IntMegaBytes > value2IntMegabytes {
		msg := fmt.Errorf("memory value is %v. Want something lower than %v", value1IntMegaBytes, value2IntMegabytes)
		return "", msg
	} else {
		msg := ("OK: Memory " + strconv.Itoa(value1IntMegaBytes) + " is lower than " + strconv.Itoa(value2IntMegabytes) + " in MBs")
		return msg, nil
	}
}

func comparePodsCpu(value1 string, value2 string) (string, error) {

	value1Float, errConver1 := strconv.ParseFloat(value1, 32)
	if errConver1 != nil {
		return "", errConver1
	}
	value1FloatMilicores := value1Float * coresToMilicores

	value2FloatMilicores, errConver2 := strconv.ParseFloat(value2, 32)
	if errConver2 != nil {
		return "", errConver1
	}

	if value1FloatMilicores > value2FloatMilicores {
		msg := fmt.Errorf("cpu value is %v. Want something lower than %v", value1FloatMilicores, value2FloatMilicores)
		return "", msg
	} else {
		msg := ("OK: CPU " + fmt.Sprintf("%f", value1FloatMilicores) + " is lower than " + fmt.Sprintf("%f", value2FloatMilicores) + " in Milicores")
		return msg, nil
	}
}

func checkFailedMetrics(fails int) (string, error) {

	if fails > 0 {
		msg := fmt.Errorf("there are %v requests failing", fails)
		return "", msg
	} else {
		msg := ("OK: No requests failing")
		return msg, nil
	}

}

func compareP95(value1 string, value2 string) (string, error) {

	value1Float, errConver1 := strconv.ParseFloat(value1, 64)
	if errConver1 != nil {
		return "", errConver1
	}

	value2Float, errConver2 := strconv.ParseFloat(value2, 64)
	if errConver2 != nil {
		return "", errConver1
	}

	if value1Float > value2Float {
		msg := fmt.Errorf("percentile 95 is %v. Want something lower than %v", value1Float, value2Float)
		return "", msg
	} else {
		msg := ("OK: Percentile 95 " + fmt.Sprintf("%f", value1Float) + " is lower than " + fmt.Sprintf("%f", value2Float))
		return msg, nil
	}

}

func execK6SyncTest(vus string, duration string, url string, test string, file string) (string, error) {
	util.Log.Info("Executing test ", test, " in ", url, " (vus/duration: ", vus, "/", duration, "s)")
	msg, err := util.ShellSilent(`k6 run --vus %s --duration %ss --env TEST_URL="%s" --summary-export %s %s/k6/%s`, vus, duration, url, file, basedir, test)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func generateSimpleTrafficLoadK6(protocol string, app string) error {

	var url string
	var err error

	if app == "bookinfo" && protocol == "http" {

		routeHost, errRoute := getRouteHost(appName, meshNamespace)
		if errRoute != nil {
			return fmt.Errorf("route %s not found in namespace %s", appName, meshNamespace)
		} else {
			url = "https://" + routeHost + "/productpage"
		}
		_, err = execK6SyncTest(testVUs, testDuration, url, "http-basic.js", reportFile)

		if err != nil {
			return err
		}

	} else {
		errorMsg := fmt.Errorf("application %s and protocol %s not supported", app, protocol)
		return errorMsg
	}

	return nil
}

func readK6File() ([]byte, error) {

	dat, err := os.ReadFile(reportFile)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func parseK6Response(response []byte) (K6Response, error) {

	var newResponse K6Response

	err := json.Unmarshal(response, &newResponse)

	if err != nil {
		return K6Response{}, err
	}

	return newResponse, nil
}

func getOCPWorkerNodesFreeResources() (int, int, error) {

	// Global Vars
	totalCPUMilicores := 0
	totalMemMegabytes := 0

	// Get pods depending on type of role
	workers, err := getOCPNodeNames("worker")
	if err != nil {
		return 0, 0, err
	}

	// Generate node list
	workersList := strings.Split(workers, ",")

	// Obtain CPU and MEM for every node
	for _, s := range workersList {

		if s != "" {
			cpu, err := getOCPNodeFreeCPU(s)
			if err != nil {
				return 0, 0, err
			}
			mem, err := getOCPNodeFreeMem(s)
			if err != nil {
				return 0, 0, err
			}
			totalCPUMilicores = totalCPUMilicores + cpu
			totalMemMegabytes = totalMemMegabytes + mem
		}

	}

	return totalCPUMilicores, totalMemMegabytes, nil
}

func calculateAppsFillCluster(app string) (int, error) {

	// Obtain free resources in the OCP Cluster
	cpu, mem, err := getOCPWorkerNodesFreeResources()
	if err != nil {
		return 0, err
	}

	// Obtain resources consumed by the application
	cpuLimit := 0
	memLimit := 0

	switch app {
	case "bookinfo":
		cpuLimit, _ = strconv.Atoi(bookinfoTotalCPU)
		memLimit, _ = strconv.Atoi(bookinfoTotalMem)
	case "jumpapp":
		cpuLimit, _ = strconv.Atoi(jumpappTotalCPU)
		memLimit, _ = strconv.Atoi(jumpappTotalMem)
	default:
		return 0, fmt.Errorf("application " + app + " not supported (bookinfo or jumpapp)")
	}

	// Calculate the number of apps by mem and cpu and return the lowest
	cpuNumApp := cpu / cpuLimit
	memNumApp := mem / memLimit
	if cpuNumApp < memNumApp {
		return cpuNumApp, nil
	} else {
		return memNumApp, nil
	}

}
