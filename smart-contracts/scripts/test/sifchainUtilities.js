const NULL_ADDRESS = "0x0000000000000000000000000000000000000000";

function getRequiredEnvironmentVariable(name) {
    const result = process.env[name];
    if (!result) {
        throw new Error(`${name} does not contain data`);
    }
    return result;
}

const bridgeBankAddressYargOptions = {
    'bridgebank_address': {
        type: "string",
        demandOption: true
    },
};

const symbolYargOption = {
    'symbol': {
        type: "string",
        coerce: addr => addr === "eth" ? NULL_ADDRESS : addr,
        demandOption: true
    },
};

const ethereumAddressYargOption = {
    'ethereum_address': {
        type: "string",
        demandOption: true
    },
};

const amountYargOption = {
    'amount': {
        describe: 'an amount',
        demandOption: true
    },
};

const ethereumNetworkYargOption = {
    'ethereum_network': {
        describe: "can be ropsten or mainnet",
        default: "http://localhost:7545",
    },
};

const transactionYargOptions = {
    ...amountYargOption,
    ...ethereumAddressYargOption,
    ...symbolYargOption,
    ...ethereumNetworkYargOption,
    'bridgebank_address': {
        type: "string",
        demandOption: true
    },
    'sifchain_address': {
        describe: "A SifChain address like sif132tc0acwt8klntn53xatchqztl3ajfxxxsawn8",
        demandOption: true
    },
}

const sharedYargOptions = {
    ...ethereumNetworkYargOption,
    'ethereum_private_key_env_var': {
        describe: "an environment variable that holds a single private key for the sender\nnot used for localnet",
        demandOption: false
    },
    'gas': {
        default: 300000
    },
    'json_path': {
        describe: 'path to the json files',
        default: "../build/contracts",
    },
};

function processArgs(context, args = {}) {
    const yargs = context.require('yargs/yargs')
    const {hideBin} = context.require('yargs/helpers')
    const result = yargs(hideBin(process.argv))
        .options(args)
        .strict()
        .argv
    return result;
}

function configureLogging(context) {
    const winston = context.require('winston');
    const logger = winston.createLogger({
        level: 'debug',
        transports: [
            new winston.transports.File({format: winston.format.simple(), filename: 'combined.log', handleExceptions: true}),
        ],
        exceptionHandlers: [
            new winston.transports.File({ filename: 'combined.log' }),
            new winston.transports.Console({
                format: winston.format.simple()
            })
        ]
    });

    logger.add(new winston.transports.Console({
        format: winston.format.simple()
    }));

    return logger;
}

module.exports = {
    processArgs,
    getRequiredEnvironmentVariable,
    sharedYargOptions,
    configureLogging,
    transactionYargOptions,
    bridgeBankAddressYargOptions,
    ethereumAddressYargOption,
    symbolYargOption,
    amountYargOption,
    NULL_ADDRESS
};