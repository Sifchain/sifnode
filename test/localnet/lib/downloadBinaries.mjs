import { $, nothrow, cd } from "zx";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const chains = require("../config/chains.json");

export async function downloadBinaries({ home = `/tmp/localnet/.bin` }) {
  await $`mkdir -p ${home}`;

  const chainsProps = Object.values(chains);
  const tempFile = `${home}/temp-file`;

  cd(home);

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

    await $`wget ${url} -O ${tempFile}`;

    if (url.endsWith(".zip")) {
      await $`unzip ${tempFile}`;
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
      cd(`${home}/${sourceRelativePath}`);
      await $`GOBIN=${home} make install`;
      cd(home);
      await $`rm -rf ${home}/${sourceRelativePath}`;
    }

    await nothrow($`chmod +x ${binary}`);
  }

  await $`rm -f ${tempFile}`;
}
