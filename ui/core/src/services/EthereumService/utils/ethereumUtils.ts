import { MetaMaskInpageProvider } from "@metamask/inpage-provider";
import JSBI from "jsbi";
import Web3 from "web3";
import {
  IpcProvider,
  provider,
  TransactionReceipt,
  WebsocketProvider,
} from "web3-core";

import { Address, Asset, AssetAmount, Network } from "../../../entities";

import erc20TokenAbi from "./erc20TokenAbi";

export function getTokenContract(web3: Web3, asset: Asset) {
  return new web3.eth.Contract(erc20TokenAbi, asset.address);
}

export async function getTokenBalance(
  web3: Web3,
  address: Address,
  asset: Asset,
) {
  const contract = getTokenContract(web3, asset);
  let tokenBalance = "0";
  try {
    tokenBalance = await contract.methods.balanceOf(address).call();
  } catch (err) {
    // console.log(`Error fetching balance for ${asset.symbol}`);
  }
  return AssetAmount(asset, tokenBalance);
}

export function isEventEmittingProvider(
  provider?: unknown,
): provider is WebsocketProvider | IpcProvider {
  if (!provider || typeof provider === "string") return false;
  return typeof (provider as any).on === "function";
}

export type EIPProvider = MetaMaskInpageProvider;

export function isMetaMaskInpageProvider(
  provider?: unknown,
): provider is EIPProvider {
  if (!provider || typeof provider === "string") return false;
  return typeof (provider as any).request === "function";
}

// Transfer token or ether
export async function transferAsset(
  web3: Web3,
  fromAddress: Address,
  toAddress: Address,
  amount: JSBI,
  asset?: Asset,
) {
  if (asset?.address) {
    return await transferToken(web3, fromAddress, toAddress, amount, asset);
  }

  return await transferEther(web3, fromAddress, toAddress, amount);
}

// Transfer token
export async function transferToken(
  web3: Web3,
  fromAddress: Address,
  toAddress: Address,
  amount: JSBI,
  asset: Asset,
) {
  const contract = getTokenContract(web3, asset);
  return new Promise<string>((resolve, reject) => {
    let hash: string;
    let receipt: boolean;

    function resolvePromise() {
      if (receipt && hash) resolve(hash);
    }

    contract.methods
      .transfer(toAddress, amount.toString())
      .send({ from: fromAddress })
      .on("transactionHash", (_hash: string) => {
        hash = _hash;
        resolvePromise();
      })
      .on("receipt", (_receipt: boolean) => {
        receipt = _receipt;
        resolvePromise();
      })
      .on("error", (err: any) => {
        reject(err);
      });
  });
}

// Transfer ether
export async function transferEther(
  web3: Web3,
  fromAddress: Address,
  toAddress: Address,
  amount: JSBI,
) {
  return new Promise<string>((resolve, reject) => {
    let hash: string;
    let receipt: TransactionReceipt;

    function resolvePromise() {
      if (receipt && hash) resolve(hash);
    }

    web3.eth
      .sendTransaction({
        from: fromAddress,
        to: toAddress,
        value: amount.toString(),
      })
      .on("transactionHash", (_hash: string) => {
        hash = _hash;
        resolvePromise();
      })
      .on("receipt", (_receipt) => {
        receipt = _receipt;
        resolvePromise();
      })
      .on("error", (err: any) => {
        reject(err);
      });
  });
}

export async function getEtheriumBalance(web3: Web3, address: Address) {
  const ethBalance = await web3.eth.getBalance(address);
  // TODO: Pull as search from supported tokens
  return AssetAmount(
    {
      symbol: "eth",
      label: "ETH",
      address: "",
      decimals: 18,
      name: "Ethereum",
      network: Network.ETHEREUM,
    },
    ethBalance,
  );
}
