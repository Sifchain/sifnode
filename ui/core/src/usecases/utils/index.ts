import { TransactionStatus } from "../../entities";
import { AppEvent } from "../../services/EventBusService";

export function isSupportedEVMChain(chainId?: string) {
  if (!chainId) return false;
  // List of supported EVM chainIds
  const supportedEVMChainIds = [
    "0x1", // 1 Mainnet
    "0x3", // 3 Ropsten
    "0x539", // 1337 Ganache/Hardhat
  ];

  return supportedEVMChainIds.includes(chainId);
}

export const ReportTransactionError = (bus: {
  dispatch: (event: AppEvent) => void;
}) => (txStatus: TransactionStatus): TransactionStatus => {
  bus.dispatch({
    type: "TransactionErrorEvent",
    payload: {
      txStatus,
      message: txStatus.memo || "There was an error with your swap",
    },
  });
  return txStatus;
};
