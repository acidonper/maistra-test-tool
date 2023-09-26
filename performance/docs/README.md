# CUSTOM - Openshif Service Mesh Performance tool

A performance tool for running specific Openshif Service Mesh use cases on an OpenShift 4.x cluster obtaining the different components' specific metrics.

The idea behind this inititative is to generate a performance testing tool that certifies Openshift Service Mesh in different scenarios. 

## Versions

| Name      | Version  |
| --------- | -------- |
| OS        | Linux    |
| Golang    | 1.13+    |
| OpenSSl   | 1.1.1+   |
| oc client | 4.10.25+ |
| k6        | v0.39.0+ |

## Prerequisites

* RedHat Service Mesh Operator has been installed on the OpenShift cluster
* Service Mesh Control Plane installed and ready 
* Service Mesh Member Role ready with 1 configured member at least
* By default, OSSM is installed with a integrated Prometheus instance configured. This tool supports [user-workload monitoring](https://docs.openshift.com/container-platform/4.13/monitoring/enabling-monitoring-for-user-defined-projects.html). For this, OSSM must be configured to use _user-workload monitoring_: [doc](https://docs.openshift.com/container-platform/4.13/service_mesh/v2x/ossm-observability.html#ossm-integrating-with-user-workload-monitoring_observability)

## Testing Prerequisite

* An `oc` client can be downloaded from [mirror openshift-v4 clients](https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/) and extract the `openshift-client-...tar.gz` file and move both `oc` and `kubectl` binaries into a local PATH directory.
* An `k6` client can be downloaded and move the binary into a local PATH directory.

## Run Performance Test Locally

* Access an OpenShift cluster from command line before running tests 

```$bash
oc login -u [user] -p [token] --server=[OCP API server]
```

* Update the environment variables 
* To run all the test cases

```$bash
cd performance/tests
go test -timeout 2h -v
```

NOTE: The `-timeout` flag is necessary when running all tests or several major test cases. Otherwise, a `go test` command falls into panic after 10 minutes.

* To run a single test case 

```$bash
cd performance/tests
go test -run A1 -timeout 2h -v
```

NOTE: Test cases shortname and mapping are in the `tests/test_cases.go` file.

## Run Performance Test from a Container Image

It is possible to generate a container image in order to execute the performance tests from a imutable container. Please follow the next procedure to execute the performance test using *podman*:

- Modify the required environment variables

```$bash
vi performance/Dockerfile
...
ENV OCP_CRED_USR=XXX
ENV OCP_CRED_PSW=XXX
ENV OCP_API_URL=XXX
...
```

- Build the container image with the environment variables, binaries and the respective code

```$bash
podman build . -t maistra-test-tool-performance -f performance/Dockerfile
```

- Run the performance tests

```$bash
podman run -it localhost/maistra-test-tool-performance
```

## Integrate other apps to the tool

If you want to use other apps besides Bookinfo, you can follow the instructions [here](./HOW_TO_INTEGRATE_APPS.md).

## Test Cases

This section tries to explain the different test cases. It is important to bear in mind that every test case has a specific prefix code that defines the kind of test:

* Ax -> Auxiliar procedures than support test cases execution
* CPx -> Define tests created for analyzing the Control Plane elements (istiod)
* DPx -> Define tests created for analyzing the Data Plane elements (Proxies, ingress or egress)

On the other hand, it is also important to take into account that multiple environment variables that define the performance test operations are defined in the following [file](../tests/test.env)

### A1 - TestSMCPInstalled

Check the Red Hat Service Mesh Control Plane is ready (All components are running correctly).

### A2 - TestSMMRInstalled

Check the Service Mesh Member Role is ready (All namespaces in the mesh are correctly added).

### CP1 - TestNSAdditionTime

Check the time for adding a new namespace in the mesh is lower than a specific time once a set of namespaces have been added previously. This time is defined in the **NSACCEPTANCETIME** environment variable in **seconds**.  

On the other hand, **NSCOUNTBUNDLE** define the number of namespaces added before the test is executed.

### A3 - TestNSAdditionTimeClean

Clean TestNSAdditionTime test namespaces deployed previously.

### A4 - CreateTestCPObjects

Create a set of applications in order to prepare the control plane tests environment. The number of application is defined by the **TESTCPAPPS** environment variable.

### CP2.1 - TestXDSPushes 

Check if the XDS pushes performed by the Control Plane pods (*istiod*) have been proper during the previous tests and are lower than a certain period of time defined in **XDSPUSHACCEPTANCETIME** environment variable in seconds. 

### CP3.1 - TestIstiodMem

Check the Control Plane's pods (*istiod*) memory in order to ensure they are lower than a certain amount of memory. This memory is defined in Megabytes in the **ISTIODACCEPTANCEMEM** environment variable.


### CP3.2 - TestIstiodCpu

Check the Control Plane's pods (*istiod*) CPU in order to ensure they are lower than a certain amount of CPU. This CPU is defined in Milicores in the **ISTIODACCEPTANCECPU** environment variable.

### A5 - DeleteTestCPObjects

Clean up the CreateTestCPObjects test objects deployed previously.

### A6 - CreateTestDPObjects

Create a set of applications in order to prepare the data plane tests environment. The number of application is defined by the **TESTDPAPPS** environment variable.

It is possible to deploy a certain number of application to fill the cluster using all the compute resoureces available. In order to enable this option, it is required to define **TESTDPAPPSFILL** a **true**. It is important to bear in mind that *TESTDPAPPS* environment variable will be ignored.

### A7 - GenerateTrafficLoadK6

Execute a k6 load test in order to generate applications' load. The idea is to emulate a real flow situation for extracting the performance metrics of the current Openshift Service Mesh architecture deployed.

Please keep in mind that the number of virtual users and the test period of time are defined by **TESTVUS** and **TESTDURATION** environment variables respectively.

### DP1.1 - TestIstioProxiesMem

Check the Data Plane's pods (*sidecar envoy proxies*) memory in order to ensure they are lower than a certain amount of memory. This memory is defined in Megabytes in the **ISTIOPROXIESACCEPTANCEMEM** environment variable.

### DP1.2 - TestIstioProxiesCpu

Check the Data Plane's pods (*sidecar envoy proxies*) CPU in order to ensure they are lower than a certain amount of CPU. This CPU is defined in Milicores in the **ISTIOPROXIESDACCEPTANCECPU** environment variable.

### DP2.1 - TestIstioIngressMem

Check the Data Plane's pods (*ingress envoy proxies*) memory in order to ensure they are lower than a certain amount of memory. This memory is defined in Megabytes in the **ISTIOINGRESSPROXIESACCEPTANCEMEM** environment variable.

### DP2.2 - TestIstioIngressCpu

Check the Data Plane's pods (*ingress envoy proxies*) CPU in order to ensure they are lower than a certain amount of CPU. This CPU is defined in Milicores in the **ISTIOINGRESSPROXIESACCEPTANCECPU** environment variable.

### DP3.1 - TestIstioEgressMem

Check the Data Plane's pods (*egress envoy proxies*) memory in order to ensure they are lower than a certain amount of memory. This memory is defined in Megabytes in the **ISTIOEGRESSPROXIESACCEPTANCEMEM** environment variable.

### DP3.2 - TestIstioEgressCpu

Check the Data Plane's pods (*egress envoy proxies*) CPU in order to ensure they are lower than a certain amount of CPU. This CPU is defined in Milicores in the **ISTIOEGRESSPROXIESACCEPTANCECPU** environment variable.

### A8 - AnalyseLoadK6Output

Check if the k6 load test output (average request time in percentil 95 metric) is lower that a certain period of time. This period of time is defined by the **REQAVG95PACCEPTANCETIME** environment variable in milliseconds.

### A9 - DeleteTestDPObjects

Clean up the CreateTestDPObjects test objects deployed previously.

## Test Cases: Environment variables
The following environment variables can be set to configure how the test will be executed. These variables are stored in the [test.env](../tests/test.env) file.

| Name                             | Value                                                         |
| -----------                      | -----------                                                   |
| MESHNAMESPACE                    | Istio Control Plane's namespace                               |
| SMCPNAME                         | ServiceMeshControlPlane's name                                |
| APPNSPREFIX                      | k8s-2                                                         |
| GODEBUG                          | Go configuration                                              |
| NSCOUNTBUNDLE                    | Number of namespaces added to the mesh                        |
| NSACCEPTANCETIME                 | Time limit of adding a new namespace to the mesh              |
| XDSPUSHACCEPTANCETIME            | Time in seconds of XDS pushes                                 |
| ISTIODACCEPTANCEMEM              | Istiod's memory limit usage in Megabytes                      |
| ISTIODACCEPTANCECPU              | Istiod's cpu limit usage in Milicores                         |
| ISTIOINGRESSPROXIESACCEPTANCEMEM | Ingress Gateways's memory limit usage in Megabytes            |
| ISTIOINGRESSPROXIESACCEPTANCECPU | Ingress Gateways's cpu limit usage in Milicores               |
| ISTIOEGRESSPROXIESACCEPTANCEMEM  | Egress Gateways's memory limit usage in Megabytes             |
| ISTIOEGRESSPROXIESACCEPTANCECPU  | Egress Gateways's cpu limit usage in Milicores                |
| ISTIOPROXIESACCEPTANCEMEM        | Istio proxies's memory limit usage in Megabytes               |
| ISTIOPROXIESDACCEPTANCECPU       | Istio proxies's cpu limit usage in Milicores                  |
| TESTCPAPPS                       | Create a set of applications to prepare the Control Plane     |
| TESTDPAPPS                       | Create a set of applications to prepare the Data Plane        |
| TESTDPAPPSFILL                   | Create a set of applications until the cluster is 100% filled |
| TESTVUS                          | K6: Number of virtual users                                   |
| TESTDURATION                     | K6: Test duration                                             |
| TESTWITHSTAGES                   | K6: Test execution with stages                                |
| REQAVG95PACCEPTANCETIME          | K6 P95 Analysis                                               |
| TRAFFICLOADAPP                   | App used in the test                                          |
| TRAFFICLOADPROTOCOL              | Protocol used in the test                                     |
| SCRIPTFILE                       | K6: Custom script used in the test                            |
| USERWORKLOADMONITORING           | User-workload monitoring disabled/enabled                     |
| USERWORKLOADNAMESPACE            | User-workload monitoring's namespace                          |
| USERWORKLOADMONITORINGSECRET     | Thanos's secret                                               |

## Author

Asier Cidon @RedHat

