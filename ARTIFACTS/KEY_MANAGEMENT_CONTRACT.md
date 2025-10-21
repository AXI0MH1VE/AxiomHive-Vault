# ARTIFACTS/KEY_MANAGEMENT_CONTRACT.md

## 1. Key Authority and Source of Trust

* **Authority:** The ACEP Secret Management Plane (SMP) is the sole authority for generating and distributing secrets related to the ACEP core.
* **Root of Trust (RoT):** The primary encryption key for the SMP must be sealed and protected by a dedicated **Hardware Security Module (HSM)**, adhering to the security principle outlined in **Phase 2.5**.

## 2. Deterministic Key Generation and Rotation

The randomness and lifecycle of all keys must be auditable and predictable, yet cryptographically strong.

| Key Type | Purpose | 11D Dimension | Rotation Policy | Enforcement Mechanism |
| :--- | :--- | :--- | :--- | :--- |
| **mTLS Certificates** | Service-to-Service Trust | Spatial | **Maximum 24 hours**. Zero-downtime rotation. | Automated by Service Mesh control plane, anchored by SMP. |
| **Data Encryption Keys (DEK)** | Data-at-Rest Encryption | Spatial | **Maximum 30 days**. Must comply with data residency (2.4). | Key Vault (e.g., HashiCorp Vault) integration, KMS. |
| **TOCTTOU Timer Seed** | Time-of-Use Validation | Temporal | **Maximum 60 seconds**. | Derived from verified TPM/Secure Enclave hardware clock. |

**Temporal Integrity Mandate:** Key rotation mechanisms must be mathematically verified via **Formal Verification (1.2)** to prove no window exists where an expired key remains active, or where a compromised key remains unrevoked for longer than the defined rotation period.

## 3. Principle of Least Privilege Enforcement

* **Scope:** Keys must be bound to a single service identity and possess the absolute minimum required permissions (e.g., a "read-only" service cannot possess a key with "write" permission).
* **Distribution:** Secrets must be injected into containers **at runtime** (e.g., via Kubernetes Secrets, environment variables via a secure proxy), never baked into the immutable container image. Once injected, the secret must be wiped from memory immediately upon container shutdown.
* **Auditing:** All key access, generation, rotation, and revocation actions must be logged immutably (4.4) and fed directly into the **SIEM (3.1)** for real-time **11D Contextual** analysis (e.g., key access pattern anomaly flagging).
