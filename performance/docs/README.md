# CUSTOM - Openshif Service Mesh Performance tool

A performance tool for running specific Openshif Service Mesh use cases on an OpenShift 4.x cluster obtaining the different components' specific metrics.

The idea behind this inititative is to generate a performance testing tool that certifies Openshift Service Mesh in different scenarios. 

## Versions

| Name      | Version  |
| --------- | -------- |
| OS        | Linux    |
| Golang    | 1.13+    |
| OpenSSl   | 1.1.1+   |
| oc client | 4.x      |
| k6        | v0.38.3+ |

## Prerequisite

* RedHat Service Mesh Operator has been installed on the OpenShift cluster
* Service Mesh Control Plane installed and ready 
* Service Mesh Member Role ready with 1 configured member at least

## Testing Prerequisite

* An `oc` client can be downloaded from [mirror openshift-v4 clients](https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/) and extract the `openshift-client-...tar.gz` file and move both `oc` and `kubectl` binaries into a local PATH directory.
* An `k6` client can be downloaded and move the binary into a local PATH directory.
* Access an OpenShift cluster from command line before running tests 

```$bash
oc login -u [user] -p [token] --server=[OCP API server]
```

## Testing Performance

* A main test is in the `performance/tests` directory. All test cases are in the `test_cases.go` and are mapped to the implementations in the `pkg` directory.
* Optionally to run all the test cases customizing the SMCP namespace and the SMCP name: A user can update the expected values in the `performance/tests/test.env`.
* To run all the test cases: `cd tests; go test -timeout 2h -v`.

NOTE: The `-timeout` flag is necessary when running all tests or several major test cases. Otherwise, a `go test` command falls into panic after 10 minutes.

* To run a single test case: e.g. `cd performance/tests; go test -run A1 -timeout 2h -v`

NOTE: Test cases shortname and mapping are in the `tests/test_cases.go` file.

## Author

Asier Cidon @RedHat

