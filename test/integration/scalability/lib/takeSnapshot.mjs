import { createArchive } from "../utils/createArchive.mjs";

export async function takeSnapshot({ home = `/tmp/localnet` }) {
  await createArchive({ src: home, name: "localnet" });
}
