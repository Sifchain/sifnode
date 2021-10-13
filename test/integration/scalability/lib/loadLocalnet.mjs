import { extractArchive } from "../utils/extractArchive.mjs";

export async function loadLocalnet({ basePath = "/tmp", name = "localnet" }) {
  await extractArchive({ basePath, name });
}
