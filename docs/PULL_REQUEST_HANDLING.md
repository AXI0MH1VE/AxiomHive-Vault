# Pull Request Handling Guide

This document outlines the complete pull request workflow, review process, merge strategy, and audit steps for the AILock project.

## 1. Pull Request Workflow

### Fork and Branch Strategy

1. **Fork the Repository**
 - Create a personal fork of the AILock repository
 - Clone your fork locally: `git clone https://github.com/YOUR_USERNAME/AILock.git`
 - Add upstream remote: `git remote add upstream https://github.com/AxiomHiveXPII/AILock.git`

2. **Create a Feature Branch**
 - Always branch from `main`: `git checkout main && git pull upstream main`
 - Use descriptive branch names following convention:
 - `feature/description` - for new features
 - `fix/description` - for bug fixes
 - `docs/description` - for documentation
 - `refactor/description` - for code refactoring
 - `security/description` - for security patches
 - Example: `git checkout -b feature/add-oauth-support`

3. **Make Your Changes**
 - Write clear, focused commits
 - Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:
 - `feat:` - new features
 - `fix:` - bug fixes
 - `docs:` - documentation changes
 - `test:` - test additions or modifications
 - `refactor:` - code refactoring
 - `security:` - security-related changes
 - `chore:` - maintenance tasks
 - Keep commits atomic and logically grouped

4. **Push and Create PR**
 - Push to your fork: `git push origin feature/your-branch`
 - Navigate to the AILock repository on GitHub
 - Click "New Pull Request" or use the prompted banner
 - Select your fork and branch as the source
 - Provide a clear, descriptive title and detailed description
 - Reference any related issues with `Closes #123` or `Fixes #456`

## 2. Pre-Merge Checklist

Before submitting a PR for review, ensure the following items are complete:

### Tests and CI
- [ ] All existing tests pass locally
- [ ] New tests added for new functionality (aim for >80% coverage)
- [ ] Edge cases and error conditions tested
- [ ] CI pipeline passes completely (GitHub Actions)
- [ ] No test flakiness or intermittent failures
- [ ] Integration tests updated if applicable

### Security Scan
- [ ] Static Application Security Testing (SAST) passes
- [ ] Dependency vulnerability scan clean (Dependabot/Snyk)
- [ ] No secrets, API keys, or credentials committed
- [ ] Security-sensitive code reviewed for vulnerabilities
- [ ] Input validation and sanitization implemented where needed

### Code Coverage
- [ ] Test coverage meets or exceeds project threshold (minimum 80%)
- [ ] Coverage report generated and reviewed
- [ ] Critical paths have comprehensive test coverage
- [ ] No significant coverage regression introduced

### Documentation Updated
- [ ] README.md updated if user-facing changes
- [ ] API documentation updated (if applicable)
- [ ] Inline code comments added for complex logic
- [ ] CHANGELOG.md prepared (draft entry)
- [ ] Configuration examples updated if needed
- [ ] Migration guide provided for breaking changes

## 3. Required Review and Approval Gates

### Reviewer Requirements
- **Minimum Reviewers**: 2 maintainers must approve
- **Code Owner Approval**: Required for changes to:
 - Security modules (`/security/**`)
 - Core infrastructure (`/infrastructure/**`)
 - API contracts (`/docs/API_CONTRACT.md`)
 - CI/CD pipelines (`.github/workflows/**`)

### Review Checklist

Reviewers should verify:

