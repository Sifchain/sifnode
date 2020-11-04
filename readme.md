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

# reset
sifnoded unsafe-reset-all

# change into the build directory
cd ./build

# scaffold
rake 'genesis:sifnode:scaffold[monkey-bars, bc849774f1daf9047ae68e7aa08f72829b7cdff4@35.166.247.98:26656, http://35.166.247.98:26657/genesis]'

# run
sifnoded start
```

You may also use the following peer addresses (in addition to the one above):

```
f5776c5a16459bce782ccdc438a4f520feab92b5@35.166.247.98:28002
7e308d56cf878d0d88780cc0d75c51b9a3aee203@35.166.247.98:28004
f8a20e44377aefd9557efc595c4a5f4d9308939b@35.166.247.98:28006
```

Additional instructions on standing up Sifnode https://www.youtube.com/watch?v=1kjdjCEcYak&feature=youtu.be&ab_channel=utx0_
Instructions on using Ethereum <> Sifchain cross-chain functionality https://youtu.be/r81NQLxMers

You can also ask questions at our Discord channel - https://discord.com/invite/zZTYnNG
