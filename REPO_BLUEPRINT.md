# AILock Repository Blueprint

Production-grade, operator-proof AILock repo blueprint — zero-mystery, clone-and-prove. Use this as your reference for implementation and audits.

## Repo layout (top level)

```
AILock/
├─ cmd/
│ ├─ ailock-api/ # API service entrypoint (authN/authZ endpoints)
│ └─ detenforce-proxy/ # L7 proxy entrypoint (policy, rate limit, TLS)
├─ internal/ # Non-public Go packages
│ ├─ auth/ # JWT/OIDC, JWKS, claim validation, clock-skew, kid rotation
│ ├─ rbac/ # Roles, permissions, policy engine, matchers
│ ├─ proxy/ # Reverse proxy core: routing, rewrite, timeouts, retries
│ ├─ ratelimit/ # Token bucket/leaky bucket, identity- and IP-scoped
│ ├─ validation/ # Request size/header limits, content-type checks
│ ├─ tls/ # TLS config, ACME (http-01/tls-alpn-01), cert storage
│ ├─ config/ # Config loader (YAML+env), schema validation
│ ├─ logging/ # Structured logs (zap/zerolog), request-id middleware
│ ├─ tracing/ # OTEL hooks (optional), metrics via Prometheus
│ ├─ health/ # Liveness/readiness endpoints
│ ├─ upstream/ # Balancer, health checks, circuit breaker
│ ├─ audit/ # Decision logs: who/what/why/latency
│ └─ security/ # Safe defaults, header normalization, SSRF guards
├─ pkg/ # Optional public packages (if you want SDK surface)
│ └─ client/ # Minimal Go client for testing (optional)
├─ api/ # API definitions (OpenAPI/Swagger), JSON schemas
│ └─ openapi.yaml
├─ deployments/
│ ├─ docker/ # Dockerfile(s), docker-compose demo
│ ├─ k8s/ # Manifests: Deployment, Service, Ingress, ConfigMap, Secret
│ └─ helm/ # (Optional) Helm chart for production installs
├─ docs/
│ ├─ ARCHITECTURE.md # Component diagram, data flows, trust boundaries
│ ├─ QUICKSTART.md # 60s path to first success (API+Proxy)
│ ├─ GETTING_STARTED.md # Step-by-step (Docker, Compose, K8s, bare-metal)
│ ├─ CONFIGURATION.md # Full config schema with sane defaults
│ ├─ POLICY_EXAMPLES.md # RBAC/policy examples (user, service, admin)
│ ├─ THREAT_MODEL.md # STRIDE-style model, mitigations, residual risk
│ ├─ SECURITY.md # Report policy, supported versions, contact
│ ├─ PERFORMANCE.md # Bench methodology + target latencies
│ ├─ OPERATIONS.md # Logging, metrics, rotation, backups, cert mgmt
│ ├─ ROADMAP.md # Near-term features, deprecations policy
│ └─ ADRs/ # Architecture Decision Records
├─ examples/
│ ├─ configs/ # Minimal + production-like configs (HS256, RS256+JWKS)
│ ├─ jwks/ # Sample keys (dev-only), JWKS static files
│ ├─ upstreams/ # Tiny echo server for offline demos
│ └─ clients/ # curl/hey/vegeta examples, Postman collection
├─ test/
│ ├─ unit/ # Go unit tests, fuzz tests (Go fuzz)
│ ├─ integration/ # Start services, hit routes, assert decisions
│ └─ e2e/ # Compose/kind-based end-to-end (auth+proxy+upstream)
├─ tools/
│ ├─ jwt_minter.go # HS256 (dev); RS256 variant + JWKS publisher
│ ├─ jwksgen.sh # Generate RS256/ES256 keys + JWKS
│ └─ devcert.sh # mkcert/step-ca script for local TLS
├─ scripts/
│ ├─ smoke.sh # 200→429 limit, 413 body, 431 headers
│ ├─ bench.sh # 'hey' 10s p50/p95/p99 smoke
│ └─ release.sh # Tag, build, SBOM, sign, create GH Release
├─ .github/
│ └─ workflows/
│ ├─ ci.yml # build+test+lint
│ ├─ security.yml # govulncheck, staticcheck, gosec, CodeQL
│ ├─ container-scan.yml # Trivy image scan + SBOM (Syft)
│ ├─ provenance.yml # SLSA provenance/cosign attest
│ └─ release.yml # versioned releases, checksums, cosign signatures
├─ ARTIFACTS/ # Ship reproducible demo config
│ └─ DETENFORCE_PROXY_CORE.yaml
├─ Dockerfile
├─ docker-compose.yaml
├─ Makefile
├─ go.mod / go.sum
├─ LICENSE (MIT)
├─ CONTRIBUTING.md
├─ CODE_OF_CONDUCT.md
├─ GOVERNANCE.md # (optional if you want community structure)
└─ CHANGELOG.md
```

