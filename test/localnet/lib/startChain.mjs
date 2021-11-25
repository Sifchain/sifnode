import { $, nothrow } from "zx";

export async function startChain(props) {
  const {
    disabled,
    chain,
    binary,
    rpcPort = 26657,
    p2pPort = 26656,
    pprofPort = 6060,
    denom,
    home = `/tmp/localnet/config/${props.chain}/${props.chainId}`,
    binPath = `/tmp/localnet/bin`,
    debug = false,
  } = props;

  if (disabled) return;

  if (!binary) throw new Error("missing requirement argument: --binary");
  if (!rpcPort) throw new Error("missing requirement argument: --rpcPort");
  if (!p2pPort) throw new Error("missing requirement argument: --p2pPort");
  if (!pprofPort) throw new Error("missing requirement argument: --pprofPort");
  if (!denom) throw new Error("missing requirement argument: --denom");
  if (!home) throw new Error("missing requirement argument: --home");
  if (!binPath) throw new Error("missing requirement argument: --binPath");

  if (debug) {
    console.log(`
chain       ${chain}
binary      ${binPath}/${binary}
rpcPort     ${rpcPort}
p2pPort     ${p2pPort}
pprofPort   ${pprofPort}
home        ${home}
binPath     ${binPath}
  `);
  }

  const proc = nothrow(
    $`${binPath}/${binary} start --home ${home} --rpc.laddr tcp://127.0.0.1:${rpcPort} --p2p.laddr tcp://127.0.0.1:${p2pPort} --grpc.enable=0 --rpc.pprof_laddr 127.0.0.1:${pprofPort}`
  );
  proc.catch(console.log);

  return {
    proc,
    ...props,
    rpcPort,
    p2pPort,
    pprofPort,
    home,
  };
}
