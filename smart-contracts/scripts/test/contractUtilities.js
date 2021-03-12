const BN = require('bn.js');

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
        case "http://localhost:7545":
            provider = new Web3.providers.HttpProvider(argv.ethereum_network);
            break;
        default:
            const privateKeyDefault = getRequiredEnvironmentVariable(argv.ethereum_private_key_env_var);
            provider = new HDWalletProvider(
                privateKeyDefault,
                argv.ethereum_network,
            );
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

/**
 * Builds a contract object at a particular address
 *
 * For interacting with deployed contracts.  If you're deploying
 * a new contract, use buildBaseContract and then call new on
 * the buildBaseContract result.
 *
 * @param context
 * @param argv
 * @param logging
 * @param name
 * @param address
 * @returns {*}
 */
function buildContract(context, argv, logging, name, address) {
    const contract = buildBaseContract(context, argv, logging, name)
    return contract.at(address);
}

async function setAllowance(context, coinDenom, amount, argv, logging, requestParameters) {
    const sifchainUtilities = context.require('./sifchainUtilities');

    logging.info(`coinDenomis: ${coinDenom}`);
    if (coinDenom != sifchainUtilities.NULL_ADDRESS) {
        const newToken = await buildContract(context, argv, logging, "BridgeToken", coinDenom);
        const currentAllowance = await newToken.allowance(argv.ethereum_address, argv.bridgebank_address, requestParameters);
        logging.info(`currentAllowance is ${currentAllowance}, amount is ${amount}, ${amount.toString(10)}`);
        if (new BN("0").lt(new BN("10"))) {
            logging.info("islt");
        } else {
            logging.info("isgt");
        }
        if (new BN(currentAllowance).lt(new BN(amount))) {
            const approveResult = await newToken.approve(argv.bridgebank_address, sifchainUtilities.SOLIDITY_MAX_INT, requestParameters);
            logging.info(`approve result is ${JSON.stringify(approveResult)}`);
        }
    }
}

module.exports = {buildProvider, buildContract, buildBaseContract, buildWeb3, setAllowance};