## File-by-file expectations (high-signal)

### cmd/ailock-api
**main.go**: wires config, logging, health, and `/api/auth/` handlers.
**Endpoints**: `/api/auth/login`, `/api/auth/refresh`, `/api/auth/keys` (JWKS read-only), `/api/auth/introspect` (optional).
**Hard edges**: issuer/audience enforced, clock-skew leeway, no `alg=none`, reject unknown `kid`, JWKS retrieval with caching and allowlist, no SSRF via `jku`.

### cmd/detenforce-proxy
**main.go**: config + server bootstrap, read/write timeouts, graceful shutdown, HSTS if TLS on.
**Middleware order**: req-id → logging → auth (optional per route) → rate-limit → validation → upstream.
**Decision log (JSON)**: ts, request_id, src_ip, method, path, subject, decision{allow|deny}, rule, reason, latency_ms.

### internal/auth
Verifiers: HS256 (dev only), RS256/ES256 with JWKS, rotate by `kid`.
Token parsing strict: disallow `none`, enforce `typ`, enforce `aud`/`iss`, validate `exp`/`nbf`/`iat` with skew.
Refresh logic guarded (if you implement it) and revocation cache (optional).

### internal/rbac
Role → permissions → resources model with examples.
Matchers for method/path, optionally headers/claims.
Policy structure documented in `docs/POLICY_EXAMPLES.md`.

### internal/proxy
Route matchers: path, path_prefix, method lists.
Rewrites, upstream selection, retries with backoff, header normalization, hop-by-hop strip.
Circuit breaker + upstream health checks.

### internal/ratelimit
Token bucket per IP and subject (JWT), sliding window metrics.
Configurable bursts; accurate 429s with Retry-After.

### internal/validation
`max_body_size`, `max_header_size`, allowed content types.
413 / 431 / 400 mapped cleanly with machine-readable body.

### internal/tls
TLS ≥1.2; sane ciphers; optional ACME (staging vs prod); cert storage path; renewal window.
CLI flags for ACME email, domains, challenge type (http-01/tls-alpn-01).

### internal/audit
Single function to emit decision logs; zero secrets in logs.
Optionally write to file/stdout/OTEL exporter.

## examples/configs/ you must ship

### detenforce-minimal-hs256.yaml
Matches the evaluation pack. Port 8080, RPM=5, httpbin upstream, `/health` open, `/api/secure/*` requires JWT.

### detenforce-rs256-jwks.yaml
Uses RS256 with local JWKS endpoint (your ailock-api on `/api/auth/keys`).

### detenforce-prod-template.yaml
TLS on, ACME enabled, rate limit identities by sub, upstream pool with health checks, unique request IDs forwarded.
Include comments and exact curl examples in each YAML.

## tools/ you must ship
`jwt_minter.go` (HS256) and `jwt_minter_rs256.go` (loads a private key; prints token + example curl).
`jwksgen.sh` (or Go tool): generates RS256/ES256 keys, jwks.json. Tags keys with kid.
`devcert.sh`: one-liner local cert via mkcert (or step-ca script).

