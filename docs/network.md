# Network Parameters

This document provides detailed information about the ABM Diora blockchain network parameters and configuration.

## 🌐 Network Overview

### Basic Information

- **Network Name**: ABM Diora Mainnet
- **Chain ID**: 1337
- **Consensus**: Proof of Stake (PoS)
- **Virtual Machine**: Ethereum Virtual Machine (EVM) Compatible
- **Block Time**: 6 seconds
- **Gas Limit**: 30,000,000 gas per block
- **Built by**: ABM Foundation

### Network Status

- **Status**: Active
- **Version**: 1.0.0
- **Last Block**: Updated in real-time
- **Active Validators**: 42
- **Total Staked**: 250,000,000 DIO

## 🔧 Core Parameters

### Block Parameters

| Parameter | Value | Description |
|-----------|--------|-------------|
| Block Time | 6 seconds | Average time between blocks |
| Block Gas Limit | 30,000,000 | Maximum gas per block |
| Block Size | ~2MB | Average block size |
| Uncle Rate | 0% | No uncle blocks in PoS |
| Finality | 1 block | Instant finality |

### Gas Parameters

| Parameter | Value | Description |
|-----------|--------|-------------|
| Base Gas Price | 20 Gwei | Minimum gas price |
| Gas Price | Dynamic | Market-based gas pricing |
| Gas Limit per Tx | 30,000,000 | Maximum gas per transaction |
| Refund Percentage | 50% | Gas refund for storage clearing |

### Consensus Parameters

| Parameter | Value | Description |
|-----------|--------|-------------|
| Consensus Algorithm | Proof of Stake | Delegated Proof of Stake |
| Validator Count | 42 | Active validators |
| Minimum Stake | 10,000 DIO | Minimum to become validator |
| Maximum Stake | 10,000,000 DIO | Maximum stake per validator |
| Unbonding Period | 21 days | Time to unstake |
| Reward Rate | 8.5% APY | Annual percentage yield |

## 💰 Economic Parameters

### Native Token (DIO)

| Parameter | Value | Description |
|-----------|--------|-------------|
| Symbol | DIO | Token symbol |
| Decimals | 18 | Number of decimal places |
| Total Supply | 1,000,000,000 | Total token supply |
| Circulating Supply | ~750,000,000 | Tokens in circulation |
| Burn Rate | 1% | Transaction burn rate |

### Staking Economics

| Parameter | Value | Description |
|-----------|--------|-------------|
| Minimum Stake | 10,000 DIO | Minimum delegation amount |
| Reward Distribution | Every block | Rewards distributed per block |
| Compound Frequency | Every block | Auto-compounding rewards |
| Penalty Rate | 1% | Early unstake penalty |

### Fee Structure

| Parameter | Value | Description |
|-----------|--------|-------------|
| Transaction Fee | Dynamic | Based on gas usage |
| Contract Deployment Fee | 100 DIO | Fixed deployment fee |
| Contract Call Fee | Dynamic | Based on gas usage |
| Transfer Fee | Dynamic | Based on gas usage |

## 🏛️ Governance Parameters

### Protocol Governance

| Parameter | Value | Description |
|-----------|--------|-------------|
| Voting Period | 7 days | Time to vote on proposal |
| Execution Delay | 2 days | Time to execute proposal |
| Quorum | 40% | Minimum participation |
| Approval Threshold | 51% | Minimum approval rate |
| Proposal Deposit | 10,000 DIO | Deposit to create proposal |

### Validator Governance

| Parameter | Value | Description |
|-----------|--------|-------------|
| Validator Selection | Stake-weighted | Based on stake amount |
| Slashing Conditions | Double signing, downtime | Penalties for misbehavior |
| Commission Rate | 0-10% | Validator commission range |
| Uptime Requirement | 95% | Minimum uptime requirement |

## 🔒 Security Parameters

### Cryptographic Parameters

| Parameter | Value | Description |
|-----------|--------|-------------|
| Hash Algorithm | Keccak-256 | Block and transaction hashing |
| Signature Algorithm | ECDSA | Transaction signatures |
| Curve | secp256k1 | Elliptic curve |
| Key Derivation | BIP-32 | Hierarchical deterministic keys |

### Network Security

| Parameter | Value | Description |
|-----------|--------|-------------|
| P2P Protocol | libp2p | Peer-to-peer networking |
| Encryption | TLS 1.3 | Network encryption |
| Rate Limiting | 1000 req/hour | API rate limiting |
| Max Connections | 1000 | Maximum peer connections |

## 🌍 Network Configuration

### Mainnet Configuration

```json
{
  "network": {
    "name": "ABM Diora Mainnet",
    "chain_id": 1337,
    "consensus": "pos",
    "block_time": 6000,
    "gas_limit": 30000000,
    "base_fee": 20000000000
  },
  "genesis": {
    "timestamp": "2026-01-01T00:00:00Z",
    "allocations": {
      "0x...": "1000000000000000000000000000"
    }
  },
  "validators": {
    "min_stake": "10000000000000000000000",
    "max_validators": 42,
    "unbonding_period": 1814400
  }
}
```

