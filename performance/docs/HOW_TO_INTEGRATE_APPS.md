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

Here is the explanation of some of the functions we used inside the installation of Bookinfo that you can use for your app too:

* `util.KubeApplySilent`: with this function you can do an `oc apply` of a YAML in the namespace you pass by parameter. It doesn't show the output of the command. We recommend it to keep it that way so all logs are handled by the logger.
* `CheckPodRunning`: this is a useful function that waits until a pod finishes running. You need to pass the namespace and a label by parameter.
* `util.KubeApplyContentSilent`: similar to the `KubeApplySilent`, but this time is useful for applying a template if you pass by parameter the result of the `util.RunTemplate` function.

You can use all of these functions or others you find useful in the maistra [`utils`](../../pkg/util/) directory.

## Add your app to the creation function

Go to the `createAppBundle` function inside the [utils](../../pkg/performance/utils.go) file and add your new app to the `case`. 

For your app you will have to:

1. Give a prefix for your namespace. For example: `nsName = bookinfoNSPrefix + strconv.Itoa(i)`
2. Create the namespaces: `err := createNSMesh(nsName)`
3. Create all the dynamic info you need using your struct.
4. Install your app with the function you created.

If you need to use some constants. We are storing them in the [`yaml_vars.go`](../../pkg/performance/yaml_vars.go) file.

## Add your app to the deletion function

Go to the `deleteAppBundle` function inside the [utils](../../pkg/performance/utils.go) file and add your new app to the `case`. The code you will need to add will be like the following example:

```go
case "bookinfo":
  prefix = bookinfoNSPrefix
```

Just add the prefix of your namespaces names and the rest of the code will delete those.

## Add your app to the filling function

Before adding your app to the filling function. You will need to calculate how much you think your app limits are in regard to CPU and Memory. And create the constants in the [`yaml_vars.go`](../../pkg/performance/yaml_vars.go) file. 

> **IMPORTANT:** If your app has several microservices, make the sum of the limits of all of them.

Now go to the `calculateAppsFillCluster` function inside the same [utils](../../pkg/performance/utils.go) file and again add your app to the `switch` using the CPU and Memory limit you calculated before. 

As in this example:

```go
// Case in utils.go
	case "bookinfo":
		cpuLimit, _ = strconv.Atoi(bookinfoTotalCPU)
		memLimit, _ = strconv.Atoi(bookinfoTotalMem)

// Constants in yaml_vars.go
bookinfoTotalCPU       = "600"  // (50 + 50 + 50 + 50 + 50 + 50) * 2
bookinfoTotalMem       = "1134" // 250 + 250 + 250 + 128 + 128 + 128
```

As you can see in the comments of the constants, those are the calculations of all the microservices' limits.

## Add your app to the traffic load function

In the same [utils](../../pkg/performance/utils.go) file, go to the `generateSimpleTrafficLoadK6` function and add your app to the `if` with the protocols your app support. 

> **IMPORTANT:** We only support doing the tests on one protocol at a time. So if your app supports `grpc` and `http`, you will have to create both conditionals in the `if`

In this function, you will need to create the URL to which the `k6` will do the traffic loading and then call the `execK6SyncTest` function. 

For example:

```go
if app == "bookinfo" && protocol == "http" {

		routeHost, errRoute := getRouteHost(appName, meshNamespace)
		if errRoute != nil {
			return fmt.Errorf("route %s not found in namespace %s", appName, meshNamespace)
		} else {
			url = "http://" + routeHost + "/productpage"
		}
		_, err = execK6SyncTest(testVUs, testDuration, url, "http-basic.js", reportFile)

		if err != nil {
			return err
		}

	} else {
		errorMsg := fmt.Errorf("application %s and protocol %s combination not supported", app, protocol)
		return errorMsg
	}
```

In this case, we get the route, construct the URL adding the `http://` and the `/productpage` endpoint needed for the traffic load to work. Then, we pass that URL to the `execK6SyncTest` along with other parameters like the number of virtual users (set by env var), duration of the test (also set by env var), the reportFile (constant in `yaml_vars.go`) and the K6 script.

For the K6 script that you need to execute the load testing: If your application is `http` you can use the already created [`http-basic.js`](../../objects/k6/http-basic.js) file or you can create one inside the [k6](../../objects/k6/) objects folder and pass it as a parameter to the `execK6SyncTest` function. If you need more information about `k6` and their scripts, you can find it in their official [documentation](https://k6.io/docs/).

## Test your app

Now you just need to test your new app by changing the `TRAFFICLOADAPP` and `TRAFFICLOADPROTOCOL` environment variables to your app name and the protocol you want to use respectively.