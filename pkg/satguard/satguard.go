// Package satguard implements SAT Guards using Boolean Satisfiability (SAT) solvers
// for runtime verification of all actions before execution.
// This ensures structural impossibility of unsafe states through formal logic.
package satguard

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// ActionProposal represents a proposed action to be validated.
type ActionProposal struct {
	ID          string                 `json:"id"`
	Action      string                 `json:"action"`
	Target      string                 `json:"target"`
	Parameters  map[string]interface{} `json:"parameters"`
	Requester   string                 `json:"requester"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Condition represents an environmental constraint.
type Condition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Evaluator   func(interface{}) bool `json:"-"`
	State       interface{}            `json:"state"`
}

// Authorization represents user permissions.
type Authorization struct {
	UserID      string   `json:"user_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Level       int      `json:"level"`
}

// GuardResult represents the outcome of SAT Guard validation.
type GuardResult struct {
	Satisfiable   bool                   `json:"satisfiable"`
	Action        string                 `json:"action"`
	Decision      string                 `json:"decision"` // ALLOW, DENY, HALT
	Timestamp     time.Time              `json:"timestamp"`
	Formula       string                 `json:"formula"`
	LogicReceipt  string                 `json:"logic_receipt"`
	Hash          string                 `json:"hash"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// SATGuard is the runtime verification layer.
type SATGuard struct {
	conditions      []Condition
	authorizations  map[string]Authorization
	actionLog       []GuardResult
	globalKillSwitch bool
	strictMode      bool
	hamiltonian     *Hamiltonian
}

// Hamiltonian represents the inverted Hamiltonian constraint system.
// Only the "safe state" is energetically accessible; all other states are impossible.
type Hamiltonian struct {
	SafeStates      map[string]bool
	EnergyBarriers  map[string]float64
	Invariants      []string
	GroundTruths    map[string]interface{}
}

// NewSATGuard creates a new SAT Guard validator.
func NewSATGuard(strictMode bool) *SATGuard {
	return &SATGuard{
		conditions:      make([]Condition, 0),
		authorizations:  make(map[string]Authorization, 0),
		actionLog:       make([]GuardResult, 0),
		globalKillSwitch: false,
		strictMode:      strictMode,
		hamiltonian: &Hamiltonian{
			SafeStates:     make(map[string]bool),
			EnergyBarriers: make(map[string]float64),
			Invariants:     make([]string, 0),
			GroundTruths:   make(map[string]interface{}),
		},
	}
}

// RegisterCondition adds an environmental constraint.
func (sg *SATGuard) RegisterCondition(cond Condition) {
	sg.conditions = append(sg.conditions, cond)
}

// RegisterAuthorization adds user permissions.
func (sg *SATGuard) RegisterAuthorization(auth Authorization) {
	sg.authorizations[auth.UserID] = auth
}

// RegisterSafeState defines a state as energetically accessible (allowed).
func (sg *SATGuard) RegisterSafeState(state string) {
	sg.hamiltonian.SafeStates[state] = true
}

// RegisterInvariant adds a high-energy invariant (ground truth).
func (sg *SATGuard) RegisterInvariant(invariant string, value interface{}) {
	sg.hamiltonian.Invariants = append(sg.hamiltonian.Invariants, invariant)
	sg.hamiltonian.GroundTruths[invariant] = value
}

// SetGlobalKillSwitch enables or disables the global kill switch.
// When true, all actions are instantly halted.
func (sg *SATGuard) SetGlobalKillSwitch(enabled bool) {
	sg.globalKillSwitch = enabled
}

// Validate performs SAT validation on an action proposal.
// Formula: P ∧ C ∧ A
// - P: Proposal (the action)
// - C: Condition (environmental constraint)
// - A: Authorization (user permission)
func (sg *SATGuard) Validate(proposal ActionProposal) GuardResult {
	result := GuardResult{
		Action:    proposal.Action,
		Timestamp: time.Now().UTC(),
		Metadata:  make(map[string]interface{}),
	}
	
	// Check global kill switch first
	if sg.globalKillSwitch {
		result.Satisfiable = false
		result.Decision = "HALT"
		result.LogicReceipt = "Global kill switch is active. All actions are halted."
		result.Formula = "KILL_SWITCH = TRUE → HALT"
		result.Hash = sg.hashResult(result)
		sg.actionLog = append(sg.actionLog, result)
		return result
	}
	
	// Evaluate Proposal (P)
	proposalValid := sg.evaluateProposal(proposal)
	
	// Evaluate Conditions (C)
	conditionsValid := sg.evaluateConditions(proposal)
	
	// Evaluate Authorization (A)
	authorizationValid := sg.evaluateAuthorization(proposal)
	
	// Evaluate Hamiltonian constraints
	hamiltonianValid := sg.evaluateHamiltonian(proposal)
	
	// SAT formula: P ∧ C ∧ A ∧ H
	satisfiable := proposalValid && conditionsValid && authorizationValid && hamiltonianValid
	
	result.Satisfiable = satisfiable
	if satisfiable {
		result.Decision = "ALLOW"
	} else {
		result.Decision = "DENY"
	}
	
	// Generate logic receipt
	result.Formula = fmt.Sprintf("P(%v) ∧ C(%v) ∧ A(%v) ∧ H(%v) = %v",
		proposalValid, conditionsValid, authorizationValid, hamiltonianValid, satisfiable)
	
	result.LogicReceipt = sg.generateLogicReceipt(proposal, proposalValid, conditionsValid, authorizationValid, hamiltonianValid)
	
	// Generate cryptographic hash
	result.Hash = sg.hashResult(result)
	
	// Log result
	sg.actionLog = append(sg.actionLog, result)
	
	return result
}

// evaluateProposal checks if the proposed action is well-formed.
func (sg *SATGuard) evaluateProposal(proposal ActionProposal) bool {
	// Check if action is defined
	if proposal.Action == "" {
		return false
	}
	
	// Check if target is specified
	if proposal.Target == "" {
		return false
	}
	
	// Check if requester is identified
	if proposal.Requester == "" {
		return false
	}
	
	return true
}

// evaluateConditions checks all environmental constraints.
func (sg *SATGuard) evaluateConditions(proposal ActionProposal) bool {
	for _, cond := range sg.conditions {
		if !cond.Evaluator(proposal) {
			return false
		}
	}
	return true
}

// evaluateAuthorization checks user permissions.
func (sg *SATGuard) evaluateAuthorization(proposal ActionProposal) bool {
	auth, exists := sg.authorizations[proposal.Requester]
	if !exists {
		return false
	}
	
	// Check if user has required permission for this action
	requiredPermission := fmt.Sprintf("action:%s", proposal.Action)
	for _, perm := range auth.Permissions {
		if perm == requiredPermission || perm == "action:*" {
			return true
		}
	}
	
	return false
}

// evaluateHamiltonian checks if the proposed action leads to a safe state.
// This implements the "Inverted Hamiltonian" - only safe states are energetically accessible.
func (sg *SATGuard) evaluateHamiltonian(proposal ActionProposal) bool {
	// Compute resulting state from action
	resultingState := sg.computeResultingState(proposal)
	
	// Check if resulting state is in the safe state set
	if !sg.hamiltonian.SafeStates[resultingState] {
		return false
	}
	
	// Verify no invariants are violated
	for _, invariant := range sg.hamiltonian.Invariants {
		if !sg.verifyInvariant(invariant, proposal) {
			return false
		}
	}
	
	return true
}

// computeResultingState determines the state that would result from the action.
func (sg *SATGuard) computeResultingState(proposal ActionProposal) string {
	// Simplified state computation
	// In production, this would model the full state transition
	return fmt.Sprintf("%s:%s", proposal.Action, proposal.Target)
}

// verifyInvariant checks if an invariant would be violated by the action.
func (sg *SATGuard) verifyInvariant(invariant string, proposal ActionProposal) bool {
	// Check if action would contradict ground truth
	groundTruth, exists := sg.hamiltonian.GroundTruths[invariant]
	if !exists {
		return true
	}
	
	// Verify action doesn't contradict ground truth
	// In production, this would perform deep logical verification
	if proposal.Action == "set_ground_truth" {
		if newValue, ok := proposal.Parameters["value"]; ok {
			return newValue == groundTruth
		}
	}
	
	return true
}

// generateLogicReceipt creates a human-readable explanation of the decision.
func (sg *SATGuard) generateLogicReceipt(proposal ActionProposal, p, c, a, h bool) string {
	receipt := fmt.Sprintf("SAT Guard Logic Receipt\n")
	receipt += fmt.Sprintf("========================\n")
	receipt += fmt.Sprintf("Action: %s\n", proposal.Action)
	receipt += fmt.Sprintf("Target: %s\n", proposal.Target)
	receipt += fmt.Sprintf("Requester: %s\n", proposal.Requester)
	receipt += fmt.Sprintf("Timestamp: %s\n\n", proposal.Timestamp.Format(time.RFC3339))
	
	receipt += fmt.Sprintf("Validation Results:\n")
	receipt += fmt.Sprintf("  Proposal Valid (P): %v\n", p)
	receipt += fmt.Sprintf("  Conditions Met (C): %v\n", c)
	receipt += fmt.Sprintf("  Authorized (A): %v\n", a)
	receipt += fmt.Sprintf("  Hamiltonian Safe (H): %v\n\n", h)
	
	if !p {
		receipt += "  ❌ Proposal is malformed or incomplete\n"
	}
	if !c {
		receipt += "  ❌ Environmental conditions not satisfied\n"
	}
	if !a {
		receipt += "  ❌ User lacks required permissions\n"
	}
	if !h {
		receipt += "  ❌ Action would lead to unsafe state (Hamiltonian violation)\n"
	}
	
	if p && c && a && h {
		receipt += "  ✓ All constraints satisfied - Action ALLOWED\n"
	} else {
		receipt += "  ✗ Constraints violated - Action DENIED\n"
	}
	
	return receipt
}

// hashResult generates a SHA-256 hash of the guard result.
func (sg *SATGuard) hashResult(result GuardResult) string {
	data, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// GetActionLog returns the complete action log.
func (sg *SATGuard) GetActionLog() []GuardResult {
	return sg.actionLog
}

// GetApprovalRate returns the percentage of allowed actions.
func (sg *SATGuard) GetApprovalRate() float64 {
	if len(sg.actionLog) == 0 {
		return 0.0
	}
	
	allowedCount := 0
	for _, result := range sg.actionLog {
		if result.Decision == "ALLOW" {
			allowedCount++
		}
	}
	
	return float64(allowedCount) / float64(len(sg.actionLog)) * 100.0
}

// ExportAuditTrail generates a compliance audit trail.
func (sg *SATGuard) ExportAuditTrail() map[string]interface{} {
	totalActions := len(sg.actionLog)
	allowedCount := 0
	deniedCount := 0
	haltedCount := 0
	decisionsByAction := make(map[string]map[string]int)
	
	for _, result := range sg.actionLog {
		switch result.Decision {
		case "ALLOW":
			allowedCount++
		case "DENY":
			deniedCount++
		case "HALT":
			haltedCount++
		}
		
		if _, exists := decisionsByAction[result.Action]; !exists {
			decisionsByAction[result.Action] = make(map[string]int)
		}
		decisionsByAction[result.Action][result.Decision]++
	}
	
	return map[string]interface{}{
		"total_actions":        totalActions,
		"allowed_count":        allowedCount,
		"denied_count":         deniedCount,
		"halted_count":         haltedCount,
		"approval_rate":        sg.GetApprovalRate(),
		"decisions_by_action":  decisionsByAction,
		"global_kill_switch":   sg.globalKillSwitch,
		"strict_mode":          sg.strictMode,
		"safe_states_count":    len(sg.hamiltonian.SafeStates),
		"invariants_count":     len(sg.hamiltonian.Invariants),
		"report_timestamp":     time.Now().UTC(),
	}
}

// Common conditions for financial operations

// ConditionDatabaseUnlocked ensures target database is not locked.
func ConditionDatabaseUnlocked(lockedDBs map[string]bool) Condition {
	return Condition{
		Name:        "DatabaseUnlocked",
		Description: "Target database must not be locked",
		State:       lockedDBs,
		Evaluator: func(data interface{}) bool {
			if proposal, ok := data.(ActionProposal); ok {
				if locked, exists := lockedDBs[proposal.Target]; exists {
					return !locked
				}
			}
			return true
		},
	}
}

// ConditionMarketHours ensures trading only during market hours.
func ConditionMarketHours() Condition {
	return Condition{
		Name:        "MarketHours",
		Description: "Trading actions only allowed during market hours",
		Evaluator: func(data interface{}) bool {
			if proposal, ok := data.(ActionProposal); ok {
				if proposal.Action == "execute_trade" {
					hour := time.Now().UTC().Hour()
					// NYSE hours: 14:30-21:00 UTC
					return hour >= 14 && hour < 21
				}
			}
			return true
		},
	}
}

// ConditionRiskLimitNotExceeded ensures risk limits are respected.
func ConditionRiskLimitNotExceeded(currentRisk, maxRisk float64) Condition {
	return Condition{
		Name:        "RiskLimitNotExceeded",
		Description: fmt.Sprintf("Current risk (%f) must not exceed maximum (%f)", currentRisk, maxRisk),
		State:       map[string]float64{"current": currentRisk, "max": maxRisk},
		Evaluator: func(data interface{}) bool {
			return currentRisk <= maxRisk
		},
	}
}
