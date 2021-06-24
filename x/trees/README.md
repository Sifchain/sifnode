# Tree Module

Tree module accepts  a create transaction, which registers the transaction sender as a seller/owner of a tree at a price seller is willing to sell them at. A limit order transaction then allows a potential buyer to register a price at which they are willing to purchase tree.

### Messages

1. MsgCreateTree registers a tree with a property category. Price attribute is in 'rowan' denomination. Whoever is signing the transaction is considered as the seller.
   `MsgCreateTree {
    	Name     string
   	Seller   sdk.AccAddress 
   	Price    sdk.Coins
   	Category string }`

2. MsgBuyTree registers a limit order at which the signer of the message is willing to buy the tree. Buyer has to send as treeId in the transaction.

   `MsgBuyTree {
    	Buyer sdk.AccAddress 
   	Price    sdk.Coins
   	Id string }`

### Keeper

Tree Keeper has access to store keys, codec and other module keepers. This keeper is registered with app module.

`Keeper {
storeKey     sdk.StoreKey
cdc          *codec.Codec
bankKeeper   types.BankKeeper
}`

### Prefixes

`TreeKey       = "tree-value-"
TreeCountKey  = "tree-count-"
OrderKey      = "order-value-"
OrderCountKey = "order-count-"`

### EndBlock

- Limit orders have to be iterated each and every block to check if the price of tree is less than limit order price, then the owner of the tree will be signer who placed limit order.
- These limit orders have to be implemented using an advanced data structure like AVL tree which have very fast search,addition,removal properties.

