const BN = require('bn.js');

module.exports = async cb => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const logging = sifchainUtilities.configureLogging(this);

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        ...sifchainUtilities.symbolYargOption,
        ...sifchainUtilities.amountYargOption,
        ...sifchainUtilities.ethereumAddressYargOption,
        ...sifchainUtilities.bridgeBankAddressYargOptions,
        'operator_address': {
            type: "string",
            demandOption: true,
        },
    });

    try {
        const amount = new BN(argv.amount, 10);

        const transactionParameters = {
            from: argv.operator_address,
            value: 0,
        }

        const newToken = await contractUtilites.buildContract(this, argv, logging, "BridgeToken", argv.symbol);

        const token_destination = argv.operator_address;

        await newToken.mint(token_destination, amount, transactionParameters)

        const bridgeBankContract = await contractUtilites.buildContract(this, argv, logging, "BridgeBank", argv.bridgebank_address);

        const result = {
            destination: token_destination,
            amount: amount.toString(10),
            token_address: newToken.address,
            symbol: await newToken.symbol(),
            name: await newToken.name(),
            decimals: (await newToken.decimals()).toString(10),
        }

        console.log(JSON.stringify(result));
    } catch (e) {
        console.error(`mintTestnetTokens.js error: ${e} ${e.message}`);
        throw(e);
    }

    return cb();
}
