# Chains

This folder is for setting up and running the backing chains the frontend uses. Currently we use:

- `sif` - sifnode
- `eth` - ganache
- `peggy` - ebrelayer and ../smart-contracts

There is also a folder for quickstart snapshots

- `snapshots`

Each chain has a set of scripts to control it:

| script          | description                                                                                                                                                                                 |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `./build.sh`    | Build the chain so it can be ready to run. For example `sifnoded` gets built from source. For ethereum this runs truffle compile. For peggy this doesn't exist as there is nothing to build |
| `./config.sh`   | Config env vars used in other scripts for this chain.                                                                                                                                       |
| `./migrate.sh`  | Send built changes to the running chain. For example push special contracts etc.                                                                                                            |
| `./pause.sh`    | Stop the chain.                                                                                                                                                                             |
| `./revert.sh`   | Revert the chain to the default snapshot.                                                                                                                                                   |
| `./launch.sh`   | Run the chain from a blank starting point. This calls `start.sh`                                                                                                                            |
| `./snapshot.sh` | Save a snapshot of the chains data to the snapshot location. (Assumes the chain is in the paused state)                                                                                     |
| `./start.sh`    | Start the chain without initializing data.                                                                                                                                                  |
