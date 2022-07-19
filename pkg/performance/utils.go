package performance

import (
	"net/http"
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
		if !strings.Contains(err.Error(), "AlreadyExists") {
			return err
		}
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

func GetWithJWT(url, token, host string) (*http.Response, error) {
	// Declare http client
	client := &http.Client{}

	// Declare HTTP Method and Url
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Host = host
	req.Header.Set("Host", req.Host)
	req.Header.Add("Authorization", "Bearer "+token)
	return client.Do(req)
}
