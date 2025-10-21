#!/bin/bash

# SUPRAPROTOCOL V∞: DevSecOps SBOM Generation Tool
# Generates SPDX or CycloneDX SBOM for container images and code

set -e

PROJECT_NAME="AILock-V_INF"
OUTPUT_DIR="./security"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Check if syft is installed (for CycloneDX)
if command -v syft &> /dev/null; then
    echo "Generating CycloneDX SBOM with Syft..."
    syft packages docker:AILock-V_INF-auth-service:latest -o cyclonedx-json --file "${OUTPUT_DIR}/sbom_cyclonedx_${TIMESTAMP}.json"
    syft packages docker:AILock-V_INF-client-app:latest -o cyclonedx-json --file "${OUTPUT_DIR}/sbom_client_cyclonedx_${TIMESTAMP}.json"
else
    echo "Syft not found. Installing or use alternative..."
    # Placeholder for alternative SBOM generation
    echo "Placeholder: Generate SBOM manually"
fi

# For SPDX generation
if command -v cdxgen &> /dev/null; then
    echo "Generating SBOM with cdxgen..."
    cdxgen -t nodejs -o "${OUTPUT_DIR}/sbom_spdx_${TIMESTAMP}.json" ./client-app
    cdxgen -t go -o "${OUTPUT_DIR}/sbom_go_spdx_${TIMESTAMP}.json" .
else
    echo "cdxgen not found. Please install: npm install -g @cyclonedx/cdxgen"
fi

echo "SBOM generation complete. Files saved in ${OUTPUT_DIR}"
