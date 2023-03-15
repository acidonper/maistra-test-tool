package performance

import (
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
)

func (b *Bookinfo) BookinfoInstall(mtls bool) error {
	util.Log.Info("Deploying Bookinfo in ", b.Namespace)
	err := util.KubeApplySilent(b.Namespace, bookinfoYaml)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(5) * time.Second)
	err = CheckPodRunning(b.Namespace, "app=details")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=ratings")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=reviews,version=v1")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=reviews,version=v2")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=reviews,version=v3")
	if err != nil {
		return err
	}
	err = CheckPodRunning(b.Namespace, "app=productpage")
	if err != nil {
		return err
	}

	util.Log.Debug("Creating Gateway from template")
	err = util.KubeApplyContentSilent(b.Namespace, util.RunTemplate(bookinfoGatewayTemplate, b))
	if err != nil {
		return err
	}

	util.Log.Debug("Creating VirtualService from template")
	err = util.KubeApplyContentSilent(b.Namespace, util.RunTemplate(bookinfoVirtualServiceTemplate, b))
	if err != nil {
		return err
	}

	util.Log.Debug("Creating Route from template")
	r := Route{Name: b.Namespace, Namespace: meshNamespace, Host: b.Host, Service: "istio-ingressgateway"}
	err = util.KubeApplyContentSilent(meshNamespace, util.RunTemplate(bookinfoSimpleHTTPRouteTemplate, r))
	if err != nil {
		return err
	}

	util.Log.Debug("Creating destination rules all")
	if mtls {
		err = util.KubeApplySilent(b.Namespace, bookinfoRuleAllTLSYaml)
		if err != nil {
			return err
		}
	} else {
		err = util.KubeApplySilent(b.Namespace, bookinfoRuleAllYaml)
		if err != nil {
			return err
		}
	}
	time.Sleep(time.Duration(10) * time.Second)
	return nil
}
