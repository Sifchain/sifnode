function buildProvider(context, argv, logging) {
    const HDWalletProvider = context.require("@truffle/hdwallet-provider");
    const Web3 = context.require("web3");
    const {getRequiredEnvironmentVariable} = context.require('./sifchainUtilities');

    let provider;
    if (!argv.ethereum_network)
        throw "Must supply argv.ethereum_network";

    switch (argv.ethereum_network) {
        case "ropsten":
        case "mainnet":
            let netConnectionString = `https://${argv.ethereum_network}.infura.io/v3/${process.env.INFURA_PROJECT_ID}`;
            if (argv.ethereum_private_key_env_var) {
                const privateKey = getRequiredEnvironmentVariable(argv.ethereum_private_key_env_var);
                provider = new HDWalletProvider(
                    privateKey,
                    netConnectionString
                );
            } else {
                provider = new Web3(netConnectionString);
            }
            break;
        default:
            provider = new Web3.providers.HttpProvider(argv.ethereum_network);
            break;
    }
    return provider;
}

function buildBridgeBank(context, argv) {
    return buildContract(context, argv, logging, "BridgeBank", argv.bridgebank_address)
}

let web3 = undefined;

function buildWeb3(context, argv, logging) {
    if (web3) {
        return web3;
    } else {
        const provider = buildProvider(context, argv, logging);
        const Web3 = context.require("web3");
        web3 = new Web3(provider);
        return web3;
    }
}

const truffleContractProvider = require("@truffle/contract");

function buildBaseContract(context, argv, logging, name) {
    const web3 = buildWeb3(context, argv, logging);
    const js = `${argv.json_path}/${name}.json`;
    let solidityContractJson = require(js);
    const contract = truffleContractProvider(solidityContractJson);
    contract.setProvider(web3.currentProvider);
    return contract;
}

function buildContract(context, argv, logging, name, address) {
    const contract = buildBaseContract(context, argv, logging, name)
    return contract.at(address);
}

module.exports = {buildProvider, buildContract, buildBaseContract, buildWeb3};