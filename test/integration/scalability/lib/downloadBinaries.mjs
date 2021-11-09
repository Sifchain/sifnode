const chains = require("../config/chains.json");

export async function downloadBinaries({ home = `/tmp/localnet` }) {
  await $`mkdir -p ${home}`;

  const chainsProps = Object.values(chains);
  const binaryFile = `${home}/binary-file`;

  for (const { binary, binaryUrl } of chainsProps) {
    console.log(`download ${binaryUrl}`);

    await $`wget ${binaryUrl} -O ${binaryFile}`;

    cd(home);
    if (binaryUrl.endsWith(".zip")) {
      await $`unzip ${binaryFile}`;
    } else {
      await $`cp -a ${binaryFile} ${binary}`;
    }
  }
}
