# ladder-app

Tiny Go HTTP service used by the CI/CD exercise ladder. `GET /`, `GET /healthz`, `GET /version`.
`LADDER_FAIL_RATE=0.5` makes half of the `/` requests fail (for canary drills).
