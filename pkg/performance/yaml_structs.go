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

type K6Response struct {
	RootGroup struct {
		Groups struct {
		} `json:"groups"`
		Checks struct {
		} `json:"checks"`
		Name string `json:"name"`
		Path string `json:"path"`
		ID   string `json:"id"`
	} `json:"root_group"`
	Metrics struct {
		HTTPReqTLSHandshaking struct {
			Med int     `json:"med"`
			Max float64 `json:"max"`
			P90 int     `json:"p(90)"`
			P95 int     `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min int     `json:"min"`
		} `json:"http_req_tls_handshaking"`
		Iterations struct {
			Count int     `json:"count"`
			Rate  float64 `json:"rate"`
		} `json:"iterations"`
		HTTPReqReceiving struct {
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
		} `json:"http_req_receiving"`
		HTTPReqDuration struct {
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
		} `json:"http_req_duration"`
		HTTPReqs struct {
			Rate  float64 `json:"rate"`
			Count int     `json:"count"`
		} `json:"http_reqs"`
		DataReceived struct {
			Count int     `json:"count"`
			Rate  float64 `json:"rate"`
		} `json:"data_received"`
		HTTPReqConnecting struct {
			Max float64 `json:"max"`
			P90 int     `json:"p(90)"`
			P95 int     `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min int     `json:"min"`
			Med int     `json:"med"`
		} `json:"http_req_connecting"`
		HTTPReqSending struct {
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
		} `json:"http_req_sending"`
		DataSent struct {
			Count int     `json:"count"`
			Rate  float64 `json:"rate"`
		} `json:"data_sent"`
		HTTPReqDurationExpectedResponseTrue struct {
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
		} `json:"http_req_duration{expected_response:true}"`
		IterationDuration struct {
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
		} `json:"iteration_duration"`
		HTTPReqFailed struct {
			Passes int `json:"passes"`
			Fails  int `json:"fails"`
			Value  int `json:"value"`
		} `json:"http_req_failed"`
		VusMax struct {
			Value int `json:"value"`
			Min   int `json:"min"`
			Max   int `json:"max"`
		} `json:"vus_max"`
		HTTPReqWaiting struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
		} `json:"http_req_waiting"`
		Vus struct {
			Value int `json:"value"`
			Min   int `json:"min"`
			Max   int `json:"max"`
		} `json:"vus"`
		HTTPReqBlocked struct {
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
		} `json:"http_req_blocked"`
	} `json:"metrics"`
}
