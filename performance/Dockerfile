FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
WORKDIR /bin

RUN microdnf install --nodocs tar gcc gzip git bind-utils sudo \
    && curl -Lo ./oc.tar.gz https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable/openshift-client-linux.tar.gz \
    && tar -xf oc.tar.gz \
    && rm -f oc.tar.gz \
    && curl -Lo ./golang.tar.gz https://go.dev/dl/go1.16.15.linux-amd64.tar.gz \
    && tar -xf golang.tar.gz -C / \
    && rm -f golang.tar.gz \
    && curl -Lo ./k6.tar.gz https://github.com/grafana/k6/releases/download/v0.39.0/k6-v0.39.0-linux-amd64.tar.gz \
    && tar -xf k6.tar.gz \
    && cp ./k6-v0.39.0-linux-amd64/k6 . \
    && rm -rf k6-v0.39.0-linux-amd64 \
    && rm -f k6.tar.gz \
    && microdnf update \
    && microdnf clean all

ENV GOROOT=/go
ENV PATH=$GOROOT/bin:$PATH
ENV OCP_TOKEN=${OCP_TOKEN}
ENV OCP_CRED_USR=XXX
ENV OCP_CRED_PSW=XXX
ENV OCP_API_URL=XXX
ENV MESHNAMESPACE=istio-system
ENV SMCPNAME=basic
ENV APPNSPREFIX=app-perf-test
ENV GODEBUG="x509ignoreCN=0"
ENV NSCOUNTBUNDLE="5"
ENV NSACCEPTANCETIME=3
ENV XDSPUSHACCEPTANCETIME=1
ENV ISTIODACCEPTANCEMEM=1024
ENV ISTIODACCEPTANCECPU=1000
ENV ISTIOINGRESSPROXIESACCEPTANCEMEM=1025
ENV ISTIOINGRESSPROXIESACCEPTANCECPU=1001
ENV ISTIOEGRESSPROXIESACCEPTANCEMEM=1026
ENV ISTIOEGRESSPROXIESACCEPTANCECPU=1002
ENV ISTIOPROXIESACCEPTANCEMEM=1027
ENV ISTIOPROXIESDACCEPTANCECPU=1003
ENV TESTCPAPPS=2
ENV TESTDPAPPS=5
ENV TESTVUS=2
ENV TESTDURATION=60
ENV REQAVG95PACCEPTANCETIME=500

COPY . /opt/maistra-test-tool
WORKDIR /opt/maistra-test-tool/performance/tests

ENTRYPOINT ["/bin/bash", "-c", "../../scripts/pipeline/run_all_tests.sh"]
