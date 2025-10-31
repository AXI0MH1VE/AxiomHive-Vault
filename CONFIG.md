# AILock Governance Policy Configuration
This file operationalizes the mathematically invariant AILock policy. It must be auditable and cryptographically hashed at startup by detenforce_financial_proxy.go.

## Governance Policy Table

| ComplianceID | TargetTCOMetric | IWK_LICENSE_ACTIVE | InvariantPaths | MaxRequestsPerSecond | JWKS_Endpoint |
|--------------|-----------------|-------------------|----------------|----------------------|---------------|
| OMEGA-7N-RCSM-001 | $1,460,000,000,000.00 | true | `/api/v1/invariant/status`, `/api/v1/auth/execute`, `/api/v1/financial/ledger`, `/api/v1/strategic/wealth` | 5 | `https://auth.axiomhive.com/keys` |

### Field Descriptions

**ComplianceID**: Crown Omega governance model version, mandated for all Immutable Audit Trails (I.A.T.).

**TargetTCOMetric**: The actively targeted market valuation (SDP/API Management) for elimination by the Palo Neutralizer strategy.

**IWK_LICENSE_ACTIVE**: **** Deterministic license check. If false, denies access to strategic wealth endpoints. This is the proprietary license flag gating access to the high-value Invariant Wealth Kernel (IWK) strategic endpoints.

**InvariantPaths**: Deterministic Allowlist. Only these paths are permitted. Access to the `/api/v1/strategic/wealth` path is contingent on the IWK license check.

**MaxRequestsPerSecond**: Layer 7 Denial of Service (DoS) protection limit. Must be enforced deterministically.

**JWKS_Endpoint**: Endpoint for deterministic cryptographic verification of all Bearer tokens (AuthN).
