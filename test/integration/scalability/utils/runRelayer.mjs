import { $, nothrow } from "zx";

export async function runRelayer({ home }) {
  const relayerHome = `${home}/relayer`;

  const proc = nothrow($`ibc-relayer start -v --poll 10 --home ${relayerHome}`);

  return {
    proc,
  };
}
