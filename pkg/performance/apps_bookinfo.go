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
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func (b *Bookinfo) BookinfoInstall(mtls bool) {
	util.Log.Info("Deploying Bookinfo in ", b.Namespace)
	util.KubeApplySilent(b.Namespace, bookinfoYaml)
	time.Sleep(time.Duration(5) * time.Second)
	CheckPodRunning(b.Namespace, "app=details")
	CheckPodRunning(b.Namespace, "app=ratings")
	CheckPodRunning(b.Namespace, "app=reviews,version=v1")
	CheckPodRunning(b.Namespace, "app=reviews,version=v2")
	CheckPodRunning(b.Namespace, "app=reviews,version=v3")
	CheckPodRunning(b.Namespace, "app=productpage")

	util.Log.Debug("Creating Gateway")
	util.KubeApplySilent(b.Namespace, bookinfoGateway)

	util.Log.Debug("Creating destination rules all")
	if mtls {
		util.KubeApplySilent(b.Namespace, bookinfoRuleAllTLSYaml)
	} else {
		util.KubeApplySilent(b.Namespace, bookinfoRuleAllYaml)
	}
	time.Sleep(time.Duration(10) * time.Second)
}

func (b *Bookinfo) BookinfoUninstall() {
	util.Log.Info("Cleanup Bookinfo in ", b.Namespace)
	util.KubeDelete(b.Namespace, bookinfoRuleAllYaml)
	util.KubeDelete(b.Namespace, bookinfoGateway)
	util.KubeDelete(b.Namespace, bookinfoYaml)
	time.Sleep(time.Duration(10) * time.Second)
}
