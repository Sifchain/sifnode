import { getAddress } from "../lib/getAddress.mjs";
import { getChainProps } from "../utils/getChainProps.mjs";

const BASE_PATH = "/tmp/localnet";
const CHAIN_A_NODE = "http://localhost:11000";
const CHAIN_B_NODE = "http://localhost:11001";

const chainAProps = getChainProps({ chain: "sifchain", network: "testnet-1" });
const chainBProps = getChainProps({ chain: "cosmos" });

const chainAHome = `${BASE_PATH}/${chainAProps.chain}/${chainAProps.chainId}`;
const chainBHome = `${BASE_PATH}/${chainBProps.chain}/${chainBProps.chainId}`;

const connection =
  await $`cat ${chainBHome}/relayer/app.yaml | grep srcConnection | cut -d ' ' -f 2`;
const channel =
  await $`ibc-setup channels --chain sifchain --home ${chainBHome}/relayer/ | grep ${connection} | cut -d ' ' -f 1`;

const chainAAddress = await getAddress({
  binary: chainAProps.binary,
  name: `${chainAProps.chain}-source`,
  home: chainAHome,
});
const chainBAddress = await getAddress({
  binary: chainBProps.binary,
  name: `${chainBProps.chain}-source`,
  home: chainBHome,
});

await $`${chainAProps.binary} q bank balances ${chainAAddress} --node ${CHAIN_A_NODE}`;
await $`${chainBProps.binary} q bank balances ${chainBAddress} --node ${CHAIN_B_NODE}`;

await $`
${chainAProps.binary} \
    tx \
    ibc-transfer \
    transfer \
    transfer \
    ${channel} \
    ${chainBAddress} \
    100000000rowan \
    --from ${chainAAddress} \
    --node ${CHAIN_A_NODE} \
    --keyring-backend test \
    --home ${chainAHome} \
    --chain-id ${chainAProps.chainId} \
    --broadcast-mode block \
    -y
`;

await sleep(10000);

await $`${chainAProps.binary} q bank balances ${chainAAddress} --node ${CHAIN_A_NODE}`;
await $`${chainBProps.binary} q bank balances ${chainBAddress} --node ${CHAIN_B_NODE}`;
