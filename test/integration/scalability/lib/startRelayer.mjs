import { runRelayer } from "../utils/runRelayer.mjs";

export async function startRelayer(props) {
  const { home = `/tmp/localnet/${props.chain}/${props.chainId}` } = props;

  // 1) start relayer
  return runRelayer({ home });
}
