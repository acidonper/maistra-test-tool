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

package performance

import (
	"github.com/maistra/maistra-test-tool/pkg/util"
)

var (
	smcpName              string            = util.Getenv("SMCPNAME", "basic")
	meshNamespace         string            = util.Getenv("MESHNAMESPACE", "istio-system")
	appNSPrefix           string            = util.Getenv("APPNSPREFIX", "app-perf-test")
	nsCountBundle         string            = util.Getenv("NSCOUNTBUNDLE", "10,50,100,500")
	nsAcceptanceTime      string            = util.Getenv("NSACCEPTANCETIME", "5")
	xdsPushAcceptanceTime string            = util.Getenv("XDSPUSHACCEPTANCETIME", "1")
	istiodAcceptanceMem   string            = util.Getenv("ISTIODACCEPTANCEMEM", "1024")
	istiodAcceptanceCpu   string            = util.Getenv("ISTIODACCEPTANCECPU", "1000")
	prometheusAPIMap      map[string]string = map[string]string{
		"xds_ppctc":  "pilot_proxy_convergence_time_count",
		"xds_ppctb":  "pilot_proxy_convergence_time_bucket{le=\"" + xdsPushAcceptanceTime + "\"}",
		"istiod_mem": "process_virtual_memory_bytes{app=\"istiod\"}",
		"istiod_cpu": "irate(process_cpu_seconds_total{app=\"istiod\"}[1m])",
	}
)
