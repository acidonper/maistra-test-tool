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

func (b *JumpApp) JumpappInstall() {
	util.Log.Info("Deploying JumpAppin ", b.Namespace)
	util.KubeApply(b.Namespace, jumpappYaml)
	time.Sleep(time.Duration(5) * time.Second)
	CheckPodRunning(b.Namespace, "app=back-golang")
	CheckPodRunning(b.Namespace, "app=back-python")
	CheckPodRunning(b.Namespace, "app=back-quarkus")
	CheckPodRunning(b.Namespace, "app=back-springboot")

	util.KubeApply(b.Namespace, jumpappNetworking)

	time.Sleep(time.Duration(10) * time.Second)
}

func (b *JumpApp) JumpappUninstall() {
	util.Log.Info("Cleanup JumpApp in ", b.Namespace)
	util.KubeDelete(b.Namespace, jumpappNetworking)
	util.KubeDelete(b.Namespace, jumpappYaml)
	time.Sleep(time.Duration(10) * time.Second)
}
