## Sifgen

### Initialize a new network

To initialize a new network, run:

```
sifgen network create <chain-id> <node-count> <output-dir> <seed-ip-address>
```

Where:

| Parameter   | Description     |
| -----------| ----------------|
| `chain-id` | The Chain ID.   |
| `node-count` | The number of nodes (validators) to generate config for. |
| `output-dir` | The root directory of where all the config files should be written to. |
| `seed-ip-address` | The IP (v4) address of the seed node. |


This will then output all the necessary keys and genesis files that are required to boot a new network.

### Create a new node

In addition to creating a new network, you can also create a node (listener/witness) to connect to an existing network, by running:

```
sifgen node create <chain-id> <peer-address> <genesis-url>
```

Where:

| Parameter   | Description     |
| -----------| ----------------|
| `chain-id` | The Chain ID.   |
| `peer-address` | The peer address of the existing network node/validator to connect to. |
| `genesis-url` | The URL of the genesis file to use. |
