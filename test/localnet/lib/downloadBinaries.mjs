import { $, nothrow, cd } from "zx";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const chains = require("../config/chains.json");

export async function downloadBinaries({ home = `/tmp/localnet` }) {
  await $`mkdir -p ${home}`;

  const chainsProps = Object.values(chains);
  const binaryFile = `${home}/binary-file`;

  cd(home);

  for (const { binary, binaryUrl } of chainsProps) {
    if (!binaryUrl) {
      continue;
    }

    console.log(`download ${binaryUrl}`);

    await $`wget ${binaryUrl} -O ${binaryFile}`;

    if (binaryUrl.endsWith(".zip")) {
      await $`unzip ${binaryFile}`;
    } else {
      await $`cp -a ${binaryFile} ${binary}`;
    }

    await nothrow($`chmod +x ${binary}`);
  }
}
