# Contributing to ABM Diora Go Blockchain

Thank you for your interest in contributing to the ABM Diora Go Blockchain! This document provides guidelines and information for contributors.

## 🚀 How to Contribute

### Reporting Issues

- Use [GitHub Issues](https://github.com/DioraFund/diora-go-blockchain/issues) to report bugs
- Provide detailed information about the issue
- Include steps to reproduce
- Add relevant logs and screenshots

### Submitting Pull Requests

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/diora-go-blockchain.git
   cd diora-go-blockchain
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **Make your changes**
   - Follow the coding standards
   - Add tests for new functionality
   - Update documentation

4. **Commit your changes**
   ```bash
   git commit -m "Add amazing feature"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/amazing-feature
   ```

6. **Create a Pull Request**
   - Provide a clear description
   - Reference relevant issues
   - Include test results

## 📝 Coding Standards

### Go Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format code
- Use `golint` to check for style issues
- Add comments for public functions
- Write unit tests

### Example Code Structure

```go
package blockchain

import (
    "context"
    "fmt"
)

// Block represents a block in the blockchain
type Block struct {
    Index     int64
    Timestamp int64
    Data      string
    Hash      string
    PrevHash  string
}

// NewBlock creates a new block with the given data
func NewBlock(index int64, prevHash, data string) *Block {
    block := &Block{
        Index:     index,
        Timestamp: time.Now().Unix(),
        Data:      data,
        PrevHash:  prevHash,
    }
    block.Hash = block.CalculateHash()
    return block
}

// CalculateHash computes the hash of the block
func (b *Block) CalculateHash() string {
    // Implementation
    return ""
}
```

### Testing

```go
package blockchain

import (
    "testing"
)

func TestNewBlock(t *testing.T) {
    block := NewBlock(1, "prevhash", "test data")
    
    if block.Index != 1 {
        t.Errorf("Expected block index 1, got %d", block.Index)
    }
    
    if block.PrevHash != "prevhash" {
        t.Errorf("Expected prev hash %s, got %s", "prevhash", block.PrevHash)
    }
}
```

## 🧪 Development Setup

### Prerequisites

- Go 1.21+
- Docker (optional)
- Git

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/DioraFund/diora-go-blockchain.git
   cd diora-go-blockchain
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

4. **Build the project**
   ```bash
   make build
   ```

5. **Start development node**
   ```bash
   ./build/diora start --dev
   ```

### Development Tools

```bash
# Format code
make fmt

# Run linter
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark
```

## 📋 Project Structure

```
diora-go-blockchain/
├── core/                 # Core blockchain logic
│   ├── blockchain.go     # Main blockchain implementation
│   ├── state.go          # State management
│   └── types.go          # Core types and interfaces
├── consensus/            # Consensus algorithms
│   └── pos.go           # Proof of Stake implementation
├── vm/                   # Virtual Machine
│   └── evm.go           # EVM implementation
├── crypto/               # Cryptographic primitives
│   └── crypto.go        # Hashing, signatures, keys
├── api/                  # API endpoints
│   └── server.go        # HTTP/WebSocket server
├── cli/                  # Command line interface
│   └── main.go          # CLI commands
├── tests/                # Test suite
│   ├── unit/            # Unit tests
│   ├── integration/     # Integration tests
│   └── benchmark/       # Performance tests
└── docs/                 # Documentation
```

## 🔧 Contribution Guidelines

### Before Contributing

1. **Check existing issues** - Look for similar issues or PRs
2. **Discuss changes** - Open an issue for major changes
3. **Follow coding standards** - Ensure code follows project conventions
4. **Add tests** - Include tests for new functionality
5. **Update docs** - Update relevant documentation

### Pull Request Process

1. **Create descriptive title** - Summarize changes clearly
2. **Provide detailed description** - Explain what and why
3. **Link to issues** - Reference related issues
4. **Include screenshots** - For UI changes
5. **Test thoroughly** - Ensure all tests pass

### Code Review Process

- **Maintainer review** - All PRs require maintainer approval
- **Automated checks** - CI/CD pipeline runs automatically
- **Community review** - Community members can review and comment
- **Security review** - Security-sensitive changes require extra review

## 🧪 Testing Guidelines

### Unit Tests

- Test individual functions and methods
- Use table-driven tests for multiple cases
- Mock external dependencies
- Aim for >80% code coverage

### Integration Tests

- Test component interactions
- Test API endpoints
- Test blockchain operations
- Test consensus mechanisms

### Benchmark Tests

- Performance-critical code needs benchmarks
- Test memory usage
- Test CPU usage
- Test network performance

### Test Examples

```go
// Table-driven test
func TestValidateAddress(t *testing.T) {
    tests := []struct {
        name    string
        address string
        want    bool
    }{
        {"valid address", "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb", true},
        {"invalid address", "0xinvalid", false},
        {"empty address", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ValidateAddress(tt.address)
            if got != tt.want {
                t.Errorf("ValidateAddress() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 🐛 Bug Reports

### Bug Report Template

```markdown
## Bug Description
Brief description of the bug

## Steps to Reproduce
1. Go to...
2. Click on...
3. See error

## Expected Behavior
What you expected to happen

## Actual Behavior
What actually happened

## Environment
- OS: [e.g. macOS, Windows, Linux]
- Go version: [e.g. 1.21.0]
- Diora version: [e.g. v1.0.0]

## Additional Context
Add any other context about the problem here

## Logs
[Include relevant logs]
```

## 💡 Feature Requests

### Feature Request Template

```markdown
## Feature Description
Brief description of the feature

## Problem Statement
What problem does this feature solve?

## Proposed Solution
How should this feature work

## Alternatives Considered
What other approaches did you consider

## Additional Context
Add any other context about the feature here
```

## 🏆 Recognition

Contributors will be recognized in:

- README.md contributors section
- Release notes
- Project website
- Community announcements

## 📞 Getting Help

- Create an issue for questions
- Join our [Telegram](https://t.me/DioraFund)
- Follow us on [Twitter](https://twitter.com/DioraCrypto)
- Check [Discord](#) (coming soon)

## 🔒 Security

- Report security vulnerabilities privately
- Follow responsible disclosure
- Don't commit sensitive data
- Use secure coding practices

## 📄 License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You

Thank you for contributing to ABM Diora Go Blockchain! Your contributions help make the project better for everyone.

---

**Remember:** Be respectful, be helpful, and follow the code of conduct.
