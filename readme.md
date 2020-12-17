# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Prerequisites / Dependencies

- [Ruby 2.6.x](https://www.ruby-lang.org/en/documentation/installation)
- [Golang](https://golang.org/doc/install)
- [jq](https://stedolan.github.io/jq/download/)
- [curl](https://curl.haxx.se/download.html)

## Getting started

### Setup

#### Connect to Sifchain

##### New Validators

1. Ensure you've installed the dependencies listed above.

2. Clone the repository:

```
git clone ssh://git@github.com/Sifchain/sifnode && cd sifnode
```

3. Checkout the latest testnet release:

```
git checkout tags/merry-go-round-1
```

4. Build:

```
make clean install
```

5. Generate a mnemonic:

```
rake "keys:generate:mnemonic"
```

6. Scaffold your node:

```
rake "genesis:sifnode:scaffold[sandpit,<moniker>,'<mnemonic>','',8930d0119e345cb4de10290d83dfdc4d251096b4@52.26.159.121:26656, http://52.26.159.121:26657/genesis]"
```

* Replace `<moniker>` with the moniker (name) of your node. 
* Replace `<mnemonic>` with the mnemonic phrase generated in the previous step.

This step will also output the keyring password, so please record this and the moniker somewhere secure.

7. Connect:

```
sifnoded start
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

##### Existing Validators

1. Checkout the latest testnet release:

```
git fetch && git checkout tags/merry-go-round-1
```

2. Build:

```
make install
```

3. Reset your local state (please take a backup of your keyring first):

```
sifnodecli unsafe-reset-all
```

4. Download the new genesis file:

```
curl http://52.26.159.121:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "8930d0119e345cb4de10290d83dfdc4d251096b4@52.26.159.121:26656,b8fbfc271516fa03018f95ac7511d55ded83c64c@3.248.60.114:26656,513528259bfcb572d9254da8f414822774081de9@13.236.54.177:26656,fc10e8d1ff247929c2d18746c9b0e982d7115eab@54.179.83.98:26656"
```

6. Start your node:

```
sifnoded start
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

##### Verify

You can verify that you're connected by running:

```
sifnodecli q tendermint-validator-set
```

and you should see the following primary validator node/s for Sifchain:

```
validators:
- address: sifvalcons1rud0nxutjt6h8ne7pwhln22t2left09p3py3np
  pubkey: sifvalconspub1zcjduepq0dglw8h63rrpzwfw05f7x3g2024hzr559m20zduy494t2xnewpqssjz3r2
  proposerpriority: -7692
  votingpower: 5000
- address: sifvalcons1esq6rv5sp0d4s906lq9x783eek5zcccg59rz27
  pubkey: sifvalconspub1zcjduepqyha0kq25kmsq4j06x8dkpkludxmrgvj24saxeqsn2ta0tnvz7qlslt9law
  proposerpriority: -1442
  votingpower: 5000
- address: sifvalcons1uwe53gyap7mt8s8gax0vpeprq3z5twl2fnfz2k
  pubkey: sifvalconspub1zcjduepq9xlhvg8pcpg6jmeqxj0v7fxwehs46y5x0ty97phuxkp9xdnf0jpszqyzyz
  proposerpriority: -182
  votingpower: 5000
- address: sifvalcons1aelywdgrqf32rayf4cdhf3hq2zjv2xxzpqnvps
  pubkey: sifvalconspub1zcjduepq36hqmp7u2tlkj8yehdcqucu4l57etgr095kqrge5krl5sj2xu0jq379243
  proposerpriority: 6068
  votingpower: 5000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
8930d0119e345cb4de10290d83dfdc4d251096b4@52.26.159.121:26656
b8fbfc271516fa03018f95ac7511d55ded83c64c@3.248.60.114:26656
513528259bfcb572d9254da8f414822774081de9@13.236.54.177:26656
fc10e8d1ff247929c2d18746c9b0e982d7115eab@54.179.83.98:26656
```

#### Become a Validator

You won't be able to participate in consensus until you become a validator.

1. Reach out to us on [Discord](https://discord.gg/3gQsRvjsRx) to request some tokens.

2. Obtain your node moniker (if you don't already know it):

```
cat ~/.sifnoded/config/config.toml | grep moniker
```

3. Run the following command to become a validator: 

```
sifnodecli tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount=1000000000rowan \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=merry-go-round \
    --min-self-delegation="1" \
    --gas="auto" \
    --from=<moniker> \
    --keyring-backend=file
```

* Replace `<moniker>` with the moniker (name) of your node. 

## Additional Resources

- [Additional instructions on standing up Sifnode](https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_).
- [Instructions on using Ethereum <> Sifchain cross-chain functionality](https://youtu.be/r81NQLxMers).

You can also ask questions on Discord [here](https://discord.com/invite/zZTYnNG).
