# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Get started

1. Install golang and set GOPATH in your env. https://golang.org/doc/install
2. Install golangci-lint. https://golangci-lint.run/usage/install/#local-installation

```
# clone
git clone git@github.com:Sifchain/sifnode.git && cd sifnode

# checkout the latest release
git checkout tags/monkey-bars-testnet-3

# build
make install

# Reset state for existing nodes
sifnoded  unsafe-reset-all

# change into the build directory
cd ./build

# scaffold
rake 'genesis:sifnode:scaffold[monkey-bars, bd17ce50e4e07b5a7ffc661ed8156ac8096f57ce@35.166.247.98:26656, http://35.166.247.98:26657/genesis]'

# run
sifnoded start
```

You may also use the following peer addresses (in addition to the one above):

```
f8f5d01fdc73e1b536084bbe42d0a81479f882b3@35.166.247.98:28002
f27548f03a4179b7a4dc3c8a62fcfc5f84be15ff@35.166.247.98:28004
dd35505768be507af3c76f5a4ecdb272537e398f@35.166.247.98:28006
```

Additional instructions on standing up Sifnode https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_
Instructions on using Ethereum <> Sifchain cross-chain functionality https://youtu.be/r81NQLxMers

You can also ask questions at our Discord channel - https://discord.com/invite/zZTYnNG
