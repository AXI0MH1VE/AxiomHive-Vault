// Package audit provides Proof of Execution (POE) logging for AILock.
// All security-sensitive operations are logged with immutable audit trails.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// EventType represents the type of audit event
type EventType string

const (
	EventAuthSuccess       EventType = "AUTH_SUCCESS"
	EventAuthFailure       EventType = "AUTH_FAILURE"
	EventPolicyExecution   EventType = "POLICY_EXECUTION"
	EventPolicyDenial      EventType = "POLICY_DENIAL"
	EventIWKExecution      EventType = "IWK_EXECUTION"
	EventIWKDenial         EventType = "IWK_DENIAL"
	EventRateLimitExceeded EventType = "RATE_LIMIT_EXCEEDED"
	EventInvalidRequest    EventType = "INVALID_REQUEST"
	EventSystemError       EventType = "SYSTEM_ERROR"
)

// Outcome represents the result of an operation
type Outcome string

const (
	OutcomeAllow Outcome = "ALLOW"
	OutcomeDeny  Outcome = "DENY"
	OutcomeError Outcome = "ERROR"
)

// POEEvent represents a Proof of Execution audit event
type POEEvent struct {
	Timestamp    time.Time         `json:"timestamp"`
	EventType    EventType         `json:"event_type"`
	ComplianceID string            `json:"compliance_id"`
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	Outcome      Outcome           `json:"outcome"`
	Message      string            `json:"message"`
	UserID       string            `json:"user_id,omitempty"`
	IP           string            `json:"ip,omitempty"`
	RequestID    string            `json:"request_id,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Logger is the interface for audit logging
type Logger interface {
	Log(event *POEEvent) error
	Close() error
}

// FileLogger writes audit logs to a file
type FileLogger struct {
	file          *os.File
	encoder       *json.Encoder
	textMode      bool
	mu            sync.Mutex
	complianceID  string
}

// NewFileLogger creates a new file-based audit logger
func NewFileLogger(path string, complianceID string, jsonFormat bool) (*FileLogger, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	logger := &FileLogger{
		file:         file,
		textMode:     !jsonFormat,
		complianceID: complianceID,
	}

	if jsonFormat {
		logger.encoder = json.NewEncoder(file)
	}

	return logger, nil
}

// Log writes a POE event to the audit log
func (l *FileLogger) Log(event *POEEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Set compliance ID if not already set
	if event.ComplianceID == "" {
		event.ComplianceID = l.complianceID
	}

	// Set timestamp if not already set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	if l.textMode {
		return l.logText(event)
	}
	return l.logJSON(event)
}

// logJSON writes event as JSON
func (l *FileLogger) logJSON(event *POEEvent) error {
	return l.encoder.Encode(event)
}

// logText writes event in text format matching CONTRACT.md specification
func (l *FileLogger) logText(event *POEEvent) error {
	line := fmt.Sprintf("[%s] POE: %s | ID: %s | Path: %s | Outcome: %s",
		event.Timestamp.Format(time.RFC3339),
		event.EventType,
		event.ComplianceID,
		event.Path,
		event.Outcome,
	)

	if event.Message != "" {
		line += fmt.Sprintf(" | Message: %s", event.Message)
	}

	if event.UserID != "" {
		line += fmt.Sprintf(" | User: %s", event.UserID)
	}

	if event.IP != "" {
		line += fmt.Sprintf(" | IP: %s", event.IP)
	}

	line += "\n"
	_, err := l.file.WriteString(line)
	return err
}

// Close closes the log file
func (l *FileLogger) Close() error {
	return l.file.Close()
}

// MultiLogger writes to multiple audit loggers
type MultiLogger struct {
	loggers []Logger
	mu      sync.Mutex
}

// NewMultiLogger creates a logger that writes to multiple destinations
func NewMultiLogger(loggers ...Logger) *MultiLogger {
	return &MultiLogger{
		loggers: loggers,
	}
}

// Log writes to all configured loggers
func (m *MultiLogger) Log(event *POEEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errors []error
	for _, logger := range m.loggers {
		if err := logger.Log(event); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple logging errors occurred: %v", errors)
	}

	return nil
}

// Close closes all loggers
func (m *MultiLogger) Close() error {
	var errors []error
	for _, logger := range m.loggers {
		if err := logger.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("error closing loggers: %v", errors)
	}

	return nil
}

// StdoutLogger writes audit logs to stdout
type StdoutLogger struct {
	encoder      *json.Encoder
	writer       io.Writer
	textMode     bool
	complianceID string
	mu           sync.Mutex
}

// NewStdoutLogger creates a new stdout audit logger
func NewStdoutLogger(complianceID string, jsonFormat bool) *StdoutLogger {
	logger := &StdoutLogger{
		writer:       os.Stdout,
		textMode:     !jsonFormat,
		complianceID: complianceID,
	}

	if jsonFormat {
		logger.encoder = json.NewEncoder(os.Stdout)
	}

	return logger
}

// Log writes event to stdout
func (l *StdoutLogger) Log(event *POEEvent) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if event.ComplianceID == "" {
		event.ComplianceID = l.complianceID
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	if l.textMode {
		line := fmt.Sprintf("[%s] POE: %s | ID: %s | Path: %s | Outcome: %s\n",
			event.Timestamp.Format(time.RFC3339),
			event.EventType,
			event.ComplianceID,
			event.Path,
			event.Outcome,
		)
		_, err := l.writer.Write([]byte(line))
		return err
	}

	return l.encoder.Encode(event)
}

// Close is a no-op for stdout logger
func (l *StdoutLogger) Close() error {
	return nil
}

// Helper functions for creating POE events

// NewAuthSuccessEvent creates an authentication success event
func NewAuthSuccessEvent(path, userID, ip string) *POEEvent {
	return &POEEvent{
		EventType: EventAuthSuccess,
		Path:      path,
		Outcome:   OutcomeAllow,
		UserID:    userID,
		IP:        ip,
		Message:   "Authentication successful",
	}
}

// NewAuthFailureEvent creates an authentication failure event
func NewAuthFailureEvent(path, ip, reason string) *POEEvent {
	return &POEEvent{
		EventType: EventAuthFailure,
		Path:      path,
		Outcome:   OutcomeDeny,
		IP:        ip,
		Message:   fmt.Sprintf("Authentication failed: %s", reason),
	}
}

// NewIWKExecutionEvent creates an IWK execution event
func NewIWKExecutionEvent(userID string, marketCapture float64) *POEEvent {
	return &POEEvent{
		EventType: EventIWKExecution,
		Path:      "/api/v1/strategic/wealth",
		Outcome:   OutcomeAllow,
		UserID:    userID,
		Message:   fmt.Sprintf("Autonomous Wealth Generation Initiated. Market Capture: $%.2f", marketCapture),
	}
}

// NewIWKDenialEvent creates an IWK denial event
func NewIWKDenialEvent(userID string) *POEEvent {
	return &POEEvent{
		EventType: EventIWKDenial,
		Path:      "/api/v1/strategic/wealth",
		Outcome:   OutcomeDeny,
		UserID:    userID,
		Message:   "IWK License Inactive. Strategic tactic access denied.",
	}
}

// NewRateLimitEvent creates a rate limit exceeded event
func NewRateLimitEvent(path, ip string) *POEEvent {
	return &POEEvent{
		EventType: EventRateLimitExceeded,
		Path:      path,
		Outcome:   OutcomeDeny,
		IP:        ip,
		Message:   "Rate limit exceeded",
	}
}

// NewPolicyExecutionEvent creates a policy execution event
func NewPolicyExecutionEvent(path, userID string) *POEEvent {
	return &POEEvent{
		EventType: EventPolicyExecution,
		Path:      path,
		Outcome:   OutcomeAllow,
		UserID:    userID,
		Message:   "Policy execution successful",
	}
}

// NewPolicyDenialEvent creates a policy denial event
func NewPolicyDenialEvent(path, userID, reason string) *POEEvent {
	return &POEEvent{
		EventType: EventPolicyDenial,
		Path:      path,
		Outcome:   OutcomeDeny,
		UserID:    userID,
		Message:   fmt.Sprintf("Policy denied: %s", reason),
	}
}
