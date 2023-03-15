package performance

import (
	"strconv"
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func CreateTestCPObjects(t *testing.T) {
	util.Log.Info("** TEST: CreateTestCPObjects")
	apps, _ := strconv.Atoi(testCPApps)
	err := createAppBundle(trafficLoadApp, apps, "controlplane")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
}

func DeleteTestCPObjects(t *testing.T) {
	util.Log.Info("** TEST: DeleteTestCPObjects")
	apps, _ := strconv.Atoi(testCPApps)
	err := deleteAppBundle(trafficLoadApp, apps, "controlplane")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
}

func CreateTestDPObjects(t *testing.T) {
	util.Log.Info("** TEST: CreateTestDPObjects")
	apps, _ := strconv.Atoi(testDPApps)
	err := createAppBundle(trafficLoadApp, apps, "dataplane")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
}

func DeleteTestDPObjects(t *testing.T) {
	util.Log.Info("** TEST: DeleteTestDPObjects")
	apps, _ := strconv.Atoi(testDPApps)
	err := deleteAppBundle(trafficLoadApp, apps, "dataplane")
	if err != nil {
		util.Log.Error(err)
		t.FailNow()
	}
}
