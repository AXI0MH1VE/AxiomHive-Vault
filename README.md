<!-- Project Banner / Logo (Black/White/Red theme) -->
<p align="center">
  <!-- If you add an image later, replace the placeholder below -->
  <img src="docs/assets/banner.png" alt="AILock – Secure Auth for Modern Systems" width="960" onerror="this.outerHTML='[ Project Banner Placeholder: place image at docs/assets/banner.png ]'" />
</p>

<h1 align="center">AILock 🔐</h1>
<p align="center"><b>Zero‑trust authentication and authorization for high‑scale APIs, microservices, and platforms.</b></p>

<p align="center">
  <a href="https://github.com/AXI0MH1VE/AILock/actions"><img src="https://img.shields.io/github/actions/workflow/status/AXI0MH1VE/AILock/ci.yml?logo=github&label=CI&color=black"></a>
  <a href="https://goreportcard.com/report/github.com/AXI0MH1VE/AILock"><img src="https://goreportcard.com/badge/github.com/AXI0MH1VE/AILock"></a>
  <a href="https://img.shields.io/codecov/c/github/AXI0MH1VE/AILock"><img src="https://img.shields.io/codecov/c/github/AXI0MH1VE/AILock?label=coverage&color=red"></a>
  <a href="https://github.com/AXI0MH1VE/AILock/security"><img src="https://img.shields.io/badge/security-verified-black"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-red"></a>
  <a href="https://github.com/AXI0MH1VE/AILock/blob/main/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/AXI0MH1VE/AILock?label=Go&color=black"></a>
</p>


Why AILock Wins
- Blazing‑fast Go core: low‑latency, high‑throughput auth and policy checks.
- Defense‑in‑depth: app‑layer DetEnforce proxy, rate‑limit, anomaly and DDOS hardening.
- Enterprise‑ready: RBAC, audit trails, revocation, rotation, short‑lived tokens.
- Drop‑in integration: clean REST, middleware hooks, and simple config.
- Cloud‑native: Docker, Compose, CI, tests, and production build profiles.

Sponsor • Contact • CTA
- Sponsor: https://buymeacoffee.com/AXI0MH1VE • PayPal: https://www.paypal.com/paypalme/yourname
- Commercial support, pilots, or security reviews: open an issue or email listed in CLAIM.md
- Star and watch this repo to follow releases.

Ownership & Licensing
- Ownership: see CLAIM.md for declaration, timestamping, and chain‑of‑custody notes.
- License: MIT (see LICENSE). Use commercially with attribution per license.

At‑a‑Glance Features
| Category | AILock | Notes |
|---|---|---|
| Auth | OAuth2, JWT, sessions | Short‑lived access + refresh rotation |
| Authorization | RBAC + permissions | Role and action scoped guards |
| Tokens | Refresh/revoke/rotate | Compromise blast radius minimization |
| Audit | Structured logs | Compliance‑ready trails |
| Proxy | DetEnforce | Request filtering, rate‑limit, WAF‑like rules |
| Scale | Go performance | Low CPU and memory overhead |
| Ops | Docker/CI | Coverage + security posture badges |

Quick Start
- Prereqs: Go 1.19+, Git, DB (Postgres/MySQL/SQLite), optional Redis
- Clone: git clone https://github.com/AXI0MH1VE/AILock && cd AILock
- deps: go mod download
- config: cp config.example.yaml config.yaml and edit
- migrate: go run cmd/migrate/main.go up
- run: go run cmd/server/main.go -> http://localhost:8080

Core Concepts
- Authentication: OAuth2/JWT with strong crypto, secret rotation, and expiry discipline.
- Authorization: flexible RBAC and permission checks via middleware.
- DetEnforce Proxy: application‑layer shield with filters, rate limiting, and logging.

API Highlights
- POST /api/v1/auth/register – create account
- POST /api/v1/auth/login – obtain access/refresh tokens
- POST /api/v1/auth/refresh – rotate access token
- GET  /api/v1/user/profile – protected example
- See README sections below for full tables.

Competitive Positioning
- Beats DIY and ad‑hoc auth by standardizing policies and reducing security debt.
- Lighter, faster footprint vs. heavyweight identity suites for embedded use.
- Designed for Go microservices and API platforms needing sub‑ms checks.

Troubleshooting
- Database connection failed
  - Check host/port/creds in config.yaml; ensure DB up and reachable; run migrations.
- Invalid JWT token
  - Ensure signing secret matches; check clock skew; verify token type and audience.
- Permission denied
  - Confirm role/permission mapping; verify middleware order and route guards.
- Port already in use
  - Change server.port in config; or free the port (lsof/kill) then retry.
- Proxy won’t start
  - Validate proxy-config.yaml path; ensure required ports available; check logs for rule parse errors.
- CORS/auth header missing
  - Add Authorization Bearer header; configure CORS allowlist in config.

Project Structure
- Core Auth Module: user auth + token lifecycle
- Authorization Layer: RBAC and permission management
- API Gateway: REST endpoints
- Configuration: environment and YAML mapping
- DetEnforce Proxy: app‑layer security enforcement

Usage Snippets
- RBAC example and middleware usage available in existing README code blocks and pkg docs.

Security Best Practices
- Always use HTTPS; terminate TLS properly.
- Rotate JWT secrets; prefer short‑lived access tokens.
- Enable audit logging and rate limiting; monitor anomalies.
- Never commit secrets; use env/secret managers.

Development
- Tests: go test ./... | coverage: go test -cover ./... | verbose: go test -v ./...
- Build: go build -o bin/ailock cmd/server/main.go
- Proxy: go build -o bin/ailock-proxy cmd/proxy/main.go
- Optimized: CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -o bin/ailock cmd/server/main.go
- Docker: docker build -t ailock:latest . && docker-compose up -d

Roadmap
- [ ] MFA (TOTP/Authenticator)
- [ ] OAuth provider logins (Google, GitHub, etc.)
- [ ] WebAuthn/FIDO2
- [ ] GraphQL endpoints
- [ ] Metrics/observability
- [ ] Kubernetes manifests
- [ ] Admin CLI

Contributors & Credit
- Core: @AXI0MH1VE (Alexis M. Adams), @ericadamsai, and community contributions.
- Thanks: dependabot, upstream libraries, and Go security ecosystem.

Support & CTA
- Issues, discussions, and security reports via GitHub Security tab and Issues.
- Sponsors welcome; integration pilots available—see CLAIM.md for contact.

Legal & IP
- License: MIT (LICENSE)
- Ownership: CLAIM.md
- Code of Conduct and Contributing guidelines included in repository.

Note: Security‑critical software—review configuration before production deploys. Ensure secrets, TLS, and network policies are correctly set in your environment.
