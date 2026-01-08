# 🚀 ABM Diora Go Blockchain

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/DioraFund/diora-go-blockchain)

A high-performance, EVM-compatible Layer 1 blockchain written in Go, built by ABM Foundation. Features Proof of Stake consensus, smart contract support, and enterprise-grade security.

## ✨ Features

### 🏗️ Core Blockchain
- **EVM Compatible**: Full Ethereum Virtual Machine compatibility
- **Proof of Stake**: Energy-efficient DPoS consensus mechanism
- **6-Second Block Times**: Fast finality and high throughput
- **1000+ TPS**: Scalable transaction processing
- **Low Gas Fees**: Sub-cent transaction costs

### 🔐 Security
- **BLS Signatures**: Advanced cryptographic primitives
- **Slashing Mechanism**: Validator accountability
- **Rate Limiting**: DDoS protection
- **Audit Logging**: Comprehensive security monitoring

### 🌐 Network
- **P2P Networking**: libp2p-based peer discovery
- **WebSocket API**: Real-time event streaming
- **REST API**: Full JSON-RPC compatibility
- **CLI Tools**: Command-line management interface

### 💎 Smart Contracts
- **Solidity Support**: Full Solidity compiler compatibility
- **ERC-20/721/1155**: Standard token contracts
- **Gas Optimization**: Efficient execution environment
- **Contract Verification**: On-chain source code verification

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional)

### Installation

```bash
# Clone the repository
git clone https://github.com/DioraFund/diora-go-blockchain.git
cd diora-go-blockchain

# Build the blockchain
make build

# Start a development node
./build/diora start --dev
```

### Docker Installation

```bash
# Build Docker image
docker build -t diora-blockchain .

# Run with Docker Compose
docker-compose up -d
```

## 📋 Commands

### Node Management

```bash
# Start node
./build/diora start

# Start with custom config
./build/diora start --config config.json

# Stop node
./build/diora stop

# Check status
./build/diora status
```

### Account Management

```bash
# Create new account
./build/diora account create

# Import account
./build/diora account import <private-key>

# Check balance
./build/diora account balance <address>

# Send transaction
./build/diora account send <from> <to> <amount>
```

### Smart Contracts

```bash
# Compile contract
./build/diora contract compile Contract.sol

# Deploy contract
./build/diora contract deploy Contract.sol

# Call contract
./build/diora contract call <address> <method> <args>
```

### Staking

```bash
# Stake tokens
./build/diora stake <amount> <validator>

# Unstake tokens
./build/diora unstake <amount> <validator>

# Check rewards
./build/diora stake rewards <validator>
```

## 🏗️ Architecture

### Core Components

```
diora-go-blockchain/
├── core/                 # Blockchain core
│   ├── blockchain.go     # Main blockchain logic
│   ├── state.go          # State management
│   └── types.go          # Core types
├── consensus/            # Proof of Stake
│   └── pos.go           # DPoS implementation
├── vm/                   # EVM implementation
│   └── evm.go           # Virtual Machine
├── crypto/               # Cryptography
│   └── crypto.go        # BLS, ECDSA, hashing
├── api/                  # HTTP/WebSocket API
│   └── server.go        # REST server
├── cli/                  # Command line interface
│   └── main.go          # CLI commands
├── contracts/            # Smart contracts
├── tests/                # Test suite
└── security/             # Security modules
```

### Consensus Algorithm

The blockchain uses a Delegated Proof of Stake (DPoS) consensus mechanism:

- **Validator Set**: 42 active validators
- **Block Production**: 6-second intervals
- **Finality**: Instant block finality
- **Rewards**: 8.5% APY for stakers
- **Slashing**: Penalty for misbehavior

### Virtual Machine

Full EVM compatibility with optimizations:

- **Gas Metering**: Accurate gas calculation
- **Precompiles**: Ethereum precompiled contracts
- **Opcode Support**: All standard EVM opcodes
- **State Trie**: Efficient state storage

## 🔧 Configuration

### Network Configuration

```json
{
  "network": {
    "chain_id": 1337,
    "block_time": 6000,
    "gas_limit": 30000000,
    "min_stake": "10000000000000000000000"
  },
  "consensus": {
    "validators": 42,
    "unbonding_period": 1814400,
    "reward_rate": "850000000000000000"
  },
  "api": {
    "rpc_port": 8545,
    "ws_port": 8546,
    "enable_cors": true
  }
}
```

