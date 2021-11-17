import { fs } from "zx";

export async function createDenomsFile({ home, denoms }) {
  await fs.writeFile(
    `${home}/config/denoms.json`,
    JSON.stringify(denoms, null, 2)
  );
}
