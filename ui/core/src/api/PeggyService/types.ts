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
} & TxEventBase<any>;

export type TxEventEthTxConfirmed = {
  type: "EthTxConfirmed";
} & TxEventBase<any>;

export type TxEventSifTxInitiated = {
  type: "SifTxInitiated";
} & TxEventBase<any>;

export type TxEventSifTxConfirmed = {
  type: "SifTxConfirmed";
} & TxEventBase<any>;

export type TxEventComplete = {
  type: "Complete";
} & TxEventBase<any>;

export type TxEvent =
  | TxEventEthConfCountChanged
  | TxEventSifConfCountChanged
  | TxEventEthTxInitiated
  | TxEventEthTxConfirmed
  | TxEventSifTxInitiated
  | TxEventSifTxConfirmed
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
};
