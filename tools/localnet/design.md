# LOCALNET

`localnet` is a tool allowing us to set up the entire local ecosystem required to build and test new features.

## Motivation

By providing new features and integrating new blockchains, more and more funds are collected in our ecosystem. It means that
responsibility of dev team grows every day. That's why we need good tools helping us to build our confidence that new
feature works well before pull request is even created by the author.

Before IBC era, testing new features locally was pretty easy. The only required component was a single `sifnode`.
Now, after releasing IBC it is much more complicated, and many interconnected components are needed.
For ease of use it seems reasonable to create a tool starting and managing all of them.

## Requirements

To be able to play with all th features of `sifnode` (including IBC) we need these components:
- at least two `sifnodes` running two separate chains
- relayer between them
- local eth testnet
- peggy between eth and sifchain
- block explorer for every chain


Every time we add new supported blockchain `localnet` is extended by adding testnet, relayer  and block explorer.

### Sifchain

To support sifchains, a single `sifnode` per chain is started. `sifnode` should be configured to not mine empty blocks,
to make output in block explorer more relevant and easier to inspect. IIRC there is an option for this in tendermint but
3 years ago when I tested it, consensus was failing because of that. Maybe it has been fixed, have to test it again.

It would be nice to implement some logs filter in `sifnode` to prevent regular `tendermint` consensus-related logs
from being printed in `localnet` unless they are errors or panics. By doing this, relevant logs printed by our app are
more meaningful and take less space on screen.

[`hermes`](https://hermes.informal.systems/) is configured to relay tokens between two sifchains.

[`Big Dipper`](https://github.com/forbole/big-dipper) is used for exploring blocks.

`sifnode` and `hermes` are single-file pieces of software. Both may be run locally easily, so they don't require
containers. On the other hand `Big Dipper` requires frontend and `mongo` instance, that's why it's better to execute them
in containers (ready to use images are available there).

### Ethereum

Single [`geth`](https://github.com/ethereum/go-ethereum) instance is used to create ethereum testnet. It's a single executable so may be used directly,
without container. 

`peggy` is started to maintain token transmission between ethereum and sifchain. Have to discuss it with peggy
team.

[`BlockScout`](https://github.com/blockscout/blockscout) block explorer is started for ethereum. Again, it's better to containerize it.

## Tools

### tmate

[`tmate`](https://tmate.io/) is a great tool based on [`tmux`](https://github.com/tmux/tmux/) allowing us to run many
applications in a single terminal window. Keyboard shortcuts are used to navigate between started apps so all the logs
may be inspected manually. `tmate` is an extension over `tmux`, supporting remote access to sessions created by others.

So instead of running many terminals, one per component, all of them may be started in a single one making our life easier.
At the same time, `tmate` provides the list of all running apps, marking those which exited (due to panic for ex.).

### podman

`podman` is a docker replacement with 100% compatible CLI interface. The difference is that `podman` doesn't require
daemon and supports containers created in user-space (no `root` required). It's much lighter than `docker` and
world is moving toward this.

## CLI interface

`localnet` exposes these commands:

- `localnet install` - downloads all the required components, unpacks, builds them etc.
- `localnet update` - updates components
- `localnet start` - starts all the components in `tmate` session
- `localnet stop` - stops `tmate` session, all components are stopped but may be (re)started again - state is preserved  between restarts  
- `localnet destroy` - drops all the state generated/collected by the components (blocks, keys, etc.)

## Block explorers

Each block explorer has statically assigned IP in `127.0.0.0/24` with frontend listening on port `8080` making it easy to bookmark each of them.

## Future work

Once stuff described above is done, we could think about using this environment to run integration tests.
By investing reasonable effort, tool could be created for sending transactions to different blockchains (sifchain, ethereum etc.)
and querying them to check if result is correct.

We could also use `localnet` in more complex scenarios like starting old version of `sifnode` then stopping it and replacing with a new one.
It would allow us to check if state is correctly maintained between upgrades. 