# EVM flows

Lock and burn operations on tokens and currency native to an EVM chain operate differently than on Cosmos and double
pegged assets. When looking at the design flow of the native EVM tokens and assets, all tokens and currencies on the EVM
side will receive a lock when imported in through the bridgebank or unlock when exported out of the bridgebank. On the
sifnode side, the created assets will be minted when the bridgebank locks or burned when the bridgebank unlocks.

@TODO@ It would be useful to name scenarios in the same way as they are named in user interface

## Lock

This is for moving EVM-native assets (either EVM native currency or ERC20 tokens) from EVM chain to Sifchain.

1. User initiates the scenario by calling `lock()` function on the [BridgeBank](SmartContracts#BridgeBank) smart contract.
1. @TODO@ Describe parameters to lock(). Describe what happens to user's ether/tokens. Mention that the user needs to approve tokens for the BridgeBank.
1. In turn, BridgeBank will emit a [LogLock](Events#LogLock) event on EVM network with the following data:
   - `_from`: Ethereum address that initiated the lock 
   - `_to`: the sifchain address that the imported assets should be credited to (UTF-8 encoded string)
   - `_token`: the token's contract address or the null address for EVM-native currency
   - `_value`: the quantity of asset being transferred (a uint256 representing the smallest unit of the base value)
   - `_nonce`: the current transaction sequence number which is indexed as a topic (_nonce) (this value increments automatically for each `lock`)
   - `_decimals`: the decimals of the asset which defaults to 18 if not found
   - `_symbol`: the symbol of the asset which defaults to empty string if not found
   - `_name`: the name of the asset which defaults to empty string if not found (_name)
   - `_networkDescriptor`: the network descriptor for the chain this asset is on
1. Upon seeing the LogLock event, [relayers](Components#relayer) and [witnesses](Components#witness) will:
   - @TODO@ Is it witness or relayer that does this?
   - [calculate the denom hash](Concepts) and add it along with the other fields to the
   - create a [NewEthBridgeClaim](Events/NewEthBridgeClaim) and broadcast it to sifnode claim, and then broadcast the event to sifnode. The relayers/witnesses will then update the sequence number they stored
for the last processed block.
1. When sifnode handles the lock, it will query the tokens metadata from the token registry;
  - if the data does not currently exist in the token registry, it will write the data from the first claim it sees into the
    token registry/metadata module.
1. Once the prophecy is complete on the sifnode side, the bank module will credit the _to account with the _value of a
   coin with the denomHash as its denom.
   
@TODO@ Describe fees, minting

@TODO@ Include sequence diagram

@TODO@ Include examples for calling lock programmatically (via web3/API/SDK)


## Burn

This is for moving EVM-native assets (either EVM native currency or ERC20 tokens) from Sifnode to their originating EVM chain.
Precondition: assets have been moved to Sifnode with a `lock` scenario.

When users initiate a burn on sifnode for either the native asset or a token on the EVM chain, they export out to the
following steps:

1. The user either mapped from the UI or direct in the cli specify the denomHash for the token to burn.
1. Sifnode pulls up the metadata on the denom hash specified. Sifnode gets the network descriptor it's exporting out to,
verifies the cross-chain fee can be paid, credits the cross-chain fee account, then burns the coins on the user's
account for that denomHash.
1. Sifnode will then emit a new event of type [EventTypeBurn)(Events/EventTypeBurn).
1. Witnesses, while watching for events, observes an EventTypeBurn
1. Witness they will sign the prophecyID of the event with their EVM native keys and then send that signature back to
   sifnode (@TODO@ details).
1. When a relayer sees that m of n signatures are available from the witnesses, it will relay those signatures in a
   call to the sumbitProphecyClaimAggregatedSigs function of the CosmosBridge contract.
1. After submitting the call to the EVM chain, the relayer will increment the sequence number on the sifnode side (@TODO@ how?)
1. Once the smart contract (@TODO@ which?) verifies the signatures are valid, it will unlock the funds for the user.
1. @TODO@ What happens to the funds? How they are transferred to the wallet? Where does it get the address?

@TODO@ Include diagram

@TODO@ Include example for calling burn on the command line or via API/SDK


## Double pegging

Preconditions:
- There was a [lock](#Lock) operation on the token on the native chain before
- A suitable amount of tokens is locked in BridgeBank on the token's native EVM chain for the users's address (@TODO@ can it be more?)
- There is a sufficient amount of twin tokens / twin currency in the user's sifnode account
- The user has enough currency to pay the gas/fees @TODO@ Describe how much, where, and when it is used

Steps:
1. User initiates a burn operation with sifnoded (either from UI or via CLI), specifying:
   - The source sifnode address
   - The destination EVM address
   - The denom hash of the currency  
   - Amount of currency/tokens to be transferred
   ```
   @todo@ sifnoded ... (include the actual command)
   ```
1. [`sifnoded`](Components#Sifnoded) invokes GetTokenMetadata(denomHash) in the Token Registry module to get the token
   metadata. From the obtained metadata, it identifies the network to export to.
    - It sees that the network descriptor is not the token's native network - @TODO@ Is this the key difference from "plain" burn?
    - It calls the token metadata module (@TODO@ Which function) to determine if the token has been exported to that network before
      - If not, it uses the "first time export fee" which is higher because it has to cover smart contract creation on the target network.
        @TODO@ Describe how/where is this fee calculated
      - If yes, it uses the standard cross-chain fee @TODO@ Describe how/where is this fee calculated 
1. sifnoded verifies that the fee can be paid and credits the cross-chain fee account for the fee amount.
1. sifnoded burns the tokens in the user's account for the given `denomHash` and `amount`. @TODO@ more details needed
1. sifnoded emits event of type [EventTypeBurn](Events#EventTypeBurn) with the following data:
    - `prophecyId`: @TODO@ How is it calculated?
    - @TODO@   
1. When witnesses see event of type `EventTypeBurn`, they will
    - sign the prophecyId of the event with their own EVM native private key
    - broadcast the event back to sifnoded (@TODO@ how?)
1. When a relayer (which is listening to events on sifchain) sees the signed event coming from a witness, it will check
   the signature. After seeing at least `m` of `n` valid signatures (@TODO@ how are `m` and `n` set?) for the same
   prophecyID, but signed with different keys, it will call a `submitProphecyClaimAggregatedSigs()` @TODO@ parameters
   function on target chain's [CosmosBridge][SmartContracts#CosmosBridge] smart contract. @TODO@ What happens if the signature is invalid, or if a long time passes? Where are they stored in the meantime and for how long? 
1. Next, the relayer will increment the sequence number of sifnode side. (@TODO@ How? Potential bottleneck?)
1. The `CosmosBridge` smart contract will verify that the signatures are valid. (@TODO@ How?)
    - If not all signatures are valid, it will ignore it 
    - If the signature is valid, it will mint the bridge token representing the asset. (@TODO@ more details needed - we're potentially doing it the first time etc.)
    - The value of `cosmosDenom` in the token will be set to `denomHash` from originating chain 
    - The corresponding amount will be sent to the user's target ETH address.

(@TODO - not clear what this means) When a user burns the asset, it follows the same flow as burning an IBC bridgetoken on the EVM side, but
instead of being credited as an IBC token, it's restored as a standard imported EVM asset.
