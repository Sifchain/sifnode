# Connecting to the Sifchain BetaNet. 

## Scaffold and run your node

1. Clone the respository.

```
git clone ssh://git@github.com/Sifchain/sifnode && cd sifnode
```

2. Checkout the mainnet release:

```
git checkout tags/mainnet-genesis
```

3. Build:

```
make clean install
```

4. Generate a mnemonic:

```
rake "keys:generate:mnemonic"
```

5. Scaffold your node:

```
rake "genesis:sifnode:scaffold[sifchain, <moniker>, '<mnemonic>', 0d4981bdaf4d5d73bad00af3b1fa9d699e4d3bc0@44.235.108.41:26656, http://44.235.108.41:26657/genesis]"
```

Where:

|Param|Description|
|-----|----------|
|`<moniker>`|A name for your node.|
|`<mnemonic>`|The mnemonic phrase generated in the previous step.|

This step will also output the keyring password, so please record this and the moniker somewhere secure.

7. Connect:

```
rake "genesis:sifnode:boot[<gas_price>]"
```

Where:

|Param|Description|
|-----|----------|
|`<gas_price>`|The minimum gas price (e.g.: 0.5rowan).|

e.g.:

```
rake "genesis:sifnode:boot[0.5rowan]"
```

and your node will start synchronizing with the network. Please note that this may take several hours or more.

## Verify

You can verify that you're connected by running:

```
sifnodecli q tendermint-validator-set
```

and you should see the following primary validator node/s for Sifchain:

```
validators:
- address: sifvalcons1qv28dvpgue9vlwzncpc75t3l3l7apcee423tem
  pubkey: sifvalconspub1zcjduepqx0jdvxtyx8fd9aff3fr4g946azapz9zujm0mtf8gqx92f0uts90skzrfws
  proposerpriority: -875
  votingpower: 1000
- address: sifvalcons18q4fh3g748d7krq4gnx0lktxlr8l6czzvvp7p6
  pubkey: sifvalconspub1zcjduepqsymd2qtgqtt5vhdzc2dphnr6ulr2eszvyre8rzzgwva232f76h7svryp06
  proposerpriority: 625
  votingpower: 1000
- address: sifvalcons12gwn2fgatqappspxevja8ry65t0rmv8k8xtgme
  pubkey: sifvalconspub1zcjduepqk2jktuqwgvs6k0xy6fg6972pu956476x5wtwtjx4al4gns2wx59sgd4kky
  proposerpriority: -125
  votingpower: 1000
- address: sifvalcons1dv83vy7k0zmezpkzqw7q95tht7fgwj5q2hz97g
  pubkey: sifvalconspub1zcjduepqw8zehuezpsse9f0pe5su0faxteqgvsa7j074s674e0pu8jrf3cyqt9frej
  proposerpriority: -1625
  votingpower: 1000
- address: sifvalcons1wn97nf5e80n0avr736a5p3sqwgf9ng6dgvctn7
  pubkey: sifvalconspub1zcjduepq8fcqvd6x3m74zdckqqsfaq5gdnd9y4ypc724v4alyyl33e5pr7fqqzae69
  proposerpriority: -1875
  votingpower: 1000
- address: sifvalcons1nz9ehhaxw6s79v5c46a2mg7q3a4p2mk8xkwkyj
  pubkey: sifvalconspub1zcjduepqkcsxq9gu5w8j32x9w28vga0d33hcasaa726c22agp892sxu4g5eqrlxm8j
  proposerpriority: 375
  votingpower: 1000
- address: sifvalcons1k6f2u93hjnn9khw5flj9sa6fvf05vzfpsyjjat
  pubkey: sifvalconspub1zcjduepqd3x4ryy8e4wnn6gxzagk3sz355gu725tx0a260xhnaa76pz3whesfyaz6f
  proposerpriority: 3875
  votingpower: 1000
- address: sifvalcons1awm72sjma7fphp0mtsfc6szyg055h2k8hdwsnn
  pubkey: sifvalconspub1zcjduepq7f72cfve29dwn09r8z3hdss9n05hhqpsj653nsrvl2t66mnnfe0s44phkl
  proposerpriority: -375
  votingpower: 1000
```

Congratulations. You are now connected to the network.

## Additional Peers

The following can be used as additional peers on the network:

```
0d4981bdaf4d5d73bad00af3b1fa9d699e4d3bc0@44.235.108.41:26656
bcc2d07a14a8a0b3aa202e9ac106dec0bef91fda@13.55.247.60:26656
663dec65b754aceef5fcccb864048305208e7eb2@34.248.110.88:26656
0120f0a48e7e81cc98829ef4f5b39480f11ecd5a@52.76.185.17:26656
6535497f0152293d773108774a705b86c2249a9c@44.238.121.65:26656
fdf5cffc2b20a20fab954d3b6785e9c382762d14@34.255.133.248:26656
8c240f71f9e060277ce18dc09d82d3bbb05d1972@13.211.43.177:26656
9fbcb6bd5a7f20a716564157c4f6296d2faf5f64@18.138.208.95:26656
```

## Become a Validator

You won't be able to participate in consensus until you become a validator.

1. You will need to have tokens (rowan) on your account in order to become a validator.

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
    --amount="<amount>" \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=sifchain \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=<moniker> \
    --keyring-backend=file
```

Where:

|Param|Description|
|-----|----------|
|`<amount>`|The amount of rowan you wish to stake (the more the better).|
|`<moniker>`|The moniker (name) of your node.|

e.g.:

```
sifnodecli tx staking create-validator \
    --commission-max-change-rate="0.1" \
    --commission-max-rate="0.1" \
    --commission-rate="0.1" \
    --amount="1000000000000000000000rowan" \
    --pubkey=$(sifnoded tendermint show-validator) \
    --moniker=<moniker> \
    --chain-id=sifchain \
    --min-self-delegation="1" \
    --gas-prices="0.5rowan" \
    --from=my-node \
    --keyring-backend=file
```

## Additional Resources

### Endpoints

|Description|Address|
|-----------|-------|
|Block Explorer|https://blockexplorer.sifchain.finance|
|RPC|https://rpc.sifchain.finance|
|API|https://api.sifchain.finance|