### Testnet Configuration

```json
{
  "network": {
    "name": "ABM Diora Testnet",
    "chain_id": 1338,
    "consensus": "pos",
    "block_time": 3000,
    "gas_limit": 30000000,
    "base_fee": 10000000000
  },
  "faucet": {
    "amount": "100000000000000000000000",
    "cooldown": 86400
  }
}
```

## 📊 Performance Metrics

### Throughput

| Metric | Value | Description |
|--------|--------|-------------|
| TPS | 1000+ | Transactions per second |
| Block Finality | 6 seconds | Time to finality |
| Network Latency | <100ms | Average network latency |
| Sync Time | ~1 hour | Full node sync time |

### Resource Requirements

| Component | Minimum | Recommended |
|-----------|---------|------------|
| CPU | 2 cores | 4 cores |
| RAM | 4GB | 8GB |
| Storage | 100GB | 500GB SSD |
| Network | 10 Mbps | 100 Mbps |

## 🔄 Upgrade Parameters

### Protocol Upgrades

| Parameter | Value | Description |
|-----------|--------|-------------|
| Upgrade Schedule | Monthly | Regular protocol upgrades |
| Backward Compatibility | 6 months | Support period |
| Fork Block | Specified | Block number for fork |
| Migration Period | 30 days | Time to migrate |

### Feature Flags

| Feature | Status | Description |
|---------|--------|-------------|
| EIP-1559 | Active | Gas market mechanism |
| EIP-2930 | Active | Access lists |
| Sharding | Planned | Future scalability |
| ZK-Rollups | Planned | Privacy features |

## 🌐 Network Endpoints

### RPC Endpoints

| Network | RPC URL | WebSocket |
|---------|---------|-----------|
| Mainnet | https://rpc.diora.io | wss://ws.diora.io |
| Testnet | https://testnet-rpc.diora.io | wss://testnet-ws.diora.io |

### API Endpoints

| Network | API URL | Documentation |
|---------|---------|-------------|
| Mainnet | https://api.diora.io | https://docs.diora.io |
| Testnet | https://testnet-api.diora.io | https://testnet-docs.diora.io |

### Explorer Endpoints

| Network | Explorer URL |
|---------|--------------|
| Mainnet | https://explorer.diora.io |
| Testnet | https://testnet-explorer.diora.io |

## 📈 Network Monitoring

### Key Metrics

| Metric | Target | Alert Threshold |
|--------|--------|-----------------|
| Block Time | 6s | >10s |
| TPS | 1000+ | <500 |
| Gas Price | 20 Gwei | >100 Gwei |
| Validator Uptime | 95% | <90% |

### Monitoring Tools

- **Prometheus**: Metrics collection
- **Grafana**: Visualization dashboard
- **Alertmanager**: Alert management
- **Jaeger**: Distributed tracing

## 🔧 Development Parameters

### Development Network

| Parameter | Value |
|-----------|--------|
| Chain ID | 1337 |
| Block Time | 2 seconds |
| Gas Price | 1 Gwei |
| Pre-funded Accounts | 10 |
| Faucet Amount | 1000 DIO |

### Testing Parameters

| Parameter | Value |
|-----------|--------|
| Test Timeout | 30 seconds |
| Gas Limit | 10,000,000 |
| Block Confirmation | 1 |
| Mining Interval | 2 seconds |

## 📝 Configuration Files

### Node Configuration

```toml
[node]
data_dir = "/var/lib/diora"
log_level = "info"
metrics_enabled = true

[network]
listen_addr = "0.0.0.0:30303"
bootnodes = [
  "/ip4/127.0.0.1/tcp/30303/p2p/12D3KooW..."
]

[consensus]
validator = false
private_key = ""
reward_address = ""

[jsonrpc]
enable = true
listen_addr = "127.0.0.1:8545"
cors_domains = ["*"]
```

### Wallet Configuration

```toml
[wallet]
keystore_dir = "/var/lib/diora/keystore"
unlock_timeout = 300
auto_unlock = false

[account]
default_account = ""
passphrase = ""
```

## 🚨 Emergency Parameters

### Circuit Breakers

| Parameter | Trigger | Action |
|-----------|---------|--------|
| High Gas Price | >1000 Gwei | Reject transactions |
| Network Congestion | >90% utilization | Increase gas price |
| Validator Downtime | >10% offline | Slash validators |
| Chain Reorg | >3 blocks | Emergency halt |

### Recovery Procedures

| Issue | Recovery Time | Procedure |
|-------|---------------|----------|
| Network Partition | 5 minutes | Auto-reconnect |
| Validator Slashing | Immediate | Manual intervention |
| Chain Reorg | 1 hour | Manual rollback |
| Data Corruption | 24 hours | Restore from backup |

---

For more information, join our [Telegram](https://t.me/DioraFund) or follow us on [Twitter](https://twitter.com/DioraCrypto).
