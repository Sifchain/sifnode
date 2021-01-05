module.exports = {
  networks: {
    development: {
      host: "127.0.0.1", // Localhost (default: none)
      port: 7545, // Standard Ethereum port (default: none)
      network_id: "5777", // Any network (default: none)
      gas: 6721975, // Truffle default development block gas limit
      gasPrice: 200000000000,
    },
  },
  compilers: {
    solc: {
      version: "^0.5.0",
    },
  },
};
