# AILock Governance Policy Configuration

This file operationalizes the mathematically invariant AILock policy. It must be auditable and cryptographically hashed at startup by detenforce_financial_proxy.go.

## Governance Policy Table

| ComplianceID | TargetTCOMetric | InvariantPaths | MaxRequestsPerSecond | JWKS_Endpoint |
|--------------|-----------------|----------------|----------------------|---------------|
| AIL-001 | <$0.01/request | `/api/v1/validate`, `/api/v1/execute` | 1000 | `https://auth.ailock.internal/.well-known/jwks.json` |

### Field Descriptions

**ComplianceID**: Unique identifier for this governance policy version. Used for audit trail and policy versioning.

**TargetTCOMetric**: Maximum cost per request threshold. This metric ensures economic efficiency and prevents cost overruns in the AILock execution model.

**InvariantPaths**: API endpoints that must maintain deterministic behavior and are subject to strict policy enforcement. These paths are protected by the proxy and require cryptographic verification.

**MaxRequestsPerSecond**: Rate limit for API requests to ensure system stability and prevent abuse. This limit is enforced at the proxy level.

**JWKS_Endpoint**: JSON Web Key Set endpoint for cryptographic signature verification. All requests to InvariantPaths must be signed with keys from this endpoint.
