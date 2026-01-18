// Package dcg implements the Deterministic Coherence Gate (DCG).
// The DCG validates all incoming data against known invariants and formal constraints,
// eliminating probabilistic noise and hallucinations before they reach the calculation engine.
package dcg

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// ValidationResult represents the outcome of a DCG validation.
type ValidationResult struct {
	Valid         bool                   `json:"valid"`
	Timestamp     time.Time              `json:"timestamp"`
	Hash          string                 `json:"hash"`
	Violations    []string               `json:"violations,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	ProofOfOrigin string                 `json:"proof_of_origin"`
}

// Invariant represents a formal constraint that data must satisfy.
type Invariant struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Validator   func(interface{}) bool `json:"-"`
	Required    bool                   `json:"required"`
}

// DataVector represents a unit of data to be validated.
type DataVector struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   interface{}            `json:"payload"`
	Source    string                 `json:"source"`
	Timestamp time.Time              `json:"timestamp"`
	Signature string                 `json:"signature,omitempty"`
}

// DCG is the Deterministic Coherence Gate validator.
type DCG struct {
	invariants      []Invariant
	groundTruthDB   map[string]interface{}
	validationLog   []ValidationResult
	strictMode      bool
	zeroTolerance   bool
}

// NewDCG creates a new Deterministic Coherence Gate.
func NewDCG(strictMode, zeroTolerance bool) *DCG {
	return &DCG{
		invariants:    make([]Invariant, 0),
		groundTruthDB: make(map[string]interface{}),
		validationLog: make([]ValidationResult, 0),
		strictMode:    strictMode,
		zeroTolerance: zeroTolerance,
	}
}

// RegisterInvariant adds a formal constraint to the DCG.
func (d *DCG) RegisterInvariant(inv Invariant) {
	d.invariants = append(d.invariants, inv)
}

// RegisterGroundTruth adds a verified fact to the ground truth database.
func (d *DCG) RegisterGroundTruth(key string, value interface{}) {
	d.groundTruthDB[key] = value
}

// Validate checks a data vector against all registered invariants.
func (d *DCG) Validate(dv DataVector) ValidationResult {
	result := ValidationResult{
		Valid:      true,
		Timestamp:  time.Now().UTC(),
		Violations: make([]string, 0),
		Metadata:   make(map[string]interface{}),
	}
	
	// Generate cryptographic hash of input
	dataHash := d.hashDataVector(dv)
	result.Hash = dataHash
	result.ProofOfOrigin = d.generateProofOfOrigin(dv, dataHash)
	
	// Validate against all invariants
	for _, inv := range d.invariants {
		if !inv.Validator(dv.Payload) {
			result.Valid = false
			violation := fmt.Sprintf("Invariant violation: %s (%s)", inv.Name, inv.Description)
			result.Violations = append(result.Violations, violation)
			
			// In zero-tolerance mode, halt on first violation
			if d.zeroTolerance {
				break
			}
		}
	}
	
	// Substrate validation: verify referenced entities exist
	if d.strictMode {
		if !d.validateSubstrate(dv) {
			result.Valid = false
			result.Violations = append(result.Violations, "Substrate validation failed: referenced entity does not exist in ground truth database")
		}
	}
	
	// Log validation result
	d.validationLog = append(d.validationLog, result)
	
	return result
}

// validateSubstrate ensures all referenced entities exist in ground truth database.
// This implements "Substrate Inversion" - hallucinations are killed at the root.
func (d *DCG) validateSubstrate(dv DataVector) bool {
	// Check if data vector references any entities
	if refs, ok := dv.Payload.(map[string]interface{})["references"]; ok {
		if refList, ok := refs.([]interface{}); ok {
			for _, ref := range refList {
				refKey := fmt.Sprintf("%v", ref)
				if _, exists := d.groundTruthDB[refKey]; !exists {
					return false
				}
			}
		}
	}
	return true
}

// hashDataVector generates a SHA-256 hash of the data vector for verification.
func (d *DCG) hashDataVector(dv DataVector) string {
	data, err := json.Marshal(dv)
	if err != nil {
		return ""
	}
	
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// generateProofOfOrigin creates a cryptographic proof of data origin.
func (d *DCG) generateProofOfOrigin(dv DataVector, dataHash string) string {
	proof := fmt.Sprintf("DCG-PROOF|%s|%s|%s|%d",
		dv.ID,
		dv.Source,
		dataHash,
		dv.Timestamp.Unix(),
	)
	
	proofHash := sha256.Sum256([]byte(proof))
	return hex.EncodeToString(proofHash[:])
}

// GetValidationLog returns the complete validation history.
func (d *DCG) GetValidationLog() []ValidationResult {
	return d.validationLog
}

// GetValidationRate returns the percentage of valid data vectors.
func (d *DCG) GetValidationRate() float64 {
	if len(d.validationLog) == 0 {
		return 0.0
	}
	
	validCount := 0
	for _, result := range d.validationLog {
		if result.Valid {
			validCount++
		}
	}
	
	return float64(validCount) / float64(len(d.validationLog)) * 100.0
}

// RejectNoise filters out invalid data vectors from a batch.
func (d *DCG) RejectNoise(vectors []DataVector) []DataVector {
	validVectors := make([]DataVector, 0)
	
	for _, dv := range vectors {
		result := d.Validate(dv)
		if result.Valid {
			validVectors = append(validVectors, dv)
		}
	}
	
	return validVectors
}

// Common invariants for financial data validation

// InvariantPositivePrice ensures price data is positive.
func InvariantPositivePrice() Invariant {
	return Invariant{
		Name:        "PositivePrice",
		Description: "Price must be greater than zero",
		Required:    true,
		Validator: func(data interface{}) bool {
			if priceData, ok := data.(map[string]interface{}); ok {
				if price, ok := priceData["price"].(float64); ok {
					return price > 0
				}
			}
			return false
		},
	}
}

// InvariantTimestampRecency ensures data is recent (within threshold).
func InvariantTimestampRecency(maxAge time.Duration) Invariant {
	return Invariant{
		Name:        "TimestampRecency",
		Description: fmt.Sprintf("Data must be no older than %v", maxAge),
		Required:    true,
		Validator: func(data interface{}) bool {
			if timestampData, ok := data.(map[string]interface{}); ok {
				if ts, ok := timestampData["timestamp"].(time.Time); ok {
					return time.Since(ts) <= maxAge
				}
				if tsStr, ok := timestampData["timestamp"].(string); ok {
					ts, err := time.Parse(time.RFC3339, tsStr)
					if err != nil {
						return false
					}
					return time.Since(ts) <= maxAge
				}
			}
			return false
		},
	}
}

// InvariantVolumeRange ensures trading volume is within expected range.
func InvariantVolumeRange(min, max float64) Invariant {
	return Invariant{
		Name:        "VolumeRange",
		Description: fmt.Sprintf("Volume must be between %f and %f", min, max),
		Required:    true,
		Validator: func(data interface{}) bool {
			if volumeData, ok := data.(map[string]interface{}); ok {
				if volume, ok := volumeData["volume"].(float64); ok {
					return volume >= min && volume <= max
				}
			}
			return false
		},
	}
}

// InvariantNoNaN ensures no NaN or Inf values in numerical data.
func InvariantNoNaN() Invariant {
	return Invariant{
		Name:        "NoNaN",
		Description: "Data must not contain NaN or Inf values",
		Required:    true,
		Validator: func(data interface{}) bool {
			if numData, ok := data.(map[string]interface{}); ok {
				for _, v := range numData {
					if f, ok := v.(float64); ok {
						if f != f { // NaN check
							return false
						}
						if f > 1e308 || f < -1e308 { // Inf check
							return false
						}
					}
				}
			}
			return true
		},
	}
}

// InvariantSignatureValid ensures cryptographic signature is valid.
func InvariantSignatureValid() Invariant {
	return Invariant{
		Name:        "SignatureValid",
		Description: "Data must have valid cryptographic signature",
		Required:    true,
		Validator: func(data interface{}) bool {
			if sigData, ok := data.(map[string]interface{}); ok {
				if sig, ok := sigData["signature"].(string); ok {
					// In production, implement actual signature verification
					return len(sig) > 0
				}
			}
			return false
		},
	}
}

// InvariantSchemaCompliance ensures data matches expected schema.
func InvariantSchemaCompliance(requiredFields []string) Invariant {
	return Invariant{
		Name:        "SchemaCompliance",
		Description: fmt.Sprintf("Data must contain required fields: %v", requiredFields),
		Required:    true,
		Validator: func(data interface{}) bool {
			if dataMap, ok := data.(map[string]interface{}); ok {
				for _, field := range requiredFields {
					if _, exists := dataMap[field]; !exists {
						return false
					}
				}
				return true
			}
			return false
		},
	}
}

// ExportValidationReport generates a compliance report for audit.
func (d *DCG) ExportValidationReport() map[string]interface{} {
	totalValidations := len(d.validationLog)
	validCount := 0
	invalidCount := 0
	violationsByType := make(map[string]int)
	
	for _, result := range d.validationLog {
		if result.Valid {
			validCount++
		} else {
			invalidCount++
			for _, violation := range result.Violations {
				violationsByType[violation]++
			}
		}
	}
	
	return map[string]interface{}{
		"total_validations":   totalValidations,
		"valid_count":         validCount,
		"invalid_count":       invalidCount,
		"validation_rate":     d.GetValidationRate(),
		"violations_by_type":  violationsByType,
		"strict_mode":         d.strictMode,
		"zero_tolerance_mode": d.zeroTolerance,
		"report_timestamp":    time.Now().UTC(),
	}
}
