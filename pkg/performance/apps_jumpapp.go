package performance

import (
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func (b *JumpApp) JumpappInstall() error {
	util.Log.Info("Deploying JumpAppin ", b.Namespace)
	err := util.KubeApplySilent(b.Namespace, jumpappYaml)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(5) * time.Second)
	err = CheckPodRunning(b.Namespace, "app=back-golang")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=back-python")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=back-quarkus")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=back-springboot")
	if err != nil {
		return err
	}

	err = util.KubeApplySilent(b.Namespace, jumpappNetworking)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(10) * time.Second)
	return nil
}

func (b *JumpApp) JumpappUninstall() error {
	util.Log.Info("Cleanup JumpApp in ", b.Namespace)
	err := util.KubeDelete(b.Namespace, jumpappNetworking)
	if err != nil {
		return err
	}
	err = util.KubeDelete(b.Namespace, jumpappYaml)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(10) * time.Second)
	return nil
}
