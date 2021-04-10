module.exports = async (cb) => {
    const Web3 = require("web3");
    const BN = require('bn.js');
    const HDWalletProvider = require("@truffle/hdwallet-provider");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.symbolYargOption,
        'ethereum_address': {
            type: "string",
            demandOption: true
        },
    });

    const web3x = new Web3(new HDWalletProvider(
        "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f",
        "http://ganache:7545"
    ));
    logging.info(`getBalance: ${(await web3x.eth.getBalance("0x2191ef87e392377ec08e7c08eb105ef5448eced5"))}`);
    logging.info(`getBalance: ${(await web3x.eth.getBalance("0x2191ef87e392377ec08e7c08eb105ef5448eced5"))}`);
    let balanceWei, balanceEth;
    const result = {};
    try {
        const web3instance = contractUtilites.buildWeb3(this, argv, logging);
        if (argv.symbol === sifchainUtilities.NULL_ADDRESS) {
            balanceWei = await web3instance.eth.getBalance(argv.ethereum_address);
            result.symbol = "eth";
            logging.info(`qinethgettokenbalnace: ${balanceWei.toString(10)} ${argv.ethereum_address}`);
        } else {
            const addr = argv.symbol;
            const tokenContract = await contractUtilites.buildContract(this, argv, logging, "BridgeToken", argv.symbol);
            result["symbol"] = await tokenContract.symbol();
            balanceWei = new BN(await tokenContract.balanceOf(argv.ethereum_address))
            logging.info(`outhgettokenbalnace: ${balanceWei}`);
        }
        balanceEth = web3instance.utils.fromWei(balanceWei.toString());
        const finalResult = {
            ...result,
            balanceWei: balanceWei.toString(10),
            balanceEth: balanceEth.toString(10),
        }

        console.log(JSON.stringify(finalResult, undefined, 0));
        return cb();
    } catch (error) {
        console.error({error});
    }

    return cb();
};
