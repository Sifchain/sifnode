 # Peggy Smart Contract Architecture

The following will contain the architecture of Peggy's smart contracts.


## Smart Contracts

First we need a list of all smart contracts and their roles in Peggy.

The BridgeToken Contract is an ERC20 token that the BridgeBank creates when a sifchain native asset gets transferred to ethereum for the first time. BridgeTokens are minted by the BridgeBank whenever a user transfers a sifchain native asset to ethereum. BridgeTokens are burned when a user transfers their BridgeTokens to sifchain.

The Valset Contract Maintains a whitelist of validator addresses as well as their respective voting power. These whitelisted addresses are then allowed to submit and sign off on prophecy claims. The CosmosBridge contract will call this contract to check if a user is in the active validator set.

The Oracle Contract Maintains all prophecy claims as well as how many validators have signed off on them. When enough validators sign off on a prophecy, the oracle contract will let you know that the prophecy has become valid. Validators used to call the oracle contract directly, but now they will call the newProphecyClaim() function on CosmosBridge to both sign off on and create the prophecy.

The BridgeBank Contract stores all tokens that get locked in peggy to move across the bridge. Bridgebank will also unlock tokens that are coming from sifchain to ethereum. Bridgebank has admin roles to mint tokens that are sifchain native. Bridgebank will also burn sifchain native tokens that are moving from ethereum back to sifchain. BridgeBank will only unlock or mint at the command of the CosmosBridge contract. Only whitelisted ethereum native tokens are allowed to move across the bridge. Four functions on Bridgebank control the flow of assets to and from sifchain.

These two functions deal with giving users tokens on ethereum in exchange for their sifchain assets or pegged assets.
```
mint(
    address payable _intendedRecipient, // address of the user receiving tokens
    uint256 _amount // amount of tokens they will receive
)

unlock(
    address payable _recipient, // address where tokens will be sent
    string memory _symbol, // symbol of the token to be sent
    uint256 _amount // amount of tokens to send
)
```

These two functions handle moving assets from sifchain to ethereum. The burn function is called for cosmos native assets, and the lock function is called for ethereum native assets. Cosmos native assets are burned because the token will have certain interfaces implemented to burn others tokens with their approvals and this token is essentially just an IOU for an asset on sifchain. This token that gets burned was minted by the BridgeBank contract calling the BridgeToken contract. Ethereum native assets are locked as our contracts do not have the ability to mint tokens such as DAI. These locked ethereum assets then generate IOU's on sifchain so that the user can trade on our DEX. 

Before making a call to either burn or lock, you will have to call the token contract (except when locking eth) and approve the BridgeBank contract to spend the amount of tokens that you are transferring over the bridge. If eth is being transferred over the bridge, there is no approval to make, simply send the eth with the lock transaction and make sure that the amount of eth in wei exactly matches the _amount you pass in to the function. If the amount of eth does not match the _amount variable, the function will revert.
```
burn(
    bytes memory _recipient, // sifaddress
    address _token, // token address
    uint256 _amount uint256
)

lock(
    bytes memory _recipient, // sifaddress
    address _token, // token address, 0x if ethereum
    uint256 _amount // amount of tokens 
)
```

The BridgeBank Contract will maintain a whitelist of both sifchain native ERC20 tokens and ethereum native ERC20 tokens that are whitelisted. The reason that this whitelist needs to exist is because the relayer broadcasts transactions based off of token symbols in the ERC20. This created a condition where an attacker could spin up a new ERC20 token, give it a name of a token that was already in existence on sifchain, lock the token in the BridgeBank, then receive that asset on sifchain. Once they had that sifchain asset, they would transfer it across the bridge back to ethereum where they would receive real tokens and not the fake ones that they had locked up. The whitelist removes the possibility of this behavior occuring.

The CosmosBridge Contract creates and stores new prophecies, but delegates the signing off on prophecies to the Oracle contract. Once the oracle contract says that the prophecy claim has been signed off on, then the CosmosBridge makes a call to the bridgebank to unlock or mint the funds. Only ethereum addresses that are whitelisted in the Valset contract will be able to create new Prophecies and sign off on them.

## Admin API

While designing Peggy, there were some unique design challenges that were faced. The main one being that the ERC20 token, ```erowan```, that represented a claim to rowan on sifchain, was going to be created before peggy and sifchain was going to be deployed to mainnet. Because of this, a solution needed to be implemented to make erowan a cosmos native token that the BridgeBank controlled without using the ```createNewBridgeToken``` function as the token would already exist. To remedy this, an Admin API was built to wire in erowan as a sifchain native even though the BridgeBank did not create the token. This led to the creation of the function ```addExistingBridgeToken```, where an admin can wire in erowan as a sifchain native asset. Only the admin can call this function. Calling this function adds the token to the cosmos native token whitelist so that it can be burned.

Addtionally, a script was created to call this admin api and wire in erowan to the bridgebank. This script is called ```setup_eRowan.js``` and resides inside the scripts folder in the smart-contracts directory. There are instructions on how to use this script inside of the Deployment.md file.

# Smart contract Architecture
Currently, the smart contracts are upgradeable which allows us to fix them should any bugs arise. In time, the control of these smart contracts will be handed over to SifDAO for control or the admin abilities completely removed for full decentralization.

