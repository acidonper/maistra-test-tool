package tests

import (
	"testing"

	performance "github.com/maistra/maistra-test-tool/pkg/performance"
)

var performanceCases = []testing.InternalTest{
	{
		Name: "A1",
		F:    performance.TestSMCPInstalled,
	},
	{
		Name: "A2",
		F:    performance.TestSMMRInstalled,
	},
	{
		Name: "CP1",
		F:    performance.TestNSAdditionTime,
	},
	{
		Name: "A3",
		F:    performance.TestNSAdditionTimeClean,
	},
	{
		Name: "A4",
		F:    performance.CreateTestCPObjects,
	},
	{
		Name: "CP2.1",
		F:    performance.TestXDSPushes,
	},
	{
		Name: "CP2.2",
		F:    performance.TestXDSErrors,
	},
	{
		Name: "CP3.1",
		F:    performance.TestIstiodMem,
	},
	{
		Name: "CP3.2",
		F:    performance.TestIstiodCpu,
	},
	{
		Name: "A5",
		F:    performance.DeleteTestCPObjects,
	},
	{
		Name: "A6",
		F:    performance.CreateTestDPObjects,
	},
	{
		Name: "A7",
		F:    performance.GenerateTrafficLoadK6,
	},
	{
		Name: "DP1.1",
		F:    performance.TestIstioProxiesMem,
	},
	{
		Name: "DP1.2",
		F:    performance.TestIstioProxiesCpu,
	},
	{
		Name: "DP2.1",
		F:    performance.TestIstioIngressMem,
	},
	{
		Name: "DP2.2",
		F:    performance.TestIstioIngressCpu,
	},
	{
		Name: "DP3.1",
		F:    performance.TestIstioEgressMem,
	},
	{
		Name: "DP3.2",
		F:    performance.TestIstioEgressCpu,
	},
	{
		Name: "A8",
		F:    performance.AnalyseLoadK6Output,
	},
	{
		Name: "A9",
		F:    performance.DeleteTestDPObjects,
	},
}
