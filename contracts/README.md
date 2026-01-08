# ABM Diora NFT Collection

Official NFT collection for ABM Diora blockchain - 15 unique digital assets representing core components of our next-generation Layer 1 blockchain ecosystem.

## 🎨 Collection Overview

### **Genesis Collection**
- **Total Supply**: 15 unique NFTs
- **Standard**: ERC-721
- **Blockchain**: ABM Diora (Chain ID: 1337)
- **Smart Contract**: ABMNFT
- **Creator**: ABM Foundation

### **Collection Categories**

1. **Foundation** (1 NFT) - Genesis Block - Legendary
2. **Consensus** (1 NFT) - Hybrid PoS - Epic
3. **Virtual Machine** (1 NFT) - EVM Compatibility - Epic
4. **Security** (1 NFT) - Enterprise Security - Legendary
5. **Performance** (2 NFTs) - 1000+ TPS, 6-Second Blocks - Rare
6. **Economics** (1 NFT) - DIO Token - Epic
7. **Network** (1 NFT) - 42 Validators - Rare
8. **Ecosystem** (1 NFT) - Developer Tools - Common
9. **Scalability** (1 NFT) - Layer 2 Ready - Rare
10. **Interoperability** (1 NFT) - Cross-Chain Bridge - Epic
11. **Finance** (1 NFT) - DeFi Integration - Rare
12. **Governance** (1 NFT) - Governance Model - Epic
13. **Infrastructure** (1 NFT) - API Gateway - Common
14. **Roadmap** (1 NFT) - Future Vision - Legendary

## 🏗️ Smart Contract Features

### **Core Functionality**
- **ERC-721 Standard**: Full compliance with Ethereum NFT standard
- **Batch Minting**: Mint multiple NFTs in single transaction
- **Presale Support**: Whitelisted presale with discounted pricing
- **Royalty Support**: Built-in royalty distribution
- **Security**: Reentrancy protection and access controls

### **Minting Mechanics**
- **Public Sale**: 100 DIO per NFT
- **Presale**: 50 DIO per NFT (whitelisted addresses)
- **Max per Transaction**: 3 NFTs
- **Max per Wallet**: 5 NFTs total
- **Supply Limit**: 15 NFTs total

### **Security Features**
- **Ownable**: Secure ownership management
- **Reentrancy Guard**: Protection against reentrancy attacks
- **Pausable**: Emergency pause functionality
- **Whitelist**: Presale access control
- **Supply Validation**: Maximum supply enforcement

## 📋 Technical Specifications

### **Contract Details**
- **Solidity Version**: 0.8.26
- **License**: MIT
- **Audited**: Yes (by ABM Foundation Security Team)
- **Gas Optimized**: Yes
- **Upgradeable**: No (immutable contract)

### **Token Standards**
- **ERC-721**: Fully compliant
- **ERC-721Metadata**: Token URI support
- **ERC-2981**: NFT royalty standard
- **EIP-165**: Interface detection

### **Metadata Structure**
```json
{
  "name": "ABM Diora #1 - Genesis Block",
  "description": "The genesis block of ABM Diora blockchain...",
  "image": "https://diorafund.github.io/diora-blockchain/images/nft/1.png",
  "attributes": [
    {
      "trait_type": "Type",
      "value": "Genesis Block"
    },
    {
      "trait_type": "Rarity",
      "value": "Legendary"
    }
  ]
}
```

## 🚀 Deployment

### **Prerequisites**
- Node.js 16+
- Hardhat framework
- ABM Diora RPC endpoint
- Private key with DIO tokens

### **Installation**
```bash
# Clone repository
git clone https://github.com/DioraFund/diora-go-blockchain.git
cd diora-go-blockchain/contracts

# Install dependencies
npm install

# Compile contracts
npm run compile
```

### **Environment Setup**
```bash
# Create .env file
touch .env

# Add environment variables
PRIVATE_KEY=your_private_key_here
DIORA_RPC_URL=https://mainnet-rpc.diora.io
DIORA_API_KEY=your_api_key_here
```

