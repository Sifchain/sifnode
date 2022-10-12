# Ledger Command Line Usage

1. Install the cosmos app onto ledger device.

2. Add a ledger key to the keyring
```shell
sifnoded keys add ledger --ledger
```
3. Send a ledger signed transaction
```shell
sifnoded tx bank send ledger toAddress 1000rowan \
  --from ledger \
  --sign-mode amino-json \
  --ledger \
  --fees 100000000000000000rowan
```