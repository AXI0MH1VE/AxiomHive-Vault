# ARTIFACTS/API_CONTRACT.md

## Asymmetric Computational Exhaustion Engine - Internal API Contract

This contract defines the strict, auditable interface for interacting with the core security engine. All inputs and outputs are JSON-serialized for verifiability.

### 1. `POST /api/v1/enumerate/vulnerabilities` (Phase 1: Enumeration)

* **Description:** Initiates the **Asymmetric Fuzzing & Symbolic Execution** process on a specified target, subject to the 11D contextual filter.
* **Request Body (JSON):**
    ```json
    {
      "target_module": "string",  // E.g., "DNS_Resolver_v2"
      "resource_factor": "float", // Override for asymmetry_factor (optional)
      "11D_context": {
        "Temporal": "string",     // E.g., "Time-lock applied: 2025-10-20"
        "Spatial": "string",      // E.g., "Geo-fenced to DC-A"
        "Contextual": "string"    // E.g., "High-risk asset"
      }
    }
    ```
* **Response Body (JSON):**
    ```json
    {
      "status": "success/failure",
      "vulnerabilities": [
        {
          "vulnerability_id": "string",
          "severity": "CRITICAL/HIGH/MEDIUM/LOW",
          "11D_vector": "Temporal/Spatial/Contextual",
          "code_path_hash": "string" // Verifiable path hash
        }
      ]
    }
    ```
---
*(Additional API endpoints for `/api/v1/remediate`, `/api/v1/monitor`, and `/api/v1/lockdown` follow a similar, strict 11D-context-gated contract.)*
