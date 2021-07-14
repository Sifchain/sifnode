# Pay from Treasury (BetaNet)

## Setup

Create a pull request, with a file named `payment-<date>`, in the `./payments` folder, where `<date>` is the current date in `YYYYMMDD` format.

E.g.:

`payment-20210420`

The file must contain the following:

```
recipient=
amount=
```

Where:

|Key|Description|
|---|-----------|
|`recipient`| The recipient wallet address. |
|`amount`| The amount (including the denomination) to send. |

E.g.:

```.env
recipient=sif18qcnjcy3hrp6svzmxaegh3vz96vwwn4augs0z9
amount=10000000000000000000rowan
```

## Execute the Payment

Before running the commands below, please ensure that your PR has been merged.

1. Ensure that the treasury address has been set by running (from your shell):

```bash
treasury_address=<address>
```

E.g.:

```bash
treasury_address=sif1gaej9rvg99xnn8zecznj2vf2tnf87gx60hdkja
```

2. Load the contents of your payment file into your shell:

```bash
source payments/<file>
```

E.g.:

```bash
source payments/payment-20210420
```

3. Transfer the funds:

```
sifnoded tx send $treasury_address $recipient $amount --node tcp://rpc.sifchain.finance:80 --gas-prices 0.5rowan --keyring-backend file --chain-id sifchain
```

4. Verify that the transfer was successful:

```
sifnoded q account $recipient --node tcp://rpc.sifchain.finance:80 --trust-node
```
