module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.transactionYargOptions,
        'approve': {
            type: 'boolean',
            default: true,
            describe: 'approve the amount before burning'
        }
    });

    logging.info(`sendBurnTx: ${JSON.stringify(argv, undefined, 2)}`);

    const bridgeBankContract = await contractUtilites.buildContract(this, argv, "BridgeBank", argv.bridgebank_address);

    const result = {};

    if (argv.approve) {
        const tokenContract = await contractUtilites.buildContract(this, argv, "BridgeToken", argv.symbol);

        result.approve = await tokenContract.approve(argv.bridgebank_address, argv.amount, {
            from: argv.ethereum_address,
            value: 0,
            gas: argv.gas
        });
    }

    result.burn = await bridgeBankContract.burn(
        Web3.utils.utf8ToHex(argv.sifchain_address),
        argv.symbol,
        argv.amount,
        {
            from: argv.ethereum_address,
            value: 0,
            gas: argv.gas
        }
    );

    console.log(JSON.stringify(result, undefined, 0));

    return cb();
};
