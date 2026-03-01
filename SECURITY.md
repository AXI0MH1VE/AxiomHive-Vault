# Security Policy

## Reporting a Vulnerability

**DO NOT** open public issues for security vulnerabilities.

To report a security issue, please email: **aciomhive@gmail.com** (or your preferred contact)

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.

## Supported Versions

| Version | Supported |
| ------- | ------------------ |
| main | :white_check_mark: |

## Security Measures

### Credential Management
- **All credentials must be stored in environment variables**
- Never commit `.env`, `config.yaml`, or files with secrets
- Use `config.example.yaml` as template only
- Real configs are git-ignored automatically

### Authentication
- JWT-based authentication with RS256
- JWKS endpoint for key rotation
- Token expiration enforced
- Refresh tokens expire after 7 days

### API Security
- Rate limiting: 5 req/sec per IP (configurable)
- Request validation (size limits, method allowlisting)
- TLS 1.2+ required in production
- Security headers enforced

### Audit & Compliance
- Proof of Execution (POE) logging enabled
- All API requests logged with compliance ID
- Immutable audit trail
- GDPR/CCPA compliant data handling

### Infrastructure
- Kubernetes security policies defined
- Network policies restrict pod communication
- Secrets managed via Kubernetes secrets (never in code)
- Container scanning in CI/CD pipeline

## Best Practices for Contributors

### Before Committing
1. Never commit files containing:
 - API keys or tokens
 - Database credentials
 - Private keys or certificates
 - Real configuration files

2. Always use:
 - Environment variables for secrets
 - `.example` suffix for config templates
 - Git pre-commit hooks (recommended)

3. Review changes:
 ```bash
 git diff # Review what you're about to commit
 ```

### Setting Up Pre-commit Hooks

Install gitleaks to prevent secret leaks:

```bash
# Install gitleaks
winget install gitleaks

# Or use pre-commit framework
pip install pre-commit
pre-commit install
```

### Creating Config Files

```bash
# Copy template
cp config.example.yaml config.yaml

# Edit with real credentials (this file is git-ignored)
# NEVER commit config.yaml
```

## Dependency Security

- Dependabot enabled for automatic security updates
- Regular dependency audits (monthly)
- Vulnerable dependencies updated within 7 days

## Incident Response

If a security breach occurs:

1. **Immediate Actions**
 - Revoke all compromised credentials
 - Rotate API keys and tokens
 - Change database passwords
 - Update JWT signing keys

2. **Investigation**
 - Review audit logs
 - Identify scope of breach
 - Document timeline

3. **Communication**
 - Notify affected users within 72 hours
 - Provide transparency report
 - Document lessons learned

## Security Audit Schedule

- **Weekly**: Automated dependency scanning
- **Monthly**: Manual security review
- **Quarterly**: Full penetration testing
- **Annually**: Third-party security audit

## Compliance

This project maintains compliance with:
- GDPR (General Data Protection Regulation)
- CCPA (California Consumer Privacy Act)
- SOC 2 Type II controls (planned)

## Contact

Security Team: aciomhive@gmail.com

---

**Last Updated:** 2025-11-08 
**Policy Version:** 1.0
