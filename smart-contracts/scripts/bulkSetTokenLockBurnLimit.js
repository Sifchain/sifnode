const whitelistLimitData = require("./" + process.argv[6]);

module.exports = async (cb) => {

    const err = () => {
        console.log("\nUsage: \nBRIDGEBANK_ADDRESS='0x9201903991991...' truffle exec scripts/bulkSetTokenLockBurnLimit.js --network develop PATH_TO_WHITELIST_FILE.json\n\n\n");
    }

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

    const addresses = whitelistLimitData.array.map(e => {return e.address})
    const limits = whitelistLimitData.array.map(e => {return e.limit})

    if (!addresses || !addresses.length) {
        err();
        throw new Error("Please provide valid address array")
    }

    if (addresses.length !== limits.length) {
        err();
        throw new Error("Address array must equal the amount array");
    }

    const web3 = new Web3(provider);

    contract.setProvider(web3.currentProvider);
    bridgeBank.setProvider(web3.currentProvider);
    BridgeBank.setProvider(web3.currentProvider);

    try {
        const accounts = await web3.eth.getAccounts();

        bridgeBank = await BridgeBank.at(process.env.BRIDGEBANK_ADDRESS)
        console.log(await bridgeBank.bulkWhitelistUpdateLimits(addresses, limits, {
            from: accounts[0],
            gas: 300000 // 300,000 gas
        }));

        console.log("\n\n~~~~ New Token Limits Set ~~~~\n\n");

        for (let i = 0; i < limits.length; i++) {
            console.log(`Token address ${addresses[i]} now has limit ${limits[i]}`);
        }

        cb();
    } catch (error) {
        err()
        console.error({ error });
        cb();
    }
}