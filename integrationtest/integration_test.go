package integrationtest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	clpkeeper "github.com/Sifchain/sifnode/x/clp/keeper"
	"github.com/Sifchain/sifnode/x/clp/test"
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	marginkeeper "github.com/Sifchain/sifnode/x/margin/keeper"
	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/proto/tendermint/types"
)

type TestCase struct {
	Setup struct {
		Accounts []banktypes.Balance
		Margin   *margintypes.GenesisState
	}
	Messages []sdk.Msg
}

func TestIntegration(t *testing.T) {
	overwriteFlag := flag.Bool("overwrite", false, "Overwrite test output")
	flag.Parse()

	externalAsset := "cusdc"
	address := "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd"

	tt := []struct {
		name               string
		externalAsset      string
		address            string
		createPoolMsg      clptypes.MsgCreatePool
		openPositionMsg    margintypes.MsgOpen
		swapMsg            clptypes.MsgSwap
		addLiquidityMsg    clptypes.MsgAddLiquidity
		removeLiquidityMsg clptypes.MsgRemoveLiquidity
		closePositionMsg   margintypes.MsgClose
	}{
		{
			externalAsset: externalAsset,
			address:       address,
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
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, app := test.CreateTestAppClpFromGenesis(false, func(app *sifapp.SifchainApp, genesisState sifapp.GenesisState) sifapp.GenesisState {
				// Initialise token registry
				trGs := &tokenregistrytypes.GenesisState{
					Registry: &tokenregistrytypes.Registry{
						Entries: []*tokenregistrytypes.RegistryEntry{
							{Denom: externalAsset, BaseDenom: externalAsset, Decimals: 6, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
							{Denom: "rowan", BaseDenom: "rowan", Decimals: 18, Permissions: []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}},
						},
					},
				}
				bz, _ := app.AppCodec().MarshalJSON(trGs)
				genesisState["tokenregistry"] = bz

				// Initialise wallet
				externalAssetBalance, ok := sdk.NewIntFromString("1000000000000000")
				require.True(t, ok)
				nativeAssetBalance, ok := sdk.NewIntFromString("1000000000000000000000000000")
				require.True(t, ok)
				balances := []banktypes.Balance{
					{
						Address: address,
						Coins: sdk.Coins{
							sdk.NewCoin(externalAsset, externalAssetBalance),
							sdk.NewCoin("rowan", nativeAssetBalance),
						},
					},
				}
				bankGs := banktypes.DefaultGenesisState()
				bankGs.Balances = append(bankGs.Balances, balances...)
				bz, _ = app.AppCodec().MarshalJSON(bankGs)
				genesisState["bank"] = bz

				// Set enabled margin pools
				marginGs := margintypes.DefaultGenesis()
				marginGs.Params.Pools = append(marginGs.Params.Pools, externalAsset)
				bz, _ = app.AppCodec().MarshalJSON(marginGs)
				genesisState["margin"] = bz

				return genesisState
			})

			clpSrv := clpkeeper.NewMsgServerImpl(app.ClpKeeper)
			marginSrv := marginkeeper.NewMsgServerImpl(app.MarginKeeper)

			ctx = ctx.WithBlockHeight(1)
			//app.BeginBlock(abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})
			app.BeginBlocker(ctx, abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})

			// Create pool
			_, err := clpSrv.CreatePool(sdk.WrapSDKContext(ctx), &tc.createPoolMsg)
			require.NoError(t, err)

			endBlock(t, app, ctx, ctx.BlockHeight())

			ctx = ctx.WithBlockHeight(2)
			app.BeginBlocker(ctx, abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})

			// Add liquidity
			_, err = clpSrv.AddLiquidity(sdk.WrapSDKContext(ctx), &tc.addLiquidityMsg)
			require.NoError(t, err)

			endBlock(t, app, ctx, ctx.BlockHeight())
			ctx = ctx.WithBlockHeight(3)
			app.BeginBlocker(ctx, abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})

			// Open position
			_, err = marginSrv.Open(sdk.WrapSDKContext(ctx), &tc.openPositionMsg)
			require.NoError(t, err)

			endBlock(t, app, ctx, ctx.BlockHeight())
			ctx = ctx.WithBlockHeight(4)
			app.BeginBlocker(ctx, abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})

			// Swap
			_, err = clpSrv.Swap(sdk.WrapSDKContext(ctx), &tc.swapMsg)
			require.NoError(t, err)

			endBlock(t, app, ctx, ctx.BlockHeight())
			ctx = ctx.WithBlockHeight(5)
			app.BeginBlocker(ctx, abci.RequestBeginBlock{Header: types.Header{Height: ctx.BlockHeight()}})

			// Remove liquidity
			_, err = clpSrv.RemoveLiquidity(sdk.WrapSDKContext(ctx), &tc.removeLiquidityMsg)
			require.NoError(t, err)

			// Close position
			_, err = marginSrv.Close(sdk.WrapSDKContext(ctx), &tc.closePositionMsg)
			require.NoError(t, err)

			// Check balances
			results := getResults(t, app, ctx, address, externalAsset)
			if *overwriteFlag {
				writeResults(t, results)
			} else {
				expected, err := getExpected()
				require.NoError(t, err)
				require.EqualValues(t, expected, &results)
			}
		})
	}

}

func endBlock(t *testing.T, app *sifapp.SifchainApp, ctx sdk.Context, height int64) {
	app.EndBlocker(ctx, abci.RequestEndBlock{Height: height})
	//app.Commit()

	// Check invariants
	res, stop := app.ClpKeeper.BalanceModuleAccountCheck()(ctx)
	require.False(t, stop, res)
}

type TestResults struct {
	Accounts map[string]sdk.Coins `json:"accounts"`
	Pools    map[string]clptypes.Pool
	LPs      map[string]clptypes.LiquidityProvider
}

func getResults(t *testing.T, app *sifapp.SifchainApp, ctx sdk.Context, address, externalAsset string) TestResults {
	// Lookup account balances
	addr, err := sdk.AccAddressFromBech32(address)
	require.NoError(t, err)
	balances := app.BankKeeper.GetAllBalances(ctx, addr)
	// Lookup pool
	pool, err := app.ClpKeeper.GetPool(ctx, externalAsset)
	require.NoError(t, err)
	// Lookup LP
	lp, err := app.ClpKeeper.GetLiquidityProvider(ctx, externalAsset, address)
	require.NoError(t, err)
	lpKey := clptypes.GetLiquidityProviderKey(externalAsset, address)

	return TestResults{
		Accounts: map[string]sdk.Coins{
			address: balances,
		},
		Pools: map[string]clptypes.Pool{
			externalAsset: pool,
		},
		LPs: map[string]clptypes.LiquidityProvider{
			string(lpKey): lp,
		},
	}
}

func writeResults(t *testing.T, results TestResults) {
	bz, err := json.Marshal(results)
	fmt.Printf("%s", bz)
	require.NoError(t, err)
	err = os.WriteFile("output/results.json", bz, 0644)
	require.NoError(t, err)
}

func getExpected() (*TestResults, error) {
	bz, err := os.ReadFile("output/results.json")
	if err != nil {
		return nil, err
	}
	var results TestResults
	err = json.Unmarshal(bz, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
