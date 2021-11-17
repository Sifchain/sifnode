import { fs } from "zx";

export async function createGenesisFiles({
  home,
  genesis,
  remoteGenesis,
  defaultGenesis,
}) {
  await fs.writeFile(
    `${home}/config/genesis.remote.json`,
    JSON.stringify(remoteGenesis, null, 2)
  );
  await fs.writeFile(
    `${home}/config/genesis.default.json`,
    JSON.stringify(defaultGenesis, null, 2)
  );
  await fs.writeFile(
    `${home}/config/genesis.json`,
    JSON.stringify(genesis, null, 2)
  );
}
