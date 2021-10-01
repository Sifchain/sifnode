const axios = require("axios");
const { print, cacheBuster } = require("./utils");

const OFAC_URL = "https://www.treasury.gov/ofac/downloads/sdnlist.txt";

async function getList() {
  const finalUrl = cacheBuster(OFAC_URL);
  const response = await axios.get(finalUrl).catch((e) => {
    throw e;
  });

  const addresses = extractAddresses(response.data);
  return addresses;
}

async function extractAddresses(rawFileContents) {
  const list = rawFileContents.match();
}

module.exports = {
  getList,
  extractAddresses,
};

getList();
