package tests

import (
	"testing"

	performance "github.com/maistra/maistra-test-tool/pkg/performance"
)

var t = &testing.T{}

var performanceCases = []testing.InternalTest{
	testing.InternalTest{
		Name: "A1",
		F:    performance.TestSMCPInstalled,
	},
	testing.InternalTest{
		Name: "A2",
		F:    performance.TestSMMRInstalled,
	},
	testing.InternalTest{
		Name: "CP1",
		F:    performance.TestNSAdditionTime,
	},
	testing.InternalTest{
		Name: "A3",
		F:    performance.TestNSAdditionTimeClean,
	},
	testing.InternalTest{
		Name: "A4",
		F:    performance.CreateTestCPObjects,
	},
	testing.InternalTest{
		Name: "CP2.1",
		F:    performance.TestXDSPushes,
	},
	testing.InternalTest{
		Name: "CP3.1",
		F:    performance.TestIstiodMem,
	},
	testing.InternalTest{
		Name: "CP3.2",
		F:    performance.TestIstiodCpu,
	},
	testing.InternalTest{
		Name: "A5",
		F:    performance.DeleteTestCPObjects,
	},
	testing.InternalTest{
		Name: "A6",
		F:    performance.CreateTestDPObjects,
	},
	testing.InternalTest{
		Name: "A7",
		F:    performance.GenerateTrafficLoadK6,
	},
	testing.InternalTest{
		Name: "DP1.1",
		F:    performance.TestIstioProxiesMem,
	},
	testing.InternalTest{
		Name: "DP1.2",
		F:    performance.TestIstioProxiesCpu,
	},
	testing.InternalTest{
		Name: "DP2.1",
		F:    performance.TestIstioIngressMem,
	},
	testing.InternalTest{
		Name: "DP2.2",
		F:    performance.TestIstioIngressCpu,
	},
	testing.InternalTest{
		Name: "DP3.1",
		F:    performance.TestIstioEgressMem,
	},
	testing.InternalTest{
		Name: "DP3.2",
		F:    performance.TestIstioEgressCpu,
	},
	// testing.InternalTest{
	// 	Name: "A8",
	// 	F:    performance.analyseLoadOutput,  TODO
	// },
	testing.InternalTest{
		Name: "A9",
		F:    performance.DeleteTestDPObjects,
	},
}
