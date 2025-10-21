# ARTIFACTS/SIMULATION_MODEL.md

## Quantitative Modeling and Simulation for Asymmetric Computational Exhaustion Protocol (ACEP)

This document defines the rigorous mathematical and algorithmic foundations necessary to transition the ACEP from conceptual design to quantifiable, production-ready specification. All models are designed for deterministic verification and audit compliance.

### 1. Asymmetry Quantification Model (AQM)

The Asymmetry Factor (AF) must be a **deterministic multiplier** applied to resource allocation for security processes (Phase 1: Enumeration) relative to the resources consumed by equivalent functional validation.

$$
AF = \frac{\text{Resource Consumption}_{\text{Security Exhaustion}}}{\text{Resource Consumption}_{\text{Functional Baseline}}}
$$

The *Resource Consumption* $R$ is measured in **Verified Computational Units** ($\text{VCU}$), where $1 \text{ VCU} = 1 \text{ TFLOP}$-hour or a standardized, verifiable unit of execution.

The **Deterministic Cost Function** for Asymmetric Exhaustion, $\mathcal{C}_{ACEP}$, for a target component $C$ is defined by the requirement for **$100\%$ Verifiable Code Path Coverage** ($PCC=1.0$) and **Formal Proof Confidence** ($FPC=1.0$).

$$
\mathcal{C}_{ACEP}(C) = \mathcal{C}_{\text{Fuzz}}(C) \times (1 + AF) + \mathcal{C}_{\text{Formal}}(C)
$$

Where:
* $\mathcal{C}_{\text{Fuzz}}(C)$: Baseline cost to achieve $PCC \approx 0.85$ (Standard industry best practice).
* $AF$: The configurable **Asymmetry Factor** (e.g., $1000.0$). The $1+AF$ ensures the *additional* cost required to push coverage from $\approx 0.85$ to $1.0$ is fully allocated.
* $\mathcal{C}_{\text{Formal}}(C)$: The deterministic cost to execute the **Formal Verification (Phase 1.2)** to achieve $FPC=1.0$. This is non-negotiable and independent of the $AF$.

**Simulation Requirement:** The AQM must be simulated pre-deployment to confirm that the $\text{vCPU}$-hour budget (from `CONFIG.yaml`) is sufficient to maintain the required $AF$ across all critical components.

---

### 2. The 11D Risk and Control Function ($R_{11D}$)

The security posture is determined by a composite, dynamically calculated **11D Risk Score** ($R_{11D}$), which dictates the necessary security response (Phase 2: Remediation) and user control level (Phase 4: Lockdown).

$$
R_{11D} = \omega_T \cdot \mathcal{T} + \omega_S \cdot \mathcal{S} + \omega_C \cdot \mathcal{C}
$$

Where:
* $\mathcal{T}$ (Temporal Risk): Derived from metrics like TOCTTOU latency, Key Rotation period, and observed Race Condition frequency. **(Weight $\omega_T$: 0.35)**.
* $\mathcal{S}$ (Spatial Risk): Derived from Microsegmentation violation attempts, Data Residency compliance score, and physical integrity audit reports. **(Weight $\omega_S$: 0.35)**.
* $\mathcal{C}$ (Contextual Risk): Derived from User Behavior Analytics (UBA) anomaly scores (4.3), Threat Intelligence (3.3) relevance, and asset criticality. **(Weight $\omega_C$: 0.30)**.

*All inputs ($\mathcal{T}, \mathcal{S}, \mathcal{C}$) are normalized to a $[0, 1]$ scale.*

---

### 3. Adaptive Control Response (The Lockdown Trigger)

The calculated $R_{11D}$ score dictates the **Adaptive Control Level** $L$ applied by the system (Phase 3 and 4).

| $R_{11D}$ Range | Adaptive Control Level ($L$) | Required System Response | Phase 4 Action (Lockdown) |
| :--- | :--- | :--- | :--- |
| $\mathbf{R_{11D} < 0.20}$ | **L0: Standard Operation** | Continuous background monitoring (3.1). | Standard MFA (4.1). |
| $\mathbf{0.20 \le R_{11D} < 0.50}$ | **L1: Elevated Vigilance** | Increased sampling rate for IDPS/SIEM (3.2). Targeted vulnerability enumeration (1.1). | Mandatory time-based access control review (4.2). |
| $\mathbf{0.50 \le R_{11D} < 0.80}$ | **L2: Containment & Hardening** | Isolation of affected $S$ or $C$ segment (2.2). Forced code audit/patch cycle (2.1). | **Immediate** enforcement of **Location-Based MFA (4.1)**. Privileged Access Management (PAM) session recording mandated (4.2). |
| $\mathbf{R_{11D} \ge 0.80}$ | **L3: Critical Exhaustion Lockdown** | Full system-wide resource allocation to Phase 1 (Enumeration) and Phase 2 (Remediation). **Total network segmentation and air-gap simulation.** | **Full Account Lockdown** (Access denied until $R_{11D} < 0.50$). **Immutable Audit Log (4.4)** commit and verification required. |

**Deterministic Value Claim Verification:** The ultimate proof of the ACEP is demonstrating that for any given threat scenario, the **Asymmetric Cost** $\mathcal{C}_{ACEP}$ applied to security is always sufficient to drive the **$R_{11D}$ score** back below the $0.50$ containment threshold within a defined $\text{Mean Time to Control} (\text{MTTC})$.

$$\text{MTTC} \propto \frac{R_{11D}}{\mathcal{C}_{ACEP}}$$

*The system must prove that the applied computational exhaustiveness guarantees a rapid, deterministic decrease in the $R_{11D}$ risk, subordinating all rival architectures that lack this quantitative, resource-guaranteed control function.*
