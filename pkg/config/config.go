// Package config provides configuration management for AILock.
// It supports loading from YAML files and environment variables.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete AILock system configuration
type Config struct {
	Governance     GovernanceConfig     `yaml:"governance"`
	IWK            IWKConfig            `yaml:"iwk"`
	Server         ServerConfig         `yaml:"server"`
	Auth           AuthConfig           `yaml:"auth"`
	Database       DatabaseConfig       `yaml:"database"`
	Proxy          ProxyConfig          `yaml:"proxy"`
	Audit          AuditConfig          `yaml:"audit"`
	Security       SecurityConfig       `yaml:"security"`
	Observability  ObservabilityConfig  `yaml:"observability"`
	Development    DevelopmentConfig    `yaml:"development"`
}

// GovernanceConfig contains compliance and governance settings
type GovernanceConfig struct {
	ComplianceID    string  `yaml:"compliance_id"`
	TargetTCOMetric float64 `yaml:"target_tco_metric"`
}

// IWKConfig contains Invariant Wealth Kernel configuration
type IWKConfig struct {
	LicenseActive   bool   `yaml:"license_active"`
	PayoutInvariant string `yaml:"payout_invariant"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	ListenPort      int           `yaml:"listen_port"`
	ListenAddress   string        `yaml:"listen_address"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// AuthConfig contains authentication and authorization settings
type AuthConfig struct {
	JWKS  JWKSConfig  `yaml:"jwks"`
	JWT   JWTConfig   `yaml:"jwt"`
	OAuth OAuthConfig `yaml:"oauth"`
}

// JWKSConfig contains JSON Web Key Set configuration
type JWKSConfig struct {
	Endpoint        string        `yaml:"endpoint"`
	RefreshInterval time.Duration `yaml:"refresh_interval"`
	Timeout         time.Duration `yaml:"timeout"`
}

// JWTConfig contains JWT token configuration
type JWTConfig struct {
	Issuer            string        `yaml:"issuer"`
	Audience          string        `yaml:"audience"`
	Expiration        time.Duration `yaml:"expiration"`
	RefreshExpiration time.Duration `yaml:"refresh_expiration"`
	SigningAlgorithm  string        `yaml:"signing_algorithm"`
}

// OAuthConfig contains OAuth provider configuration
type OAuthConfig struct {
	Enabled   bool                       `yaml:"enabled"`
	Providers map[string]OAuthProvider `yaml:"providers"`
}

// OAuthProvider represents an OAuth provider configuration
type OAuthProvider struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Driver     string               `yaml:"driver"`
	URI        string               `yaml:"uri"`
	Pool       DatabasePoolConfig   `yaml:"pool"`
	Migrations DatabaseMigrations   `yaml:"migrations"`
}

// DatabasePoolConfig contains connection pool settings
type DatabasePoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

// DatabaseMigrations contains migration settings
type DatabaseMigrations struct {
	AutoMigrate    bool   `yaml:"auto_migrate"`
	MigrationsPath string `yaml:"migrations_path"`
}

// ProxyConfig contains DetEnforce proxy settings
type ProxyConfig struct {
	Enabled    bool                   `yaml:"enabled"`
	RateLimit  RateLimitConfig        `yaml:"rate_limit"`
	Validation ValidationConfig       `yaml:"validation"`
	TLS        TLSConfig              `yaml:"tls"`
	Allowlist  []string               `yaml:"allowlist"`
}

// RateLimitConfig contains rate limiting settings
type RateLimitConfig struct {
	MaxRequestsPerSecond int           `yaml:"max_requests_per_second"`
	BurstSize            int           `yaml:"burst_size"`
	CleanupInterval      time.Duration `yaml:"cleanup_interval"`
}

// ValidationConfig contains request validation settings
type ValidationConfig struct {
	MaxBodySize     int      `yaml:"max_body_size"`
	MaxHeaderSize   int      `yaml:"max_header_size"`
	AllowedMethods  []string `yaml:"allowed_methods"`
}

// TLSConfig contains TLS settings
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	MinVersion string `yaml:"min_version"`
}

// AuditConfig contains audit and logging settings
type AuditConfig struct {
	Enabled bool              `yaml:"enabled"`
	POE     POEConfig         `yaml:"poe"`
	Logging LoggingConfig     `yaml:"logging"`
	Storage StorageConfig     `yaml:"storage"`
}

// POEConfig contains Proof of Execution logging settings
type POEConfig struct {
	Enabled             bool   `yaml:"enabled"`
	LogPath             string `yaml:"log_path"`
	Format              string `yaml:"format"`
	IncludeRequestBody  bool   `yaml:"include_request_body"`
	IncludeResponseBody bool   `yaml:"include_response_body"`
}