### **Deployment Commands**
```bash
# Deploy to local network
npm run deploy:local

# Deploy to testnet
npm run deploy:testnet

# Deploy to mainnet
npm run deploy:mainnet
```

### **Verification**
```bash
# Verify contract on block explorer
npm run verify -- --network diora <contract-address>
```

## 📄 Metadata Generation

### **Generate All Metadata**
```bash
# Generate metadata for all 15 NFTs
python3 generate_metadata.py
```

### **Metadata Structure**
- **1.json**: Genesis Block
- **2.json**: Hybrid PoS
- **3.json**: EVM Compatibility
- ...
- **15.json**: Future Vision

### **Rarity Distribution**
- **Legendary**: 4 NFTs (26.7%)
- **Epic**: 5 NFTs (33.3%)
- **Rare**: 4 NFTs (26.7%)
- **Common**: 2 NFTs (13.3%)

## 🎨 Visual Design

### **Art Style**
- **Modern Digital Art**: Clean, professional aesthetic
- **Blockchain Theme**: Technology and network motifs
- **Color Schemes**: Gradient-based unique colors
- **Animations**: Subtle motion effects
- **Resolution**: High-quality 1024x1024 PNG

### **Color Coding**
- **Legendary**: Gold, Rainbow gradients
- **Epic**: Purple, Violet gradients
- **Rare**: Blue, Green, Orange gradients
- **Common**: Gray, Teal gradients

## 🔗 Important Links

### **Contract Addresses**
- **Mainnet**: [TBD after deployment]
- **Testnet**: [TBD after deployment]
- **Explorer**: https://diorafund.github.io/diora-blockchain/explorer/

### **Documentation**
- **Whitepaper**: https://diorafund.github.io/diora-whitepaper/
- **Developer Docs**: https://github.com/DioraFund/diora-go-blockchain/tree/main/docs
- **API Reference**: https://github.com/DioraFund/diora-go-blockchain/blob/main/docs/api.md

### **Community**
- **Telegram**: https://t.me/DioraFund
- **Twitter**: https://twitter.com/DioraCrypto
- **GitHub**: https://github.com/DioraFund

## 🛠️ Development

### **Testing**
```bash
# Run all tests
npm test

# Run with gas reporting
npm run gas-report

# Run coverage
npm run coverage
```

### **Linting**
```bash
# Format code
npm run prettier

# Check Solidity style
npm run solhint

# Run all linting
npm run lint
```

### **Local Development**
```bash
# Start local Hardhat node
npm run node:local

# Deploy to local network
npm run deploy:local
```

## 📊 Gas Costs

### **Estimated Gas Usage**
- **Mint Single**: ~150,000 gas
- **Mint Batch (3)**: ~400,000 gas
- **Presale Mint**: ~160,000 gas
- **Transfer**: ~50,000 gas
- **Approve**: ~45,000 gas

### **Gas Optimization**
- **Batch Operations**: Reduced gas per NFT
- **Storage Optimization**: Efficient state management
- **Loop Optimization**: Minimized iteration costs

## 🔒 Security

### **Audit Status**
- **Internal Audit**: ✅ Completed by ABM Foundation
- **Third-Party Audit**: 🔄 In progress
- **Bug Bounty**: 📋 Active program
- **Security Team**: 🛡️ Monitoring 24/7

### **Security Measures**
- **Reentrancy Protection**: OpenZeppelin guards
- **Access Control**: Role-based permissions
- **Input Validation**: Parameter checking
- **Emergency Controls**: Pause and withdraw functions

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.

**Copyright (c) 2026 ABM Foundation**

## 🤝 Contributing

We welcome contributions from the community! Please see our [Contributing Guide](../CONTRIBUTING.md) for details.

## 📞 Support

For questions or support:
- **Email**: info@diora.io
- **Telegram**: https://t.me/DioraFund
- **GitHub Issues**: https://github.com/DioraFund/diora-go-blockchain/issues

---

**Built with ❤️ by ABM Foundation for the Diora ecosystem**
