# AILock Internal API Specification

This contract defines the deterministic API surface for AILock, enforced by `detenforce_financial_proxy.go` and governed by `CONFIG.md`.

## Overview

The AILock API provides mathematically invariant endpoints for policy validation and execution. All requests must be cryptographically signed and conform to the governance constraints defined in CONFIG.md.

## Security Schemes

**Bearer Authentication (JWT)**
- Type: HTTP Bearer
- Format: JWT (JSON Web Token)
- Verification: All tokens must be validated against the JWKS endpoint specified in CONFIG.md
- Required Claims: `sub` (subject), `iat` (issued at), `exp` (expiration), `scope` (permissions)

## Paths Allowlist

### POST /api/v1/validate

**Description**: Validates a policy request against AILock governance rules without executing.

**Security**: Requires Bearer JWT authentication

**Request Body** (application/json):
```json
{
  "policy_id": "string (required)",
  "parameters": "object (required)",
  "dry_run": "boolean (optional, default: true)"
}
```

**Expected Responses**:
- `200 OK`: Policy validation successful
  ```json
  {
    "valid": true,
    "compliance_id": "AIL-001",
    "estimated_cost": 0.008,
    "deterministic_hash": "sha256:..."
  }
  ```
- `400 Bad Request`: Invalid policy or parameters
- `401 Unauthorized`: Missing or invalid JWT
- `429 Too Many Requests`: Rate limit exceeded (see CONFIG.md MaxRequestsPerSecond)

### POST /api/v1/execute

**Description**: Executes a validated policy with full audit trail and proof of execution.

**Security**: Requires Bearer JWT authentication

**Request Body** (application/json):
```json
{
  "policy_id": "string (required)",
  "parameters": "object (required)",
  "validation_hash": "string (required, from /validate response)",
  "idempotency_key": "string (required)"
}
```

**Expected Responses**:
- `200 OK`: Policy executed successfully
  ```json
  {
    "execution_id": "uuid",
    "result": "object",
    "proof_of_execution": "string (cryptographic proof)",
    "cost": 0.009,
    "timestamp": "ISO8601"
  }
  ```
- `400 Bad Request`: Invalid execution request
- `401 Unauthorized`: Missing or invalid JWT
- `409 Conflict`: Idempotency key already used
- `422 Unprocessable Entity`: Policy validation failed
- `429 Too Many Requests`: Rate limit exceeded
- `503 Service Unavailable`: Target system unavailable

## Proof of Execution and Audit Trail

Every execution generates a cryptographic proof that includes:
- Request hash (SHA-256)
- Execution timestamp
- Result hash
- JWT signature
- Compliance policy version (ComplianceID)

All proofs are logged to an immutable audit trail for compliance verification.

## Enforcement

This API contract is enforced by:
- **Proxy**: `detenforce_financial_proxy.go` (deterministic request routing and policy enforcement)
- **Governance**: `CONFIG.md` (policy parameters and constraints)
- **Hash Verification**: All configuration changes require cryptographic hash validation at startup

## References

- Governance Policy: See `CONFIG.md`
- Proxy Implementation: See `detenforce_financial_proxy.go`
- Strategic Amplification Engine (SAE): AILock's core deterministic execution framework
