import http from "k6/http";
import { check } from "k6";

export const options = {
    thresholds: {
        http_req_duration: ["p(99) < 3000"],
    },
    stages: [
        { duration: "30s", target: 10 },
        { duration: "1m", target: 10 },
        { duration: "20s", target: 0 },
    ],
};

const payload = JSON.stringify({
    acommodation_id: "019356b3-cce2-7bb6-a7b7-40d96b8dc233",
    user_id: "019356b3-e1a4-755d-b7ed-6135ac8d05d3",
    start_at: "2021-12-01T00:00:00Z",
    end_at: "2021-12-02T00:00:00Z",
});

export default function() {
    let res = http.post("http://localhost:8080/api/v1/bookings", payload);
    check(res, { "status was 202": (r) => r.status == 202 });
}
