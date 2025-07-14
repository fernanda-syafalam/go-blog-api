import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

export let loginTrend = new Trend('login_response_time');
export let loginSuccessRate = new Rate('login_success');

export let options = {
  stages: [
    { duration: '5s', target: 10 },     // Warm-up
    { duration: '10s', target: 50 },    // Ramp-up
    { duration: '30s', target: 50 },    // Sustained
    { duration: '10s', target: 0 },     // Ramp-down
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],   // 95% requests < 200ms
    login_success: ['rate>0.95'],       // 95% success rate
  }
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'; // Bisa di-set via env
const USERS = [
  { id: 'username1', password: 'password' },
  { id: 'username2', password: 'password' },
  { id: 'username3', password: 'password' },
];

export default function () {
  const user = USERS[Math.floor(Math.random() * USERS.length)];

  const payload = JSON.stringify({
    email: 'john.doe@example.com',
    password: 'strongpassword123',
  });

  const headers = {
    'Content-Type': 'application/json',
  };

  const res = http.post(`${BASE_URL}/api/auth/login`, payload, { headers });

  loginTrend.add(res.timings.duration);
  loginSuccessRate.add(res.status === 200);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response has token': (r) => r.json('data.token') !== undefined,
  });

  sleep(1); 
}
