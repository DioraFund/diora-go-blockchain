require("@nomicfoundation/hardhat-toolbox");

/**
 * Hardhat configuration for ABM Diora NFT deployment
 * @dev Configured for ABM Diora blockchain (Chain ID: 1337)
 */

const PRIVATE_KEY = process.env.PRIVATE_KEY || "";
const DIORA_RPC_URL = process.env.DIORA_RPC_URL || "http://localhost:8545";

module.exports = {
  solidity: {
    version: "0.8.26",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
      viaIR: true,
    },
  },
  networks: {
    // ABM Diora Mainnet
    diora: {
      url: DIORA_RPC_URL,
      chainId: 1337,
      gasPrice: 20000000000, // 20 Gwei
      gas: 6000000,
      accounts: PRIVATE_KEY ? [PRIVATE_KEY] : [],
      timeout: 60000,
      blockGasLimit: 30000000,
    },
    // ABM Diora Testnet
    dioraTestnet: {
      url: "https://testnet-rpc.diora.io",
      chainId: 1338,
      gasPrice: 10000000000, // 10 Gwei
      gas: 6000000,
      accounts: PRIVATE_KEY ? [PRIVATE_KEY] : [],
      timeout: 60000,
      blockGasLimit: 30000000,
    },
    // Local development
    localhost: {
      url: "http://127.0.0.1:8545",
      chainId: 31337,
      gasPrice: 20000000000,
      gas: 6000000,
      timeout: 60000,
    },
    // Hardhat network
    hardhat: {
      chainId: 31337,
      gasPrice: 20000000000,
      gas: 6000000,
      timeout: 60000,
      blockGasLimit: 30000000,
    },
  },
  etherscan: {
    apiKey: {
      diora: process.env.DIORA_API_KEY || "",
      dioraTestnet: process.env.DIORA_TESTNET_API_KEY || "",
    },
    customChains: [
      {
        network: "diora",
        chainId: 1337,
        urls: {
          apiURL: "https://api.diora.io/api",
          browserURL: "https://diorafund.github.io/diora-blockchain/explorer",
        },
      },
      {
        network: "dioraTestnet",
        chainId: 1338,
        urls: {
          apiURL: "https://testnet-api.diora.io/api",
          browserURL: "https://testnet.diorafund.github.io/diora-blockchain/explorer",
        },
      },
    ],
  },
  paths: {
    sources: "./contracts",
    tests: "./test",
    cache: "./cache",
    artifacts: "./artifacts",
  },
  mocha: {
    timeout: 60000,
  },
  gasReporter: {
    enabled: process.env.REPORT_GAS !== undefined,
    currency: "USD",
    gasPrice: 20,
    coinmarketcap: process.env.COINMARKETCAP_API_KEY,
  },
};
