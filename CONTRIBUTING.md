# Contributing to AILock

Thank you for your interest in contributing to AILock! We welcome contributions from the community and are grateful for your support in making this project better.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Pull Request Process](#pull-request-process)
- [Development Setup](#development-setup)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Getting Help](#getting-help)

## Code of Conduct

By participating in this project, you agree to maintain a respectful, inclusive, and professional environment. We are committed to providing a welcoming experience for everyone, regardless of background or identity.

### Our Standards

- **Be Respectful**: Treat all contributors with respect and courtesy
- **Be Collaborative**: Work together constructively and help others
- **Be Professional**: Keep discussions focused and productive
- **Be Inclusive**: Welcome diverse perspectives and contributions

## How Can I Contribute?

There are many ways to contribute to AILock:

### Reporting Bugs

- Check if the issue has already been reported
- Use the issue template and provide detailed information
- Include steps to reproduce, expected behavior, and actual behavior
- Add relevant logs, screenshots, or error messages

### Suggesting Enhancements

- Open an issue with the "enhancement" label
- Clearly describe the feature and its benefits
- Explain the use case and motivation
- Be open to discussion and feedback

### Code Contributions

- Fix bugs or implement features from open issues
- Improve documentation or examples
- Add or improve tests
- Refactor code for better maintainability

### Documentation

- Fix typos or improve clarity
- Add examples or tutorials
- Translate documentation
- Update outdated information

## Pull Request Process

We follow a structured pull request workflow to ensure code quality and security. For detailed information on our complete PR process, including review requirements, merge strategies, and audit procedures, please see:

** [Pull Request Handling Guide](docs/PULL_REQUEST_HANDLING.md)**

This comprehensive guide covers:
- Fork and feature branch workflow
- Pre-merge checklist (tests, security, coverage, docs)
- Review and approval requirements
- Responding to feedback
- Merge strategies (squash vs. rebase)
- Conventional commits and semantic versioning
- Final audit steps (SLSA, SBOM, signatures)
- Security and revert protocols

### Quick Summary

1. **Fork and Branch**: Create a fork and feature branch from `main`
2. **Make Changes**: Follow coding standards and write tests
3. **Commit**: Use [Conventional Commits](https://www.conventionalcommits.org/) format
4. **Push and PR**: Submit a pull request with clear description
5. **Review**: Address feedback and keep PR updated
6. **Merge**: Maintainers will merge once all checks pass

## Development Setup

### Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For version control
- **Docker**: For containerized testing (optional but recommended)
- **Make**: For build automation

### Getting Started

1. **Clone your fork**:
 ```bash
 git clone https://github.com/YOUR_USERNAME/AILock.git
 cd AILock
 ```

2. **Add upstream remote**:
 ```bash
 git remote add upstream https://github.com/AxiomHiveXPII/AILock.git
 ```

3. **Install dependencies**:
 ```bash
 go mod download
 ```

4. **Run tests**:
 ```bash
 make test
 ```

5. **Build the project**:
 ```bash
 make build
 ```

## Coding Standards

### Go Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html) conventions
- Use `gofmt` to format code
- Run `golint` and address warnings
- Keep functions small and focused
- Use meaningful variable and function names

### Code Organization

- Group related code in packages
- Keep package dependencies minimal
- Use interfaces for abstraction
- Avoid circular dependencies

### Comments and Documentation

- Document all exported functions, types, and packages
- Use GoDoc format for documentation comments
- Explain "why" in comments, not just "what"
- Keep comments up-to-date with code changes

### Error Handling

- Always check and handle errors
- Provide context with error messages
- Use custom error types when appropriate
- Don't panic in library code

## Testing Guidelines

### Test Requirements

- Write tests for all new functionality
- Maintain minimum 80% code coverage
- Include unit tests for functions and methods
- Add integration tests for component interactions
- Test edge cases and error conditions

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./path/to/package

# Run with verbose output
go test -v ./...
```

### Writing Good Tests

- Use descriptive test function names
- Follow the Arrange-Act-Assert pattern
- Test one thing per test function
- Use table-driven tests for multiple scenarios
- Mock external dependencies appropriately

## Documentation

### Documentation Standards

- Keep README.md updated with major changes
- Update API documentation for interface changes
- Add examples for new features
- Document configuration options
- Include troubleshooting information

### Key Documentation Files

- [README.md](README.md) - Project overview and quick start
- [API_CONTRACT.md](docs/API_CONTRACT.md) - API specifications
- [CONFIGURATION.md](docs/CONFIGURATION.md) - Configuration guide
- [PULL_REQUEST_HANDLING.md](docs/PULL_REQUEST_HANDLING.md) - PR workflow and audit

## Getting Help

### Questions and Discussions

- **Issues**: For bugs, features, and technical questions
- **Pull Requests**: For code review and implementation discussions
- **Documentation**: Check existing docs first

### Best Practices for Getting Help

- Search existing issues before creating new ones
- Provide context and examples
- Be specific about your environment and setup
- Follow up on responses and mark issues resolved

---

## License

By contributing to AILock, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to AILock! 
