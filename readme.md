# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Get started

1. Install golang and set GOPATH in your env. https://golang.org/doc/install

2. Setup your node.

Please perform the following:

```
# Clone:
git clone git@github.com:Sifchain/sifnode.git && cd sifnode

# Checkout the latest release:
git checkout tags/monkey-bars-testnet-3

# Build:
make install
```

And then, if you're a new node operator:

```
# Change into the build directory:
cd ./build

# Scaffold a new node on the network (Only run once. Creates your private keys, etc):
rake 'genesis:sifnode:scaffold[monkey-bars, bd17ce50e4e07b5a7ffc661ed8156ac8096f57ce@35.166.247.98:26656, http://35.166.247.98:26657/genesis]'

# run:
sifnoded start
```

Or if you're an existing node operator:

```
# Reset state for existing nodes:
sifnoded unsafe-reset-all

# Update the genesis file:
wget -O ~/.sifnoded/config/genesis.json https://raw.githubusercontent.com/Sifchain/networks/feature/genesis/testnet/monkey-bars-testnet-3/genesis.json

# Change the value for "persistent_peers" in:
vim ~/.sifnoded/config/config.toml 

To:

"bd17ce50e4e07b5a7ffc661ed8156ac8096f57ce@35.166.247.98:26656,f8f5d01fdc73e1b536084bbe42d0a81479f882b3@35.166.247.98:28002,f27548f03a4179b7a4dc3c8a62fcfc5f84be15ff@35.166.247.98:28004,dd35505768be507af3c76f5a4ecdb272537e398f@35.166.247.98:28006"

# Start:
sifnoded start
```

New node operators may also use the following peer addresses:

```
f8f5d01fdc73e1b536084bbe42d0a81479f882b3@35.166.247.98:28002
f27548f03a4179b7a4dc3c8a62fcfc5f84be15ff@35.166.247.98:28004
dd35505768be507af3c76f5a4ecdb272537e398f@35.166.247.98:28006
```

# Verfiy
sifnodecli q tendermint-validator-set
```
blockheight: 285
validators:
- address: sifvalcons1zdzdqejfsn49ntnzynl5aukctaj66mpkt5e6vj
  pubkey: sifvalconspub1zcjduepqa0ams0c3d0n0f4jadfnreh0dxlmknk46x39ngumv7hgkzahgregs99qjpz
  proposerpriority: -15000
  votingpower: 5000
- address: sifvalcons1g7gw770x0qnpd8y86sr6ggwp9t84dvgrff9jaw
  pubkey: sifvalconspub1zcjduepqw4f7kgju5uh4c8vu6zmgwp9f5nmqtgrjaqcm28ymjv7e9p0vqxxq0t6ujv
  proposerpriority: 5000
  votingpower: 5000
- address: sifvalcons1daxr5v7kv2fy6wfzr3nrgajhaa995zz37ag6f4
  pubkey: sifvalconspub1zcjduepqg6ueqp8ev30wskud7jcgaet632c4n8qzq7s8yyp5xmgr43x9x69s397kpy
  proposerpriority: 5000
  votingpower: 5000
- address: sifvalcons1kgc7jvs2azzx8jjm97sn0vwnyk7kl6treeg5t4
  pubkey: sifvalconspub1zcjduepq9skuxclrd5z2q8f8le0xlpe0pd9uz2uvh8e9c4l7as3ul6a7y86qtlhvzr
  proposerpriority: 5000
  votingpower: 5000
```


Additional instructions on standing up Sifnode https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_
Instructions on using Ethereum <> Sifchain cross-chain functionality https://youtu.be/r81NQLxMers

You can also ask questions at our Discord channel - https://discord.com/invite/zZTYnNG
