## Sifgen

### Setup

#### Seed

When initialising a completely new network, we need to setup a new seed - or master - node:

```
sifgen node create <chain-id>
```

Where:

| Parameter   | Description     |
| -----------| ----------------|
| `chain-id` | The Chain ID.   |


This will output the relevant information about the newly created seed node:

```
chain_id: sifchain
moniker: dawn-river
key_address: sif12tqnz6lfwdvzurvw55h24rgwzrdth7wph9g7qx
key_password: 9SnuJBifdAFoQqE4RxbYsvjO0tly8K1G
peer_address: c3119bfc796f39ad9034854fe3112f2ce80ff580@10.0.22.5:26656
validator_public_key_address: |
    sifvalconspub1zcjduepq0rqzk2gn55k5m2897gn30ayuwdn5fx2c6g592f6duqwlp9d9crmqhwnn4n
```

To start the node, run:

```
sifnoded start
```


#### Validators

There are several steps required in order to setup a new validator. The node is first setup as a listener (or `witness` in Cosmos nomenclature), then funds need to be transferred to the node, and finally the node is promoted by bonding/staking the necessary funds.

##### Listener

Much like how the seed node was setup, to configure this new node as a listener, run (on the new node in question):

```
sifgen node create <chain-id> <peer-address> <genesis-url>
```

Where:

| Parameter       | Description     |
| ---------------| ----------------|
| `chain-id`     | The Chain ID.   |
| `peer-address` | The address of the peer (seed) node. Format is `<node-id>@<ip-address>:<p2p-port>`. Use the `peer_address` from the seed node output. |
| `genesis-url`  | The URL of the genesis RPC endpoint (e.g.: `http://<ip-address>:<rpc-port>/genesis`). |

As with the seed node, this too will result in similar output:

```
chain_id: sifchain
moniker: summer-fire
key_address: sif1qqfmsluekeec93xsuly3656geaahc9sjhq2pdu
key_password: cpVmdNl4fOJaXKgQq7Ak1bZM6rPniw3h
peer_address: ""
validator_public_key_address: |
    sifvalconspub1zcjduepqx0pfm7e0aar92mr77phfeuu0d4pvd3w8u83kc2ehdu05ss09eywqpgxu4l
```

Then start the node by running:

```
sifnoded start
```

and it'll then connect to the seed and sync accordingly.

##### Transfer Funds

The seed node currently contains a hardcoded amount of funds, which can be transferred to would-be validators.

On the seed node, run:

```
sifgen faucet transfer <chain-id> <faucet-password> <faucet-address> <validator-address> <amount>
```

Where:

| Parameter          | Description     |
| ------------------| ----------------|
| `chain-id`        | The Chain ID.   |
| `faucet-password` | The faucet password. Use the `key_password` from the seed node output. |
| `faucet-address`  | The faucet address. Use the `key_address` from the seed node output. |
| `validator-address`      | The validator address. Use the `key_address` from the validator node output. |
| `amount`          | The amount of stake to transfer. E.g.: `1000000000stake`. |

##### Promotion

Once funds have been transferred, the node can be promoted to a full validator. Run (on the node itself):

```
sifgen node promote <chain-id> <moniker> <validator-public-key-address> <key-password> <bond-amount>
```

Where:

| Parameter          | Description     |
| ------------------| ----------------|
| `chain-id`        | The Chain ID.   |
| `moniker` | The validator moniker. Use the `moniker` from the validator node output. |
| `validator-public-key-address`  | The validator public key address. Use the `validator_public_key_address` from the validator node output. |
| `key-password`      | The key password. Use the `key_password` from the validator node output. |
| `bond-amount`          | The amount to bond. E.g.: `1000000000stake`. |

and then the node will be promoted to a validator. To confirm that it was promoted, run:

```
sifnodecli q tendermint-validator-set
```

### Operations

#### Seeds

Nodes can be updated with additional seeds by running (on the node itself):

```
sifgen node update-peers <chain-id> <moniker> <peer-addresses>
```

| Parameter          | Description     |
| ------------------| ----------------|
| `chain-id`        | The Chain ID.   |
| `moniker` | The validator moniker. Use the `moniker` from the validator node output. |
| `peer-addresses` | A comma separated list of peers. |

Once run, the node will require a restart.
