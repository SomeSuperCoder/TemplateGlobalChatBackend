import http from "k6/http";
import { sleep, check } from "k6";

export const options = {
  discardResponseBodies: false,

  // SIMPLE, BRUTAL, AND EFFECTIVE: Constant Virtual Users
  scenarios: {
    constant_vus_siege: {
      executor: "constant-vus", // The simplest executor
      vus: 10_000, // MAX OUT YOUR LOCAL HARDWARE (2000 VUs)
      duration: "10m", // Siege for 10 minutes straight
    },
  },

  thresholds: {
    http_req_failed: ["rate<0.95"],
    http_req_duration: ["p(99)<10000"],
  },
};

const sessionToken = "jF1r5OhhNVia9ztFQtgnefak6rES1K_xor_gnFEd7no=";
const csrfToken = "BRP803aW1mY2xwyGOe3nO1-CAH0OryW77Xjz-t_VAGY=";

export default function () {
  let res = http.get("http://localhost:8090/messages/?page=1&limit=50", {
    headers: { "X-CSRF-Token": csrfToken },
    cookies: {
      session_token: sessionToken,
    },
  });
  check(res, { "status is 200": (res) => res.status === 200 });
  sleep(1);
}
