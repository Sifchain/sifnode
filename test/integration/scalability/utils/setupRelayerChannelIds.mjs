import { $, nothrow } from "zx";

export async function setupRelayerChannelIds({ home }) {
  const relayerHome = `${home}/relayer`;

  await nothrow($`ibc-setup ics20 --home ${relayerHome}`);
}
