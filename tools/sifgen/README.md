## Sifgen

### Initialize a net network

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
