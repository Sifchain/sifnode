const BN = require('bn.js');

module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.symbolYargOption,
        ...sifchainUtilities.amountYargOption,
        ...sifchainUtilities.ethereumAddressYargOption,
        ...sifchainUtilities.bridgeTokenAddressYargOptions,
        'operator_address': {
            type: "string",
            demandOption: true,
        },
    });

    const amount = new BN(argv.amount, 10);

    const standardOptions = {
        from: argv.operator_address
    }

    const newToken = await contractUtilites.buildContract(this, argv, logging, "BridgeToken", argv.bridgetoken_address);

    logging.info(`newToken is ${newToken}`);
    const token_destination = argv.operator_address;

    await newToken.mint(token_destination, amount, standardOptions);

    const result = {
        destination: token_destination,
        "amount": amount.toString(),
    }
    console.log(JSON.stringify(result));

    return cb();
};
