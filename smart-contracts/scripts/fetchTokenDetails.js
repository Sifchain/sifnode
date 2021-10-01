/**
 * Given a list of token addresses, fetches metadata for each token.
 * This script is part of the whitlisting process.
 * Please read Whitelist_Update.md for instructions.
 */

require("dotenv").config();
const fs = require('fs');
const axios = require('axios');
const { ethers } = require("hardhat");

const addressListFile = process.env.ADDRESS_LIST_SOURCE;
const destinationFolder = 'data';
const destinationFile = generateDestinationFilename();

async function main() {
  print('yellow', 'Starting...', true);

  const ERC20Factory = await ethers.getContractFactory("BridgeToken");

  const data = fs.readFileSync(addressListFile, 'utf8');
  const addressList = JSON.parse(data);
  
  print('yellow', `Will fetch data for the following addresses:\n${addressList.join(', ')}`, true);

  const finalList = [];

  let address;
  for (let i = 0; i < addressList.length; i++) {
    try {
      address = addressList[i];
      console.log(`Processing token ${address}. Please wait...`);
      const instance = await ERC20Factory.attach(address);
      const name = await instance.name();
      const decimals = await instance.decimals();
      const symbol = await instance.symbol();

      if(!isValidSymbol(symbol)) {
        print('red', `Skipping token ${address} (${name}) because it's symbol has spaces or special characters: ${symbol}`);
        continue;
      }

      const iconUrl = await getTokenMetadata(address);
      
      finalList.push({
        address,
        name,
        symbol,
        decimals,
        // below, properties that  UI cares for:
        network: "ethereum",
        homeNetwork: "ethereum",
        imageUrl: iconUrl,
      });

      print('green', `--> Processed token "${name}" (${symbol}) successfully: ${decimals} decimals.`, true);
    } catch(e) {
      print('red', `--> Failed to fetch details of token ${address}: ${e.message}`);
    }
  }

  // The output file expects this format:
  const output = {
    array: finalList
  }

  fs.writeFileSync(destinationFile, JSON.stringify(output, null, 2));

  print('cyan', `DONE! These results have been written to ${destinationFile}:`);
  print('cyan', JSON.stringify(finalList, null, 2));
}

const colors = {
  green: '\x1b[42m\x1b[37m',
  red: '\x1b[41m\x1b[37m',
  yellow: '\x1b[33m',
  cyan: '\x1b[36m',
  close: '\x1b[0m'
}
function print(color, message, breakLine) {
  const lb = breakLine ? '\n' : '';
  console.log(`${colors[color]}${message}${colors.close}${lb}`);
}

/**
* Will return false for a symbol that has spaces and/or special characters in it
* @param {string} symbol 
* @returns {bool} does the symbol match the RegExp?
*/
function isValidSymbol(symbol) {
  const regexp = new RegExp('^[a-zA-Z0-9]+$');
  return regexp.test(symbol);
}


function generateDestinationFilename() {
  // setup month names
  const monthNames = [
    "Jan", "Feb", "Mar", "Apr", "May", "Jun",
    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
  ];

  // get current date (we do it manually so that it's not dependant on user's locale)
  const today = new Date();
  const day = String(today.getDate()).padStart(2, '0');
  const month = monthNames[today.getMonth()];
  const year = today.getFullYear();

  // transform it in a string with the following format:
  // whitelist_mainnet_update_14_sep_2021.json
  const filename = `${destinationFolder}/whitelist_mainnet_update_${day}_${month}_${year}.json`;

  return filename;
}

async function getTokenMetadata(address) {
  const response = await axios.post(process.env.MAINNET_URL, {
    "jsonrpc":"2.0",
    "method":"alchemy_getTokenMetadata",
    "params":[address],
    "id":1
  }).catch(e => {
    print('red', `-> Cannot find imageUrl. Setting imageUrl to null.`);
    return null;
  });

  return response?.data?.result?.logo;
}

main()
  .catch((error) => {
    console.error({ error });
  })
  .finally(() => process.exit(0))