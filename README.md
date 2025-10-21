# AILock

## Overview

AILock is a modern authentication and authorization system designed to provide secure access control mechanisms for applications. Built with a focus on security, flexibility, and ease of integration, AILock offers a comprehensive solution for managing user authentication flows.

Whether you're building a microservices architecture, a monolithic application, or need to add authentication to an existing system, AILock provides the tools and flexibility to secure your applications with industry-standard protocols and best practices.

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
- A supported database (PostgreSQL, MySQL, or SQLite)
- Redis (optional, for session management)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/AXI0MH1VE/AILock.git
cd AILock
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the application:
```bash
cp config.example.yaml config.yaml
# Edit config.yaml with your settings
```

4. Run database migrations:
```bash
go run cmd/migrate/main.go up
```

5. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` by default.

## Quick Start

### Basic Authentication Flow

1. **Register a new user**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "SecurePass123!"
  }'
```

2. **Login to get access token**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

3. **Access protected resources**:
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Usage

### Configuration

AILock uses a YAML configuration file. Key configuration options include:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

auth:
  jwt_secret: "your-secret-key-change-this"
  access_token_duration: 1h
  refresh_token_duration: 168h  # 7 days
  bcrypt_cost: 12

database:
  type: "postgres"
  host: "localhost"
  port: 5432
  name: "ailock"
  user: "ailock_user"
  password: "your_password"
  ssl_mode: "disable"

logging:
  level: "info"
  format: "json"
  audit_enabled: true
```

### Role-Based Access Control (RBAC)

AILock implements a flexible RBAC system:

```go
// Define roles
roles := []string{"admin", "user", "moderator"}

// Assign permissions to roles
permissions := map[string][]string{
    "admin": {"read", "write", "delete", "admin"},
    "moderator": {"read", "write", "moderate"},
    "user": {"read"},
}

// Check permissions in your handlers
if ailock.HasPermission(userID, "write") {
    // Allow action
}
```

### API Integration

Integrate AILock into your Go application:

```go
package main

import (
    "github.com/AXI0MH1VE/AILock/pkg/auth"
    "github.com/AXI0MH1VE/AILock/pkg/middleware"
)

func main() {
    // Initialize AILock
    authService := auth.NewService(config)
    
    // Use middleware for protected routes
    router.Use(middleware.AuthRequired(authService))
    router.Use(middleware.RoleRequired("admin"))
    
    // Your application logic
}
```

### DetEnforce Proxy

The DetEnforce proxy provides an additional security layer:

```bash
# Start the proxy
go run cmd/proxy/main.go --config proxy-config.yaml
```

Proxy features:
- Request filtering and validation
- Rate limiting
- DDoS protection
- Request/response logging
- Custom security rules

## API Documentation

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login and receive tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/logout` | Logout and invalidate tokens |
| POST | `/api/v1/auth/verify` | Verify email address |
| POST | `/api/v1/auth/reset-password` | Request password reset |

### User Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/user/profile` | Get user profile |
| PUT | `/api/v1/user/profile` | Update user profile |
| DELETE | `/api/v1/user/account` | Delete user account |
| GET | `/api/v1/user/sessions` | List active sessions |
| DELETE | `/api/v1/user/sessions/:id` | Revoke specific session |

### Admin Endpoints (Requires Admin Role)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/users` | List all users |
| GET | `/api/v1/admin/users/:id` | Get user details |
| PUT | `/api/v1/admin/users/:id/role` | Update user role |
| DELETE | `/api/v1/admin/users/:id` | Delete user |
| GET | `/api/v1/admin/audit-logs` | View audit logs |

## Security Best Practices

1. **Always use HTTPS in production** - Never transmit tokens over unencrypted connections
2. **Rotate JWT secrets regularly** - Update your JWT secret keys periodically
3. **Implement rate limiting** - Protect against brute force attacks
4. **Enable audit logging** - Monitor authentication events for suspicious activity
5. **Use strong password policies** - Enforce minimum complexity requirements
6. **Implement token refresh** - Use short-lived access tokens with refresh token rotation
7. **Secure your configuration** - Never commit secrets to version control

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Building

```bash
# Build the main server
go build -o bin/ailock cmd/server/main.go

# Build the proxy
go build -o bin/ailock-proxy cmd/proxy/main.go

# Build for production (with optimizations)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags=\"-s -w\" -o bin/ailock cmd/server/main.go
```

### Docker Support

```bash
# Build Docker image
docker build -t ailock:latest .

# Run with Docker Compose
docker-compose up -d
```

Example `docker-compose.yaml`:
```yaml
version: '3.8'
services:
  ailock:
    image: ailock:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
    depends_on:
      - postgres
  
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=ailock
      - POSTGRES_USER=ailock_user
      - POSTGRES_PASSWORD=your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Contributing

We welcome contributions to AILock! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- All tests pass
- Code follows Go best practices
- Documentation is updated
- Commit messages are clear and descriptive

## Troubleshooting

### Common Issues

**Issue**: "Database connection failed"
- **Solution**: Verify database credentials in `config.yaml` and ensure the database server is running

**Issue**: "Invalid JWT token"
- **Solution**: Check that the JWT secret in your config matches the one used to sign tokens

**Issue**: "Permission denied"
- **Solution**: Verify the user has the required role/permissions for the requested resource

**Issue**: "Port already in use"
- **Solution**: Change the port in `config.yaml` or stop the process using the current port

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with Go and industry-standard security libraries
- Inspired by modern authentication best practices
- Community contributions and feedback

## Support

For questions, issues, or feature requests:
- Open an issue on GitHub
- Check existing documentation
- Review closed issues for solutions

## Roadmap

- [ ] Multi-factor authentication (MFA)
- [ ] OAuth 2.0 provider support (Google, GitHub, etc.)
- [ ] WebAuthn/FIDO2 support
- [ ] GraphQL API endpoints
- [ ] Enhanced monitoring and metrics
- [ ] Kubernetes deployment templates
- [ ] CLI tool for management

---

**Note**: This is a security-focused project. Always review the code and configuration before deploying to production environments.
