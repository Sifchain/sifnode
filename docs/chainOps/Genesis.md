# Sifnode Genesis

Sifnode includes the functionality to instantiate a four(4) node local testnet cluster, to perform a genesis ceremony.

## Requirements

- [Ruby 2.6.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker)

## Setup

1. Ensure you've installed the dependencies listed above.

2. Clone the repository:

```
git clone ssh://git@github.com/Sifchain/sifnode && cd sifnode
```

3. Compile:

```
make clean install
```

4. Switch into the `build` directory:

```
cd build
```

5. Scaffold your network (this will create all the keys and associated config files for your cluster):

```
rake "genesis:network:scaffold[chainnet]"
```

where:

| Parameter | Description |
|-----------|-------------|
|chainnet | The name of your chain/network |


6. Once the scaffold process has completed, your new cluster can be booted by running:

```
rake "genesis:network:boot[chainnet,eth_address,eth_keys,eth_websocket]"
```

where:

| Parameter | Description | Example |
|-----------|-------------|---------|
| chainnet | The name of your chain/network | |
| eth_address | The Ethereum contract address for peggy. | |
| eth_keys | A list of Ethereum private keys for each of the four validators. | 'key1 key2 key3 key4' |
| eth_websocket | The Ethereum websocket address to connect to. | |

and your local cluster will start accordingly.

The nodes in the cluster will map to the following TCP ports on your environment:

| Node | RPC Port | P2P Port | Rest Server |
|------|----------|----------|-------------|
| sifnode1 | 26657 | 26656 | 1317 |
| sifnode2 | 28003 | 28002 | 1502 |
| sifnode3 | 28005 | 28004 | 1503 |
| sifnode4 | 28007 | 28006 | 1504 |
