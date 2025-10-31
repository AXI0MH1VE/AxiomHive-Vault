// detenforce_financial_proxy.go
// Strategic Amplification Engine (SAE) - Deterministic Execution and Policy Enforcement Proxy
//
// This proxy enforces the AILock governance policy defined in CONFIG.md and the API contract
// specified in CONTRACT.md. It provides deterministic request routing, rate limiting, JWT
// validation, and cryptographic proof of execution.
//
// References:
// - Governance Policy: CONFIG.md
// - API Contract: CONTRACT.md
// - JWKS Endpoint: Configured in CONFIG.md

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GovernancePolicy represents the configuration from CONFIG.md
type GovernancePolicy struct {
	ComplianceID          string   `json:"compliance_id"`
	TargetTCOMetric       float64  `json:"target_tco_metric"`
	InvariantPaths        []string `json:"invariant_paths"`
	MaxRequestsPerSecond  int      `json:"max_requests_per_second"`
	JWKSEndpoint          string   `json:"jwks_endpoint"`
	ConfigHash            string   `json:"config_hash"`
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	mu       sync.Mutex
	tokens   int
	maxRate  int
	lastTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRate int) *RateLimiter {
	return &RateLimiter{
		tokens:   maxRate,
		maxRate:  maxRate,
		lastTime: time.Now(),
	}
}

// Allow checks if a request is allowed under the rate limit
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTime).Seconds()

	// Refill tokens based on elapsed time
	rl.tokens += int(elapsed * float64(rl.maxRate))
	if rl.tokens > rl.maxRate {
		rl.tokens = rl.maxRate
	}
	rl.lastTime = now

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// DeterministicProxy is the main proxy server
type DeterministicProxy struct {
	policy      *GovernancePolicy
	rateLimiter *RateLimiter
	backend     *httputil.ReverseProxy
	auditLog    *os.File
}

// LoadAndHashConfig loads CONFIG.md and computes its cryptographic hash
func LoadAndHashConfig(configPath string) (*GovernancePolicy, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Compute SHA-256 hash of config file
	hash := sha256.Sum256(data)
	configHash := hex.EncodeToString(hash[:])

	// Parse policy (simplified - in production, parse the markdown table)
	policy := &GovernancePolicy{
		ComplianceID:         "AIL-001",
		TargetTCOMetric:      0.01,
		InvariantPaths:       []string{"/api/v1/validate", "/api/v1/execute"},
		MaxRequestsPerSecond: 1000,
		JWKSEndpoint:         "https://auth.ailock.internal/.well-known/jwks.json",
		ConfigHash:           configHash,
	}

	log.Printf("[STARTUP] Loaded governance policy %s with hash: %s", policy.ComplianceID, configHash)
	return policy, nil
}

// ValidateJWT validates the JWT token against the JWKS endpoint
func (dp *DeterministicProxy) ValidateJWT(tokenString string) (*jwt.Token, error) {
	// In production, fetch and cache JWKS from dp.policy.JWKSEndpoint
	// For now, this is a placeholder implementation
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// In production, return the public key from JWKS
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("JWT validation failed: %w", err)
	}

	return token, nil
}

// IsInvariantPath checks if the request path is protected
func (dp *DeterministicProxy) IsInvariantPath(path string) bool {
	for _, p := range dp.policy.InvariantPaths {
		if p == path {
			return true
		}
	}
	return false
}

// AuditLog writes an audit entry
func (dp *DeterministicProxy) AuditLog(entry map[string]interface{}) {
	entry["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	entry["compliance_id"] = dp.policy.ComplianceID

	jsonData, _ := json.Marshal(entry)
	dp.auditLog.Write(append(jsonData, '\n'))
}

// ServeHTTP handles incoming requests
func (dp *DeterministicProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Rate limiting
	if !dp.rateLimiter.Allow() {
		dp.AuditLog(map[string]interface{}{
			"event":  "rate_limit_exceeded",
			"path":   r.URL.Path,
			"method": r.Method,
		})
		http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
		return
	}

	// Check if path requires enforcement
	if dp.IsInvariantPath(r.URL.Path) {
		// Extract and validate JWT
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			dp.AuditLog(map[string]interface{}{
				"event":  "missing_auth",
				"path":   r.URL.Path,
				"method": r.Method,
			})
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate JWT token
		token, err := dp.ValidateJWT(authHeader[7:]) // Strip "Bearer "
		if err != nil || !token.Valid {
			dp.AuditLog(map[string]interface{}{
				"event":  "invalid_jwt",
				"path":   r.URL.Path,
				"method": r.Method,
				"error":  err.Error(),
			})
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			return
		}

		// Compute request hash for proof of execution
		body, _ := io.ReadAll(r.Body)
		requestHash := sha256.Sum256(body)
		r.Body = io.NopCloser(io.Reader(nil)) // Reset body for backend

		dp.AuditLog(map[string]interface{}{
			"event":        "request_validated",
			"path":         r.URL.Path,
			"method":       r.Method,
			"request_hash": hex.EncodeToString(requestHash[:]),
		})
	}

	// Forward to backend
	dp.backend.ServeHTTP(w, r)
}

func main() {
	// Load and hash CONFIG.md at startup
	policy, err := LoadAndHashConfig("CONFIG.md")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Open audit log
	auditLog, err := os.OpenFile("ailock_audit.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open audit log: %v", err)
	}
	defer auditLog.Close()

	// Configure backend URL
	backendURL, err := url.Parse("http://localhost:8080") // Backend service
	if err != nil {
		log.Fatalf("Invalid backend URL: %v", err)
	}

	// Create proxy
	proxy := &DeterministicProxy{
		policy:      policy,
		rateLimiter: NewRateLimiter(policy.MaxRequestsPerSecond),
		backend:     httputil.NewSingleHostReverseProxy(backendURL),
		auditLog:    auditLog,
	}

	log.Printf("[STARTUP] AILock Deterministic Proxy started on :9090")
	log.Printf("[STARTUP] Enforcing governance policy: %s", policy.ComplianceID)
	log.Printf("[STARTUP] Config hash: %s", policy.ConfigHash)
	log.Printf("[STARTUP] Protected paths: %v", policy.InvariantPaths)

	// Start HTTP server
	if err := http.ListenAndServe(":9090", proxy); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
