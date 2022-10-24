/**
 * Given a list of addresses, returns token details for each address
 * Details fetched are: name, symbol and decimals.
 *
 * Expected usage: first, set the .env variable ADDRESS_LIST_SOURCE with the path for
 * the file that contains a list of addresses in the following format:
 * [
 *  "0x217ddead61a42369a266f1fb754eb5d3ebadc88a",
 *  "0x9e32b13ce7f2e80a01932b42553652e053d6ed8e"
 * ]
 *
 * Then, set the .env variable ADDRESS_LIST_DESTINATION with the destination file
 * The script will write the result to that file.
 *
 * EXAMPLE (.env):
 * ADDRESS_LIST_SOURCE="data/testAddressList.json"
 * ADDRESS_LIST_DESTINATION="data/tokenData.json"
 *
 * Finally, run
 * $ npx hardhat run scripts/generateSifnodeWhitelist.js --network mainnet
 */

require("dotenv").config();
const fs = require("fs");
const { ethers } = require("hardhat");
const _ = require("lodash");

const sifnodeDS = require("../data/ds_sifnode_whitelist.json");
const addressListFile = process.env.ADDRESS_LIST_SOURCE;
const destinationFile = process.env.ADDRESS_LIST_DESTINATION;

function generateDenom(symbol) {
  const denom = "c" + symbol.toLowerCase();
  return denom;
}

async function main() {
  print("yellow", "Starting...", true);

  const ERC20Factory = await ethers.getContractFactory("BridgeToken");

  const data = fs.readFileSync(addressListFile, "utf8");
  const addressList = JSON.parse(data);

  print(
    "yellow",
    `Will fetch data for the following addresses:\n${addressList.join(", ")}`,
    true
  );

  const finalList = [];

  let address;
  for (let i = 0; i < addressList.length; i++) {
    try {
      address = addressList[i];
      console.log(`Processing token ${address}... Please wait...`);

      const instance = await ERC20Factory.attach(address);
      const symbol = generateDenom(await instance.symbol());
      const decimals = (await instance.decimals()).toString();
      const obj = _.cloneDeep(sifnodeDS);

      // required fields to whitelist
      // denom: string
      // base_denom: string
      // decimals: string
      obj.decimals = decimals;
      obj.base_denom = symbol;
      obj.denom = symbol;

      finalList.push(obj);

      print(
        "green",
        `--> Processed token ${symbol} successfully: ${decimals} decimals.`,
        true
      );
    } catch (e) {
      print(
        "red",
        `--> Failed to fetch details of token ${address}: ${e.message}`
      );
    }
  }

  fs.writeFileSync(destinationFile, JSON.stringify(finalList, null, 2));

  print("cyan", `DONE! These results have been written to ${destinationFile}:`);
  print("cyan", JSON.stringify(finalList, null, 2));
}

const colors = {
  green: "\x1b[42m\x1b[37m",
  red: "\x1b[41m\x1b[37m",
  yellow: "\x1b[33m",
  cyan: "\x1b[36m",
  close: "\x1b[0m",
};
function print(color, message, breakLine) {
  const lb = breakLine ? "\n" : "";
  console.log(`${colors[color]}${message}${colors.close}${lb}`);
}

main()
  .catch((error) => {
    console.error({ error });
  })
  .finally(() => process.exit(0));
