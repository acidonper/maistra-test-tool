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
