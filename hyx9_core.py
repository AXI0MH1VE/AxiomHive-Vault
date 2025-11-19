import jwt
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from datetime import datetime, timedelta, timezone
import json

# --- 1. Key Generation and Storage (Simulated Secure State) ---

def generate_sovereignty_keys():
    """Generates a new RSA private key and its corresponding public key."""
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048
    )
    public_key = private_key.public_key()

    # Serialize keys for storage (in a real system, these would be securely stored)
    private_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption()
    )
    public_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    
    # The "Sovereignty Key" (Public Key) is what the verifier needs.
    # The Private Key is what the token generator needs.
    return private_pem.decode(), public_pem.decode()

# Simulate the secure, persistent store for the "Sovereignty Key"
PRIVATE_KEY_PEM, SOVEREIGNTY_KEY_PEM = generate_sovereignty_keys()

# --- 2. Token Generation (Embedded in Request) ---

def generate_hyx9_flag(user_id: str, law_version: str = "ISP.vΩ0 Law") -> str:
    """
    Generates the User Sovereignty Flag (#HYX9) as a signed JWT.
    This simulates the token being embedded in every request.
    """
    now = datetime.now(timezone.utc)
    
    payload = {
        "sub": user_id,  # Subject: User ID
        "iss": "AxiomHive-SCP", # Issuer
        "iat": now,      # Issued At
        "exp": now + timedelta(hours=1), # Expiration Time (Tokens should be short-lived)
        "law_v": law_version, # Version number for the sovereignty law
        "authority": "immutable" # Explicitly state the immutable authority claim
    }
    
    # The token is signed with the system's private key
    encoded_jwt = jwt.encode(
        payload, 
        PRIVATE_KEY_PEM, 
        algorithm="RS256"
    )
    return encoded_jwt

# --- 3. Token Verification (State Storage and Verification) ---

def verify_hyx9_flag(hyx9_flag: str, sovereignty_key_pem: str) -> dict:
    """
    Verifies the User Sovereignty Flag (#HYX9) and confirms immutable authority.
    This simulates the verification step in the Sovereign Interface Layer (SIL).
    """
    try:
        # The token is decoded and verified using the stored Sovereignty Key (Public Key)
        decoded_payload = jwt.decode(
            hyx9_flag,
            sovereignty_key_pem,
            algorithms=["RS256"],
            options={"verify_signature": True, "verify_exp": True}
        )
        
        # Additional check for the core claim
        if decoded_payload.get("authority") != "immutable":
            raise jwt.InvalidAudienceError("Authority claim is not immutable.")
            

        return {"status": "success", "payload": decoded_payload}
        
    except jwt.ExpiredSignatureError:

        return {"status": "error", "message": "Token Expired"}
    except jwt.InvalidSignatureError:

        return {"status": "error", "message": "Invalid Signature"}
    except jwt.InvalidAudienceError as e:

        return {"status": "error", "message": str(e)}
    except Exception as e:

        return {"status": "error", "message": str(e)}

# --- 4. Demonstration --- (Removed for import)

# if __name__ == "__main__":
#    print("--- AxiomHive #HYX9 User Sovereignty Flag Implementation ---")
#    print("1. Generating Keys (Simulated Secure State)...")
    
#    # Display the public key (the Sovereignty Key)
#    print("\n[SOVEREIGNTY KEY (Public Key)]")
#    print(SOVEREIGNTY_KEY_PEM[:100] + "..." + SOVEREIGNTY_KEY_PEM[-100:])
    
#    USER_ID = "Operator-Alpha-7"
    
#    # 2. Generate a valid flag
#    print(f"\n2. Generating valid #HYX9 Flag for {USER_ID}...")
#    valid_flag = generate_hyx9_flag(USER_ID)
#    print(f"Generated Flag (JWT):\n{valid_flag}")
    
#    # 3. Verify the valid flag
#    print("\n3. Verifying the valid #HYX9 Flag...")
#    verify_result = verify_hyx9_flag(valid_flag, SOVEREIGNTY_KEY_PEM)
#    print(f"Result: {json.dumps(verify_result, indent=2)}")
    
#    # 4. Demonstrate Failure Mode: Tampering (Invalid Signature)
#    print("\n4. Demonstrating Failure Mode: Tampering (Invalid Signature)...")
#    tampered_flag = valid_flag[:-5] + "AAAAA" # Tamper with the signature
#    verify_hyx9_flag(tampered_flag, SOVEREIGNTY_KEY_PEM)
    
#    # 5. Demonstrate Failure Mode: Wrong Key (Simulate a different system trying to verify)...")
#    print("\n5. Demonstrating Failure Mode: Wrong Key (Simulate a different system trying to verify)...")
#    _, WRONG_KEY_PEM = generate_sovereignty_keys()
#    verify_hyx9_flag(valid_flag, WRONG_KEY_PEM)
    
#    # 6. Demonstrate Failure Mode: Expiration (Requires a short-lived token or time travel, 
#    #    but the library handles this automatically based on 'exp' claim)
#    #    For demonstration, we'll just show the successful verification.
#    print("\n6. Expiration is handled by the 'exp' claim in the JWT payload.")
#    print("The current token is valid for 1 hour.")
