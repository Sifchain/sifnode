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
rake "genesis:sifnode:scaffold[sandpit,<moniker>,'<mnemonic>','',32789411cfb76bf5cf5bbe2ee78bb7cf64085805@35.166.205.133:26656, http://35.166.205.133:26657/genesis]"
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
curl http://35.166.205.133:26657/genesis | jq '.result.genesis' > ~/.sifnoded/config/genesis.json
```

5. Update your persistent peers in the file `~/.sifnoded/config/config.toml` so that it reads: 

```
persistent_peers = "32789411cfb76bf5cf5bbe2ee78bb7cf64085805@35.166.205.133:26656,e5da71b75d064bb7649c8d0f7c77063c743427b4@108.128.168.119:26656,93fa246f7af23fb51c824fcb72d04916e3a0dd8c@13.210.44.239:26656,b6a24bfbf50a6dbee83e75353817ab106eaeffa4@18.139.106.51:26656"
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
- address: sifvalcons1r2te04yyr5ryjnrremld67kektdddwhazqtwhk
  pubkey: sifvalconspub1zcjduepqympmagddsq572s8jhumc59uu5q6s2sk3pext8txraamfwgyacyws2nmcky
  proposerpriority: 9375
  votingpower: 5000
- address: sifvalcons1dznuy3cq8ls97swx50z03qjune7gpataz24ec7
  pubkey: sifvalconspub1zcjduepq8rkz0tda3h98s74k67mllg44v04gxpt5l7hhc3jhc9wxglltmq2sft23zz
  proposerpriority: 3125
  votingpower: 5000
- address: sifvalcons1hjej2hpwhtu0jw380cd2ek294vcxvl46ufu3gf
  pubkey: sifvalconspub1zcjduepqwhsxwvr07l6357rvslp9f53trwam5rtfkn7fxs8p5nzjhlnjs4kstju0l7
  proposerpriority: -9375
  votingpower: 5000
- address: sifvalcons17cmmgphys058nhkt8sgpwv633fyj9238t4z8t0
  pubkey: sifvalconspub1zcjduepqgeau6hn4q0v3nt2agnc65yceg42dsjedcqa0d3ehe9zcvw2lpadswvf6vf
  proposerpriority: -3125
  votingpower: 5000
```

Congratulations. You are now connected to the network.

#### Additional Peers

The following can be used as additional peers on the network:

```
32789411cfb76bf5cf5bbe2ee78bb7cf64085805@35.166.205.133:26656
e5da71b75d064bb7649c8d0f7c77063c743427b4@108.128.168.119:26656
93fa246f7af23fb51c824fcb72d04916e3a0dd8c@13.210.44.239:26656
b6a24bfbf50a6dbee83e75353817ab106eaeffa4@18.139.106.51:26656
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
