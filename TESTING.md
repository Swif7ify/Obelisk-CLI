# 🧪 Testing Guide for Obelisk CLI

Comprehensive testing documentation for developers and contributors.

## 📋 Table of Contents

- [Overview](#overview)
- [Running Tests](#running-tests)
- [Test Structure](#test-structure)
- [Writing Tests](#writing-tests)
- [Coverage Reports](#coverage-reports)
- [CI/CD Integration](#cicd-integration)
- [Best Practices](#best-practices)

---

## 🎯 Overview

Obelisk CLI uses Go's built-in testing framework with comprehensive unit tests, integration tests, and CI/CD automation.

### Test Coverage

- **Unit Tests**: Test individual functions and components
- **Integration Tests**: Test complete workflows
- **Security Scans**: Automated security vulnerability detection
- **Cross-Platform**: Tests run on Windows, macOS, and Linux

---

## 🚀 Running Tests

### Quick Start

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Run benchmarks
make bench
```

### Manual Test Commands

```bash
# Run all tests
go test ./...

# Run tests in specific package
go test ./internal/scanner/...

# Run specific test
go test -run TestScanSecrets ./internal/scanner/

# Run with verbose output
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...
```

---

## 📁 Test Structure

```
Obelisk-CLI/
├── cmd/
│   └── update_test.go              # Command tests
├── internal/
│   ├── scanner/
│   │   ├── secrets_test.go         # Secret scanning tests
│   │   ├── naming_test.go          # Naming convention tests
│   │   └── gitignore_test.go       # Gitignore tests
│   └── report/
│       └── formatter_test.go       # Report formatting tests
└── .github/
    └── workflows/
        └── test.yml                # CI/CD test workflow
```

---

## ✍️ Writing Tests

### Unit Test Template

```go
package scanner

import (
    "testing"
)

func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid input",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
        {
            name:    "invalid input",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Testing with Temporary Files

```go
func TestWithFiles(t *testing.T) {
    // Create temporary directory
    tmpDir := t.TempDir()

    // Create test file
    testFile := filepath.Join(tmpDir, "test.txt")
    if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
        t.Fatal(err)
    }

    // Run your test
    result, err := YourFunction(tmpDir)

    // Assert results
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Cleanup is automatic with t.TempDir()
}
```

### Table-Driven Tests

```go
func TestMultipleScenarios(t *testing.T) {
    tests := []struct {
        name         string
        input        string
        wantFindings int
        wantSeverity string
    }{
        {
            name:         "scenario 1",
            input:        "test1",
            wantFindings: 1,
            wantSeverity: "critical",
        },
        {
            name:         "scenario 2",
            input:        "test2",
            wantFindings: 0,
            wantSeverity: "",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## 📊 Coverage Reports

### Generate Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html  # macOS
xdg-open coverage.html  # Linux
start coverage.html  # Windows
```

### Using Makefile

```bash
# Generate and open coverage report
make test-coverage
```

### Coverage Goals

- **Overall**: Aim for 80%+ coverage
- **Critical Paths**: 90%+ coverage for security-related code
- **New Code**: All new features should include tests

---

## 🔄 CI/CD Integration

### GitHub Actions Workflow

Tests run automatically on:

- Push to `main` or `develop` branches
- Pull requests
- Multiple OS (Ubuntu, Windows, macOS)
- Multiple Go versions (1.21, 1.22)

### Workflow Steps

1. **Unit Tests**: Run all unit tests
2. **Linting**: Check code quality with golangci-lint
3. **Build**: Build for all platforms
4. **Integration Tests**: Test built binaries
5. **Security Scan**: Run Gosec security scanner
6. **Coverage Upload**: Upload to Codecov

### View Test Results

```bash
# Check workflow status
gh run list

# View specific run
gh run view <run-id>

# View logs
gh run view <run-id> --log
```

---

## 🎯 Best Practices

### 1. Test Naming

```go
// Good
func TestScanSecrets_DetectsAPIKeys(t *testing.T) {}
func TestScanSecrets_IgnoresSafePatterns(t *testing.T) {}

// Bad
func TestSecrets(t *testing.T) {}
func Test1(t *testing.T) {}
```

### 2. Test Independence

```go
// Good - Each test is independent
func TestFeatureA(t *testing.T) {
    tmpDir := t.TempDir()
    // Test in isolation
}

func TestFeatureB(t *testing.T) {
    tmpDir := t.TempDir()
    // Test in isolation
}

// Bad - Tests depend on each other
var sharedState string

func TestFeatureA(t *testing.T) {
    sharedState = "value"
}

func TestFeatureB(t *testing.T) {
    // Depends on TestFeatureA running first
    if sharedState != "value" {
        t.Fatal("wrong state")
    }
}
```

### 3. Clear Assertions

```go
// Good
if got != want {
    t.Errorf("FunctionName() = %v, want %v", got, want)
}

// Bad
if got != want {
    t.Error("failed")
}
```

### 4. Use Subtests

```go
// Good
func TestFeature(t *testing.T) {
    t.Run("scenario 1", func(t *testing.T) {
        // Test scenario 1
    })

    t.Run("scenario 2", func(t *testing.T) {
        // Test scenario 2
    })
}

// Bad
func TestFeatureScenario1(t *testing.T) {}
func TestFeatureScenario2(t *testing.T) {}
```

### 5. Test Error Cases

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", "test", false},
        {"empty input", "", true},
        {"invalid input", "bad", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 🐛 Debugging Tests

### Run Specific Test

```bash
# Run single test
go test -run TestScanSecrets ./internal/scanner/

# Run with verbose output
go test -v -run TestScanSecrets ./internal/scanner/

# Run with race detector
go test -race -run TestScanSecrets ./internal/scanner/
```

### Print Debug Information

```go
func TestDebug(t *testing.T) {
    result := SomeFunction()

    // Print for debugging
    t.Logf("Result: %+v", result)

    // Only prints if test fails
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Skip Tests

```go
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode")
    }

    // Long-running test
}

// Run with: go test -short
```

---

## 📈 Benchmarking

### Write Benchmarks

```go
func BenchmarkFunction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Function("input")
    }
}
```

### Run Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkFunction ./internal/scanner/

# With memory stats
go test -bench=. -benchmem ./...

# Compare benchmarks
go test -bench=. -benchmem ./... > old.txt
# Make changes
go test -bench=. -benchmem ./... > new.txt
benchcmp old.txt new.txt
```

---

## 🔒 Security Testing

### Run Security Scan

```bash
# Install gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run scan
gosec ./...

# Generate report
gosec -fmt=json -out=results.json ./...
```

### Security Tests in CI

Security scans run automatically in GitHub Actions on every push and PR.

---

## 📚 Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Test Comments](https://github.com/golang/go/wiki/TestComments)
- [Testify Library](https://github.com/stretchr/testify) (optional)

---

## 🤝 Contributing Tests

When contributing:

1. **Write tests** for all new features
2. **Update existing tests** when modifying code
3. **Ensure tests pass** locally before pushing
4. **Check coverage** doesn't decrease
5. **Follow naming conventions**
6. **Document complex tests**

### Pull Request Checklist

- [ ] All tests pass locally
- [ ] New tests added for new features
- [ ] Coverage maintained or improved
- [ ] Tests are independent
- [ ] Clear test names and assertions
- [ ] CI/CD tests pass

---

**Happy Testing! 🧪**
