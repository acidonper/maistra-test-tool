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

	util.Log.Debug("Creating Gateway and VirtualService from template")
	util.KubeApplyContents(b.Namespace, util.RunTemplate(bookinfoGatewayTemplate, b))

	util.Log.Debug("Creating Route from template")
	r := Route{Name: b.Namespace, Namespace: meshNamespace, Host: b.Host, Service: "istio-ingressgateway"}
	util.KubeApplyContents(meshNamespace, util.RunTemplate(bookinfoSimpleHTTPRouteTemplate, r))

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
