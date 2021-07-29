function get() {
    require("dotenv").config();
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");
    const truffleContract = require("truffle-contract");

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

    const owner = process.env.OWNER;
    const pauser = process.env.PAUSER;
    const operator = process.env.OPERATOR;
    const mnemonic = process.env.MNEMONIC;
    const infuraProjectId = process.env.INFURA_PROJECT_ID;

    const consensusThreshold = process.env.CONSENSUS_THRESHOLD;
    const initialValidatorAddresses = process.env.INITIAL_VALIDATOR_ADDRESSES.split(",");
    const initialValidatorPowers = process.env.INITIAL_VALIDATOR_POWERS.split(",");
    const mainnetGasPrice = process.env.MAINNET_GAS_PRICE;
    const erowanAddress = process.env.EROWAN_ADDRESS;
    const alchemyUrl = process.env.ALCHEMY_URL;

    const ethereumPrivateKey = process.env.ETHEREUM_PRIVATE_KEY;
    const localProvider = process.env.LOCAL_PROVIDER;

    const BridgeBankContract = truffleContract(
        require("../build/contracts/BridgeBank.json")
    );
    BridgeBankContract.setProvider(web3.currentProvider);

    const CosmosBridgeContract = truffleContract(
        require("../build/contracts/CosmosBridge.json")
    );
    CosmosBridgeContract.setProvider(web3.currentProvider);

    return {
        web3,
        truffleContract,
        NETWORK_ROPSTEN,
        NETWORK_MAINNET,
        BridgeBankContract,
        CosmosBridgeContract,
        env: {
            owner,
            pauser,
            operator,
            mnemonic,
            infuraProjectId,
            consensusThreshold,
            initialValidatorAddresses,
            initialValidatorPowers,
            mainnetGasPrice,
            erowanAddress,
            alchemyUrl,
            ethereumPrivateKey,
            localProvider
        }
    }
}

module.exports = { get };