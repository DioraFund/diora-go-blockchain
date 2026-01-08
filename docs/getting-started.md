# Getting Started Guide

Welcome to ABM Diora Go Blockchain! This guide will help you get started with building and running your own blockchain node.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional)
- 4GB+ RAM
- 100GB+ storage

### Installation

#### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/DioraFund/diora-go-blockchain.git
cd diora-go-blockchain

# Build the blockchain
make build

# Start a development node
./build/diora start --dev
```

#### Option 2: Docker Installation

```bash
# Pull the image
docker pull diorafund/diora-go-blockchain

# Run with Docker
docker run -p 8545:8545 -p 8546:8546 diorafund/diora-go-blockchain

# Or use docker-compose
docker-compose up -d
```

### Verification

```bash
# Check if node is running
curl http://localhost:8545

# Check node status
./build/diora status

# Get latest block
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

## 📋 Basic Concepts

### Blockchain Architecture

ABM Diora is an EVM-compatible blockchain with the following components:

- **Core Layer**: Block processing, state management
- **Consensus Layer**: Proof of Stake with DPoS delegation
- **Virtual Machine**: Ethereum Virtual Machine compatibility
- **API Layer**: JSON-RPC and WebSocket endpoints
- **P2P Layer**: Peer-to-peer networking

### Network Parameters

| Parameter | Value | Description |
|-----------|--------|-------------|
| Chain ID | 1337 | Unique network identifier |
| Block Time | 6 seconds | Average block interval |
| Gas Limit | 30,000,000 | Maximum gas per block |
| Validators | 42 | Active validator count |
| Min Stake | 10,000 DIO | Minimum delegation amount |

## 🔧 Node Configuration

### Configuration File

Create a `config.json` file:

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
    "enable_cors": true,
    "rpc_addr": "0.0.0.0"
  },
  "p2p": {
    "listen_addr": "0.0.0.0:30303",
    "bootnodes": [
      "/ip4/127.0.0.1/tcp/30303/p2p/12D3KooWExample"
    ]
  },
  "database": {
    "path": "./data",
    "cache_size": 1000
  },
  "logging": {
    "level": "info",
    "file": "./logs/diora.log"
  }
}
```

### Environment Variables

```bash
# Node configuration
export DIORA_CONFIG_PATH=/path/to/config.json
export DIORA_DATA_DIR=/var/lib/diora
export DIORA_LOG_LEVEL=info

# API configuration
export DIORA_RPC_PORT=8545
export DIORA_WS_PORT=8546
export DIORA_ENABLE_CORS=true

# Network configuration
export DIORA_NETWORK_ID=1337
export DIORA_BOOTNODES="/ip4/127.0.0.1/tcp/30303/p2p/12D3KooWExample"
```

## 👛 Account Management

### Creating Accounts

```bash
# Create new account
./build/diora account create

# Create with password
./build/diora account create --password

# Create with custom name
./build/diora account create --name "My Validator"
```

### Importing Accounts

```bash
# Import from private key
./build/diora account import --private-key "0x..."

# Import from keystore
./build/diora account import --keystore /path/to/keystore.json

# List all accounts
./build/diora account list
```

### Managing Accounts

```bash
# Check balance
./build/diora account balance <address>

# Get account info
./build/diora account info <address>

# Export private key
./build/diora account export <address> --private-key

# Export keystore
./build/diora account export <address> --keystore /path/to/export.json
```

## 🔍 Node Operations

### Starting a Node

```bash
# Start with default settings
./build/diora start

# Start with custom config
./build/diora start --config /path/to/config.json

# Start in development mode
./build/diora start --dev

# Start as validator
./build/diora start --validator --private-key "0x..."
```

### Stopping a Node

```bash
# Graceful shutdown
./build/diora stop

# Force shutdown
./build/diora stop --force
```

### Node Status

```bash
# Check node status
./build/diora status

# Detailed status
./build/diora status --detailed

# Check network connectivity
./build/diora status --network

# Check validator status
./build/diora status --validator
```

## 💸 Transactions

### Sending Transactions

```bash
# Send DIO tokens
./build/diora account send <from> <to> <amount>

