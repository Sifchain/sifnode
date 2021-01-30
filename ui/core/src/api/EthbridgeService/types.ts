import { TransactionStatus } from "../../entities";

type TxEventBase<T> = {
  txHash: string;
  payload: T;
};

export type TxEventEthConfCountChanged = {
  type: "EthConfCountChanged";
} & TxEventBase<number>;

export type TxEventSifConfCountChanged = {
  type: "SifConfCountChanged";
} & TxEventBase<number>;

export type TxEventEthTxInitiated = {
  type: "EthTxInitiated";
} & TxEventBase<unknown>;

export type TxEventHashReceived = {
  type: "HashReceived";
} & TxEventBase<string>;

export type TxEventEthTxConfirmed = {
  type: "EthTxConfirmed";
} & TxEventBase<unknown>;

export type TxEventSifTxInitiated = {
  type: "SifTxInitiated";
} & TxEventBase<unknown>;

export type TxEventSifTxConfirmed = {
  type: "SifTxConfirmed";
} & TxEventBase<unknown>;

export type TxEventComplete = {
  type: "Complete";
} & TxEventBase<unknown>;

export type TxEventError = {
  type: "Error";
} & TxEventBase<TransactionStatus>;

export type TxEvent =
  | TxEventEthConfCountChanged
  | TxEventSifConfCountChanged
  | TxEventEthTxInitiated
  | TxEventEthTxConfirmed
  | TxEventSifTxInitiated
  | TxEventSifTxConfirmed
  | TxEventHashReceived
  | TxEventError
  | TxEventComplete;

export type TxEventPrepopulated<T extends TxEvent = TxEvent> = Omit<
  T,
  "txHash"
> & {
  txHash?: string;
};

export type PegTxEventEmitter = {
  setTxHash: (hash: string) => void;
  emit: <T extends TxEvent>(e: TxEventPrepopulated<T>) => void;
  onTxEvent: (handler: (e: TxEvent) => void) => PegTxEventEmitter;
  onTxHash: (handler: (e: TxEventHashReceived) => void) => PegTxEventEmitter;
  onEthConfCountChanged: (
    handler: (e: TxEventEthConfCountChanged) => void
  ) => PegTxEventEmitter;
  onEthTxInitiated: (
    handler: (e: TxEventEthTxInitiated) => void
  ) => PegTxEventEmitter;
  onEthTxConfirmed: (
    handler: (e: TxEventEthTxConfirmed) => void
  ) => PegTxEventEmitter;
  onSifTxInitiated: (
    handler: (e: TxEventSifTxInitiated) => void
  ) => PegTxEventEmitter;
  onSifTxConfirmed: (
    handler: (e: TxEventSifTxConfirmed) => void
  ) => PegTxEventEmitter;
  onComplete: (handler: (e: TxEventComplete) => void) => PegTxEventEmitter;
  onError: (handler: (e: TxEventError) => void) => PegTxEventEmitter;
};
