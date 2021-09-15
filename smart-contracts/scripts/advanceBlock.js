const {web3} = require("@openzeppelin/test-helpers/src/setup");
const { time } = require("@openzeppelin/test-helpers");

async function main() {
  const DEFAULT_BLOCK_TO_ADVANCE = 5
  const argv = require('yargs/yargs')(process.argv.slice(2))
                    .options({
                        'n': {
                        alias: 'blocks',
                        demandOption: false,
                        default: DEFAULT_BLOCK_TO_ADVANCE,
                        describe: 'Number of blocks to advance',
                        type: 'int'
                        },
                        'provider': {
                            demandOption: false,
                            default: process.env.LOCAL_PROVIDER,
                            describe: 'Web3 Provider. e.g. http://localhost:8545. Defaults to env LOCAL_PROVIDER',
                            type: 'string'
                        }
                    })
                    .usage("Usage: node advanceBlock.js [-n blocks] [--provider ethereum_node]")
                    .argv;

  let blocks_to_advance = argv.n;
  if (argv.provider == null || argv.provider == undefined || argv.provider.trim() == "") {
    console.error("Invalid provider. Supply with value or set LOCAL_PROVIDER env var. Received:[",argv.provider,"]")
    return
  }

  require('@openzeppelin/test-helpers/configure')({
    provider: argv.provider,
  });

  try {
    for (let i = 0; i < blocks_to_advance; i++) {
      await time.advanceBlock();
    }

    console.log(`Advanced ${blocks_to_advance} blocks`);

    let bn = await web3.eth.getBlockNumber();

    console.log(`current block number is ${bn}`)

    console.log(JSON.stringify({nBlocks: blocks_to_advance, currentBlockNumber: bn}))
  } catch (error) {
    console.error({ error });
  }
}

main()