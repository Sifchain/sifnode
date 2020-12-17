Note

We need to sketch out how we are going to get these events

# 1. lock eth in ethereum then mint ceth in sifchain

1.  TxEventEthTxInitiated `{ originTxHash: xxx, destAddress: xxx, amount: n }`
2.  TxEventEthConfCountChanged `{ originTxHash: xxx, count: n }`
3.  TxEventEthTxConfirmed `{ originTxHash: xxx }`
4.  TxEventSifTxInitiated `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }` (_optional_)
5.  TxEventSifConfCountChanged `{ originTxHash: xxx, count: n }` (_optional_)
6.  TxEventSifTxConfirmed `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }`

# 2. burn ceth in sifchain then eth back to ethereum

1.  TxEventEthTxInitiated `{ originTxHash: xxx, destAddress: xxx, amount: n, ... }`
2.  TxEventEthConfCountChanged `{ originTxHash: xxx, count: n }`
3.  TxEventEthTxConfirmed `{ originTxHash: xxx }`
4.  TxEventSifTxInitiated `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }` (_optional_)
5.  TxEventSifConfCountChanged `{ originTxHash: xxx, count: n }` (_optional_)
6.  TxEventSifTxConfirmed `{ originTxHash: xxx, destAddress: xxx, destTxHash: xxx, amount: n }`

# 3. lock rowan in sifchain transfer to ethereum

# 4. burn erowan in ethereum then transfer rowan back to sifchain

### TxEventEthTxInitiated

### TxEventEthConfCountChanged

### TxEventSifConfCountChanged

### TxEventEthTxConfirmed

### TxEventSifTxInitiated

### TxEventSifTxConfirmed

### TxEventError;

### TxEventComplete;
