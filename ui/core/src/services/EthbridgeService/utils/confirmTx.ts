// Some of this was from https://medium.com/pixelpoint/track-blockchain-transactions-like-a-boss-with-web3-js-c149045ca9bf

import Web3 from "web3";

export async function getConfirmations(web3: Web3, txHash: string) {
  try {
    // Get transaction details
    const trx = await web3.eth.getTransaction(txHash);

    // Get current block number
    const currentBlock = await web3.eth.getBlockNumber();

    // When transaction is unconfirmed, its block number is null.
    // In this case we return 0 as number of confirmations
    return trx.blockNumber === null ? 0 : currentBlock - trx.blockNumber;
  } catch (error) {
    console.log(error);
    return 0;
  }
}

export function confirmTx({
  web3,
  txHash,
  confirmations = 10,
  onSuccess = () => {},
  onCheckConfirmation = () => {},
}: {
  web3: Web3;
  txHash: string;
  confirmations: number;
  onSuccess?: () => void;
  onCheckConfirmation?: (confirmations: number) => void;
}) {
  let currentCount = 0;

  setTimeout(async () => {
    const confirmationCount = await getConfirmations(web3, txHash);

    if (currentCount !== confirmationCount) {
      onCheckConfirmation && onCheckConfirmation(confirmationCount);
    }

    currentCount = confirmationCount;

    if (currentCount >= confirmations) {
      onSuccess && onSuccess();
      return;
    }
    confirmTx({
      web3,
      txHash,
      confirmations,
      onSuccess,
      onCheckConfirmation,
    });
  }, 500);
}
