package performance

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func testNamespaceAdditionTime(index int, acceptanceTime int) error {
	// Convert second to milisecond
	acceptanceTimeMS := acceptanceTime * 1000
	util.Log.Info("Adding a new namespace to the mesh (Acceptance Time: ", acceptanceTimeMS, " milliseconds)")

	// Create namespace
	nsName := appNSPrefix + strconv.Itoa(index) + "-measure"
	err := createNS(nsName)
	if err != nil {
		return err
	}

	// Measure namespace addition to the mesh time
	startT := time.Now()
	errAddNS := addNSToMesh(nsName)
	if errAddNS != nil {
		return err
	}
	duration := int(time.Since(startT) / (time.Second / time.Microsecond))

	// Check addition time
	if duration > acceptanceTimeMS {
		return fmt.Errorf("acceptance time exceeded")
	}
	util.Log.Info("Duration OK: ", duration, " milliseconds")
	return nil
}

func TestNSAdditionTime(t *testing.T) {
	util.Log.Info("** TEST: TestNSAdditionTime")

	// Test SMMR and SMCP are Ready
	TestSMCP(t)
	TestSMMR(t)

	// Iterate namespace bundles
	nsCounts := strings.Split(nsCountBundle, ",")
	for i, s := range nsCounts {
		// Define required variables
		count, _ := strconv.Atoi(s)
		acceptanceTime, _ := strconv.Atoi(nsAcceptanceTime)

		// Create required namespaces taking into account previous namespaces creation
		if i > 0 {
			tmp, _ := strconv.Atoi(nsCounts[i-1])
			err := createNSBundle(tmp, count, appNSPrefix)
			if err != nil {
				t.Error(err.Error())
				break
			}
		} else {
			err := createNSBundle(0, count, appNSPrefix)
			if err != nil {
				t.Error(err.Error())
				break
			}
		}

		// Wait for SMMR is Ready
		TestSMMR(t)

		// Measure a new namespace creation
		err := testNamespaceAdditionTime(count, acceptanceTime)
		if err != nil {
			t.Error(err.Error())
			break
		}
	}
}

func TestNSAdditionTimeClean(t *testing.T) {
	util.Log.Info("** TEST: TestNSAdditionTimeClean")

	// Test SMMR and SMCP are Ready
	TestSMCP(t)
	TestSMMR(t)

	// Iterate namespace bundles
	nsCounts := strings.Split(nsCountBundle, ",")

	// clean namespaces
	util.Log.Info("Cleaning up TestNSAdditionTime objects")
	max, _ := strconv.Atoi(nsCounts[len(nsCounts)-1])
	err := deleteNSBundle(0, max, appNSPrefix)
	if err != nil {
		t.Error(err.Error())
	}
	for _, s := range nsCounts {
		util.Log.Info("Deleting namespace " + appNSPrefix + s + "-measure")
		nsName := appNSPrefix + s + "-measure"
		deleteNSMesh(nsName)
	}

	// Test SMMR and SMCP are Ready
	TestSMCP(t)
	TestSMMR(t)
}
