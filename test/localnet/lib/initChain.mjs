import { $ } from "zx";
import { cleanUpGenesisState } from "../utils/cleanUpGenesisState.mjs";
import { createDenomsFile } from "../utils/createDenomsFile.mjs";
import { createGenesisFiles } from "../utils/createGenesisFiles.mjs";
import { getDenoms } from "../utils/getDenoms.mjs";
import { getRemoteGenesis } from "../utils/getRemoteGenesis.mjs";

import { createRequire } from "module";
const require = createRequire(import.meta.url);

export async function initChain(props) {
  const {
    disabled,
    chain,
    binary,
    chainId,
    node,
    amount = 10e18,
    denom,
    home = `/tmp/localnet/config/${props.chain}/${props.chainId}`,
    binPath = `/tmp/localnet/bin`,
    debug = false,
  } = props;

  if (disabled) return;

  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!node) throw new Error("missing requirement argument: --node");
  if (!amount) throw new Error("missing requirement argument: --amount");
  if (!denom) throw new Error("missing requirement argument: --denom");
  if (!home) throw new Error("missing requirement argument: --home");
  if (!binPath) throw new Error("missing requirement argument: --binPath");

  const validatorAccountName = `${chain}-validator`;
  const sourceAccountName = `${chain}-source`;

  if (debug) {
    console.log(`
chain                 ${chain}
binary                ${binPath}/${binary}
chainId               ${chainId}
node                  ${node}
validatorAccountName  ${validatorAccountName}
sourceAccountName     ${sourceAccountName}
amount                ${amount}
denom                 ${denom}
home                  ${home}
binPath               ${binPath}
  `);
  }

  const { remoteGenesis } = await getRemoteGenesis({ node });

  await $`rm -rf ${home}`;
  await $`mkdir -p ${home}`;
  await $`${binPath}/${binary} init ${chainId} --chain-id ${chainId} --home ${home}`;
  await $`${binPath}/${binary} keys add ${validatorAccountName} --keyring-backend test --home ${home}`;
  await $`${binPath}/${binary} keys add ${sourceAccountName} --keyring-backend test --home ${home}`;
  await $`${binPath}/${binary} add-genesis-account ${validatorAccountName} ${amount}${denom} --keyring-backend test --home ${home}`;
  await $`${binPath}/${binary} add-genesis-account ${sourceAccountName} ${amount}${denom} --keyring-backend test --home ${home}`;

  const defaultGenesis = require(`${home}/config/genesis.json`);
  const genesis = cleanUpGenesisState({ remoteGenesis, defaultGenesis });
  await createGenesisFiles({ home, genesis, remoteGenesis, defaultGenesis });

  if (chain === "sifchain") {
    const denoms = getDenoms();
    await createDenomsFile({ home, denoms });
    await $`${binPath}/${binary} set-gen-denom-whitelist ${home}/config/denoms.json --home ${home}`;
  }

  await $`${binPath}/${binary} gentx ${validatorAccountName} ${amount}${denom} --chain-id ${chainId} --keyring-backend test --home ${home}`;
  // $.verbose = false;
  await $`${binPath}/${binary} collect-gentxs --home ${home}`;
  // $.verbose = true;
}
