# How to integrate other apps into the tool

As of this moment, the only application integrated with the Performance Test Tool is [Bookinfo](https://github.com/Maistra/istio/tree/maistra-2.4/samples/bookinfo). This is the most used app to demo Service Mesh capabilities, but it is true that probably a lot of customers' apps won't have the same effects in performance that the Bookinfo app has. That's why this guide is created with the idea to make easier the integration of new apps, but it will require some code.

## Fork the repository

The first step as always is to fork the repository. Our upstream right now is https://github.com/acidonper/maistra-test-tool/tree/feature/performance.

> **IMPORTANT:** Since this tool is not yet integrated with the `maistra-test-tool` we are still using the `feature/performance` branch as our main branch.

## Create a directory for your app in the objects folder

If you have some immutable manifests (eg: deployment, service) to deploy your application, store all of them in a directory inside the [objects](../../objects/) folder. Create it with the name of your app.

## Create some templates for the dynamically adjustable files

If some of your files need to be dynamically changed, you will need to create some templates. For that, go to the [templates](../../pkg/performance/yaml_tmpl.go) file and add them there as variables. For example:

```go
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
```

In this case, the `.Name` and `.Host` are created dynamically, and for them to be referenced here, they are part of the following struct in the [structs](../../pkg/performance/yaml_structs.go):

```go
type Bookinfo struct {
    Namespace string `json:"namespace,omitempty"`
    Name      string `json:"name,omitempty"`
    Host      string `json:"host,omitempty"`
}
```

So you will have to create something similar for your app to dynamically create some of the Kubernetes resources.

## Add the code to install your new app

Create a new Go file in the [performance](../../pkg/performance/) package with the following name `apps_<your app name>.go`. Use the [`apps_bookinfo.go`](../../pkg/performance/apps_bookinfo.go) as an example of the code you need to add. 

> **TIP:** You only need to install the app, since the tool is going to uninstall it when deleting the namespace.

## Add your app to the creation method

Go to the `createAppBundle` method inside the [utils](../../pkg/performance/utils.go) file and add your new app to the `case`. 

For your app you will have to:

1. Give a prefix for your namespace. For example: `nsName = bookinfoNSPrefix + strconv.Itoa(i)`
2. Create the namespaces: `err := createNSMesh(nsName)`
3. Create all the dynamic info you need using your struct.
4. Install your app with the method you created.

If you need to use some constants. We are storing them in the [`yaml_vars.go`](../../pkg/performance/yaml_vars.go) file.

## Add your app to the traffic load method

In the same [utils](../../pkg/performance/utils.go) file, go to the `generateSimpleTrafficLoadK6` method and add your app to the `if` with the protocols your app support. 

> **IMPORTANT:** We only support doing the tests on one protocol at a time. So if your app supports `grpc` and `http`, you will have to create both conditionals in the `if`

In this method, you will need to create the URL to which the `k6` will do the traffic loading and then call the `execK6SyncTest` function. 

You will also need a script to execute the load testing. If your application is `http` you can use the already created [`http-basic.js`](../../objects/k6/http-basic.js) file or you can create one inside the [k6](../../objects/k6/) objects folder and pass it as a parameter to the `execK6SyncTest` function. If you need more information about `k6` and their scripts, you can find it in their official [documentation](https://k6.io/docs/).

## Test your app

Now you just need to test your new app by changing the `TRAFFICLOADAPP` and `TRAFFICLOADPROTOCOL` environment variables to your app name and the protocol you want to use respectively.