import http from "k6/http";
import { Rate } from "k6/metrics";
import { check, sleep, group } from "k6";

export let errorRate = new Rate("errors");

export const options = {
    discardResponseBodies: true,
    stages: [
        { duration: '10s', target: 100 }, // simulate ramp-up of traffic from 1 to 100 users over 5 minutes.
        { duration: '10s', target: 500 }, // ramp-up to 500 users for 10 minutes
        { duration: '10s', target: 10 }, // ramp-down to 0 users
    ],
    noConnectionReuse: true,
    tags: {
      benchmark_runtime: new Date().toISOString(),
      job_name: 'http-bookinfo-stages',
    },
    insecureSkipTLSVerify: true,
};

export default function() {
    group("http-bookinfo", () => {
    const myArray = __ENV.TEST_URL.split(" ");
    myArray.forEach(element => {
        let res = http.get(`${element}`, {
        tags: { url: `${element}` },
        });
        let success = check(res, {
        'HTTP_SUCCESS': (r) => r.status === 200,
        }, {
        url: res.request.url,
        });
        errorRate.add(!success);
    });
    sleep(0.1)
    });
}