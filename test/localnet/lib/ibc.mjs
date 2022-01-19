import { $ } from "zx";
import { getChainProps } from "../utils/getChainProps.mjs";
import { getAddress } from "./getAddress.mjs";

export async function ibc(props) {
  const sifChainProps = getChainProps({ chain: "sifchain" });

  const {
    node,
    chainId,
    type,
    channelId,
    counterpartyChannelId,
    binary,
    name,
    denom,
    amount = 1,
    fees = undefined,
    gas = undefined,
    dryRun = false,
    times = 1,
    broadcast = "async",
    timeout = 600000000000,
    memo = undefined,
  } = getChainProps({
    ...props,
    chain: props.type === "issuer" ? props.chain : "sifchain",
    node: props.type === "issuer" ? props.node : undefined,
    chainId: props.type === "issuer" ? props.chainId : undefined,
    binary: props.type === "issuer" ? props.binary : undefined,
    denom: props.type === "issuer" ? props.denom : undefined,
  });

  if (!node) throw new Error("missing requirement argument: --node");
  if (!chainId) throw new Error("missing requirement argument: --chain-id");
  if (!type) throw new Error("missing requirement argument: --type");
  if (type !== "issuer" && type !== "receiver")
    throw new Error("wrong argument value: --type [issuer|receiver]");
  if (!channelId) throw new Error("missing requirement argument: --channelId");
  if (!counterpartyChannelId)
    throw new Error("missing requirement argument: --counterpartyChannelId");
  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!name) throw new Error("missing requirement argument: --name");
  if (!denom) throw new Error("missing requirement argument: --denom");
  if (!amount) throw new Error("missing requirement argument: --amount");
  if (!times) throw new Error("missing requirement argument: --times");
  if (!broadcast) throw new Error("missing requirement argument: --broadcast");
  if (!timeout) throw new Error("missing requirement argument: --timeout");

  console.log(`
node                      ${node}
chainId                   ${chainId}
type                      ${type}
channelId                 channel-${channelId}
counterpartyChannelId     channel-${counterpartyChannelId}
binary                    ${binary}
name                      ${name}
denom                     ${props.denom}
amount                    ${amount}
fees denom                ${denom}
fees                      ${fees}
gas                       ${gas}
times                     ${times}
broadcast                 ${broadcast}
timeout                   ${timeout}
memo                      ${memo}
`);

  const srcAddr = await getAddress({
    binary: type === "receiver" ? sifChainProps.binary : props.binary,
    name,
  });
  const dstAddr = await getAddress({
    binary: type === "receiver" ? props.binary : sifChainProps.binary,
    name,
  });

  const response =
    await $`${binary} q auth account ${srcAddr} --node ${node} --chain-id ${chainId} --output json | jq -r "{account_number: .account_number, sequence: .sequence}"`;
  const { account_number: accountNumber, sequence } = JSON.parse(response);

  let seq = sequence;
  for (let i = 0; i < times; i++) {
    console.log(`tx ${i} processing`);
    await $`
${binary} \
  tx \
  ibc-transfer \
  transfer \
  transfer \
  channel-${type === "receiver" ? channelId : counterpartyChannelId} \
  ${dstAddr} \
  ${amount}${props.denom} \
  --from ${srcAddr} \
  --keyring-backend test \
  ${fees ? `--fees` : ``} ${fees ? `${fees}${denom}` : ``} \
  ${gas ? `--gas` : ``} ${gas ? `${gas}` : ``} \
  --chain-id ${chainId} \
  --node ${node} \
  --broadcast-mode ${broadcast ? broadcast : `async`} \
  --packet-timeout-timestamp ${timeout} \
  ${broadcast === "async" ? `--offline` : ``} \
  --sequence ${seq} \
  --account-number ${accountNumber} \
  ${!dryRun ? "--yes" : ""}
`;
    console.log(`tx ${i} done`);
    seq++;
  }
}

// TODO: memo does not seem to have worked when tested with IRIS on Fri 3rd Sep
// --memo ${memo ? `[tx-${i}]: ${memo}` : `[tx-${i}]`} \
