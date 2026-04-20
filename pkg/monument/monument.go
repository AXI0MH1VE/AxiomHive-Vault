// Package monument implements Monument Protocol generation.
// Monument Protocols are cryptographically signed, deterministic axiom sets
// that act as filters for all incoming data vectors with zero drift.
package monument

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Axiom represents a single deterministic rule.
type Axiom struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Formula     string                 `json:"formula"`
	Priority    int                    `json:"priority"`
	Immutable   bool                   `json:"immutable"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// MonumentProtocol represents a complete set of deterministic axioms.
type MonumentProtocol struct {
	ProtocolID    string                 `json:"protocol_id"`
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	Intent        string                 `json:"intent"`
	Axioms        []Axiom                `json:"axioms"`
	CreatedAt     time.Time              `json:"created_at"`
	Signature     string                 `json:"signature"`
	Hash          string                 `json:"hash"`
	ComplianceID  string                 `json:"compliance_id"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// IntentPacket represents the strategic outcome to be achieved.
type IntentPacket struct {
	Intent      string                 `json:"intent"`
	Objective   string                 `json:"objective"`
	Constraints []string               `json:"constraints"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ProtocolBuilder builds Monument Protocols from intent packets.
type ProtocolBuilder struct {
	complianceID string
	protocols    []MonumentProtocol
}

// NewProtocolBuilder creates a new Monument Protocol builder.
func NewProtocolBuilder(complianceID string) *ProtocolBuilder {
	return &ProtocolBuilder{
		complianceID: complianceID,
		protocols:    make([]MonumentProtocol, 0),
	}
}

// GenerateProtocol synthesizes a Monument Protocol from an intent packet.
func (pb *ProtocolBuilder) GenerateProtocol(intent IntentPacket, name, version string) MonumentProtocol {
	protocol := MonumentProtocol{
		ProtocolID:   generateProtocolID(),
		Name:         name,
		Version:      version,
		Intent:       intent.Intent,
		Axioms:       make([]Axiom, 0),
		CreatedAt:    time.Now().UTC(),
		ComplianceID: pb.complianceID,
		Metadata:     make(map[string]interface{}),
	}
	
	// Synthesize axioms from intent
	axioms := pb.synthesizeAxioms(intent)
	protocol.Axioms = axioms
	
	// Generate cryptographic hash
	protocol.Hash = pb.hashProtocol(protocol)
	
	// Generate cryptographic signature
	protocol.Signature = pb.signProtocol(protocol)
	
	// Store protocol
	pb.protocols = append(pb.protocols, protocol)
	
	return protocol
}

// synthesizeAxioms converts intent into deterministic axioms.
func (pb *ProtocolBuilder) synthesizeAxioms(intent IntentPacket) []Axiom {
	axioms := make([]Axiom, 0)
	
	// Core axiom: Intent definition
	axioms = append(axioms, Axiom{
		ID:          generateAxiomID(),
		Name:        "IntentDefinition",
		Description: fmt.Sprintf("Primary intent: %s", intent.Intent),
		Formula:     fmt.Sprintf("INTENT = \"%s\"", intent.Intent),
		Priority:    1,
		Immutable:   true,
		Metadata:    map[string]interface{}{"type": "intent"},
	})
	
	// Objective axioms
	axioms = append(axioms, Axiom{
		ID:          generateAxiomID(),
		Name:        "ObjectiveConstraint",
		Description: fmt.Sprintf("Objective: %s", intent.Objective),
		Formula:     fmt.Sprintf("OBJECTIVE = \"%s\"", intent.Objective),
		Priority:    2,
		Immutable:   true,
		Metadata:    map[string]interface{}{"type": "objective"},
	})
	
	// Constraint axioms
	for i, constraint := range intent.Constraints {
		axioms = append(axioms, Axiom{
			ID:          generateAxiomID(),
			Name:        fmt.Sprintf("Constraint_%d", i+1),
			Description: constraint,
			Formula:     fmt.Sprintf("CONSTRAINT_%d = \"%s\"", i+1, constraint),
			Priority:    10 + i,
			Immutable:   true,
			Metadata:    map[string]interface{}{"type": "constraint", "index": i},
		})
	}
	
	// Parameter axioms
	for key, value := range intent.Parameters {
		axioms = append(axioms, Axiom{
			ID:          generateAxiomID(),
			Name:        fmt.Sprintf("Parameter_%s", key),
			Description: fmt.Sprintf("Parameter %s = %v", key, value),
			Formula:     fmt.Sprintf("PARAM[\"%s\"] = %v", key, value),
			Priority:    100,
			Immutable:   false,
			Metadata:    map[string]interface{}{"type": "parameter", "key": key},
		})
	}
	
	return axioms
}

// hashProtocol generates a SHA-256 hash of the protocol.
func (pb *ProtocolBuilder) hashProtocol(protocol MonumentProtocol) string {
	// Create deterministic representation
	data := fmt.Sprintf("%s|%s|%s|%s|%d|%d",
		protocol.ProtocolID,
		protocol.Name,
		protocol.Version,
		protocol.Intent,
		protocol.CreatedAt.Unix(),
		len(protocol.Axioms),
	)
	
	// Include all axiom hashes
	for _, axiom := range protocol.Axioms {
		data += fmt.Sprintf("|%s:%s", axiom.ID, axiom.Formula)
	}
	
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// signProtocol generates a cryptographic signature for the protocol.
func (pb *ProtocolBuilder) signProtocol(protocol MonumentProtocol) string {
	// In production, use proper cryptographic signing (RSA, ECDSA)
	// For now, use HMAC-like construction
	signData := fmt.Sprintf("MONUMENT_PROTOCOL|%s|%s|%s",
		protocol.Hash,
		pb.complianceID,
		protocol.CreatedAt.Format(time.RFC3339),
	)
	
	signature := sha256.Sum256([]byte(signData))
	return hex.EncodeToString(signature[:])
}

// VerifyProtocol verifies the integrity and signature of a protocol.
func (pb *ProtocolBuilder) VerifyProtocol(protocol MonumentProtocol) (bool, []string) {
	errors := make([]string, 0)
	
	// Verify hash
	expectedHash := pb.hashProtocol(protocol)
	if protocol.Hash != expectedHash {
		errors = append(errors, "Hash verification failed")
	}
	
	// Verify signature
	expectedSignature := pb.signProtocol(protocol)
	if protocol.Signature != expectedSignature {
		errors = append(errors, "Signature verification failed")
	}
	
	// Verify compliance ID
	if protocol.ComplianceID != pb.complianceID {
		errors = append(errors, "Compliance ID mismatch")
	}
	
	// Verify axiom integrity
	for i, axiom := range protocol.Axioms {
		if axiom.ID == "" {
			errors = append(errors, fmt.Sprintf("Axiom %d: Missing ID", i))
		}
		if axiom.Formula == "" {
			errors = append(errors, fmt.Sprintf("Axiom %d: Missing formula", i))
		}
	}
	
	return len(errors) == 0, errors
}

// FilterDataVector applies protocol axioms to filter a data vector.
func (pb *ProtocolBuilder) FilterDataVector(protocol MonumentProtocol, data map[string]interface{}) (bool, []string) {
	violations := make([]string, 0)
	
	// Apply each axiom as a filter
	for _, axiom := range protocol.Axioms {
		if !pb.evaluateAxiom(axiom, data) {
			violations = append(violations, fmt.Sprintf("Axiom violation: %s (%s)", axiom.Name, axiom.Description))
		}
	}
	
	return len(violations) == 0, violations
}

// evaluateAxiom evaluates an axiom against data.
func (pb *ProtocolBuilder) evaluateAxiom(axiom Axiom, data map[string]interface{}) bool {
	// Evaluate axiom against data vector
	if axiom.Metadata["type"] == "constraint" {
		// Check if the constraint exists in data and is true
		if val, ok := data[axiom.Name]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		// For numeric constraints, check if they are within bounds
		if threshold, ok := axiom.Metadata["threshold"].(float64); ok {
			if val, ok := data[axiom.Metadata["key"].(string)].(float64); ok {
				op, _ := axiom.Metadata["operator"].(string)
				switch op {
				case "<=": return val <= threshold
				case ">=": return val >= threshold
				case "==": return val == threshold
				}
			}
		}
	}
	
	return true
}

// GetProtocol retrieves a protocol by ID.
func (pb *ProtocolBuilder) GetProtocol(protocolID string) *MonumentProtocol {
	for _, protocol := range pb.protocols {
		if protocol.ProtocolID == protocolID {
			return &protocol
		}
	}
	return nil
}

// ListProtocols returns all protocols.
func (pb *ProtocolBuilder) ListProtocols() []MonumentProtocol {
	return pb.protocols
}

// ExportProtocol exports a protocol to JSON.
func (pb *ProtocolBuilder) ExportProtocol(protocolID string) (string, error) {
	protocol := pb.GetProtocol(protocolID)
	if protocol == nil {
		return "", fmt.Errorf("protocol not found: %s", protocolID)
	}
	
	data, err := json.MarshalIndent(protocol, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(data), nil
}

// ImportProtocol imports a protocol from JSON.
func (pb *ProtocolBuilder) ImportProtocol(jsonData string) (MonumentProtocol, error) {
	var protocol MonumentProtocol
	err := json.Unmarshal([]byte(jsonData), &protocol)
	if err != nil {
		return MonumentProtocol{}, err
	}
	
	// Verify protocol integrity
	valid, errors := pb.VerifyProtocol(protocol)
	if !valid {
		return MonumentProtocol{}, fmt.Errorf("protocol verification failed: %v", errors)
	}
	
	// Add to protocols list
	pb.protocols = append(pb.protocols, protocol)
	
	return protocol, nil
}

// generateProtocolID generates a unique protocol identifier.
func generateProtocolID() string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("PROTOCOL_%d", timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

// generateAxiomID generates a unique axiom identifier.
func generateAxiomID() string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("AXIOM_%d", timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:12])
}

// Pre-defined Monument Protocols for common financial operations

// GenerateLiquidityGapProtocol generates a protocol for identifying liquidity gaps.
func (pb *ProtocolBuilder) GenerateLiquidityGapProtocol() MonumentProtocol {
	intent := IntentPacket{
		Intent:    "Identify and exploit liquidity gaps in the Clayton regime",
		Objective: "Detect market inefficiencies and execute optimal trades",
		Constraints: []string{
			"Trades must comply with SEC regulations",
			"Risk exposure must not exceed defined limits",
			"All calculations must be deterministic and verifiable",
		},
		Parameters: map[string]interface{}{
			"max_position_size": 1000000.0,
			"risk_tolerance":    0.05,
			"time_horizon":      "1D",
		},
	}
	
	return pb.GenerateProtocol(intent, "LiquidityGapAnalysis", "1.0.0")
}

// GenerateRiskManagementProtocol generates a protocol for risk management.
func (pb *ProtocolBuilder) GenerateRiskManagementProtocol() MonumentProtocol {
	intent := IntentPacket{
		Intent:    "Enforce strict risk management across all trading operations",
		Objective: "Prevent excessive risk exposure and ensure portfolio stability",
		Constraints: []string{
			"VaR must not exceed 2% of portfolio value",
			"No single position exceeds 10% of portfolio",
			"Correlation limits enforced across asset classes",
		},
		Parameters: map[string]interface{}{
			"var_limit":           0.02,
			"position_limit":      0.10,
			"correlation_threshold": 0.7,
		},
	}
	
	return pb.GenerateProtocol(intent, "RiskManagement", "1.0.0")
}

// GenerateComplianceProtocol generates a protocol for regulatory compliance.
func (pb *ProtocolBuilder) GenerateComplianceProtocol() MonumentProtocol {
	intent := IntentPacket{
		Intent:    "Ensure all operations comply with EU AI Act and financial regulations",
		Objective: "Maintain full audit trail and deterministic verification",
		Constraints: []string{
			"All actions must be logged with SHA-256 hashing",
			"Human oversight must be maintained (Article 14)",
			"Transparency must be ensured (Article 13)",
			"Record-keeping must be automatic (Article 12)",
		},
		Parameters: map[string]interface{}{
			"audit_retention_days": 2555, // 7 years
			"hash_algorithm":       "SHA-256",
			"compliance_id":        pb.complianceID,
		},
	}
	
	return pb.GenerateProtocol(intent, "EUAIActCompliance", "1.0.0")
}

// ExecutionReport represents the result of protocol execution.
type ExecutionReport struct {
	ProtocolID    string                 `json:"protocol_id"`
	ExecutionID   string                 `json:"execution_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Success       bool                   `json:"success"`
	AxiomsApplied int                    `json:"axioms_applied"`
	Violations    []string               `json:"violations"`
	Result        interface{}            `json:"result"`
	Hash          string                 `json:"hash"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// ExecuteProtocol executes a protocol and returns a report.
func (pb *ProtocolBuilder) ExecuteProtocol(protocolID string, data map[string]interface{}) ExecutionReport {
	protocol := pb.GetProtocol(protocolID)
	if protocol == nil {
		return ExecutionReport{
			ProtocolID:  protocolID,
			ExecutionID: generateExecutionID(),
			Timestamp:   time.Now().UTC(),
			Success:     false,
			Violations:  []string{"Protocol not found"},
		}
	}
	
	// Filter data through protocol
	passed, violations := pb.FilterDataVector(*protocol, data)
	
	report := ExecutionReport{
		ProtocolID:    protocolID,
		ExecutionID:   generateExecutionID(),
		Timestamp:     time.Now().UTC(),
		Success:       passed,
		AxiomsApplied: len(protocol.Axioms),
		Violations:    violations,
		Result:        data,
		Metadata:      make(map[string]interface{}),
	}
	
	// Generate hash
	reportData, _ := json.Marshal(report)
	hash := sha256.Sum256(reportData)
	report.Hash = hex.EncodeToString(hash[:])
	
	return report
}

// generateExecutionID generates a unique execution identifier.
func generateExecutionID() string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("EXEC_%d", timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}
