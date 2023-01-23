package performance

var (
	bookinfoGatewayTemplate = `
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: {{ .Name }}
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - {{ .Host }}
`

	bookinfoVirtualServiceTemplate = `
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: bookinfo
spec:
  hosts:
  - {{ .Host }}
  gateways:
  - {{ .Name }}
  http:
  - match:
    - uri:
        exact: /productpage
    - uri:
        prefix: /static
    - uri:
        exact: /login
    - uri:
        exact: /logout
    - uri:
        prefix: /api/v1/products
    route:
    - destination:
        host: productpage
        port:
          number: 9080
`

	bookinfoSimpleHTTPRouteTemplate = `
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  host: {{ .Host }}
  to:
    kind: Service
    name: {{ .Service }}
  tls:
    termination: edge
    insecureEdgeTerminationPolicy: Redirect
  port:
    targetPort: http2
`
)