## deployments/
**Dockerfile**: static Go build, nonroot user, minimal base (distroless/alpine), healthcheck.
**docker-compose.yaml**: API + Proxy + Echo upstream. Wires ports, volumes for configs and keys.
**k8s/**: Deployment (readiness/liveness), Service, ConfigMap (config YAML), Secret (keys), Ingress, PodSecurityContext, Resources/limits, Probes tuned.
**helm/** (optional but helpful): values for auth, ratelimit, upstreams, TLS.

## docs/ (operator trust)
**ARCHITECTURE.md**: diagram of API, Proxy, JWKS flow, decision path; trust boundaries; failure modes.
**THREAT_MODEL.md**: STRIDE table, specific mitigations (JWT confusion, SSRF, header smuggling, slowloris, ACME pitfalls).
**SECURITY.md**: supported versions, report email/PGP, response SLAs.
**CONFIGURATION.md**: every key with defaults, examples; precedence (flags > env > file).
**POLICY_EXAMPLES.md**: copy-paste policies (IP only; JWT subject role; service-to-service).
**PERFORMANCE.md**: how to run hey/vegeta, expected p50/p95/p99 under sample hardware.
**OPERATIONS.md**: log rotation, metrics endpoints, cert renewals, backup/restore, upgrading safely.
**ROADMAP.md + ADRs/**: where it's going, and why past decisions were made.

## Testing & quality gates
**Unit tests**: auth parsing, JWKS caching/rotation, rate-limit math, header normalization.
**Fuzz tests**: token parser, header parser, path matcher.
**Integration**: bring up API+Proxy+Upstream (compose or go test with testcontainers-go); assert 401→200 with valid JWT; 429 under RPM; 413/431 boundaries.
**E2E**: kind-based workflow kicks on PR; uploads junit + coverage.
**Linters**: golangci-lint or staticcheck; formatting in CI.
**Security**: govulncheck, gosec, CodeQL.
**Container**: build image, generate SBOM (Syft), scan (Trivy), sign (cosign).
**Provenance**: SLSA provenance attest on release.
**Coverage target**: ≥80% for core packages (auth, proxy, ratelimit).

## CI/CD (.github/workflows/)
**ci.yml**: matrix (Go 1.21–1.22), build, test, race, cover, lint.
**security.yml**: govulncheck, staticcheck, gosec, CodeQL.
**container-scan.yml**: build, SBOM, Trivy; upload SARIF.
**provenance.yml**: cosign attest with GitHub OIDC.
**release.yml**: version tag → build multi-arch images, upload binaries + checksums, cosign sign, attach SBOM, generate release notes from conventional commits.

## Developer ergonomics
**Makefile**: `make build`, `test`, `lint`, `run-proxy`, `run-api`, `smoke`, `bench`.
Dev containers or `.tool-versions` (asdf) for consistent toolchains.
Pre-commit with `go fmt`/`vet`/`lint`, YAML schema check.
Postman collection in `examples/clients/`.

## OSS hygiene
LICENSE (MIT), CONTRIBUTING.md, CODE_OF_CONDUCT.md.
CHANGELOG.md (Conventional Commits + semantic versioning).
GOVERNANCE.md (optional) if you invite maintainers.

## Minimal "first-run" you must prove in README (front-and-center)

```bash
# 1) Run API (provides /api/auth/keys for JWKS) and Proxy via docker-compose
docker compose up -d

# 2) Health and rate-limit (expect 200×5 then 429)
bash scripts/smoke.sh

# 3) AuthN: mint token and call protected route
(cd tools && go run jwt_minter.go) # prints token
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/secure/test # expect 200 echo

# 4) Bench (optional)
export JWT_TOKEN=<TOKEN>
bash scripts/bench.sh
```

**Promise**: anyone following those four blocks should see success without editing a single line.

## Final sanity checklist (before you announce)
- `docs/QUICKSTART.md` mirrors README steps exactly.
- `ARTIFACTS/DETENFORCE_PROXY_CORE.yaml` matches docs verbatim.
- Logs are JSON and include: request_id, decision, rule, reason, latency_ms.
- HS256 examples clearly labeled dev-only; RS256/JWKS path is the default in docs.
- CI badges green; security scans visible in PRs.
- A 10-second clip (GIF/webm) of rate limit triggering and JWT allow/deny ready for social.

If you want, I'll extend your evaluation pack with the RS256/JWKS config, a tiny offline echo upstream, and a docker-compose that boots all three—so the repo's first run is fully air-gapped and deterministic.
