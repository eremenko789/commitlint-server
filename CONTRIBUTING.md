# Contributing to Commitlint

Thank you for your interest in contributing to Commitlint! This guide will help you get started with contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Commit Message Guidelines](#commit-message-guidelines)

## Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/your-username/commitlint.git
   cd commitlint
   ```
3. Add the upstream repository as a remote:
   ```bash
   git remote add upstream https://github.com/conventionalcommit/commitlint.git
   ```

## Development Setup

### Prerequisites

- Go 1.23.0 or later
- Make (optional, but recommended)
- Docker (optional, for testing Docker builds)

### Building the Project

```bash
# Build both CLI and webhook server
make build

# Build only the CLI
make build-cli

# Build only the webhook server
make build-webhook-server

# Run tests
go test ./...

# Run with race detector
go test -race ./...
```

## Project Structure

```
.
├── cmd/
│   ├── cli/                    # CLI application
│   └── webhook-server/         # Webhook server application
├── internal/
│   ├── webhook/               # Webhook server implementation
│   │   ├── config.go         # Configuration management
│   │   ├── server.go         # HTTP server
│   │   ├── gitea.go          # Gitea API client
│   │   ├── lint.go           # Commitlint integration
│   │   └── payload.go        # Webhook payload types
│   └── ...                    # Other internal packages
├── .github/
│   └── workflows/             # GitHub Actions workflows
├── .goreleaser.yml           # GoReleaser configuration
├── Makefile                  # Build automation
└── ...
```

## Making Changes

1. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes, following the coding standards:
   - Use `gofmt` to format your code
   - Follow Go best practices and idioms
   - Add tests for new functionality
   - Update documentation as needed

3. Run tests to ensure everything works:
   ```bash
   go test ./...
   ```

4. Validate GoReleaser configuration if you've made changes to it:
   ```bash
   goreleaser check
   ```

## Testing

### Unit Tests

Write unit tests for new functionality. Place test files in the same package as the code they test, with the suffix `_test.go`.

### Integration Tests

For integration tests that require external services (like Gitea), use build tags to exclude them from regular test runs:

```go
//go:build integration
// +build integration

package webhook_test
```

Run integration tests with:
```bash
go test -tags=integration ./...
```

### Manual Testing

1. **CLI Testing:**
   ```bash
   go run cmd/cli/main.go --help
   echo "feat: add new feature" | go run cmd/cli/main.go
   ```

2. **Webhook Server Testing:**
   ```bash
   # Start the server
   go run cmd/webhook-server/main.go -c webhook-server.yml
   
   # In another terminal, send a test webhook
   curl -X POST http://localhost:3000/webhook \
     -H "Content-Type: application/json" \
     -H "X-Gitea-Signature: your-signature" \
     -d @test-payload.json
   ```

## Submitting Changes

1. Commit your changes using conventional commit format (see below)
2. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
3. Open a Pull Request against the `main` branch
4. Ensure all CI checks pass
5. Wait for review and address any feedback

## Commit Message Guidelines

This project uses [Conventional Commits](https://www.conventionalcommits.org/). Each commit message should have the format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to CI configuration files and scripts
- `chore`: Other changes that don't modify src or test files

### Examples

```
feat(webhook): add support for GitHub webhooks

Add a new webhook handler that can process GitHub webhook events
in addition to Gitea events.

Closes #123
```

```
fix(cli): correct version flag output

The version flag was showing incorrect build information.
This commit fixes the version string construction.
```

```
docs: update installation instructions

Add instructions for installing via Homebrew and update
the Docker installation section.
```

## Release Process

Releases are automated using GitHub Actions and GoReleaser. When a tag is pushed:

1. The `release.yml` workflow triggers
2. GoReleaser builds binaries for all supported platforms
3. Docker images are built and pushed to GitHub Container Registry
4. A GitHub Release is created with all artifacts

To create a release:

```bash
git tag -a v1.2.3 -m "Release version 1.2.3"
git push upstream v1.2.3
```

## Questions?

If you have questions, feel free to:

1. Open an issue for discussion
2. Ask in existing issues or pull requests
3. Reach out to the maintainers

Thank you for contributing!
