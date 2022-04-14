/**
 * This will parse the OFAC list, extracting EVM addresses
 * It will also convert addresses to their checksum version
 * And remove any duplicate addresses found in OFAC's list
 */

import { ethers } from "ethers";
import axios from "axios";
import {
  print,
  cacheBuster,
  removeDuplicates,
} from "../../scripts/helpers/utils";

export async function getList(url: string) {
  print("yellow", "Fetching and parsing OFAC blocklist. Please wait...");

  const finalUrl = cacheBuster(url);
  const response = await axios.get(finalUrl).catch((e) => {
    throw e;
  });

  const addresses = extractAddresses(response.data);

  return addresses;
}

export function extractAddresses(rawFileContents: string) {
  const list = rawFileContents.match(/0x[a-fA-F0-9]{40}/g);
  // Handle condition of no addresses in the list
  if (list == null) {
    return []
  }
  const checksumSet = new Set<string>() 
  for (const address of list) {
    try {
      checksumSet.add(ethers.utils.getAddress(address));
    } catch (error) {
      // If address is not a valid Ethereum address, just return an empty string
      console.error(`Detected address ${address} by regex, but ethers says its not valid`, error);
    }
  }

  const finalList = [...checksumSet];

  print("magenta", `Found ${finalList.length} unique EVM addresses.`);

  return finalList;
}