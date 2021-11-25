import { runRelayer } from "../utils/runRelayer.mjs";

export async function startRelayer(props) {
  const { home = `/tmp/localnet/config/${props.chain}/${props.chainId}` } =
    props;

  // 1) start relayer
  const { proc } = await runRelayer({ home });

  for await (let chunk of proc.stdout) {
    if (chunk.includes("next heights to relay")) break;
  }
  proc.kill();
}
