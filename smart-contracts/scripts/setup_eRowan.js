module.exports = async (cb) => {
    try {
        const HDWalletProvider = require("@truffle/hdwallet-provider");
        const Web3 = require("web3");

        // Contract abstraction
        const truffleContract = require("truffle-contract");
        const contract = truffleContract(
            require("../build/contracts/BridgeToken.json")
        );
        let bridgeBank = truffleContract(
            require("../build/contracts/BridgeBank.json")
        );

        const BridgeBank = artifacts.require("BridgeBank")

        const address = process.env.EROWAN_ADDRESS
        if (!address || address.length !== 42) {
          throw new Error("Please provide valid eRowan token address")
        }

        const NETWORK_ROPSTEN =
          process.argv[4] === "--network" && process.argv[5] === "ropsten";

        const NETWORK_MAINNET =
          process.argv[4] === "--network" && process.argv[5] === "mainnet";

        let provider;
        if (NETWORK_ROPSTEN) {
          provider = new HDWalletProvider(
            process.env.ETHEREUM_PRIVATE_KEY,
            "https://ropsten.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
          );
        } else if (NETWORK_MAINNET) {
          provider = new HDWalletProvider(
            process.env.ETHEREUM_PRIVATE_KEY,
            "https://mainnet.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
          );
        } else {
          provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
        }

        const web3 = new Web3(provider);

        contract.setProvider(web3.currentProvider);
        bridgeBank.setProvider(web3.currentProvider);
        BridgeBank.setProvider(web3.currentProvider);

        try {
          const accounts = await web3.eth.getAccounts();
          const bridgeToken = await contract.at(address);

          bridgeBank = await BridgeBank.deployed()

          await bridgeBank.addExistingBridgeToken(bridgeToken.address, {
            from: accounts[0],
            gas: 300000, // 300,000 Gwei
            gasPrice: 190000000000 // web3.utils.toWei("50", "gwei"),
          });
          console.log("Finished")
          cb()
        } catch (error) {
          console.error({ error })
          cb()
        }
    } catch(error) {
        console.log("error: ", error)
    }
}
