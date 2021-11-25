import { $, sleep } from "zx";
import { getAddress } from "../lib/getAddress";
import { getChainProps } from "../utils/getChainProps";
import { loadLocalNet } from "../lib/loadLocalNet";

const CONFIG_PATH = "/tmp/localnet/config";
const BIN_PATH = "/tmp/localnet/bin";
const CHAIN_A_NODE = "http://localhost:11000";
const CHAIN_B_NODE = "http://localhost:11001";

const chainAProps = getChainProps({ chain: "sifchain", network: "testnet-1" });
const chainBProps = getChainProps({ chain: "cosmos" });

const chainAHome = `${CONFIG_PATH}/${chainAProps.chain}/${chainAProps.chainId}`;
const chainBHome = `${CONFIG_PATH}/${chainBProps.chain}/${chainBProps.chainId}`;

let chainsProps;
let relayersProps;

beforeEach(async () => {
  $.verbose = false;
  const result = await loadLocalNet({
    configPath: CONFIG_PATH,
    archivePath: "/tmp/localnet/config.tbz",
  });
  chainsProps = result.chainsProps;
  relayersProps = result.relayersProps;
  $.verbose = true;
}, 10000);

afterEach(async () => {
  await Promise.all(
    relayersProps.map(async ({ proc }) => {
      proc.kill();
    })
  );
  await Promise.all(
    chainsProps.map(async ({ proc }) => {
      proc.kill();
    })
  );
}, 10000);

test("ibc transfer between two chains", async () => {
  $.verbose = false;

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

  const [{ amount: preBalanceA }] = JSON.parse(
    await $`${BIN_PATH}/${chainAProps.binary} q bank balances ${chainAAddress} --node ${CHAIN_A_NODE} --output json`
  ).balances;
  const [{ amount: preBalanceB }] = JSON.parse(
    await $`${BIN_PATH}/${chainBProps.binary} q bank balances ${chainBAddress} --node ${CHAIN_B_NODE} --output json`
  ).balances;

  expect(preBalanceA).toBe("9999999799999800000");
  expect(preBalanceB).toBe("9999999900000000000");

  await sleep(5000);

  await $`
${BIN_PATH}/${chainAProps.binary} \
    tx \
    ibc-transfer \
    transfer \
    transfer \
    ${channel} \
    ${chainBAddress} \
    100000000rowan \
    --fees ${chainAProps.fees}${chainAProps.denom} \
    --from ${chainAAddress} \
    --node ${CHAIN_A_NODE} \
    --keyring-backend test \
    --home ${chainAHome} \
    --chain-id ${chainAProps.chainId} \
    --broadcast-mode block \
    -y
`;

  // await sleep(10000);

  const [{ amount: postBalanceA }] = JSON.parse(
    await $`${BIN_PATH}/${chainAProps.binary} q bank balances ${chainAAddress} --node ${CHAIN_A_NODE} --output json`
  ).balances;
  const [{ amount: postBalanceB }] = JSON.parse(
    await $`${BIN_PATH}/${chainBProps.binary} q bank balances ${chainBAddress} --node ${CHAIN_B_NODE} --output json`
  ).balances;

  expect(postBalanceA).toBe("9999999799899700000");
  expect(postBalanceB).toBe(preBalanceB);

  $.verbose = true;
}, 20000);
