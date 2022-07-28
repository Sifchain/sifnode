import { spinner } from "zx/experimental";

export async function getPools() {
  const {
    result: { pools },
  } = await spinner("loading pools                              ", () =>
    within(
      async () =>
        // (await fetch("https://api.sifchain.finance/clp/getPools")).json()
        await fs.readJson(`./pools.json`)
    )
  );
  return pools;
}
