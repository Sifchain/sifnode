require("dotenv").config();

var HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
  compilers: {
    solc: {
      version: "0.6.9",
      optimizer: {
          enabled: true,
          runs: 1000
      }
    }
  },
  networks: {
    sifdocker: {
      host: "localhost",
      port: 7546, // Match default network 'ganache'
      network_id: 5777,
      gas: 6721975, // Truffle default development block gas limit
      gasPrice: 200000000000
    },
    develop: {
      host: "localhost",
      port: 7545, // Match default network 'ganache'
      network_id: 5777,
      gas: 6721975, // Truffle default development block gas limit
      gasPrice: 150000000000
    },
    ropsten: {
      provider: function () {
        return new HDWalletProvider(
          process.env.ETHEREUM_PRIVATE_KEY,
          "https://eth-ropsten.alchemyapi.io/v2/r5vDxb5RNM2NyXGs6A59_sYmH-8C5js6"
        );
      },
      network_id: 3,
      gas: 6000000
    },
    mainnet: {
      provider: function () {
        return new HDWalletProvider(
          process.env.ETHEREUM_PRIVATE_KEY,
          "https://mainnet.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
        );
      },
      network_id: 1,
      gas: 6000000,
      gasPrice: 150000000000
    },
    xdai: {
      provider: function () {
        return new HDWalletProvider(
          process.env.MNEMONIC,
          "https://dai.poa.network"
        );
      },
      network_id: 100,
      gas: 6000000
    }
  },
  rpc: {
    host: "localhost",
    post: 8080
  },
  mocha: {
    useColors: true
  },
  plugins: ["truffle-contract-size", "solidity-coverage"]
};
