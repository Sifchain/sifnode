import { fetch } from "zx";

export async function getRemoteGenesis({ node }) {
  const response = await fetch(`${node}/genesis`);
  const data = await response.json();

  if (!data.result || !data.result.genesis) {
    console.error("wrong genesis content");

    return { remoteGenesis: {} };
  }

  const {
    result: { genesis: remoteGenesis },
  } = data;

  return { remoteGenesis };
}
