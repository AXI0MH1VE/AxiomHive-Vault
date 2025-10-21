# AILock

## Overview

AILock is a modern authentication and authorization system designed to provide secure access control mechanisms for applications. Built with a focus on security, flexibility, and ease of integration, AILock offers a comprehensive solution for managing user authentication flows.

## Project Structure

The project is organized as a Go workspace with the following key components:

- **Core Authentication Module**: Handles user authentication and token management
- **Authorization Layer**: Implements role-based access control (RBAC) and permission management
- **API Gateway**: Provides RESTful endpoints for integration
- **Configuration Management**: Flexible configuration system for deployment scenarios
- **DetEnforce Proxy**: Application-layer security proxy for enhanced protection

## Features

- 🔐 **Secure Authentication**: Industry-standard authentication protocols including OAuth 2.0 and JWT
- 🛡️ **Authorization Framework**: Fine-grained access control with role and permission management
- 🔄 **Token Management**: Automatic token refresh and revocation capabilities
- 📊 **Audit Logging**: Comprehensive logging of authentication and authorization events
- 🚀 **High Performance**: Optimized for scalability and low-latency operations
- 🔧 **Easy Integration**: Clean API design for seamless integration with existing systems
- 🔒 **Proxy Protection**: DetEnforce proxy for advanced security enforcement at the application layer

## Installation

### Prerequisites

- Go 1.19 or higher
- Git

### Setup

1. Clone the repository:
```bash
git clone https://github.com/AXI0MH1VE/AILock.git
cd AILock
```

2. Install dependencies:
```bash
go work sync
go mod download
```

3. Build the project:
```bash
go build ./...
```

## Proxy Configuration

AILock now includes **DetEnforce Proxy**, a powerful application-layer security proxy that provides advanced protection mechanisms including request filtering, rate limiting, and security policy enforcement.

### Proxy Overview

The DetEnforce Proxy acts as a security gateway between clients and your AILock services, offering:

- **Request Validation**: Comprehensive validation of incoming requests
- **Rate Limiting**: Protection against abuse and DoS attacks
- **Security Policies**: Configurable security rules and enforcement
- **Logging & Monitoring**: Detailed audit trails of all proxy activity
- **TLS Termination**: Secure communication with automatic certificate management

### Proxy Setup

The proxy configuration is located in `ARTIFACTS/DETENFORCE_PROXY_CORE.yaml`. This YAML file contains all the necessary settings for deploying and operating the proxy.

#### Configuration File Structure

The proxy configuration includes the following key sections:

1. **Server Settings**: Listen address, ports, and TLS configuration
2. **Upstream Targets**: Backend services that the proxy forwards requests to
3. **Security Policies**: Rules for request filtering and validation
4. **Rate Limiting**: Thresholds and rate limit configurations
5. **Logging**: Audit log settings and output destinations

#### Example Configuration

```yaml
# Basic proxy configuration example
server:
  listen_address: "0.0.0.0"
  port: 8443
  tls:
    enabled: true
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"

upstreams:
  - name: "ailock-api"
    address: "localhost:8080"
    health_check:
      enabled: true
      interval: 30s

security:
  rate_limiting:
    enabled: true
    requests_per_minute: 100
  request_validation:
    enabled: true
    max_body_size: 1048576

logging:
  level: "info"
  format: "json"
  output: "/var/log/detenforce-proxy.log"
```

### Deploying the Proxy

1. **Configure the proxy**: Edit `ARTIFACTS/DETENFORCE_PROXY_CORE.yaml` to match your environment:
   - Set the appropriate listen address and port
   - Configure upstream backend services
   - Adjust security policies and rate limits as needed
   - Set up TLS certificates for production use

2. **Start the proxy** (example command, adjust based on your deployment method):
```bash
# Using the configuration file
./detenforce-proxy --config ARTIFACTS/DETENFORCE_PROXY_CORE.yaml
```

3. **Verify proxy operation**:
```bash
# Check proxy health
curl -k https://localhost:8443/health

# Test proxied request to AILock
curl -k https://localhost:8443/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "pass"}'
```

### Proxy Integration with AILock

To integrate the proxy with your AILock deployment:

1. **Update client configurations** to point to the proxy endpoint instead of directly to AILock
2. **Configure the proxy upstream** to point to your AILock service backend
3. **Set up security policies** that align with your authentication requirements
4. **Monitor proxy logs** for security events and performance metrics

### Security Considerations

- **TLS/SSL**: Always enable TLS in production environments
- **Rate Limiting**: Adjust rate limits based on your expected traffic patterns
- **Request Validation**: Configure appropriate request size limits and validation rules
- **Logging**: Ensure audit logs are properly stored and rotated
- **Upstream Health Checks**: Enable health checking to ensure high availability

### Troubleshooting

- **Connection Issues**: Verify the proxy is listening on the correct address/port
- **Certificate Errors**: Ensure TLS certificates are valid and properly configured
- **Rate Limiting**: Check logs if requests are being blocked due to rate limits
- **Upstream Failures**: Verify backend services are running and accessible

For detailed configuration options, refer to the `DETENFORCE_PROXY_CORE.yaml` file in the ARTIFACTS directory.

## Usage

Start the AILock service:
```bash
go run main.go
```

### API Examples

Authentication endpoint:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "pass"}'
```

## Contributing

We welcome contributions to AILock! Here's how you can help:

1. **Fork the Repository**: Create your own fork of the project
2. **Create a Feature Branch**: `git checkout -b feature/your-feature-name`
3. **Make Your Changes**: Implement your feature or bug fix
4. **Write Tests**: Ensure your code is well-tested
5. **Commit Your Changes**: Use clear and descriptive commit messages
6. **Push to Your Fork**: `git push origin feature/your-feature-name`
7. **Submit a Pull Request**: Open a PR with a clear description of your changes

### Development Guidelines

- Follow Go best practices and idioms
- Write comprehensive tests for new functionality
- Update documentation as needed
- Ensure all tests pass before submitting PRs

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For questions, issues, or suggestions, please open an issue on GitHub or contact the maintainers.

---

Built with ❤️ by the AILock team
