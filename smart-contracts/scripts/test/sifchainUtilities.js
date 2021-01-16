function getRequiredEnvironmentVariable(name) {
    const result = process.env[name];
    if (!result) {
        throw new Error(`${name} does not contain data`);
    }
    return result;
}

const sharedYargOptions = {
    'ethereum_network': {
        describe: "can be ropsten or mainnet",
        default: "http://localhost:7545",
    },
    'ethereum_private_key_env_var': {
        describe: "an environment variable that holds a single private key for the sender\nnot used for localnet",
        demandOption: false
    },
}

const transactionYargOptions = {
    'symbol': {
        describe: 'eth, erowan, etc',
        default: "eth",
    },
    'amount': {
        describe: 'an amount',
        demandOption: true
    },
    'gas': {
        default: 300000
    },
    'json_path': {
        describe: 'path to the json files',
        default: "../build/contracts",
    },
    'bridgebank_address': {
        type: "string",
        demandOption: true
    },
    'sifchain_address': {
        describe: "A SifChain address like sif132tc0acwt8klntn53xatchqztl3ajfxxxsawn8",
        demandOption: true
    },
    'ethereum_address': {
        type: "string",
        demandOption: true
    },
    'ethereum_network': {
        describe: "can be ropsten or mainnet",
        default: "http://localhost:7545",
    },
}

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
            new winston.transports.File({filename: 'combined.log', handleExceptions: true}),
        ],
    });

    logger.add(new winston.transports.Console({
        format: winston.format.simple()
    }));

    return logger;
}

module.exports = {processArgs, getRequiredEnvironmentVariable, sharedYargOptions, configureLogging, transactionYargOptions};