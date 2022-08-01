// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"

	"github.com/maistra/maistra-test-tool/pkg/performance"
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
	// testing.InternalTest{
	// 	Name: "A4",
	// 	F:    performance.createTestCPObjects,  TODO
	// },
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
	// testing.InternalTest{
	// 	Name: "A5",
	// 	F:    performance.deleteTestCPObjects,  TODO
	// },
	// testing.InternalTest{
	// 	Name: "A6",
	// 	F:    performance.createTestDPObjects,  TODO
	// },
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
		F:    performance.TestIstioIngressProxiesMem,
	},
	testing.InternalTest{
		Name: "DP2.2",
		F:    performance.TestIstioIngressProxiesCpu,
	},
	testing.InternalTest{
		Name: "DP3.1",
		F:    performance.TestIstioEgressProxiesMem,
	},
	testing.InternalTest{
		Name: "DP3.2",
		F:    performance.TestIstioEgressProxiesCpu,
	},
	// testing.InternalTest{
	// 	Name: "A6",
	// 	F:    performance.deleteTestDPObjects,  TODO
	// },
}
