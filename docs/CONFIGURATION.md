# Configuration Guide

## Overview

This document describes the configuration options for AILock's detection enforcement proxy, with particular focus on secure TLS implementation.

## TLS Configuration

### Security Requirements

The proxy implements robust TLS security according to industry best practices:

1. **Protocol Versions**: Only TLS 1.2 and TLS 1.3 are supported. SSLv3, TLS 1.0, and TLS 1.1 are explicitly disabled.

2. **Cipher Suites**: Modern, secure cipher suites are configured:
 - TLS 1.3: `TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`, `TLS_CHACHA20_POLY1305_SHA256`
 - TLS 1.2: `ECDHE-RSA-AES128-GCM-SHA256`, `ECDHE-RSA-AES256-GCM-SHA384`, `ECDHE-RSA-CHACHA20-POLY1305`
 - Forward secrecy enabled through ECDHE key exchange
 - No CBC-mode ciphers (vulnerable to padding oracle attacks)
 - No RC4, 3DES, or other weak ciphers

3. **ACME Certificate Management**: Automated certificate provisioning and renewal via Let's Encrypt or other ACME providers:
 - Supported challenge types: `http-01`, `tls-alpn-01`
 - Automatic renewal starting 30 days before expiration
 - Staging environment support for testing

4. **Certificate Storage**: Certificates are stored securely with proper permissions:
 - Private keys: 0600 (owner read/write only)
 - Certificates: 0644 (world-readable)
 - Configurable storage path

### CLI Flags for TLS

The following command-line flags control TLS behavior:

#### Required Flags

- `--tls-enabled`: Enable TLS/HTTPS on the proxy (boolean, default: `false`)
- `--tls-port`: Port for TLS listener (default: `8443`)

#### Certificate Configuration

**Manual Certificate Mode:**
- `--tls-cert-file`: Path to TLS certificate file (PEM format)
- `--tls-key-file`: Path to TLS private key file (PEM format)

**ACME/Let's Encrypt Mode:**
- `--acme-enabled`: Enable ACME automatic certificate provisioning (boolean)
- `--acme-email`: Email address for ACME account registration (required)
- `--acme-domains`: Comma-separated list of domains for certificate (required, e.g., `api.example.com,www.example.com`)
- `--acme-challenge`: Challenge type to use: `http-01` or `tls-alpn-01` (default: `http-01`)
- `--acme-storage-path`: Directory for ACME certificate storage (default: `./acme-certs`)
- `--acme-staging`: Use Let's Encrypt staging environment for testing (boolean, default: `false`)

#### Advanced Options

- `--tls-min-version`: Minimum TLS version (default: `1.2`, options: `1.2`, `1.3`)
- `--tls-cipher-suites`: Comma-separated list of allowed cipher suites (overrides defaults)
- `--tls-renewal-days`: Days before expiration to start renewal (default: `30`)

### Environment Separation

#### Production Configuration

```bash
./ailock-proxy \
 --tls-enabled \
 --tls-port 443 \
 --acme-enabled \
 --acme-email admin@example.com \
 --acme-domains api.example.com \
 --acme-challenge http-01 \
 --acme-storage-path /var/lib/ailock/certs \
 --tls-min-version 1.2
```

#### Staging/Testing Configuration

```bash
./ailock-proxy \
 --tls-enabled \
 --tls-port 8443 \
 --acme-enabled \
 --acme-email test@example.com \
 --acme-domains test.example.com \
 --acme-challenge http-01 \
 --acme-storage-path ./test-certs \
 --acme-staging \
 --tls-min-version 1.2
```

**Note**: The `--acme-staging` flag uses Let's Encrypt's staging environment, which has higher rate limits and issues untrusted certificates for testing.

### Example Valid Configuration

#### YAML Configuration File (config.yaml)

```yaml
tls:
 enabled: true
 port: 443
 min_version: "1.2"

 # ACME configuration
 acme:
 enabled: true
 email: "security@example.com"
 domains:
 - "api.example.com"
 - "proxy.example.com"
 challenge_type: "http-01"
 storage_path: "/var/lib/ailock/acme-certs"
 staging: false
 renewal_days: 30

 # Or use manual certificates
 # manual:
 # cert_file: "/etc/ailock/tls/cert.pem"
 # key_file: "/etc/ailock/tls/key.pem"

 # Advanced cipher configuration (optional)
 cipher_suites:
 - "TLS_AES_128_GCM_SHA256"
 - "TLS_AES_256_GCM_SHA384"
 - "TLS_CHACHA20_POLY1305_SHA256"
 - "ECDHE-RSA-AES128-GCM-SHA256"
 - "ECDHE-RSA-AES256-GCM-SHA384"
```

#### Docker Compose Example

```yaml
version: '3.8'
services:
 ailock-proxy:
 image: ailock/proxy:latest
 command:
 - --tls-enabled
 - --tls-port=443
 - --acme-enabled
 - --acme-email=admin@example.com
 - --acme-domains=api.example.com
 - --acme-challenge=http-01
 - --acme-storage-path=/certs
 - --tls-min-version=1.2
 ports:
 - "443:443"
 - "80:80" # Required for http-01 challenge
 volumes:
 - ./certs:/certs
 environment:
 - LOG_LEVEL=info
```

### Security Best Practices

1. **Use ACME for Production**: Automated certificate management reduces the risk of expired certificates
2. **Separate Staging and Production**: Always test configuration changes in staging first
3. **Monitor Certificate Expiry**: Set up alerts for certificates expiring within 7 days
4. **Secure Storage Paths**: Ensure certificate storage directories have appropriate permissions (0700)
5. **Regular Updates**: Keep TLS libraries and the proxy software up to date
6. **Rate Limits**: Be aware of Let's Encrypt rate limits (50 certificates per registered domain per week)

### Troubleshooting

#### Certificate Provisioning Fails

1. Verify DNS records point to the proxy server
2. Ensure ports 80 (http-01) or 443 (tls-alpn-01) are accessible
3. Check firewall rules
4. Review logs for specific ACME errors
5. Try staging environment first: `--acme-staging`

#### TLS Handshake Errors

1. Verify client supports TLS 1.2+
2. Check cipher suite compatibility
3. Ensure certificate is valid and not expired
4. Verify certificate domain matches requested hostname

## Additional Configuration

For other configuration options (rate limiting, detection thresholds, logging), see:
- [ARTIFACTS/DETENFORCE_PROXY_CORE.yaml](../ARTIFACTS/DETENFORCE_PROXY_CORE.yaml) - Core configuration schema
- [API_CONTRACT.md](./API_CONTRACT.md) - API endpoints and authentication