# Send with custom gas price
./build/diora account send <from> <to> <amount> --gas-price 20000000000

# Send with password
./build/diora account send <from> <to> <amount> --password

# Estimate gas
./build/diora account estimate <from> <to> <amount>
```

### Transaction History

```bash
# Get transaction history
./build/diora account history <address>

# Get pending transactions
./build/diora tx pending

# Get transaction details
./build/diora tx info <tx-hash>
```

## 🏛️ Staking and Validators

### Becoming a Validator

```bash
# Register as validator
./build/diora validator register --private-key "0x..." --commission 5

# Stake tokens
./build/diora stake <amount> <validator-address>

# Unstake tokens
./build/diora unstake <amount> <validator-address>

# Check rewards
./build/diora stake rewards <validator-address>
```

### Validator Management

```bash
# List validators
./build/diora validator list

# Get validator info
./build/diora validator info <address>

# Update commission
./build/diora validator update-commission <address> <new-commission>

# Withdraw rewards
./build/diora validator withdraw-rewards <address>
```

## 📜 Smart Contracts

### Deploying Contracts

```bash
# Compile Solidity contract
./build/diora contract compile MyContract.sol

# Deploy contract
./build/diora contract deploy MyContract.sol

# Deploy with constructor arguments
./build/diora contract deploy MyContract.sol --args "arg1" "arg2"

# Deploy with custom gas
./build/diora contract deploy MyContract.sol --gas 2000000
```

### Interacting with Contracts

```bash
# Call contract method
./build/diora contract call <address> <method> <args>

# Send transaction to contract
./build/diora contract send <address> <method> <args> --value 1000000000000000000

# Get contract ABI
./build/diora contract abi <address>

# Get contract source
./build/diora contract source <address>
```

## 🔌 API Usage

### JSON-RPC API

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

### WebSocket API

```javascript
const ws = new WebSocket('ws://localhost:8546');

// Subscribe to new blocks
ws.send(JSON.stringify({
  "id": 1,
  "method": "eth_subscribe",
  "params": ["newHeads"]
}));

// Subscribe to pending transactions
ws.send(JSON.stringify({
  "id": 2,
  "method": "eth_subscribe",
  "params": ["newPendingTransactions"]
}));

// Listen for events
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
};
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test
go test ./core/blockchain_test.go

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmarks
make benchmark
```

### Test Examples

```bash
# Test blockchain functionality
go test -v ./core/...

# Test consensus mechanism
go test -v ./consensus/...

# Test API endpoints
go test -v ./api/...

# Test CLI commands
go test -v ./cli/...
```

## 🐛 Troubleshooting

### Common Issues

**Node won't start**
```bash
# Check logs
./build/diora logs

# Check configuration
./build/diora config validate

# Check ports
netstat -tulpn | grep :8545
```

**Transactions not confirming**
```bash
# Check gas price
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}'

# Check network status
./build/diora status --network
```

**API not responding**
```bash
# Check if API is running
curl http://localhost:8545

# Check API logs
tail -f logs/diora.log | grep API

# Restart API service
./build/diora restart --api
```

### Getting Help

```bash
# Show help
./build/diora --help

# Show command help
./build/diora account --help

# Show version
./build/diora version
```

## 📚 Next Steps

1. **Explore the API**: Read the [API Reference](./api.md)
2. **Deploy Smart Contracts**: See [Smart Contracts Guide](./contracts.md)
3. **Run a Validator**: Check [Network Parameters](./network.md)
4. **Secure Your Node**: Read [Security Guide](./security.md)
5. **Join the Community**: Connect with other developers

## 🔗 Additional Resources

- [API Reference](./api.md)
- [Smart Contracts](./contracts.md)
- [Network Parameters](./network.md)
- [Security Guide](./security.md)
- [Contributing Guide](../CONTRIBUTING.md)

---

For more information, join our [Telegram](https://t.me/DioraFund) or follow us on [Twitter](https://twitter.com/DioraCrypto).

**Built with ❤️ by ABM Foundation**
