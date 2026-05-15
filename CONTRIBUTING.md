# Contributing to Obelisk CLI

First off, thank you for considering contributing to Obelisk CLI! 🎉

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Style Guide](#style-guide)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/Obelisk-CLI.git
   cd Obelisk-CLI
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/Swif7ify/Obelisk-CLI.git
   ```

## Development Setup

### Prerequisites

- **Go 1.26+** installed ([download](https://golang.org/dl/))
- **Git** installed
- A **Google Gemini API key** (for AI features)

### Build & Run

```bash
# Install dependencies
go mod tidy

# Build the binary
make build

# Run tests
make test

# Run the tool
./bin/obelisk check
```

## How to Contribute

### Reporting Bugs

- Use the [Bug Report](https://github.com/Swif7ify/Obelisk-CLI/issues/new?template=bug_report.md) issue template
- Include steps to reproduce, expected behavior, and actual behavior
- Include your Go version and OS

### Suggesting Features

- Use the [Feature Request](https://github.com/Swif7ify/Obelisk-CLI/issues/new?template=feature_request.md) issue template
- Describe the problem your feature would solve
- Suggest a possible implementation

### Adding a New Framework Adapter

Obelisk uses an adapter pattern. To add support for a new framework:

1. Create a new file in `adapters/` (e.g., `python.go`)
2. Implement the `Adapter` interface
3. Register it in `internal/detector/detector.go`
4. Add tests

## Pull Request Process

1. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/my-feature
   ```
2. Make your changes with clear, descriptive commits
3. Ensure all tests pass: `make test`
4. Ensure the build succeeds: `make build`
5. Update documentation if needed
6. Submit a Pull Request using our [PR template](.github/PULL_REQUEST_TEMPLATE.md)

### PR Requirements

- [ ] Code compiles without errors
- [ ] All existing tests pass
- [ ] New features include tests
- [ ] Documentation is updated
- [ ] Commit messages are clear and descriptive

## Style Guide

### Go Code

- Follow standard [Go conventions](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Keep functions focused and small
- Add comments for exported functions and types
- Handle errors explicitly — no silent swallows

### Commit Messages

Use conventional commit format:
```
type(scope): description

feat(scanner): add entropy-based secret detection
fix(ui): resolve spinner alignment on narrow terminals
docs(readme): update installation instructions
```

---

Thank you for helping make Obelisk better! 🏛️
