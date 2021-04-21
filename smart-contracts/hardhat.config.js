require('@openzeppelin/hardhat-upgrades');
require("@nomiclabs/hardhat-waffle");
require('hardhat-local-networks-config-plugin')
require("hardhat-typechain");

// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task("accounts", "Prints the list of accounts", async () => {
  const accounts = await ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  localNetworksConfig: '~/.hardhat/networks.ts',
  networks: {
    hardhat: {
      allowUnlimitedContractSize: false,
    },
    localhost: {
      url: "http://localhost:7545",
    }
  },
  solidity: {
    version: "0.8.0",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200
      },
    },
  },
  typechain: {
    outDir: "build",
    target: "ethers-v5"
  },
};
