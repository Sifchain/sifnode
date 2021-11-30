import { $, nothrow, cd } from "zx";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const chains = require("../config/chains.json");

export async function downloadBinaries({ binPath = `/tmp/localnet/bin` }) {
  await $`rm -rf ${binPath}`;
  await $`mkdir -p ${binPath}`;

  const chainsProps = Object.values(chains);
  const tempFile = `${binPath}/.tempFile`;

  cd(binPath);

  for (const {
    disabled,
    binary,
    binaryUrl,
    binaryRelativePath,
    sourceUrl,
    sourceRelativePath,
  } of chainsProps) {
    if (disabled) {
      continue;
    }

    const url = binaryUrl || sourceUrl;

    console.log(`download ${url}`);

    await $`wget -q ${url} -O ${tempFile}`;

    if (url.endsWith(".zip")) {
      await $`unzip -q ${tempFile}`;
    } else if (binaryUrl.endsWith(".tar.gz")) {
      await $`tar xvzf ${tempFile}`;
    } else {
      if (binaryUrl) {
        await $`cp -a ${tempFile} ${binary}`;
      }
    }

    if (binaryUrl && binaryRelativePath) {
      await $`mv ${binaryRelativePath} ${binary}`;
    } else if (sourceUrl && sourceRelativePath) {
      cd(`${binPath}/${sourceRelativePath}`);
      await $`GOBIN=${binPath} make install`;
      cd(binPath);
      await $`rm -rf ${binPath}/${sourceRelativePath}`;
    }

    await nothrow($`chmod +x ${binary}`);
  }

  await $`rm -f ${tempFile}`;
}
