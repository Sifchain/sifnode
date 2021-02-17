# Peggy architecture

The following document will explain the architecture of peggy from a 10,000 foot view. Other documentation will drill down more into the details of how this works.


## Glossary
Relayer: A piece of middleware that listens to transactions on one chain and submits them to another chain. This relayer will listen to events on both the ethereum and sifchain blockchain.

BridgeBank: A smart contract on ethereum where users will unlock, lock, mint and burn funds to transfer them across the chains.

BridgeToken: An ERC20 token that is created by the BridgeBank to represent a sifchain native asset on ethereum. BridgeTokens are minted by the BridgeBank whenever a user transfers a sifchain native asset to ethereum.

LogLock: An event that is triggered when an ethereum native asset is locked in the BridgeBank contract.

LogBurn: An event that is triggered when a sifchain native asset is burned from the BridgeBank contract.

MsgLock: A sifchain event that signals that a sifchain native asset has been locked.

MsgBurn: A sifchain event that signals that an ethereum native asset has been burned.

ProphecyClaim: A transaction that tells us that a certain amount of coins should be sent to someone. This event is triggered by a lock or burn transaction on one chain, then the relayer submits this prophecy claim to the receiving chain.

Validators: Whitelisted ethereum addresses who submit new prophecy claims.

Valset: A smart contract that stores the whitelist of validators and their powers.

Oracle: A smart contract that stores the current amount of sign off on a given prohpecy claim.

CosmosBridge: A smart contract on ethereum where validators will submit new prophecy claims to unlock or mint assets on ethereum.

Validator Power: The weight a single validator has on voting for a prophecy claim.

Consensus threshold: The percent of validators power that must sign off on a prophecy claim for it to mint or unlock assets on the ethereum side.


## Event Listener

Peggy is a cross chain bridge that currently moves assets from ethereum to sifchain, and from sifchain to ethereum. 

To move assets from ethereum to sifchain, the relayer subscribes to the BridgeBank smart contract deployed on ethereum and listens for the LogLock and LogBurn messages. When the relayer receives lock or burn messages, it waits 50 blocks to ensure that the transaction is still valid, then submits new prophecy claims to sifchain. Other relayers then sign off on that prophecy claim and then once enough relayers have approved the prophecy claim, the assets are minted and sent to that sifchain recipient.

To move assets from sifchain to ethereum, the relayer subscribes to the cosmos chain and listens for MsgLock and MsgBurn event. Once that event is heard, a new ProphecyClaim is submitted to the ethereum CosmosBridge smart contract. Once enough validators sign off on the prophecy claim such that the consensus threshold is met, the funds are unlocked or minted on the ethereum side.

# Smart contracts

Please note that only the whitelisted validators in the valset smart contract can submit or sign off on prophecy claims on ethereum.

On the ethereum side of the world, we maintain smart contracts that will lock and burn funds to move them across the bridge. There are many smart contracts, this is the high level flow from eth to sifchain:
1. User locks up funds in BridgeBank smart contract
2. Relayer hears the event generated from the BridgeBank contract
3. Relayer submits a new prophecy claim to mint assets on the cosmos side of the world.

When a user transfers value from sifchain to ethereum this is what the flow looks like:
1. User locks or burns assets on the cosmos side of the world.
2. Relayer hears this transaction and submits a new prophecy claim to the CosmosBridge smart contract
3. Other relayers sign off on this transaction.
4. Once enough relayers sign off on this prophecy claim and the consensus threshold is reached, one of two things happen. If this was a sifchain native asset being moved across the bridge, then we will mint assets for that user through the BridgeBank. If this asset being moved across the bridge was an ethereum native asset, then the BridgeBank will unlock those funds and send them to the user specified in the prohpecy claim.

# Smart contract Architecture
Currently, the smart contracts are upgradeable which allows us to fix them should any bugs arise. In time, the control of these smart contracts will be handed over to SifDAO for control or the admin abilities completely removed for full decentralization.