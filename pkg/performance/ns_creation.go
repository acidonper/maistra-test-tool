package performance

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func arrayPositionFind(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

func deleteNSBundle(min int, max int) error {
	util.Log.Info("Deleting namespaces from ", min, " to ", max)
	for i := min; i < max; i++ {
		nsName := appNSPrefix + strconv.Itoa(i)
		err := delNamespaceMesh(nsName)
		if err != nil {
			return err
		}
	}
	return nil
}

func createNSBundle(min int, max int) error {
	util.Log.Info("Creating namespaces from ", min, " to ", max)
	for i := min; i < max; i++ {
		nsName := appNSPrefix + strconv.Itoa(i)
		err := createNamespaceMesh(nsName)
		if err != nil {
			return err
		}
	}
	return nil
}

func testNamespaceAdditionTime(index int, acceptanceTime int) error {
	// Convert second to milisecond
	acceptanceTimeMS := acceptanceTime * 1000
	util.Log.Info("Adding a new namespace to the mesh (Acceptance Time: ", acceptanceTimeMS, " miliseconds)")

	// Create namespace
	nsName := appNSPrefix + strconv.Itoa(index) + "-measure"
	err := createNS(nsName)
	if err != nil {
		return err
	}

	// Measure namespace addition to the mesh time
	startT := time.Now()
	errAddNS := addNamespaceMesh(nsName)
	if errAddNS != nil {
		return err
	}
	duration := int(time.Since(startT) / (time.Second / time.Microsecond))

	// Check addition time
	if duration > acceptanceTimeMS {
		util.Log.Error("Acceptance time exceeded")
		return fmt.Errorf("Acceptance time exceeded")
	}
	util.Log.Info("Duration OK: ", duration)
	return nil
}

func TestNSAdditionTime(t *testing.T) {
	util.Log.Info("** TEST: TestNSAdditionTime")

	// Test SMMR and SMCP are Ready
	TestSMCP(t)
	TestSMMR(t)

	// Iterate namespace bundles
	nsCounts := strings.Split(nsCountBundle, ",")
	var step int
	for i, s := range nsCounts {
		step = i
		// Define required variables
		count, _ := strconv.Atoi(s)
		acceptanceTime, _ := strconv.Atoi(nsAcceptanceTime)

		// Create required namespaces taking into account previous namespaces creation
		if i > 0 {
			tmp, _ := strconv.Atoi(nsCounts[i-1])
			err := createNSBundle(tmp, count)
			if err != nil {
				t.Error(err.Error())
				break
			}
		} else {
			err := createNSBundle(0, count)
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

	// clean namespaces
	util.Log.Info("cleaning up test")
	max, _ := strconv.Atoi(nsCounts[step])
	err := deleteNSBundle(0, max)
	if err != nil {
		t.Error(err.Error())
	}
	for _, s := range nsCounts {
		nsName := appNSPrefix + s + "-measure"
		delNamespaceMesh(nsName)
	}
}
