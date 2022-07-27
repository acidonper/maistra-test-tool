package performance

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

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

func getMeshPods() ([]string, error) {
	meshPods, err := util.ShellSilent(`oc get pods -A -l istio.io/rev=%s --field-selector=status.phase==Running -o go-template='{{range .items}}{{.metadata.name}}/{{.metadata.namespace}},{{end}}'`, smcpName)
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

func delNamespaceMesh(namespace string) error {
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

func createNamespaceMesh(namespace string) error {
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

func addNamespaceMesh(namespace string) error {
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
		return "", fmt.Errorf("auth method not defined: " + auth)
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

func getMeshProxies() (map[string]string, error) {
	// Find proxies in the mesh
	podsList, err := getMeshPods()
	if err != nil {
		return nil, err
	}

	// generate the respective map excluding control plane namespace proxies
	pods := make(map[string]string)
	for _, s := range podsList {
		if !strings.Contains(s, meshNamespace) && s != "" {
			podMetadata := strings.Split(s, "/")
			pods[podMetadata[0]] = podMetadata[1]
		}
	}

	return pods, nil
}

func getMeshIngressProxies() (map[string]string, error) {
	// Find proxies in the mesh
	podsList, err := getMeshPods()
	if err != nil {
		return nil, err
	}

	// generate the respective map filtering ingress gateways
	pods := make(map[string]string)
	for _, s := range podsList {
		if strings.Contains(s, "istio-ingressgateway-") && s != "" {
			podMetadata := strings.Split(s, "/")
			pods[podMetadata[0]] = podMetadata[1]
		}
	}

	return pods, nil
}

func getMeshEgressProxies() (map[string]string, error) {
	// Find proxies in the mesh
	podsList, err := getMeshPods()
	if err != nil {
		return nil, err
	}

	// generate the respective map filtering egress gateways
	pods := make(map[string]string)
	for _, s := range podsList {
		if strings.Contains(s, "istio-egressgateway-") && s != "" {
			podMetadata := strings.Split(s, "/")
			pods[podMetadata[0]] = podMetadata[1]
		}
	}

	return pods, nil
}
