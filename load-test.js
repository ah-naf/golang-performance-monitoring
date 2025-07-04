import http from "k6/http";
import { check, sleep } from "k6";

// Base URL of your Go API
const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

export let options = {
  stages: [
    { duration: "1m", target: 100 }, // ramp-up: 0 → 100 users over 1m
    { duration: "2m", target: 300 }, // ramp-up: 100 → 300 users over 2m
    { duration: "5m", target: 500 }, // ramp-up: 300 → 1000 users over 5m
    { duration: "3m", target: 0 }, // ramp-down: 1000 → 0 users over 3m
  ],

  // (optional) some thresholds so the test fails if error rate too high
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests under 500ms
    http_req_failed: ["rate<0.01"], // error rate < 1%
  },
};

export default function () {
  // give each VU+iter a unique integer ID
  const taskId = __VU * 100000 + __ITER;

  const jsonHeaders = { headers: { "Content-Type": "application/json" } };

  // 1) Create
  let payload = JSON.stringify({
    id: taskId,
    title: `Task ${taskId}`,
    description: "Load-test task",
    status: "pending",
  });
  let res = http.post(`${BASE_URL}/task`, payload, jsonHeaders);
  check(res, { "create → 201": (r) => r.status === 201 });

  // 2) Read all
  res = http.get(`${BASE_URL}/task`);
  check(res, { "get all → 200": (r) => r.status === 200 });

  // 3) Read one
  res = http.get(`${BASE_URL}/task/${taskId}`);
  check(res, { "get one → 200": (r) => r.status === 200 });

  // 4) Update
  let updated = JSON.stringify({
    title: `Task ${taskId} ✅`,
    description: "Updated by k6",
    status: "done",
  });
  res = http.put(`${BASE_URL}/task/${taskId}`, updated, jsonHeaders);
  check(res, { "update → 200": (r) => r.status === 200 });

  // 5) Delete
  res = http.del(`${BASE_URL}/task/${taskId}`);
  check(res, { "delete → 200": (r) => r.status === 200 });

  // pause between iterations
  sleep(1);
}
