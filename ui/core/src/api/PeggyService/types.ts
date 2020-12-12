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
} & TxEventBase<unknown>;

export type TxEvent =
  | TxEventEthConfCountChanged
  | TxEventSifConfCountChanged
  | TxEventEthTxInitiated
  | TxEventEthTxConfirmed
  | TxEventSifTxInitiated
  | TxEventSifTxConfirmed
  | TxEventError
  | TxEventComplete;

export type TxEventPrepopulated = Omit<TxEvent, "txHash"> & { txHash?: string };

export type TxEventEmitter = {
  emit: (e: TxEventPrepopulated) => void;
  onTxEvent: (handler: (e: TxEvent) => void) => TxEventEmitter;
  onEthConfCountChanged: (
    handler: (e: TxEventEthConfCountChanged) => void
  ) => TxEventEmitter;
  onEthTxInitiated: (
    handler: (e: TxEventEthTxInitiated) => void
  ) => TxEventEmitter;
  onEthTxConfirmed: (
    handler: (e: TxEventEthTxConfirmed) => void
  ) => TxEventEmitter;
  onSifTxInitiated: (
    handler: (e: TxEventSifTxInitiated) => void
  ) => TxEventEmitter;
  onSifTxConfirmed: (
    handler: (e: TxEventSifTxConfirmed) => void
  ) => TxEventEmitter;
  onComplete: (handler: (e: TxEventComplete) => void) => TxEventEmitter;
  onError: (handler: (e: TxEventError) => void) => TxEventEmitter;
};
