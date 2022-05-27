# siftool

To start the local environment:

siftool run-env

It will automatically install Python dependencies upon first use. This command will detect if you are on Peggy1 or
Peggy2 branch, and will start local processes accordingly:
- For Peggy1, it will run ganache-cli, sifnoded and ebrelayer.
- For Peggy2, it will run hardhat, sifnoded and two instances of ebrelayer.

At the moment, the environment consists of Ethereum-compliant local node (ganache/hardhat), one `sifnode` validator and
a Peggy bridge implemented by `ebrelayer` binary.


Original design document: https://docs.google.com/document/d/1IhE2Y03Z48ROmTwO9-J_0x_lx2vIOFkyDFG7BkAIqCk/edit#
