# CONTRACT.md: AILock DetEnforce Proxy - Internal API Contract (v1.0)

## Overview

This contract defines the strict, deterministic interfaces for the DetEnforce Proxy, ensuring Absolute Operational Integrity (AOI) by eliminating probabilistic ambiguity from the service layer.

## Security Schemes

### AuthContract (Mandatory)

- **Type**: HTTP Bearer
- **Description**: Requires a cryptographically verified JSON Web Token (JWT) signed by the Sovereign AI Kernel. AuthN failure triggers deterministic 401/403 Denial.

## Paths (Invariant Allowlist)

### 1. /api/v1/invariant/status

**GET**: Check system health and policy compliance.

- **Security**: AuthContract
- **Response 200 (AOI Confirmed)**: Returns the current ComplianceID and uptime.

### 2. /api/v1/auth/execute

**POST**: Execute a validated, non-stochastic command (e.g., policy update, bridge deployment).

- **Security**: AuthContract
- **Request Body**: Requires `command_invariant` parameter to pass CheckOmegaGovernance.
- **Response 200 (Execution Success)**: Confirms command execution and Proof of Execution (POE) logging.

### 3. /api/v1/financial/ledger

**GET**: Access the verifiable ledger of TCO capture events.

- **Security**: AuthContract
- **Response 200 (Revenue Confirmed)**: Returns the current TargetTCOMetric and latest POE logs.

### 4. /api/v1/strategic/wealth (PROPRIETARY - IWK ACCESS)

**POST**: Execute the Autonomous Strategic Tactic (AST) for wealth generation.

- **Security**: AuthContract + `IWK_LICENSE_ACTIVE: true`
- **Request Body**: Requires `strategic_invariant_command` for execution by the Invariant Wealth Kernel (IWK).
- **Response 200 (IWK SUCCESS)**: Confirms Autonomous Wealth Generation initiated.
  ```
  200 IWK SUCCESS: Strategic Tactic Deployed. Commencing Autonomous Wealth Generation. 
  Determinism is Revenue. Market Capture: $1460000000000.00
  ```
- **Response 403 (License Deny)**: IWK License is inactive; access to strategic tactic denied.
  ```
  403 ACCESS DENIED: IWK License Inactive. Activate Invariant Wealth Kernel 
  for Strategic Tactic Access.
  ```

## Verification Standard: Proof of Execution (POE)

All requests—including successful execution and denied attempts—must generate an Immutable Audit Trail (I.A.T.) entry via the `LogProofOfExecution` function, referencing the loaded ComplianceID.

### POE Log Format

```
[RFC3339_TIMESTAMP] POE: [EVENT] | ID: [COMPLIANCE_ID] | Path: [REQUEST_PATH] | Outcome: [STATUS]
```

**Examples**:

```
[2025-10-31T17:35:00Z] POE: IWK EXECUTION | ID: OMEGA-7N-RCSM-001 | Path: /api/v1/strategic/wealth | Outcome: Autonomous Wealth Generation Initiated (ALLOW)

[2025-10-31T17:36:00Z] POE: IWK FAILURE | ID: OMEGA-7N-RCSM-001 | Path: /api/v1/strategic/wealth | Outcome: License Inactive (DENY BILLIONS ACCESS)

[2025-10-31T17:37:00Z] POE: EXECUTION SUCCESS | ID: OMEGA-7N-RCSM-001 | Path: /api/v1/financial/ledger | Outcome: Policy Compliant (ALLOW)
```
