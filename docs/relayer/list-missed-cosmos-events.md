# sub-command listMissedCosmosEvent

## algorithm behind the sub-command
Sometimes, we need check if the ebrelayer missed some cosmos events yesterday or in latest days. The command will check the cosmos side, collect all lock/burn events then search the prophecy transaction in Ethereum. All cosmos events not processed by the ebrelayer will be output in console. To get the specific block height, we compute it by one Ethereum block every 15 seconds and one Cosmos block every 6 seconds.

## Usage format and example
ebrelayer listMissedCosmosEventCmd [tendermintNode] [web3Provider] [bridgeRegistryContractAddress] [ebrelayerEthereumAddress] [days]

ebrelayer listMissedCosmosEventCmd tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 0x627306090abaB3A6e1400e9345bC60c78a8BEf57 1

### output as following:

listMissedEvent.go:116: missed cosmos event: 
Claim Type: burn
Cosmos Sender: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
Cosmos Sender Sequence: 1
Ethereum Recipient: 0x11111111262B236c9AC9A9A8C8e4276B5Cf6b2C9
Symbol: eth
Amount: 10000

listMissedEvent.go:116: missed cosmos event: 
Claim Type: lock
Cosmos Sender: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
Cosmos Sender Sequence: 2
Ethereum Recipient: 0x11111111262B236c9AC9A9A8C8e4276B5Cf6b2C9
Symbol: rowan
Amount: 10

listMissedEvent.go:116: missed cosmos event: 
Claim Type: burn
Cosmos Sender: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
Cosmos Sender Sequence: 3
Ethereum Recipient: 0x11111111262B236c9AC9A9A8C8e4276B5Cf6b2C9
Symbol: eth
Amount: 10000

listMissedEvent.go:116: missed cosmos event: 
Claim Type: lock
Cosmos Sender: sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5
Cosmos Sender Sequence: 4
Ethereum Recipient: 0x11111111262B236c9AC9A9A8C8e4276B5Cf6b2C9
Symbol: rowan
Amount: 10
