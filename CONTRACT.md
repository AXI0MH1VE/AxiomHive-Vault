# CONTRACT.md: AILock DetEnforce Proxy — Internal API Contract (v1.0)

## Overview

This contract establishes the deterministic, auditable service interfaces for the AILock DetEnforce Proxy. Its architecture enforces Absolute Operational Integrity (AOI) and eliminates all probabilistic ambiguity at the API and service layers.

## Security Schemes

### AuthContract (Mandatory)

- **Type**: HTTP Bearer Authentication
- **Requirement**: All requests must provide a cryptographically verified JSON Web Token (JWT), signed by the Sovereign AI Kernel. AuthN failure triggers deterministic 401/403 Denial.

## Invariant API Paths (Allowlist)

### /api/v1/invariant/status

**Method**: GET

**Purpose**: Check system health and policy compliance.

- **Security**: AuthContract
- **Response 200 (AOI Confirmed)**: Returns the current ComplianceID and uptime.

### /api/v1/auth/execute

**Method**: POST

**Purpose**: Execute a validated, non-stochastic command (e.g., policy update, bridge deployment).

- **Security**: AuthContract
- **Body**: Requires `command_invariant` parameter (must pass CheckOmegaGovernance).
- **Response 200 (Execution Success)**: Confirms execution and Proof of Execution (POE) logging.

### /api/v1/financial/ledger

**Method**: GET

**Purpose**: Access the verifiable ledger of TCO capture events.

- **Security**: AuthContract
- **Response 200 (Revenue Confirmed)**: Returns the current TargetTCOMetric and latest POE logs.

### /api/v1/strategic/wealth (PROPRIETARY — IWK ACCESS)

**Method**: POST

**Purpose**: Execute the Autonomous Strategic Tactic (AST) for wealth generation.

- **Security**: AuthContract and `IWK_LICENSE_ACTIVE: true`
- **Body**: Requires `strategic_invariant_command` (executed by Invariant Wealth Kernel).
- **Response 200 (IWK SUCCESS)**: Autonomous Wealth Generation initiated.
  - Example:
    ```
    200 IWK SUCCESS: Strategic Tactic Deployed. Commencing Autonomous Wealth Generation. Determinism is Revenue. Market Capture: $1,460,000,000,000.00
    ```
- **Response 403 (License Deny)**: IWK License inactive; strategic tactic access denied.
  - Example:
    ```
    403 ACCESS DENIED: IWK License Inactive. Activate Invariant Wealth Kernel for Strategic Tactic Access.
    ```

## Verification Standard: Proof of Execution (POE)

All requests—including both successful executions and denied attempts—must generate an Immutable Audit Trail (I.A.T.) entry via the LogProofOfExecution function, always referencing the loaded ComplianceID.

### POE Log Format

```
[RFC3339_TIMESTAMP] POE: [EVENT] | ID: [COMPLIANCE_ID] | Path: [REQUEST_PATH] | Outcome: [STATUS]
```

### Examples

```
[2025-10-31T17:35:00Z] POE: IWK EXECUTION | ID: OMEGA-7N-RCSM-001 | Path: /api/v1/strategic/wealth | Outcome: Autonomous Wealth Generation Initiated (ALLOW)
[2025-10-31T17:36:00Z] POE: IWK FAILURE | ID: OMEGA-7N-RCSM-001 | Path: /api/v1/strategic/wealth | Outcome: License Inactive (DENY BILLIONS ACCESS)
[2025-10-31T17:37:00Z] POE: EXECUTION SUCCESS | ID: OMEGA-7N-RCSM-001 | Path: /api/v1/financial/ledger | Outcome: Policy Compliant (ALLOW)
```

---

This contract is immutable, non-negotiable, and legally enforceable as the root legal/technical record of interface and operational law for the AXIOM HIVE DetEnforce Proxy.
