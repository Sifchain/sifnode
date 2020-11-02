# sifnode

**sifnode** is a blockchain application built using Cosmos SDK and Tendermint and generated with [Starport](https://github.com/tendermint/starport).

## Get started

1. Install golang and set GOPATH in your env. https://golang.org/doc/install
2. Install golangci-lint. https://golangci-lint.run/usage/install/#local-installation

```
# clone
git clone git@github.com:Sifchain/sifnode.git && cd sifnode

# ensure you're on the develop branch
git branch

# build
make && make install

# change into the build directory
cd ./build

# scaffold
rake 'genesis:sifnode:scaffold[monkey-bars, 190cb35265860f182e35a3bceeb297082858eebd@35.166.247.98:26656, http://35.166.247.98:26657/genesis]'

# run
sifnoded start
```

You may also use the following peer addresses (in addition to the one above):

```
c91f8892b391b4af7bf1e17ea9b7f44027913002@35.166.247.98:28002
ca3341e335cba0bc8780d49986a6b49f11f804a9@35.166.247.98:28004
a95912df79d1a256c339267af50a8222c1c3185d@35.166.247.98:28006
```

Additional instructions on standing up Sifnode https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_
Instructions on using Ethereum <> Sifchain cross-chain functionality https://youtu.be/r81NQLxMers

You can also ask questions at our Discord channel - https://discord.com/invite/zZTYnNG
