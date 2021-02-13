import { Address, TxHash } from "src/entities";

const win = window as any;
const localStorage = win.localStorage;

// data model
// { address: {[ txHash: { status, amount, symbol }]}}
// { address: {[ txHash: { tA, tB, aA, aB }]}}
// localStorage.setItem(address, JSON.stringify({[{txHash:{ tA, tB, aA, aB }]}))

// todo DRY

// todo pass in object as defined above, right now just txHash,  w hardcode "pending"
export function setTransaction(address: Address, txHash: TxHash ) {
  if (!address && !txHash) return;
  // check if key exists
  const transactions: string = localStorage.getItem(address);
  console.log(transactions)
  if (!transactions) {
    const transaction = {} as any
    transaction[txHash] = { status: "pending" }
    localStorage.setItem(address, JSON.stringify([transaction]));
    return
  }
  const parsedTransactions = JSON.parse(transactions)
  const transaction = {} as any
  transaction[txHash] = { status: "pending" }
  parsedTransactions.push(transaction)  
  localStorage.setItem(address, JSON.stringify(parsedTransactions))
  return

}

export function getTransactions(address: Address, txHash?: TxHash) {
  if (!address) return;
  const transactions = localStorage.getItem(address);
  console.log("Found localStorage txd:", transactions);
  if (!transactions) return;
  return JSON.parse(transactions);

  // then what ?
  // for each tx, query chain, create notification, setItem
}

export function removeItem(address: Address, txHash: TxHash) {
  // check if key exists,
  // if does, parse it, find txHash,
  // then remove it
  // return success
}
