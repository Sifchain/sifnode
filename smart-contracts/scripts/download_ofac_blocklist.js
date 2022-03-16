const parser = require("./helpers/ofacParser");
const { print } = require("./helpers/utils");
const fs = require("fs");

/**
 * The command line argument is expected to be the full path wehre we want save the OFAC parsed list.
 * Example:
 * node sifnode/smart-contracts/scripts/download_ofac_blocklist.js ~/sifnode/smart-contracts/data/msg-set-blacklist.json
 */

async function main() {
  if (process.argv.length < 3) {
    print("h_red", "please specify a filename to store parsed list");
  }
  const ofac = await parser.getList();
  const msg = {
    addresses: ofac,
  };
  const msgJSON = JSON.stringify(msg);

  try {
    fs.writeFileSync(process.argv[2], msgJSON);
  } catch (err) {
    print("h_red", { err });
    return;
  }

  print("magenta", "File saved.");
}

main()
  .catch((error) => {
    print("h_red", error.message);
  })
  .finally(() => process.exit(0));
