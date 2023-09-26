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
	Metrics struct {
		HTTPReqReceiving struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
		} `json:"http_req_receiving"`
		Vus struct {
			Max   int `json:"max"`
			Value int `json:"value"`
			Min   int `json:"min"`
		} `json:"vus"`
		HTTPReqs struct {
			Rate  float64 `json:"rate"`
			Count int     `json:"count"`
		} `json:"http_reqs"`
		HTTPReqDuration struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
		} `json:"http_req_duration"`
		GrpcReqDuration struct {
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
		} `json:"grpc_req_duration"`
		HTTPReqConnecting struct {
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
		} `json:"http_req_connecting"`
		Checks struct {
			Passes int     `json:"passes"`
			Fails  int     `json:"fails"`
			Value  float64 `json:"value"`
		} `json:"checks"`
		HTTPReqBlocked struct {
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
		} `json:"http_req_blocked"`
		IterationDuration struct {
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
		} `json:"iteration_duration"`
		VusMax struct {
			Min   int `json:"min"`
			Max   int `json:"max"`
			Value int `json:"value"`
		} `json:"vus_max"`
		HTTPReqFailed struct {
			Passes int     `json:"passes"`
			Fails  int     `json:"fails"`
			Value  float64 `json:"value"`
		} `json:"http_req_failed"`
		Errors struct {
			Fails  int     `json:"fails"`
			Passes int     `json:"passes"`
			Value  float64 `json:"value"`
		} `json:"errors"`
		DataReceived struct {
			Rate  float64 `json:"rate"`
			Count int     `json:"count"`
		} `json:"data_received"`
		HTTPReqWaiting struct {
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
		} `json:"http_req_waiting"`
		DataSent struct {
			Rate  float64 `json:"rate"`
			Count int     `json:"count"`
		} `json:"data_sent"`
		Iterations struct {
			Count int     `json:"count"`
			Rate  float64 `json:"rate"`
		} `json:"iterations"`
		HTTPReqSending struct {
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
		} `json:"http_req_sending"`
		HTTPReqDurationExpectedResponseTrue struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
		} `json:"http_req_duration{expected_response:true}"`
		GroupDuration struct {
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
			Avg float64 `json:"avg"`
		} `json:"group_duration"`
		HTTPReqTLSHandshaking struct {
			Avg float64 `json:"avg"`
			Min float64 `json:"min"`
			Med float64 `json:"med"`
			Max float64 `json:"max"`
			P90 float64 `json:"p(90)"`
			P95 float64 `json:"p(95)"`
		} `json:"http_req_tls_handshaking"`
	} `json:"metrics"`
}