#### Code Quality
- [ ] Code follows project style guidelines and conventions
- [ ] Logic is clear, maintainable, and well-structured
- [ ] No unnecessary complexity or over-engineering
- [ ] DRY (Don't Repeat Yourself) principle followed
- [ ] Appropriate error handling and logging

#### Security Review
- [ ] No injection vulnerabilities (SQL, command, etc.)
- [ ] Proper authentication and authorization checks
- [ ] Sensitive data properly encrypted/protected
- [ ] Input validation and sanitization adequate
- [ ] Dependencies vetted and up-to-date

#### Testing
- [ ] Tests are comprehensive and meaningful
- [ ] Test names clearly describe what they test
- [ ] Mocks/stubs used appropriately
- [ ] No hard-coded test data that could cause issues

#### Documentation
- [ ] Code changes are well-documented
- [ ] Public APIs have clear docstrings
- [ ] Breaking changes clearly called out
- [ ] Examples provided where helpful

#### Architecture
- [ ] Changes align with project architecture
- [ ] No architectural anti-patterns introduced
- [ ] Performance implications considered
- [ ] Scalability considerations addressed

## 4. Responding to Review Feedback

### Best Practices

1. **Acknowledge Feedback Promptly**
 - Respond to all review comments
 - Ask clarifying questions if feedback is unclear
 - Mark conversations as resolved only after addressing

2. **Update Your PR**
 - Make requested changes in new commits (don't force-push during review)
 - Reference which review comments each commit addresses
 - Push updates: `git push origin feature/your-branch`

3. **Re-request Review**
 - After addressing feedback, re-request review from reviewers
 - Add a comment summarizing what was changed
 - Be specific: "Updated per @reviewer's suggestion in comment #3"

4. **Handling Disagreements**
 - Engage constructively and professionally
 - Provide technical rationale for alternative approaches
 - Escalate to project lead if consensus cannot be reached
 - Be willing to compromise for the good of the project

5. **Keeping PR Updated**
 - Regularly sync with upstream `main`: `git pull upstream main`
 - Resolve merge conflicts promptly
 - Rebase if requested: `git rebase upstream/main`

## 5. Merge Strategy

### Merge Methods

**Default: Squash and Merge**
- Used for most feature branches
- Condenses all commits into a single commit
- Keeps `main` history clean and linear
- Commit message must follow Conventional Commits
- Automatically includes PR number reference

**Rebase and Merge**
- Used when commit history is already clean and valuable
- Each commit in PR must:
 - Follow Conventional Commits format
 - Be atomic and self-contained
 - Pass all tests individually
- Maintains detailed commit history
- Requires pre-approval from maintainer

**Standard Merge (Rarely)**
- Only for long-lived feature branches with multiple contributors
- Creates merge commit
- Preserves complete branch history
- Requires explicit maintainer approval

### Conventional Commits Format

All merge commits must follow this format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Example:**
```
feat(auth): add OAuth2 support for GitHub login

Implements OAuth2 authentication flow using GitHub as provider.
Includes token refresh mechanism and session management.

Closes #145
```

### Semantic Versioning

Commit types map to version bumps:
- `feat:` → Minor version bump (0.X.0)
- `fix:` → Patch version bump (0.0.X)
- `BREAKING CHANGE:` or `!` → Major version bump (X.0.0)
- Other types → No version bump (until next release)

### Changelog Update

- CHANGELOG.md is updated automatically via CI after merge
- Based on Conventional Commits messages
- Grouped by type (Features, Fixes, etc.)
- Manual curation by maintainer before release
- Draft changes visible in "Unreleased" section

## 6. Final Audit Steps Before Merge

Maintainer performing merge must verify:

### CI/CD Status
- [ ] All CI checks are green (passing)
- [ ] No workflows skipped or cancelled
- [ ] Build artifacts generated successfully
- [ ] Test reports show 100% pass rate

### Security and Provenance

#### SLSA Provenance Check
- [ ] SLSA (Supply-chain Levels for Software Artifacts) level verified
- [ ] Build provenance attestation generated
- [ ] Provenance signature validated
- [ ] No unsigned commits from unknown sources

#### SBOM Verification
- [ ] Software Bill of Materials (SBOM) generated
- [ ] All dependencies declared and tracked
- [ ] No undeclared or hidden dependencies
- [ ] SBOM format validated (SPDX or CycloneDX)
- [ ] Vulnerability scan run against SBOM

#### Signature Verification
- [ ] All commits are signed (GPG/SSH)
- [ ] Committer identity verified
- [ ] Signatures validated against known keys
- [ ] No unverified or suspicious signatures

### Badge and Status Checks
- [ ] Build status badge shows passing
- [ ] Coverage badge reflects accurate percentage
- [ ] Security scan badge shows "no vulnerabilities"
- [ ] License compliance badge passing
- [ ] All required status checks completed

### Documentation Verification
- [ ] All documentation renders correctly
- [ ] Links are valid (no 404s)
- [ ] Code examples are syntactically correct
- [ ] API documentation matches implementation

### Final Review
- [ ] PR description is complete and accurate
- [ ] All conversations marked as resolved
- [ ] Minimum required approvals obtained
- [ ] No merge conflicts with target branch
- [ ] Branch is up-to-date with `main`

## 7. Security and Revert Protocol

### Security Incident Response

If a security vulnerability is discovered in a merged PR:

1. **Immediate Actions**
 - Notify all maintainers immediately via secure channel
 - Assess severity using CVSS scoring
 - Do NOT discuss publicly until patch is ready

2. **Critical Vulnerabilities (CVSS ≥ 7.0)**
 - Revert the PR immediately
 - Create private security advisory on GitHub
 - Develop and test fix in private fork
 - Coordinate disclosure timeline

3. **Medium Vulnerabilities (CVSS 4.0-6.9)**
 - Evaluate if revert is necessary
 - Create hotfix branch from affected commit
 - Fast-track fix through expedited review
 - Document in security advisory

4. **Post-Incident**
 - Conduct retrospective on how vulnerability was introduced
 - Update review checklist to prevent recurrence
 - Enhance automated security scanning rules
 - Document lessons learned

### Revert Procedure

**When to Revert:**
- Security vulnerability discovered
- Critical bug causing production issues
- Breaking change not caught in review
- Performance regression exceeding acceptable limits
- Data corruption or integrity issues

**Revert Process:**

1. **Create Revert PR**
 ```bash
 git checkout main
 git pull upstream main
 git revert <commit-sha>
 git push origin revert-<original-pr-number>
 ```

2. **Revert PR Requirements**
 - Title: `revert: <original PR title> (#original-pr-number)`
 - Reference original PR: `Reverts #123`
 - Explain reason for revert clearly
 - Link to issue tracking the fix

3. **Expedited Review**
 - Requires only 1 maintainer approval (vs. usual 2)
 - Fast-tracked through CI (parallel execution)
 - Merge as soon as approval + green CI

4. **Communication**
 - Post in original PR explaining revert
 - Update any related issues
 - Notify affected users/contributors
 - Document in CHANGELOG under "Reverted"

5. **Follow-up**
 - Create issue to re-implement feature correctly
 - Reference reverted PR in new implementation PR
 - Add regression test before re-merging
 - Extra scrutiny during re-implementation review

### Rollback Verification

After revert is merged:
- [ ] Application starts successfully
- [ ] Critical functionality works
- [ ] No data loss or corruption
- [ ] Metrics return to baseline
- [ ] Security issue no longer present
- [ ] All tests pass

---

## Additional Resources

- [Contributing Guidelines](CONTRIBUTING.md)
- [Security Policy](../security/SECURITY.md)
- [API Contract](API_CONTRACT.md)
- [Configuration Guide](CONFIGURATION.md)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [SLSA Framework](https://slsa.dev/)

---

**Questions?** Open an issue or contact maintainers via the project's communication channels.
