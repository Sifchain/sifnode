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

    const bridgeBankContract = await contractUtilites.buildContract(this, argv, logging, "BridgeBank", argv.bridgebank_address);

    const result = {};

    let gasLimit = argv.gas;
    if (gasLimit === 'estimate') {
        gasLimit = 6000000; // we don't do an actual estimate for burns, just locks
    }

    // see if the user asked to approve the amount first
    if (argv.approve) {
        const tokenContract = await contractUtilites.buildContract(this, argv, logging,"BridgeToken", argv.symbol);

        result.approve = await tokenContract.approve(argv.bridgebank_address, argv.amount, {
            from: argv.ethereum_address,
            value: 0,
            gas: gasLimit
        });
    }

    result.burn = await bridgeBankContract.burn(
        Web3.utils.utf8ToHex(argv.sifchain_address),
        argv.symbol,
        argv.amount,
        {
            from: argv.ethereum_address,
            value: 0,
            gas: gasLimit
        }
    );

    console.log(JSON.stringify(result, undefined, 0));

    return cb();
};
