# Smart Contracts Development Guide

This guide covers smart contract development on the ABM Diora blockchain, which is fully EVM-compatible.

## 🏗️ EVM Compatibility

ABM Diora supports all standard Ethereum Virtual Machine (EVM) features:

- **Solidity** smart contracts (version 0.8.26+)
- **Vyper** smart contracts
- **ERC-20** token standard
- **ERC-721** NFT standard
- **ERC-1155** multi-token standard
- **Standard Ethereum precompiles**
- **Full opcode support**

## 📝 Writing Smart Contracts

### Solidity Example

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract DioraToken is ERC20, Ownable {
    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply
    ) ERC20(name, symbol) {
        _mint(msg.sender, initialSupply);
    }

    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }

    function burn(uint256 amount) public {
        _burn(msg.sender, amount);
    }
}
```

### Vyper Example

```vyper
# @version >=0.3.0

from vyper.interfaces import ERC20

implements: ERC20

name: public(String[32])
symbol: public(String[32])
decimals: public(uint8)
balanceOf: public(HashMap[address, uint256])
allowance: public(HashMap[address, HashMap[address, uint256]])
totalSupply: public(uint256)

event Transfer:
    sender: indexed(address)
    receiver: indexed(address)
    value: uint256

event Approval:
    owner: indexed(address)
    spender: indexed(address)
    value: uint256

def __init__(name: String[32], symbol: String[32], decimals: uint8, initialSupply: uint256):
    self.name = name
    self.symbol = symbol
    self.decimals = decimals
    self.totalSupply = initialSupply
    self.balanceOf[msg.sender] = initialSupply
    log Transfer(empty_address, msg.sender, initialSupply)

@external
def transfer(to: address, value: uint256) -> bool:
    self.balanceOf[msg.sender] -= value
    self.balanceOf[to] += value
    log Transfer(msg.sender, to, value)
    return True
```

## 🪙 Standard Token Contracts

### ERC-20 Token

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is ERC20, ERC20Burnable, Ownable {
    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply
    ) ERC20(name, symbol) {
        _mint(msg.sender, initialSupply);
    }

    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }
}
```

### ERC-721 NFT

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721, ERC721URIStorage, Ownable {
    uint256 private _tokenIdCounter;

    constructor() ERC721("MyNFT", "MNFT") {}

    function safeMint(address to, string memory uri) public onlyOwner {
        uint256 tokenId = _tokenIdCounter;
        _tokenIdCounter += 1;
        _safeMint(to, tokenId);
        _setTokenURI(tokenId, uri);
    }

    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721, ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
}
```

## 🔧 Development Tools

### Compiler Setup

```bash
# Install Solidity compiler
npm install -g solc

# Install Vyper compiler
pip install vyper

# Verify installation
solc --version
vyper --version
```

### Project Structure

```
contracts/
├── src/
│   ├── Token.sol
│   ├── NFT.sol
│   └── MultiToken.sol
├── test/
│   ├── Token.test.js
│   └── NFT.test.js
├── scripts/
│   ├── deploy.js
│   └── verify.js
├── hardhat.config.js
└── package.json
```

### Hardhat Configuration

```javascript
require("@nomicfoundation/hardhat-toolbox");

module.exports = {
  solidity: {
    version: "0.8.26",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      }
    }
  },
  networks: {
    diora: {
      url: "http://localhost:8545",
      chainId: 1337,
      gasPrice: 20000000000,
      accounts: [
        "0x..." // Your private key
      ]
    }
  }
};
```

## 🚀 Deployment

### CLI Deployment

```bash
# Compile contract
solc --bin --abi Token.sol -o build/

# Deploy contract
./build/diora contract deploy Token.sol

# Deploy with constructor arguments
./build/diora contract deploy Token.sol --args "MyToken" "MTK" "1000000000000000000000000"
```

### Hardhat Deployment

```javascript
// scripts/deploy.js
const { ethers } = require("hardhat");

async function main() {
  const Token = await ethers.getContractFactory("Token");
  const token = await Token.deploy("MyToken", "MTK", "1000000000000000000000000");
  await token.deployed();

  console.log("Token deployed to:", token.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
```

## 🧪 Testing

### Unit Tests with Hardhat

```javascript
// test/Token.test.js
const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Token", function () {
  it("Should have correct name and symbol", async function () {
    const Token = await ethers.getContractFactory("Token");
    const token = await Token.deploy("MyToken", "MTK", "1000000000000000000000000");
    await token.deployed();

    expect(await token.name()).to.equal("MyToken");
    expect(await token.symbol()).to.equal("MTK");
  });

  it("Should mint initial supply", async function () {
    const Token = await ethers.getContractFactory("Token");
    const token = await Token.deploy("MyToken", "MTK", "1000000000000000000000000");
    await token.deployed();

    const ownerBalance = await token.balanceOf(await token.owner());
    expect(await token.totalSupply()).to.equal(ownerBalance);
  });
});
```

## 🔒 Security Best Practices

### Security Checklist

- [ ] Use latest Solidity version
- [ ] Implement access controls
- [ ] Check for integer overflow/underflow
- [ ] Use Reentrancy Guard on external calls
- [ ] Validate all inputs
- [ ] Use SafeMath for arithmetic operations

### Security Example

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract SecureContract is ReentrancyGuard, Pausable, Ownable {
    mapping(address => uint256) public balances;
    
    function deposit() external payable nonReentrant whenNotPaused {
        balances[msg.sender] += msg.value;
    }
    
    function withdraw(uint256 amount) external nonReentrant whenNotPaused {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        
        balances[msg.sender] -= amount;
        
        (bool success,) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
    
    function pause() external onlyOwner {
        _pause();
    }
    
    function unpause() external onlyOwner {
        _unpause();
    }
}
```

## 📚 Additional Resources

- [Solidity Documentation](https://docs.soliditylang.org/)
- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Hardhat Framework](https://hardhat.org/docs)
- [Ethereum Smart Contract Best Practices](https://consensys.github.io/smart-contract-best-practices/)

---

For more information, join our [Telegram](https://t.me/DioraFund) or follow us on [Twitter](https://twitter.com/DioraCrypto).
