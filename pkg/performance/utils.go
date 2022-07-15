package performance

import (
	"strconv"
	"strings"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func TestSMCP(t *testing.T) {
	util.Log.Info("Checking SMCP in ", meshNamespace)
	msg, _ := util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)
	if !strings.Contains(msg, "condition met") {
		util.Log.Error("SMCP not Ready")
		t.Error("SMCP not Ready")
		t.FailNow()
	}
}

func TestSMMR(t *testing.T) {
	util.Log.Info("Checking SMMR in ", meshNamespace)
	msg, _ := util.Shell(`oc wait --for condition=Ready -n %s smmr/default --timeout 300s`, meshNamespace)
	if !strings.Contains(msg, "condition met") {
		util.Log.Error("SMMR not Ready")
		t.Error("SMMR not Ready")
		t.FailNow()
	}
}

func deleteNS(namespace string) {
	util.Log.Debug("Deleting namespace", namespace)
	util.ShellSilent(`oc delete project %s`, namespace)
}

func createNS(namespace string) {
	util.Log.Debug("Creating namespace", namespace)
	util.ShellSilent(`oc new-project %s`, namespace)
}

func delNamespaceMesh(namespace string) {
	// Find namespace in members array
	tmp, _ := util.ShellSilent(`oc get smmr default -n %s --template='{{ .spec.members }}'`, meshNamespace)
	members := strings.Split(tmp, " ")
	position := arrayPositionFind(members, namespace)

	// If namespace exists
	if position > 0 {
		util.ShellSilent(`oc patch smmr default -n %s --type='json' -p='[{"op": "remove", "path": "/spec/members/%s"}]'`, meshNamespace, strconv.Itoa(position))
		configured := false
		util.Log.Debug("Waiting for the namespace to be deleted to the mesh: ", namespace)
		for !configured {
			msgConfigured, _ := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.configuredMembers }}'`, meshNamespace)
			msgPending, _ := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.pendingMembers }}'`, meshNamespace)
			if !strings.Contains(msgConfigured, namespace+" ") && !strings.Contains(msgPending, namespace+"]") {
				if !strings.Contains(msgPending, namespace+" ") && !strings.Contains(msgPending, namespace+"]") {
					deleteNS(namespace)
					configured = true
				}
			}
		}
	}
}

func createNamespaceMesh(namespace string) {
	createNS(namespace)
	util.ShellSilent(`oc patch smmr default -n %s --type='json' -p='[{"op": "add", "path": "/spec/members/-", "value":"%s"}]'`, meshNamespace, namespace)
	configured := false
	util.Log.Debug("Waiting for the namespace to be added to the mesh: ", namespace)
	for !configured {
		msg, _ := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.configuredMembers }}'`, meshNamespace)
		if strings.Contains(msg, namespace) {
			configured = true
		}
	}
}

func addNamespaceMesh(namespace string) {
	util.ShellSilent(`oc patch smmr default -n %s --type='json' -p='[{"op": "add", "path": "/spec/members/-", "value":"%s"}]'`, meshNamespace, namespace)
	configured := false
	util.Log.Debug("Waiting for the namespace to be added to the mesh: ", namespace)
	for !configured {
		msg, _ := util.ShellSilent(`oc get smmr default -n %s --template='{{ .status.configuredMembers }}'`, meshNamespace)
		if strings.Contains(msg, namespace) {
			configured = true
		}
	}
}
