We need to sketch out how we are going to get these events

# 1. lock eth in ethereum then mint ceth in sifchain

1.  TxEventEthTxInitiated `{ originTxHash: xxx, destAddress: xxx, amount: n }`
2.  TxEventEthConfCountChanged `{ originTxHash: xxx, count: n }`
3.  TxEventEthTxConfirmed `{ originTxHash: xxx }`
4.  TxEventSifTxInitiated `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }` (_optional_)
5.  TxEventSifConfCountChanged `{ originTxHash: xxx, count: n }` (_optional_)
6.  TxEventSifTxConfirmed `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }`

# 2. burn ceth in sifchain then eth back to ethereum

1.  TxEventSifTxInitiated `{ originTxHash: xxx, destAddress: xxx, amount: n, ... }`
2.  TxEventSifConfCountChanged `{ originTxHash: xxx, count: n }`
3.  TxEventSifTxConfirmed `{ originTxHash: xxx }`
4.  TxEventEthTxInitiated `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }` (_optional_)
5.  TxEventEthConfCountChanged `{ originTxHash: xxx, count: n }` (_optional_)
6.  TxEventEthTxConfirmed `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }`

# 3. lock rowan in sifchain transfer to ethereum

1.  TxEventSifTxInitiated `{ originTxHash: xxx, destAddress: xxx, amount: n, ... }`
2.  TxEventSifConfCountChanged `{ originTxHash: xxx, count: n }`
3.  TxEventSifTxConfirmed `{ originTxHash: xxx }`
4.  TxEventEthTxInitiated `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }` (_optional_)
5.  TxEventEthConfCountChanged `{ originTxHash: xxx, count: n }` (_optional_)
6.  TxEventEthTxConfirmed `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }`

# 4. burn erowan in ethereum then transfer rowan back to sifchain

1.  TxEventEthTxInitiated `{ originTxHash: xxx, destAddress: xxx, amount: n }`
2.  TxEventEthConfCountChanged `{ originTxHash: xxx, count: n }`
3.  TxEventEthTxConfirmed `{ originTxHash: xxx }`
4.  TxEventSifTxInitiated `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }` (_optional_)
5.  TxEventSifConfCountChanged `{ originTxHash: xxx, count: n }` (_optional_)
6.  TxEventSifTxConfirmed `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }`

### TxEventEthTxInitiated

TBD. What call/event triggers this event?

### TxEventEthConfCountChanged

TBD. What call/event triggers this event?

### TxEventSifConfCountChanged

TBD. What call/event triggers this event?

### TxEventEthTxConfirmed

TBD. What call/event triggers this event?

### TxEventSifTxInitiated

TBD. What call/event triggers this event?

### TxEventSifTxConfirmed

TBD. What call/event triggers this event?

### TxEventError;

TBD. What call/event triggers this event?

### TxEventComplete;

TBD. What call/event triggers this event?
