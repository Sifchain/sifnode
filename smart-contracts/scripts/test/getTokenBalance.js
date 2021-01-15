module.exports = async (cb) => {
    const Web3 = require("web3");
    const BigNumber = require("bignumber.js")

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
        const web3instance = contractUtilites.buildWeb3(this, argv);
        if (argv.symbol === 'eth') {
            balanceWei = await web3instance.eth.getBalance(argv.ethereum_address);
            result.symbol = "eth";
        } else {
            const addr = argv.bridgetoken_address;
            if (!addr)
                throw "must provide --bridgetoken_address for non-eth"
            const bridgeTokenContract = await contractUtilites.buildContract(this, argv, "BridgeToken", argv.bridgetoken_address);
            result["symbol"] = bridgeTokenContract.symbol;
            balanceWei = new BigNumber(await bridgeTokenContract.balanceOf(argv.ethereum_address))
        }
        balanceEth = web3instance.utils.fromWei(balanceWei.toString());
        result = {
            ...result,
            balanceWei,
            balanceEth,
        }
        console.log(JSON.stringify(result, undefined, 0));
        return cb();
    } catch (error) {
        console.error({error});
    }

    return cb();
};
