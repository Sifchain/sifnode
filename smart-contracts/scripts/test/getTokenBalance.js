module.exports = async (cb) => {
    const Web3 = require("web3");

    const sifchainUtilities = require('./sifchainUtilities')
    const contractUtilites = require('./contractUtilities');

    const argv = sifchainUtilities.processArgs(this, {
        ...sifchainUtilities.sharedYargOptions,
        'bridgetoken_address': {
            type: "string",
        },
        'ethereum_address': {
            type: "string",
            demandOption: true
        },
        'symbol': {
            describe: 'eth or the address of a token',
            default: "eth",
        }
    });

    let balanceWei, balanceEth;
    let result = {};
    try {
        if (argv.symbol === 'eth') {
            const web3instance = contractUtilites.buildWeb3(this, argv);
            balanceWei = await web3instance.eth.getBalance(argv.ethereum_address);
            balanceEth = web3instance.utils.fromWei(balanceWei);
        } else {
            const addr = argv.bridgetoken_address;
            if (!addr)
                throw "must provide --bridgetoken_address for non-eth"
            const bridgeTokenContract = contractUtilites.buildContract(this, argv, "BridgeToken", argv.bridgetoken_address);
            const tokenInstance = await bridgeTokenContract.at(token);
            const name = await tokenInstance.name();
            const symbol = await tokenInstance.symbol();
            const decimals = await tokenInstance.decimals();
            balanceWei = new BigNumber(await tokenInstance.balanceOf(account));
            balanceEth = balanceWei.div(new BigNumber(10).pow(decimals.toNumber()));
        }
        result = {
            ...result,
            balanceWei,
            balanceEth,
            "symbol": argv.symbol
        }
        console.log(JSON.stringify(result, undefined, 0));
        return cb();
    } catch (error) {
        console.error({error});
    }

    return cb();
};
