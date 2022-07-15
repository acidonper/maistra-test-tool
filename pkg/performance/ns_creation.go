package performance

import (
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

func deleteNSBundle(min int, max int) {
	util.Log.Info("Deleting namespaces from ", min, " to ", max)
	for i := min; i < max; i++ {
		nsName := appNSPrefix + strconv.Itoa(i)
		delNamespaceMesh(nsName)
	}
}

func createNSBundle(min int, max int) {
	util.Log.Info("Creating namespaces from ", min, " to ", max)
	for i := min; i < max; i++ {
		nsName := appNSPrefix + strconv.Itoa(i)
		createNamespaceMesh(nsName)
	}
}

func testNamespaceAdditionTime(index int, acceptanceTime int, t *testing.T) {
	// Convert second to milisecond
	acceptanceTimeMS := acceptanceTime * 1000
	util.Log.Info("Adding a new namespace to the mesh (Acceptance Time: ", acceptanceTimeMS, " miliseconds)")

	// Create namespace
	nsName := appNSPrefix + strconv.Itoa(index) + "-measure"
	createNS(nsName)

	// Measure namespace addition to the mesh time
	startT := time.Now()
	addNamespaceMesh(nsName)
	duration := int(time.Since(startT) / (time.Second / time.Microsecond))

	// Check addition time
	if duration > acceptanceTimeMS {
		util.Log.Error("Acceptance time exceeded")
		t.Error("Acceptance time exceeded")
		t.FailNow()
	}
	util.Log.Info("Duration OK: ", duration)
}

func TestNSAdditionTime(t *testing.T) {
	util.Log.Info("** TEST: TestNSAdditionTime")

	// Iterate namespace bundles
	nsCounts := strings.Split(nsCountBundle, ",")
	for i, s := range nsCounts {
		// Define required variables
		count, _ := strconv.Atoi(s)
		acceptanceTime, _ := strconv.Atoi(nsAcceptanceTime)

		// Create required namespaces taking into account previous namespaces creation
		if i > 0 {
			tmp, _ := strconv.Atoi(nsCounts[i-1])
			createNSBundle(tmp, count)
		} else {
			createNSBundle(0, count)
		}

		// Wait for SMMR is Ready
		TestSMMR(t)

		// Measure a new namespace creation
		testNamespaceAdditionTime(count, acceptanceTime, t)
	}

	// clean namespaces
	util.Log.Info("cleaning up test")
	max, _ := strconv.Atoi(nsCounts[len(nsCounts)-1])
	deleteNSBundle(0, max)
	for _, s := range nsCounts {
		nsName := appNSPrefix + s + "-measure"
		delNamespaceMesh(nsName)
	}
}
