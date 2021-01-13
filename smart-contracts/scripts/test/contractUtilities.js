function buildProvider(context, argv) {
    const HDWalletProvider = context.require("@truffle/hdwallet-provider");
    const Web3 = context.require("web3");
    const {getRequiredEnvironmentVariable} = context.require('./sifchainUtilities');

    let provider;
    if (!argv.ethereum_network)
        throw "Must supply argv.ethereum_network";

    switch (argv.ethereum_network) {
        case "ropsten":
            let ropstenConnectionString = "https://ropsten.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID);
            if (argv.ethereum_private_key_env_var) {
                const ropstenKey = getRequiredEnvironmentVariable(argv.ethereum_private_key_env_var);
                provider = new HDWalletProvider(
                    ropstenKey,
                    ropstenConnectionString
                );
            } else {
                provider = new Web3(ropstenConnectionString);
            }
            break;
        case "mainnet":
            if (argv.ethereum_private_key_env_var) {
                const mainnetKey = getRequiredEnvironmentVariable(argv.ethereum_private_key_env_var);
                provider = new HDWalletProvider(
                    mainnetKey,
                    "https://mainnet.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
                );
            }
            break;
        default:
            provider = new Web3.providers.HttpProvider(argv.ethereum_network);
            break;
    }
    return provider;
}

function buildBridgeBank(context, argv) {
    return buildContract(context, argv, "BridgeBank", argv.bridgebank_address)
}

let web3 = undefined;

function buildWeb3(context, argv) {
    if (web3) {
        return web3;
    } else {
        const provider = buildProvider(context, argv);
        const Web3 = context.require("web3");
        web3 = new Web3(provider);
        return web3;
    }
}

function buildContract(context, argv, name, address) {
    const web3 = buildWeb3(context, argv);
    const contract = context.artifacts.require(name);
    contract.setProvider(web3.currentProvider);
    return contract.at(address);
}

module.exports = {buildProvider, buildContract, buildWeb3};