// LoggingConfig contains general logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// StorageConfig contains audit storage settings
type StorageConfig struct {
	Type     string         `yaml:"type"`
	FilePath string         `yaml:"file_path"`
	Rotation RotationConfig `yaml:"rotation"`
}

// RotationConfig contains log rotation settings
type RotationConfig struct {
	MaxSize    int  `yaml:"max_size"`
	MaxAge     int  `yaml:"max_age"`
	MaxBackups int  `yaml:"max_backups"`
	Compress   bool `yaml:"compress"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	CORS      CORSConfig      `yaml:"cors"`
	Request   RequestConfig   `yaml:"request"`
	IPControl IPControlConfig `yaml:"ip_control"`
}

// CORSConfig contains CORS settings
type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// RequestConfig contains request security settings
type RequestConfig struct {
	MaxConcurrentRequests int    `yaml:"max_concurrent_requests"`
	RequestIDHeader       string `yaml:"request_id_header"`
}

// IPControlConfig contains IP filtering settings
type IPControlConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Allowlist []string `yaml:"allowlist"`
	Blocklist []string `yaml:"blocklist"`
}

// ObservabilityConfig contains observability settings
type ObservabilityConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	Health  HealthConfig  `yaml:"health"`
	Tracing TracingConfig `yaml:"tracing"`
}

// MetricsConfig contains metrics settings
type MetricsConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Endpoint   string `yaml:"endpoint"`
	Prometheus bool   `yaml:"prometheus"`
}

// HealthConfig contains health check settings
type HealthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// TracingConfig contains tracing settings
type TracingConfig struct {
	Enabled    bool    `yaml:"enabled"`
	Provider   string  `yaml:"provider"`
	Endpoint   string  `yaml:"endpoint"`
	SampleRate float64 `yaml:"sample_rate"`
}

// DevelopmentConfig contains development settings
type DevelopmentConfig struct {
	DebugMode     bool `yaml:"debug_mode"`
	PprofEnabled  bool `yaml:"pprof_enabled"`
	PprofPort     int  `yaml:"pprof_port"`
}

// Load loads configuration from a YAML file with environment variable overrides
func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	// Load from file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	applyEnvOverrides(cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// LoadFromEnv loads configuration entirely from environment variables
func LoadFromEnv() (*Config, error) {
	cfg := &Config{}
	applyEnvOverrides(cfg)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// applyEnvOverrides applies environment variable overrides to configuration
func applyEnvOverrides(cfg *Config) {
	if val := os.Getenv("AILOCK_COMPLIANCE_ID"); val != "" {
		cfg.Governance.ComplianceID = val
	}
	if val := os.Getenv("AILOCK_IWK_LICENSE_ACTIVE"); val == "true" {
		cfg.IWK.LicenseActive = true
	} else if val == "false" {
		cfg.IWK.LicenseActive = false
	}
	if val := os.Getenv("AILOCK_IWK_PAYOUT_INVARIANT"); val != "" {
		cfg.IWK.PayoutInvariant = val
	}
	if val := os.Getenv("AILOCK_SERVER_PORT"); val != "" {
		_, _ = fmt.Sscanf(val, "%d", &cfg.Server.ListenPort)
	}
	if val := os.Getenv("AILOCK_DB_URI"); val != "" {
		cfg.Database.URI = val
	}
	if val := os.Getenv("AILOCK_JWKS_ENDPOINT"); val != "" {
		cfg.Auth.JWKS.Endpoint = val
	}
	if val := os.Getenv("AILOCK_LOG_LEVEL"); val != "" {
		cfg.Audit.Logging.Level = val
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Governance validation
	if c.Governance.ComplianceID == "" {
		return fmt.Errorf("governance.compliance_id is required")
	}

	// Server validation
	if c.Server.ListenPort <= 0 || c.Server.ListenPort > 65535 {
		return fmt.Errorf("server.listen_port must be between 1 and 65535")
	}

	// Auth validation
	if c.Auth.JWKS.Endpoint == "" {
		return fmt.Errorf("auth.jwks.endpoint is required")
	}

	// Database validation
	if c.Database.Driver == "" {
		return fmt.Errorf("database.driver is required")
	}
	if c.Database.URI == "" {
		return fmt.Errorf("database.uri is required")
	}

	// IWK validation
	if c.IWK.LicenseActive && c.IWK.PayoutInvariant == "" {
		return fmt.Errorf("iwk.payout_invariant is required when license is active")
	}

	return nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.ListenAddress, c.Server.ListenPort)
}

// IsDebugMode returns whether debug mode is enabled
func (c *Config) IsDebugMode() bool {
	return c.Development.DebugMode
}

// IsIWKActive returns whether the IWK license is active
func (c *Config) IsIWKActive() bool {
	return c.IWK.LicenseActive
}
