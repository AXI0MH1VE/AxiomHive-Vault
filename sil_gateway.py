from fastapi import FastAPI, Header, HTTPException, Depends, status
from typing import Annotated
import jwt
import json

# Import the core logic and keys from the hyx9_core module
from hyx9_core import SOVEREIGNTY_KEY_PEM, verify_hyx9_flag, generate_hyx9_flag

# --- 1. FastAPI Application Initialization ---
app = FastAPI(
    title="SIL-Gateway: Sovereign Coherence Protocol (SCP) Request Interceptor",
    description="Enforces User Sovereignty Flag (#HYX9) verification before routing to the Deterministic Engine Core (DEC)."
)

# --- 2. Dependency Function for #HYX9 Verification ---

def get_sovereign_operator(
    hyx9_flag: Annotated[str | None, Header(alias="X-Sovereignty-Flag")] = None
):
    """
    Dependency that extracts and verifies the #HYX9 token from the request header.
    This acts as the Request Interceptor, enforcing the Forge Prime mandate.
    """
    if hyx9_flag is None:
        # Rejection: No Flag provided. This is a non-negotiable failure.
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Forge Prime Mandate Violation: X-Sovereignty-Flag header missing. Immutable authority not established."
        )

    # Attempt to verify the flag using the stored Sovereignty Key
    verification_result = verify_hyx9_flag(hyx9_flag, SOVEREIGNTY_KEY_PEM)

    if verification_result["status"] == "error":
        # Rejection: Verification failed (tampered, expired, or wrong key).
        # This is the deterministic rejection of chaos/entropy.
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail=f"Sovereign Coherence Protocol Violation: #HYX9 Verification Failed. Reason: {verification_result['message']}"
        )

    # Success: Immutable authority confirmed. Return the verified payload (Operator's identity)
    return verification_result["payload"]

# --- 3. DEC-Routing Endpoint (Simulated) ---

@app.get("/dec/process")
def route_to_dec(
    operator_payload: Annotated[dict, Depends(get_sovereign_operator)]
):
    """
    Simulated endpoint for the Deterministic Engine Core (DEC).
    This route is only reachable if the #HYX9 verification is successful.
    """
    user_id = operator_payload.get("sub", "UNKNOWN_OPERATOR")
    law_version = operator_payload.get("law_v", "UNKNOWN_LAW")
    
    # In a real system, the DEC would now begin processing, knowing the
    # Operator's immutable authority is established.
    
    # The next step (DEC-Logger) would use this payload to write the
    # Cryptographic Commitment Chain entry.
    
    return {
        "status": "DEC_ROUTE_SUCCESS",
        "message": "Request successfully routed to Deterministic Engine Core.",
        "operator_id": user_id,
        "law_enforced": law_version,
        "trace_note": "Immutable authority established via #HYX9. Proceeding to Inverted Lagrangian Solver."
    }

# --- 4. Utility Endpoint for Testing (Token Generation) ---

@app.get("/utility/generate_flag")
def generate_test_flag(user_id: str = "Operator-Alpha-7"):
    """
    Utility endpoint to generate a valid #HYX9 flag for testing the SIL-Gateway.
    """
    flag = generate_hyx9_flag(user_id)
    return {
        "status": "FLAG_GENERATED",
        "user_id": user_id,
        "X-Sovereignty-Flag": flag,
        "note": "Use this flag in the 'X-Sovereignty-Flag' header for /dec/process"
    }

# --- 5. Key Exposure Endpoint (For Audit/Debugging) ---

@app.get("/utility/sovereignty_key")
def get_sovereignty_key():
    """
    Exposes the public Sovereignty Key (Public Key) for external audit/verification.
    """
    return {
        "status": "KEY_EXPOSED",
        "key_type": "Public RSA Key (Sovereignty Key)",
        "key_pem": SOVEREIGNTY_KEY_PEM
    }
