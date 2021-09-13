# relayer architecture
After the signature aggregation feature introduced in the relayer, there are two different nodes of relayer according its role in the signature aggregation. The first type of node (witness node) just listen the lock and burn events from Sifnode, and sign against the prophecy id of lock/burn message, then send to Sifnode with signature. Another type of node (relay node) listen the prophecy completed message, and get all aggregated signature, then forward the message to Ethereum in a transaction. We usually need deploy multiple nodes for first type of relayer, but just one for second type of node.

## start witness node
To start the witness node, you need run a subcommand of eblayer, the usage as following:
Use: "init-witness [networkDescriptor] [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMnemonic]",
Example: "ebrelayer init-witness 1 tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 mnemonic --chain-id=peggy".

## start relay node
To start the relay node, you need run a subcommand of eblayer, the usage as following:
Use: "init-relayer [networkDescriptor] [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [validatorMnemonic]",	
Example: "ebrelayer init-relayer 1 tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 mnemonic --chain-id=peggy".

Because in the Sifchain network, we just need single relay node for each target network. In current stage, target network means EVM-based network like Ethereum, BSC, Polygon like so on.

