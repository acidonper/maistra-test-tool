package performance

type Bookinfo struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Host      string `json:"host,omitempty"`
}

type JumpApp struct {
	Namespace string `json:"namespace,omitempty"`
}

type Route struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Host      string `json:"host,omitempty"`
	Service   string `json:"service,omitempty"`
}

type PromResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name                 string `json:"__name__"`
				App                  string `json:"app"`
				Instance             string `json:"instance"`
				Istio                string `json:"istio"`
				IstioIoRev           string `json:"istio_io_rev"`
				Job                  string `json:"job"`
				KubernetesNamespace  string `json:"kubernetes_namespace"`
				KubernetesPodName    string `json:"kubernetes_pod_name"`
				Le                   string `json:"le"`
				MaistraControlPlane  string `json:"maistra_control_plane"`
				PodTemplateHash      string `json:"pod_template_hash"`
				SidecarIstioIoInject string `json:"sidecar_istio_io_inject"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
