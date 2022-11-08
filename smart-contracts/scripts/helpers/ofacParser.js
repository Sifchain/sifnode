/**
 * This will parse the OFAC list, extracting EVM addresses
 * It will also convert addresses to their checksum version
 * And remove any duplicate addresses found in OFAC's list
 */

const Web3 = require("web3");
const web3 = new Web3();
const axios = require("axios");
const { print, cacheBuster, removeDuplicates } = require("./utils");

const OFAC_URL = "https://www.treasury.gov/ofac/downloads/sdnlist.txt";

async function getList() {
  print("yellow", "Fetching and parsing OFAC blocklist. Please wait...");

  const finalUrl = cacheBuster(OFAC_URL);
  const response = await axios.get(finalUrl).catch((e) => {
    throw e;
  });

  const addresses = extractAddresses(response.data);

  return addresses;
}

function extractAddresses(rawFileContents) {
  const list = rawFileContents.match(/0x[a-fA-F0-9]{40}/g);
  const checksumList = list.map(web3.utils.toChecksumAddress);

  print(
    "magenta",
    `Found ${checksumList.length} EVM addresses. Removing duplicates...`
  );

  const finalList = removeDuplicates(checksumList);

  print("magenta", `The final list has ${finalList.length} unique addresses.`);

  return finalList;
}

module.exports = {
  getList,
  extractAddresses,
};