### Environment Variables

```bash
# Node configuration
export DIORA_DATA_DIR=/var/lib/diora
export DIORA_LOG_LEVEL=info
export DIORA_NETWORK_ID=1337

# API configuration
export DIORA_RPC_PORT=8545
export DIORA_WS_PORT=8546
export DIORA_ENABLE_CORS=true
```

## 📊 API Documentation

### JSON-RPC Endpoints

```bash
# Get latest block
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Get balance
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x...","latest"],"id":1}'

# Send transaction
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{...}],"id":1}'
```

### WebSocket Events

```javascript
const ws = new WebSocket('ws://localhost:8546');

// Subscribe to new blocks
ws.send(JSON.stringify({
  "id": 1,
  "method": "eth_subscribe",
  "params": ["newHeads"]
}));

// Listen for events
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('New block:', data);
};
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test ./core/blockchain_test.go

# Run with coverage
make test-coverage

# Run integration tests
make test-integration
```

### Benchmark Tests

```bash
# Run benchmarks
make benchmark

# Test TPS
make benchmark-tps

# Test memory usage
make benchmark-memory
```

## 📈 Performance

### Benchmarks

| Metric | Value | Description |
|--------|-------|-------------|
| TPS | 1000+ | Transactions per second |
| Block Time | 6s | Average block interval |
| Finality | 6s | Block finality time |
| Gas Limit | 30M | Maximum gas per block |
| Sync Time | ~1 hour | Full node sync |

### Resource Requirements

| Component | Minimum | Recommended |
|-----------|---------|------------|
| CPU | 2 cores | 4 cores |
| RAM | 4GB | 8GB |
| Storage | 100GB | 500GB SSD |
| Network | 10 Mbps | 100 Mbps |

## 🔒 Security

### Security Features

- **BLS Signatures**: Threshold signature schemes
- **Rate Limiting**: Request throttling
- **Input Validation**: Comprehensive input checking
- **Audit Logging**: Security event tracking
- **Access Control**: Role-based permissions

### Security Audits

- ✅ **Code Review**: Peer-reviewed codebase
- ✅ **Penetration Testing**: Regular security audits
- ✅ **Bug Bounty**: Responsible disclosure program
- ✅ **Smart Contract Audits**: Third-party contract verification

## 🌐 Network Parameters

### Mainnet

| Parameter | Value |
|------------|-------|
| Chain ID | 1337 |
| Block Time | 6 seconds |
| Gas Price | Dynamic |
| Validators | 42 |
| Min Stake | 10,000 DIO |
| Unbonding | 21 days |

### Testnet

| Parameter | Value |
|------------|-------|
| Chain ID | 1338 |
| Block Time | 3 seconds |
| Gas Price | 1 Gwei |
| Faucet Amount | 1000 DIO |
| Faucet Cooldown | 24 hours |

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create feature branch
3. Make changes
4. Add tests
5. Submit pull request

### Code Standards

- Follow Go best practices
- Use `gofmt` for formatting
- Add comprehensive tests
- Document all public APIs

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

**Copyright (c) 2026 ABM Foundation**

## 🌍 Community

- **GitHub**: [DioraFund/diora-go-blockchain](https://github.com/DioraFund/diora-go-blockchain)
- **Telegram**: [@DioraFund](https://t.me/DioraFund)
- **Twitter**: [@DioraCrypto](https://twitter.com/DioraCrypto)

## 📚 Documentation

- [Getting Started](./docs/getting-started.md)
- [Smart Contracts](./docs/contracts.md)
- [Network Parameters](./docs/network.md)
- [Security Guide](./docs/security.md)
- [📄 Whitepaper](https://diorafund.github.io/diora-whitepaper) - Complete technical and economic documentation

## 🏆 Acknowledgments

- **Ethereum Foundation** for EVM specification
- **Go-Ethereum** for reference implementation
- **libp2p** for P2P networking
- **OpenZeppelin** for smart contract standards

---

**Built with ❤️ by ABM Foundation for the future of decentralized finance**
