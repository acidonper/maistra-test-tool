package performance

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestSMCP(t *testing.T) {
	util.Log.Info("Checking SMCP in ", meshNamespace)
	msg, _ := util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 30s`, meshNamespace, smcpName)
	if !strings.Contains(msg, "condition met") {
		util.Log.Error("SMCP not Ready")
		t.Error("SMCP not Ready")
		t.FailNow()
	}
}

func TestSMMR(t *testing.T) {
	util.Log.Info("Checking SMMR in ", meshNamespace)
	msg, _ := util.Shell(`oc wait --for condition=Ready -n %s smmr/default --timeout 30s`, meshNamespace)
	if !strings.Contains(msg, "condition met") {
		util.Log.Error("SMMR not Ready")
		t.Error("SMMR not Ready")
		t.FailNow()
	}
}

func getRouteHost(route string, namespace string) (string, error) {
	util.Log.Debug("Getting Route Host", route, " in namespace ", namespace)
	msg, err := util.ShellSilent(`oc get route %s -n %s --template='{{ .spec.host }}'`, route, namespace)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func deleteNS(namespace string) error {
	util.Log.Debug("Deleting namespace", namespace)
	_, err := util.ShellSilent(`oc delete project %s`, namespace)
	if err != nil {
		return err
	}
	return nil
}

func createNS(namespace string) error {
	util.Log.Debug("Creating namespace", namespace)
	_, err := util.ShellSilent(`oc new-project %s`, namespace)
	if err != nil {
		// if !strings.Contains(err.Error(), "AlreadyExists") {
		// 	return err
		// }
		return err
	}
	return nil
}

func delNamespaceMesh(namespace string) error {
	// Find namespace in members array
	tmp, _ := util.ShellSilent(`oc get smmr default -n %s --template='{{ .spec.members }}'`, meshNamespace)
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

func httpPostQueryAuth(url string, user string, token string) (*http.Response, error) {
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
	req.SetBasicAuth(user, token)

	return client.Do(req)
}

func obtainMeshPrometheusToken() (string, error) {
	token, err := util.ShellSilent(`oc get secret htpasswd -n istio-system --template='{{ .data.rawPassword }}' | base64 -d`)
	if err != nil {
		return "", err
	}
	return token, nil
}

func getMetricPrometheus(query string) (string, error) {
	// Create Prometheus URL
	routeHost, err := getRouteHost("prometheus", meshNamespace)
	if err != nil {
		return "", err
	}
	promUrl := "https://" + routeHost + "/api/v1/query"
	baseUrl, err := url.Parse(promUrl)
	if err != nil {
		return "", err
	}
	values := baseUrl.Query()
	values.Add("query", prometheusAPIMap[query])
	baseUrl.RawQuery = values.Encode()

	// Obtain token to connect to Prometheus
	token, err := obtainMeshPrometheusToken()
	if err != nil {
		return "", err
	}

	// HTTP Post call to Prometheus
	resp, err := httpPostQueryAuth(baseUrl.String(), "internal", token)
	if err != nil {
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

//curl -kv -u internal:$PEPE https://prometheus-istio-system.apps.meshtests.sandbox84.opentlc.com/api/v1/query --data-urlencode "query=histogram_quantile(0.99, sum(rate(pilot_proxy_convergence_time_bucket[1m])) by (le))"
