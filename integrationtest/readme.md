# Integration Test Framework

This test framework takes a test case,
which specifies how to setup the chain from scratch,
and a sequence of different types of messages to execute.

Results of the test case are compared to previous results stored.

Setup: 
* Token registry entries
* Account balances
* Margin Params
* Clp Params
* Admin accounts

Message spec:

```go
createPoolMsg: clptypes.MsgCreatePool{
    Signer:              address,
    ExternalAsset:       &clptypes.Asset{Symbol: "cusdc"},
    NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
    ExternalAssetAmount: sdk.NewUintFromString("1000000000"),             // 1000cusdc
},
openPositionMsg: margintypes.MsgOpen{
    Signer:           address,
    CollateralAsset:  "rowan",
    CollateralAmount: sdk.NewUintFromString("10000000000000000000"), // 10rowan
    BorrowAsset:      externalAsset,
    Position:         margintypes.Position_LONG,
    Leverage:         sdk.NewDec(2),
},
swapMsg: clptypes.MsgSwap{
    Signer:             address,
    SentAsset:          &clptypes.Asset{Symbol: externalAsset},
    ReceivedAsset:      &clptypes.Asset{Symbol: clptypes.NativeSymbol},
    SentAmount:         sdk.NewUintFromString("10000"),
    MinReceivingAmount: sdk.NewUint(0),
},
addLiquidityMsg: clptypes.MsgAddLiquidity{
    Signer:              address,
    ExternalAsset:       &clptypes.Asset{Symbol: externalAsset},
    NativeAssetAmount:   sdk.NewUintFromString("1000000000000000000000"), // 1000rowan
    ExternalAssetAmount: sdk.NewUintFromString("1000000000"),
},
removeLiquidityMsg: clptypes.MsgRemoveLiquidity{
    Signer:        address,
    ExternalAsset: &clptypes.Asset{Symbol: externalAsset},
    WBasisPoints:  sdk.NewInt(5000),
    Asymmetry:     sdk.NewInt(0),
},
closePositionMsg: margintypes.MsgClose{
    Signer: address,
    Id:     1,
},
```

