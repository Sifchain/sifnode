# Pay from Treasury (BetaNet/Mainnet)

## Setup

Create a pull request, with a `.payment` file, in the project root that contains:

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

1. Ensure that the treasury address has been set by running (from your shell):

```bash
treasury_address=<address>
```

E.g.:

```bash
treasury_address=sif1gaej9rvg99xnn8zecznj2vf2tnf87gx60hdkja
```

2. Once the PR has been merged, execute the payment as follows:

```bash
source .payment
sifnodecli tx send $treasury_address $recipient $amount --node tcp://rpc.sifchain.finance:80 --gas-prices 0.5rowan --keyring-backend file --chain-id sifchain
```

3. Verify that the payment has gone through:

```
sifnodecli q account $recipient --node tcp://rpc.sifchain.finance:80 --trust-node
```
