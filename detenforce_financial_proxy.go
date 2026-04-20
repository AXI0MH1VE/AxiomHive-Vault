// detenforce_financial_proxy.go - Financial Infrastructure Proxy
// This proxy enforces deterministic access control policies for financial services.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Policy defines the access control constraints.
type Policy struct {
	AllowedPaths     map[string]bool
	ComplianceID     string
	MetricValue      float64
	ServiceActive    bool
}

var currentPolicy Policy

func init() {
	currentPolicy = loadPolicyFromEnv()
	log.Println("Financial Proxy Initialized with default security policies.")
}

func loadPolicyFromEnv() Policy {
	allowed := map[string]bool{
		"/api/v1/status": true,
	}

	if strings.EqualFold(os.Getenv("ALLOW_AUTH"), "true") {
		allowed["/api/v1/auth"] = true
	}
	if strings.EqualFold(os.Getenv("ALLOW_LEDGER"), "true") {
		allowed["/api/v1/ledger"] = true
	}

	complianceID := os.Getenv("COMPLIANCE_ID")
	if complianceID == "" {
		complianceID = "AXIOM-VAULT-PROTOTYPE"
	}

	metric := 0.0
	if v := os.Getenv("METRIC_VALUE"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			metric = parsed
		}
	}

	active := strings.EqualFold(os.Getenv("SERVICE_ACTIVE"), "true")

	return Policy{
		AllowedPaths:  allowed,
		ComplianceID:  complianceID,
		MetricValue:   metric,
		ServiceActive: active,
	}
}

func LogExecution(event string, path string, outcome string) {
	fmt.Printf("[%s] AUDIT: %s | ID: %s | Path: %s | Outcome: %s\n", time.Now().Format(time.RFC3339), event, currentPolicy.ComplianceID, path, outcome)
}

func ValidatePath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("security violation: path traversal attempt")
	}
	return nil
}

func FinancialProxyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if err := ValidatePath(path); err != nil {
		LogExecution("VALIDATION FAILURE", path, fmt.Sprintf("Deny: %v", err))
		http.Error(w, fmt.Sprintf("403 Forbidden: %v", err), http.StatusForbidden)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		LogExecution("AUTH FAILURE", path, "Missing Token")
		http.Error(w, "401 Unauthorized: Authentication Required.", http.StatusUnauthorized)
		return
	}

	if !currentPolicy.AllowedPaths[path] {
		LogExecution("ACCESS FAILURE", path, "Path not in allowlist")
		http.Error(w, "403 Forbidden: Resource not in allowlist.", http.StatusForbidden)
		return
	}

	LogExecution("EXECUTION SUCCESS", path, "Policy Compliant")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "200 OK. Metric: %.2f", currentPolicy.MetricValue)
}

func main() {
	listenPort := ":8080"
	if port := os.Getenv("LISTEN_PORT"); port != "" {
		listenPort = fmt.Sprintf(":%s", strings.TrimPrefix(port, ":"))
	}
	log.Printf("Financial Proxy starting on port %s...", listenPort)

	http.HandleFunc("/", FinancialProxyHandler)

	if err := http.ListenAndServe(listenPort, nil); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
